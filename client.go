package sftp

import (
	"io"
	"sync"
	"sync/atomic"
)

// Interface composed of functions used from golang.org/x/crypto/ssh#Session
type Session interface {
	Close() error
	RequestSubsystem(subsystem string) error
	StdinPipe() (io.WriteCloser, error)
	StdoutPipe() (io.Reader, error)
}

type responder struct {
	buf *Buffer
	ch  chan<- *Buffer
}

type Client struct {
	idgen uint32

	done       chan struct{}
	responders chan responder
	ingress    chan *Buffer
	egress     chan *Buffer

	closeOnce sync.Once
	closers   [2]io.Closer
}

func (c *Client) nextPacketId() uint32 {
	return atomic.AddUint32(&c.idgen, 1)
}

func (c *Client) Close() (err error) {
	c.closeOnce.Do(func() {
		close(c.done)
		for _, x := range c.closers {
			if err != nil {
				x.Close()
			} else {
				err = x.Close()
			}
		}
	})
	return
}

// The provided io.ReadWriteCloser should be a properly set up ssh session.
// It's takes a very general interface for testing purposes. Most users will
// provide a Conn object to this constructor.
func NewClient(s Session) (*Client, error) {
	err := s.RequestSubsystem("sftp")
	if err != nil {
		return nil, err
	}
	r, err := s.StdoutPipe()
	if err != nil {
		return nil, err
	}
	wc, err := s.StdinPipe()
	if err != nil {
		return nil, err
	}
	c := &Client{
		idgen:      2, // set idgen to 2 so that the Init function sends the correct *version*.
		done:       make(chan struct{}),
		responders: make(chan responder),
		ingress:    make(chan *Buffer, 20),
		egress:     make(chan *Buffer, 20),
	}
	c.closers[0], c.closers[1] = wc, s
	go c.loopEgress(wc)
	go c.loopIngress(r)
	go c.loopMultiplex()
	return c, nil
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
		buf := NewBuffer()
		_, err := CopyPacket(buf, r)
		if err != nil {
			c.Close()
			bufPool.Put(buf)
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
	cs := make(map[uint32]chan<- *Buffer) // no locks!
	for open := true; open; {
		select {

		case _, open = <-c.done:

		case r := <-c.responders:
			id := r.buf.PeekId()
			if _, found := cs[id]; found || id < 3 {
				bufPool.Put(r.buf)
				r.ch <- NewStatus(STATUS_BAD_MESSAGE).Buffer()
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
	for _, ch := range cs {
		ch <- NewStatus(STATUS_CONNECTION_LOST).Buffer()
	}
}

// Concurrently send the packet buffer over the client. When a response is
// returned by the server, it will be delivered via the "response" channel. If
// the client is closed or lost before sending/receiving, or something goes
// wrong on the server, a status packet will be sent over the channel. Status
// packets with the code STATUS_OK indicate a successful action.
func (c *Client) send(buf *Buffer) (response <-chan *Buffer) {
	//	if err := validOutgoingPacketHeader(buf); err != nil {
	//		errstring = err.Error()
	//		errstatus = STATUS_BAD_MESSAGE
	//		goto Done
	//	}
	resp := make(chan *Buffer, 1)
	select {
	case c.responders <- responder{buf, resp}:
		// It's impossible for a sent responder to race with c.done. The c.responders
		// queue is unbuffered. Therefore any responder sent WILL be handled by the
		// loopResponder function. The slight performance penalty is worth the safety.
	case <-c.done:
		bufPool.Put(buf)
		resp <- NewStatus(STATUS_NO_CONNECTION).Buffer()
	}
	return resp
}
