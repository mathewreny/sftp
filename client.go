package sftp

import (
	"fmt"
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
	flush      chan *Buffer

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
		ingress:    make(chan *Buffer),
		egress:     make(chan *Buffer),
		flush:      make(chan *Buffer),
	}
	c.closers[0], c.closers[1] = wc, s
	go c.loopFlush(wc)
	go c.loopEgress()
	go c.loopIngress(r)
	go c.loopMultiplex(wc)
	return c, nil
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

type bufqueue struct {
	head, tail int
	queue      []*Buffer
}

func newBufqueue() bufqueue {
	return bufqueue{queue: make([]*Buffer, 30)}
}

func (bufq *bufqueue) NextIndex(i int) int {
	return (i + 1) % len(bufq.queue)
}

func (bufq *bufqueue) Peek() *Buffer {
	return bufq.queue[bufq.head]
}

func (bufq *bufqueue) Pop() (removed *Buffer) {
	removed = bufq.queue[bufq.head]
	bufq.queue[bufq.head] = nil
	if removed != nil {
		bufq.head = bufq.NextIndex(bufq.head)
	}
	return
}

func (bufq *bufqueue) Push(buf *Buffer) {
	bufq.queue[bufq.tail] = buf
	bufq.tail = bufq.NextIndex(bufq.tail)
	if bufq.tail == bufq.head {
		bufq.Grow()
	}
}

func (bufq *bufqueue) Grow() {
	head := bufq.head
	qlen := len(bufq.queue)
	q := make([]*Buffer, 2*qlen)
	if head == 0 {
		copy(q, bufq.queue)
	} else {
		copy(q, bufq.queue[head:])
		copy(q[qlen-head:], bufq.queue[:head])
	}
	bufq.head = 0
	bufq.tail = qlen
	bufq.queue = q
}

func (c *Client) loopEgress() {
	bufq := newBufqueue()
	for {
		empty := nil == bufq.Peek()
		if empty {
			bufq.Push(<-c.egress)
		} else {
			select {
			case buf := <-c.egress:
				bufq.Push(buf)
			case c.flush <- bufq.Peek():
				bufq.Pop()
			}
		}
	}
}

func (c *Client) loopFlush(w io.Writer) {
	for buf := range c.flush {
		_, err := CopyPacket(w, buf)
		if err != nil {
			panic("CLOSED C IN EGRESS????")
			c.Close()
		}
		bufPool.Put(buf)
	}
}

func (c *Client) loopMultiplex(w io.Writer) {
	cs := make(map[uint32]chan<- *Buffer) // no locks!
	for open := true; open; {

		select {

		case <-c.done:
			open = false
			panic("DONE??????")

		case r := <-c.responders:
			fmt.Println("Got responder")
			if r.buf == nil {
				panic("NIL RESPONDER???????")
				break
			}
			id := r.buf.PeekId()
			if id < 3 {
				panic("ID less than three!???")
			} else if _, found := cs[id]; found {
				panic("Found responder in responders?????")
				r.ch <- NewStatus(STATUS_BAD_MESSAGE).Buffer()
			} else {
				cs[id] = r.ch
				fmt.Println("Sending buffer to egress.")
				c.egress <- r.buf
				fmt.Println("Sent buffer to egress.")
			}

		case buf := <-c.ingress:
			if buf == nil {
				panic("NIL INGRESS???????")
			}
			id := buf.PeekId()
			if id < 3 {
				panic("ID IS LESS THAN THREE!!!???")
			} else if ch, found := cs[id]; !found {
				panic("NOT FOUND ?????")
				c.Close()
			} else {
				ch <- buf
				delete(cs, id)
			}
		}
	}
	// clean up and exit cleanly for all outstanding requests.
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
	fmt.Println("IN SEND FUNCTION")
	select {
	case c.responders <- responder{buf, resp}:
		// It's impossible for a sent responder to race with c.done. The c.responders
		// queue is unbuffered. Therefore any responder sent WILL be handled by the
		// loopResponder function. The slight performance penalty is worth the safety.
	case <-c.done:
		resp <- NewStatus(STATUS_NO_CONNECTION).Buffer()
	}
	response = resp
	fmt.Println("OUT OF SEND FUNCTION")
	return
}
