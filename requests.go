package sftp

// Automatically generated file. Do not touch.

import "errors"

// Extended is an optinal list of `name` and `data` that you want to support.
func (c *Conn) Init(extended [][2]string) ([][2]string, error) {
	id := c.nextPacketId()
	var pktLen uint32 = 4 + 1 + 4
	buf := NewBuffer()
	buf.Grow(4 + pktLen)
	buf.WriteUint32(pktLen)
	buf.WriteByte(FXP_INIT)
	buf.WriteUint32(id)
	reply := c.Send(buf)
	replyisnil := nil == reply
	// TODO Temporary
	if replyisnil {
		return nil, errors.New("Internal: Nil response channel.")
	}
	return parseVersionResponse(<-reply)
}

// Open a file which is represented by a `Handle`.
func (c *Conn) Open(path string, pflags uint32, attrs FxpAttrs) (Handle, error) {
	id := c.nextPacketId()
	var pktLen uint32 = 4 + 1 + 4 + 4 + uint32(len(path)) + 4 + attrs.Len()
	buf := NewBuffer()
	buf.Grow(4 + pktLen)
	buf.WriteUint32(pktLen)
	buf.WriteByte(FXP_OPEN)
	buf.WriteUint32(id)
	buf.WriteString(path)
	buf.WriteUint32(pflags)
	buf.WriteAttrs(attrs)
	reply := c.Send(buf)
	replyisnil := nil == reply
	// TODO Temporary
	if replyisnil {
		return Handle{}, errors.New("Internal: Nil response channel.")
	}
	return parseHandleResponse(<-reply, c)
}
func (h *Handle) Close() error {
	id := h.conn.nextPacketId()
	var pktLen uint32 = 4 + 1 + 4 + 4 + uint32(len(h.h))
	buf := NewBuffer()
	buf.Grow(4 + pktLen)
	buf.WriteUint32(pktLen)
	buf.WriteByte(FXP_CLOSE)
	buf.WriteUint32(id)
	buf.WriteString(h.h)
	reply := h.conn.Send(buf)
	replyisnil := nil == reply
	// TODO Temporary
	if replyisnil {
		return errors.New("Internal: Nil response channel.")
	}
	return parseStatusResponse(<-reply)
}
func (h *Handle) Read(offset uint64, length uint32) ([]byte, error) {
	id := h.conn.nextPacketId()
	var pktLen uint32 = 4 + 1 + 4 + 4 + uint32(len(h.h)) + 8 + 4
	buf := NewBuffer()
	buf.Grow(4 + pktLen)
	buf.WriteUint32(pktLen)
	buf.WriteByte(FXP_READ)
	buf.WriteUint32(id)
	buf.WriteString(h.h)
	buf.WriteUint64(offset)
	buf.WriteUint32(length)
	reply := h.conn.Send(buf)
	replyisnil := nil == reply
	// TODO Temporary
	if replyisnil {
		return nil, errors.New("Internal: Nil response channel.")
	}
	return parseDataResponse(<-reply)
}
func (h *Handle) Write(offset uint64, length uint32, data []byte) error {
	id := h.conn.nextPacketId()
	var pktLen uint32 = 4 + 1 + 4 + 4 + uint32(len(h.h)) + 8 + 4 + uint32(len(data))
	buf := NewBuffer()
	buf.Grow(4 + pktLen)
	buf.WriteUint32(pktLen)
	buf.WriteByte(FXP_WRITE)
	buf.WriteUint32(id)
	buf.WriteString(h.h)
	buf.WriteUint64(offset)
	buf.WriteUint32(length)
	buf.Write(data)
	reply := h.conn.Send(buf)
	replyisnil := nil == reply
	// TODO Temporary
	if replyisnil {
		return errors.New("Internal: Nil response channel.")
	}
	return parseStatusResponse(<-reply)
}
func (c *Conn) Lstat(path string) (FxpAttrs, error) {
	id := c.nextPacketId()
	var pktLen uint32 = 4 + 1 + 4 + 4 + uint32(len(path))
	buf := NewBuffer()
	buf.Grow(4 + pktLen)
	buf.WriteUint32(pktLen)
	buf.WriteByte(FXP_LSTAT)
	buf.WriteUint32(id)
	buf.WriteString(path)
	reply := c.Send(buf)
	replyisnil := nil == reply
	// TODO Temporary
	if replyisnil {
		return FxpAttrs{}, errors.New("Internal: Nil response channel.")
	}
	return parseAttrsResponse(<-reply)
}
func (h *Handle) Fstat() (
	FxpAttrs, error) {
	id := h.conn.nextPacketId()
	var pktLen uint32 = 4 + 1 + 4 + 4 + uint32(len(h.h))
	buf := NewBuffer()
	buf.Grow(4 + pktLen)
	buf.WriteUint32(pktLen)
	buf.WriteByte(FXP_FSTAT)
	buf.WriteUint32(id)
	buf.WriteString(h.h)
	reply := h.conn.Send(buf)
	replyisnil := nil == reply
	// TODO Temporary
	if replyisnil {
		return FxpAttrs{}, errors.New("Internal: Nil response channel.")
	}
	return parseAttrsResponse(<-reply)
}
func (c *Conn) Setstat(path string, flags uint32, attrs FxpAttrs) error {
	id := c.nextPacketId()
	var pktLen uint32 = 4 + 1 + 4 + 4 + uint32(len(path)) + 4 + attrs.Len()
	buf := NewBuffer()
	buf.Grow(4 + pktLen)
	buf.WriteUint32(pktLen)
	buf.WriteByte(FXP_SETSTAT)
	buf.WriteUint32(id)
	buf.WriteString(path)
	buf.WriteUint32(flags)
	buf.WriteAttrs(attrs)
	reply := c.Send(buf)
	replyisnil := nil == reply
	// TODO Temporary
	if replyisnil {
		return errors.New("Internal: Nil response channel.")
	}
	return parseStatusResponse(<-reply)
}
func (h *Handle) Fsetstat(flags uint32, attrs FxpAttrs) error {
	id := h.conn.nextPacketId()
	var pktLen uint32 = 4 + 1 + 4 + 4 + uint32(len(h.h)) + 4 + attrs.Len()
	buf := NewBuffer()
	buf.Grow(4 + pktLen)
	buf.WriteUint32(pktLen)
	buf.WriteByte(FXP_FSETSTAT)
	buf.WriteUint32(id)
	buf.WriteString(h.h)
	buf.WriteUint32(flags)
	buf.WriteAttrs(attrs)
	reply := h.conn.Send(buf)
	replyisnil := nil == reply
	// TODO Temporary
	if replyisnil {
		return errors.New("Internal: Nil response channel.")
	}
	return parseStatusResponse(<-reply)
}
func (c *Conn) Opendir(path string) (Handle, error) {
	id := c.nextPacketId()
	var pktLen uint32 = 4 + 1 + 4 + 4 + uint32(len(path))
	buf := NewBuffer()
	buf.Grow(4 + pktLen)
	buf.WriteUint32(pktLen)
	buf.WriteByte(FXP_OPENDIR)
	buf.WriteUint32(id)
	buf.WriteString(path)
	reply := c.Send(buf)
	replyisnil := nil == reply
	// TODO Temporary
	if replyisnil {
		return Handle{}, errors.New("Internal: Nil response channel.")
	}
	return parseHandleResponse(<-reply, c)
}
func (h *Handle) Readdir() (
	[]FxpName, error) {
	id := h.conn.nextPacketId()
	var pktLen uint32 = 4 + 1 + 4 + 4 + uint32(len(h.h))
	buf := NewBuffer()
	buf.Grow(4 + pktLen)
	buf.WriteUint32(pktLen)
	buf.WriteByte(FXP_READDIR)
	buf.WriteUint32(id)
	buf.WriteString(h.h)
	reply := h.conn.Send(buf)
	replyisnil := nil == reply
	// TODO Temporary
	if replyisnil {
		return nil, errors.New("Internal: Nil response channel.")
	}
	return parseNameResponse(<-reply)
}
func (c *Conn) Remove(path string) error {
	id := c.nextPacketId()
	var pktLen uint32 = 4 + 1 + 4 + 4 + uint32(len(path))
	buf := NewBuffer()
	buf.Grow(4 + pktLen)
	buf.WriteUint32(pktLen)
	buf.WriteByte(FXP_REMOVE)
	buf.WriteUint32(id)
	buf.WriteString(path)
	reply := c.Send(buf)
	replyisnil := nil == reply
	// TODO Temporary
	if replyisnil {
		return errors.New("Internal: Nil response channel.")
	}
	return parseStatusResponse(<-reply)
}
func (c *Conn) Mkdir(path string, flags uint32) error {
	id := c.nextPacketId()
	var pktLen uint32 = 4 + 1 + 4 + 4 + uint32(len(path)) + 4
	buf := NewBuffer()
	buf.Grow(4 + pktLen)
	buf.WriteUint32(pktLen)
	buf.WriteByte(FXP_MKDIR)
	buf.WriteUint32(id)
	buf.WriteString(path)
	buf.WriteUint32(flags)
	reply := c.Send(buf)
	replyisnil := nil == reply
	// TODO Temporary
	if replyisnil {
		return errors.New("Internal: Nil response channel.")
	}
	return parseStatusResponse(<-reply)
}
func (c *Conn) Rmdir(path string) error {
	id := c.nextPacketId()
	var pktLen uint32 = 4 + 1 + 4 + 4 + uint32(len(path))
	buf := NewBuffer()
	buf.Grow(4 + pktLen)
	buf.WriteUint32(pktLen)
	buf.WriteByte(FXP_RMDIR)
	buf.WriteUint32(id)
	buf.WriteString(path)
	reply := c.Send(buf)
	replyisnil := nil == reply
	// TODO Temporary
	if replyisnil {
		return errors.New("Internal: Nil response channel.")
	}
	return parseStatusResponse(<-reply)
}
func (c *Conn) Realpath(path string) ([]FxpName, error) {
	id := c.nextPacketId()
	var pktLen uint32 = 4 + 1 + 4 + 4 + uint32(len(path))
	buf := NewBuffer()
	buf.Grow(4 + pktLen)
	buf.WriteUint32(pktLen)
	buf.WriteByte(FXP_REALPATH)
	buf.WriteUint32(id)
	buf.WriteString(path)
	reply := c.Send(buf)
	replyisnil := nil == reply
	// TODO Temporary
	if replyisnil {
		return nil, errors.New("Internal: Nil response channel.")
	}
	return parseNameResponse(<-reply)
}
func (c *Conn) Stat(path string) (FxpAttrs, error) {
	id := c.nextPacketId()
	var pktLen uint32 = 4 + 1 + 4 + 4 + uint32(len(path))
	buf := NewBuffer()
	buf.Grow(4 + pktLen)
	buf.WriteUint32(pktLen)
	buf.WriteByte(FXP_STAT)
	buf.WriteUint32(id)
	buf.WriteString(path)
	reply := c.Send(buf)
	replyisnil := nil == reply
	// TODO Temporary
	if replyisnil {
		return FxpAttrs{}, errors.New("Internal: Nil response channel.")
	}
	return parseAttrsResponse(<-reply)
}
func (c *Conn) Rename(path, newpath string) error {
	id := c.nextPacketId()
	var pktLen uint32 = 4 + 1 + 4 + 4 + uint32(len(path)) + 4 + uint32(len(newpath))
	buf := NewBuffer()
	buf.Grow(4 + pktLen)
	buf.WriteUint32(pktLen)
	buf.WriteByte(FXP_RENAME)
	buf.WriteUint32(id)
	buf.WriteString(path)
	buf.WriteString(newpath)
	reply := c.Send(buf)
	replyisnil := nil == reply
	// TODO Temporary
	if replyisnil {
		return errors.New("Internal: Nil response channel.")
	}
	return parseStatusResponse(<-reply)
}
func (c *Conn) Readlink(path string) ([]FxpName, error) {
	id := c.nextPacketId()
	var pktLen uint32 = 4 + 1 + 4 + 4 + uint32(len(path))
	buf := NewBuffer()
	buf.Grow(4 + pktLen)
	buf.WriteUint32(pktLen)
	buf.WriteByte(FXP_READLINK)
	buf.WriteUint32(id)
	buf.WriteString(path)
	reply := c.Send(buf)
	replyisnil := nil == reply
	// TODO Temporary
	if replyisnil {
		return nil, errors.New("Internal: Nil response channel.")
	}
	return parseNameResponse(<-reply)
}
func (c *Conn) Symlink(path, target string) error {
	id := c.nextPacketId()
	var pktLen uint32 = 4 + 1 + 4 + 4 + uint32(len(path)) + 4 + uint32(len(target))
	buf := NewBuffer()
	buf.Grow(4 + pktLen)
	buf.WriteUint32(pktLen)
	buf.WriteByte(FXP_SYMLINK)
	buf.WriteUint32(id)
	buf.WriteString(path)
	buf.WriteString(target)
	reply := c.Send(buf)
	replyisnil := nil == reply
	// TODO Temporary
	if replyisnil {
		return errors.New("Internal: Nil response channel.")
	}
	return parseStatusResponse(<-reply)
}
func (c *Conn) Extended(request string, payload []byte) ([]byte, error) {
	id := c.nextPacketId()
	var pktLen uint32 = 4 + 1 + 4 + 4 + uint32(len(request)) + uint32(len(payload))
	buf := NewBuffer()
	buf.Grow(4 + pktLen)
	buf.WriteUint32(pktLen)
	buf.WriteByte(FXP_EXTENDED)
	buf.WriteUint32(id)
	buf.WriteString(request)
	buf.Write(payload)
	reply := c.Send(buf)
	replyisnil := nil == reply
	// TODO Temporary
	if replyisnil {
		return nil, errors.New("Internal: Nil response channel.")
	}
	return parseExtendedReplyResponse(<-reply)
}
