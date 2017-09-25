package sftp

import (
	"sync/atomic"
)

type Conn struct {
	idgen uint32
}

func (c *Conn) generatePacketId() uint32 {
	return atomic.AddUint32(&c.idgen, 1)
}

// TODO
func (c *Conn) send(id uint32, buf *PacketBuffer) {
	// send the packet asynchronously.

	// Do this when the packet is sent.
	bufPool.Put(buf)
}

// TODO
func (c *Conn) statusResponse(id uint32) chan error {
	return nil
}

// TODO
func (c *Conn) handleResponse(id uint32) chan string {
	return nil
}

// TODO
func (c *Conn) dataResponse(id uint32) chan []byte {
	return nil
}

// TODO
func (c *Conn) nameResponse(id uint32) chan []FxpName {
	return nil
}

// TODO
func (c *Conn) attrsResponse(id uint32) chan FxpAttrs {
	return nil
}

// TODO
func (c *Conn) extendedReplyResponse(id uint32) chan []byte {
	return nil
}
