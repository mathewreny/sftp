package sftp

import (
	"errors"
	"io"
	"sync/atomic"
)

type responder struct {
	id uint32
	ch chan<- *PacketBuffer
}

type Conn struct {
	idgen uint32

	done       chan struct{}
	responders chan responder
	ingress    chan *PacketBuffer
	egress     chan *PacketBuffer

	debug bool
}

func (c *Conn) nextPacketId() uint32 {
	return atomic.AddUint32(&c.idgen, 1)
}

func (c *Conn) Close() error {
	close(c.done)
	return nil
}

func NewConn(r io.Reader, w io.Writer) (c *Conn) {
	c = &Conn{
		idgen:      2, // set idgen to 2 so that the Init function sends the correct *version*.
		done:       make(chan struct{}),
		responders: make(chan responder),
		ingress:    make(chan *PacketBuffer, 20),
		egress:     make(chan *PacketBuffer, 20),
	}
	go c.loopResponder()
	go c.loopEgress(w)
	go c.loopIngress(r)
	return
}

func (c *Conn) loopEgress(w io.Writer) {
	for {
		select {
		case <-c.done:
			return
		case buf := <-c.egress:
			_, err := io.Copy(w, buf)
			if err != nil {
				c.Close()
			}
			bufPool.Put(buf)
		}
	}
}

func (c *Conn) loopIngress(r io.Reader) {
	for {
		buf, err := ConsumePacket(r)
		if err != nil || buf == nil {
			c.Close()
			return
		}
		select {
		case c.ingress <- buf:
		case <-c.done:
			return
		}
	}
}

func (c *Conn) loopResponder() {
	cs := make(map[uint32]chan<- *PacketBuffer)
	for {
		select {

		case <-c.done:
			// clean up and exit cleanly for all outstanding requests.
			for id, ch := range cs {
				ch <- BufferStatus(nil, STATUS_NO_CONNECTION, id, "Done", "en")
			}
			return

		case r := <-c.responders:
			if ch, found := cs[r.id]; found && ch != nil {
				close(ch)
			}
			cs[r.id] = r.ch

		case buf := <-c.ingress:
			id := buf.PeekId()
			if ch, found := cs[id]; found && ch != nil {
				ch <- buf
				delete(cs, id)
			}
		}
	}
}

// Concurrently send the packet buffer over the connection. When a response is
// returned by the server, it will be delivered via the "response" channel. If
// the connection is closed or lost before sending/receiving, or something goes
// wrong on the server, a status packet will be sent over the channel. Status
// packets with the code STATUS_OK indicate a successful action.
//
// Sending a packet with an invalid SFTP 3 header causes this function to panic.
// Do not send bad packets. 99% of users should not use this function. Its here
// for the 1% of users who want to form their own raw packets using buffers.
func (c *Conn) Send(buf *PacketBuffer) (response <-chan *PacketBuffer) {
	if err := validOutgoingPacketHeader(buf); err != nil {
		// Don't send bad packets.
		panic(err)
	}

	resp := make(chan *PacketBuffer, 1)

	select {
	case c.responders <- responder{buf.PeekId(), resp}:
		// It's impossible for a sent responder to race with c.done. The c.responders
		// queue is unbuffered. Therefore any responder sent WILL be handled by the
		// loopResponder function. The slight performance penalty is worth the safety.
	case <-c.done:
		resp <- BufferStatus(buf, STATUS_NO_CONNECTION, buf.PeekId(), "Conn is closed.", "en")
		bufPool.Put(buf)
		return resp
	}

	select {
	case c.egress <- buf:
	case <-c.done:
		resp <- BufferStatus(buf, STATUS_CONNECTION_LOST, buf.PeekId(), "Conn is closed.", "en")
		bufPool.Put(buf)
	}
	return resp

}

// Ugly validation logic condensed into a convenience function. This function
// should be used as a gateway. It checks for problems that might cause a panic
// in later uses of the buffer.
func validIncomingPacketHeader(buf *PacketBuffer) error {
	if buf == nil {
		return errors.New("Buffer is nil")
	}
	if !buf.IsValidLength() {
		return errors.New("Packet length doesn't match buffer.")
	}
	id := buf.PeekId()
	if id == 0 {
		return errors.New("Packet id must not be zero.")
	}
	switch buf.PeekType() {
	case FXP_STATUS,
		FXP_HANDLE,
		FXP_DATA,
		FXP_NAME,
		FXP_ATTRS,
		FXP_VERSION,
		FXP_EXTENDEDREPLY:
		return nil
	default:
		return errors.New("Packet is not a valid server response type.")
	}
}

// The packet header is the first 9 bytes in a packet. There are many ways this
// can be wrong. This function does not check packets for correctness based on
// any state. For instance, if an INIT packet is called twice, this function
// has no way of knowing. All the uglyness is condensed into one function.
func validOutgoingPacketHeader(buf *PacketBuffer) error {
	if buf == nil {
		return errors.New("Buffer is nil.")
	}
	if !buf.IsValidLength() {
		return errors.New("Packet buffer length is invalid.")
	}
	id := buf.PeekId()
	if 0 == id {
		return errors.New("Packet id must be greater than zero.")
	}
	t := buf.PeekType()
	if (t > 2 && t <= 20) || t == FXP_INIT || t == FXP_EXTENDED {
		return nil
	}
	switch t {
	case FXP_STATUS,
		FXP_HANDLE,
		FXP_DATA,
		FXP_NAME,
		FXP_ATTRS,
		FXP_VERSION,
		FXP_EXTENDEDREPLY:
		return errors.New("Packet is not an outgoing packet type.")
	default:
		return errors.New("Packet type is unknown.")
	}
}
