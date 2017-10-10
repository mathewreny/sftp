package sftp

import (
	"errors"
	"io"
	"sync"
	"sync/atomic"
)

type responder struct {
	buf *PacketBuffer
	ch  chan<- *PacketBuffer
}

type Conn struct {
	idgen uint32

	done       chan struct{}
	responders chan responder
	ingress    chan *PacketBuffer
	egress     chan *PacketBuffer

	closeOnce sync.Once
	closeRW   func() error
	debug     bool
}

func (c *Conn) nextPacketId() uint32 {
	return atomic.AddUint32(&c.idgen, 1)
}

func (c *Conn) Close() (err error) {
	c.closeOnce.Do(func() {
		close(c.done)
		err = c.closeRW()
	})
	return
}

// The provided io.ReadWriteCloser should be a properly set up ssh session.
// It's provided for both testing and highly customized ssh connections.
func NewConn(rwc io.ReadWriteCloser) (c *Conn) {
	c = &Conn{
		idgen:      2, // set idgen to 2 so that the Init function sends the correct *version*.
		done:       make(chan struct{}),
		responders: make(chan responder),
		ingress:    make(chan *PacketBuffer, 20),
		egress:     make(chan *PacketBuffer, 20),
		closeRW:    rwc.Close,
	}
	go c.loopEgress(rwc)
	go c.loopIngress(rwc)
	go c.loopMultiplex()
	return
}

func (c *Conn) loopEgress(w io.Writer) {
	for buf := range c.egress {
		_, err := io.Copy(w, buf)
		bufPool.Put(buf)
		if err != nil {
			c.Close()
		}
	}
}

func (c *Conn) loopIngress(r io.Reader) {
	for {
		buf, err := ConsumePacket(r)
		if err != nil {
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

func (c *Conn) loopMultiplex() {
	cs := make(map[uint32]chan<- *PacketBuffer) // no locks!
	for open := true; open; {
		select {

		case _, open = <-c.done:

		case r := <-c.responders:
			if id := r.buf.PeekId(); id < 3 {
				r.ch <- BufferStatusCode(r.buf, STATUS_BAD_MESSAGE)
			} else if _, found := cs[id]; found {
				r.ch <- BufferStatusCode(r.buf, STATUS_BAD_MESSAGE)
			} else {
				cs[id] = r.ch
				c.egress <- r.buf
			}

		case buf := <-c.ingress:
			if id := buf.PeekId(); id < 3 {
				c.Close()
			} else if ch := cs[id]; ch == nil {
				c.Close()
			} else {
				ch <- buf
				delete(cs, id)
				continue
			}
		}
	}
	// clean up and exit cleanly for all outstanding requests.
	close(c.responders)
	close(c.ingress)
	close(c.egress)
	for id, ch := range cs {
		ch <- BufferStatus(nil, STATUS_CONNECTION_LOST, id, "Done", "en")
	}
}

// Concurrently send the packet buffer over the connection. When a response is
// returned by the server, it will be delivered via the "response" channel. If
// the connection is closed or lost before sending/receiving, or something goes
// wrong on the server, a status packet will be sent over the channel. Status
// packets with the code STATUS_OK indicate a successful action.
func (c *Conn) Send(buf *PacketBuffer) (response <-chan *PacketBuffer) {
	//	if err := validOutgoingPacketHeader(buf); err != nil {
	//		errstring = err.Error()
	//		errstatus = STATUS_BAD_MESSAGE
	//		goto Done
	//	}
	resp := make(chan *PacketBuffer, 1)
	select {
	case c.responders <- responder{buf, resp}:
		// It's impossible for a sent responder to race with c.done. The c.responders
		// queue is unbuffered. Therefore any responder sent WILL be handled by the
		// loopResponder function. The slight performance penalty is worth the safety.
	case <-c.done:
		resp <- BufferStatus(buf, STATUS_NO_CONNECTION, buf.PeekId(), "Conn is closed.", "en")
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
	if 3 > id {
		return errors.New("Packet id must be greater than two.")
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
	if 3 > id {
		return errors.New("Packet id must be greater than two.")
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
