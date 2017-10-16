package sftp

import (
	//	"fmt"
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

// Creates a packet ID that hasn't been used by the provided client.
// Behind the hood, this function simply increments a counter within
// the Client struct using the sync/atomic library.
//
// Don't forget to use this function when creating the Init packet.
func NextId(c *Client) PacketId {
	return PacketId(atomic.AddUint32(&c.idgen, 1))
}

// SFTP version 3 client. This design is the lowest level a client API can get. Unlike many others,
// it is an inherently functional design. Init must be the first packet sent by users of this
// client. This allows the librarymust be the
// sent must be the Init packet.
// to a SSH session
//
// Its architecture is
// designed sl
type Client struct {
	// number used to create packet ids.
	idgen uint32

	done       chan struct{}
	responders chan responder
	ingress    chan *Buffer
	egress     chan *Buffer
	queue      chan *Buffer

	closeOnce sync.Once
	closers   []io.Closer
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
func NewClient(s Session) (c *Client, err error) {
	err = s.RequestSubsystem("sftp")
	if err != nil {
		return
	}
	var r io.Reader
	var wc io.WriteCloser
	if r, err = s.StdoutPipe(); err != nil {
		return
	}
	if wc, err = s.StdinPipe(); err != nil {
		return
	}
	defer func() {
		c.closers = append(c.closers, s)
	}()
	return NewRawClient(r, wc, 0, 0, 20, 0)
}

// Useful for testing purposes. The integers correspond to channel buffer sizes. Read the source
// code to see which channels they affect. This function was made public to encourage tuning.
func NewRawClient(r io.Reader, wc io.WriteCloser, nr, ne, ni, nq int) (c *Client, err error) {
	c = &Client{
		done:       make(chan struct{}),
		responders: make(chan responder, nr),
		egress:     make(chan *Buffer, ne),
		ingress:    make(chan *Buffer, ni),
		queue:      make(chan *Buffer, nq),
		closers:    make([]io.Closer, 1, 2),
	}
	c.closers[0] = wc
	go c.loopMultiplex()
	go c.loopQueue()
	go c.loopEgress(wc)
	go c.loopIngress(r)
	// Important: idgen is set to 2 so that the Init function sends the correct *version* number.
	// This design choice greatly simplifies the code and allows for extensions. Do not remove it.
	c.idgen = 2
	return c, err
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

func newBufqueue() *bufqueue {
	return &bufqueue{queue: make([]*Buffer, 30)}
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

func (c *Client) loopQueue() {
	bufq := newBufqueue()
	defer func() {
		// clean up.
		for x := bufq.Pop(); x != nil; x = bufq.Pop() {
			bufPool.Put(x)
		}
	}()
	for {
		if bufq.Peek() == nil {
			select {
			case <-c.done:
				return
			case buf := <-c.queue:
				bufq.Push(buf)
			}
		} else {
			select {
			case <-c.done:
				return
			case buf := <-c.queue:
				bufq.Push(buf)
			case c.egress <- bufq.Peek():
				bufq.Pop()
			}
		}
	}
}

func (c *Client) loopEgress(w io.Writer) {
	for {
		select {
		case <-c.done:
			return
		case buf := <-c.egress:
			_, err := CopyPacket(w, buf)
			if err != nil {
				c.Close()
			}
			bufPool.Put(buf)
		}
	}
}

func (c *Client) loopMultiplex() {
	cs := make(map[uint32]chan<- *Buffer) // no locks!
	for open := true; open; {

		select {

		case <-c.done:
			open = false

		case r := <-c.responders:
			id := uint32(PeekId(r.buf))
			if _, found := cs[id]; found {
				bufPool.Put(r.buf)
				r.ch <- NewStatus(STATUS_BAD_MESSAGE).Buffer()
			} else {
				cs[id] = r.ch
				c.queue <- r.buf
			}

		case buf := <-c.ingress:
			id := uint32(PeekId(buf))
			if ch, found := cs[id]; !found {
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
	resp := make(chan *Buffer, 1)
	response = resp
	select {
	case c.responders <- responder{buf, resp}:
		// It's impossible for a sent responder to race with c.done. The c.responders
		// queue is unbuffered. Therefore any responder sent WILL be handled by the
		// loopResponder function. The slight performance penalty is worth the safety.
	case <-c.done:
		bufPool.Put(buf)
		resp <- NewStatus(STATUS_NO_CONNECTION).Buffer()
	}
	return
}
