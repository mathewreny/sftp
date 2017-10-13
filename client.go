package sftp

import (
	"io"
	"sync"
	"sync/atomic"
)

type responder struct {
	buf *PacketBuffer
	ch  chan<- *PacketBuffer
}

type Client struct {
	idgen uint32

	done       chan struct{}
	responders chan responder
	ingress    chan *PacketBuffer
	egress     chan *PacketBuffer

	closeOnce sync.Once
	closeConn func() error
	debug     bool
}

func (c *Client) nextPacketId() uint32 {
	return atomic.AddUint32(&c.idgen, 1)
}

func (c *Client) Close() (err error) {
	c.closeOnce.Do(func() {
		close(c.done)
		err = c.closeConn()
	})
	return
}

// The provided io.ReadWriteCloser should be a properly set up ssh session.
// It's takes a very general interface for testing purposes. Most users will
// provide a Conn object to this constructor.
func NewClient(conn io.ReadWriteCloser) (c *Client) {
	c = &Client{
		idgen:      2, // set idgen to 2 so that the Init function sends the correct *version*.
		done:       make(chan struct{}),
		responders: make(chan responder),
		ingress:    make(chan *PacketBuffer, 20),
		egress:     make(chan *PacketBuffer, 20),
		closeConn:  conn.Close,
	}
	go c.loopEgress(conn)
	go c.loopIngress(conn)
	go c.loopMultiplex()
	return
}

func (c *Client) loopEgress(w io.Writer) {
	for buf := range c.egress {
		_, err := io.Copy(w, buf)
		bufPool.Put(buf)
		if err != nil {
			c.Close()
		}
	}
}

func (c *Client) loopIngress(r io.Reader) {
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

func (c *Client) loopMultiplex() {
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
	close(c.egress)
	for id, ch := range cs {
		ch <- BufferStatus(nil, STATUS_CONNECTION_LOST, id, "Done", "en")
	}
}

// Concurrently send the packet buffer over the client. When a response is
// returned by the server, it will be delivered via the "response" channel. If
// the client is closed or lost before sending/receiving, or something goes
// wrong on the server, a status packet will be sent over the channel. Status
// packets with the code STATUS_OK indicate a successful action.
func (c *Client) Send(buf *PacketBuffer) (response <-chan *PacketBuffer) {
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
		resp <- BufferStatus(buf, STATUS_NO_CONNECTION, buf.PeekId(), "Client is closed.", "en")
	}
	return resp
}
