package sftp

// Automatically generated file. Do not touch.

import "errors"

// Extended is an optinal list of `name` and `data` that you want to support.
func (c *Client) Init(extended [][2]string) ([][2]string, error) {
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
func (c *Client) Open(path string, pflags uint32, attrs Attrs) (Handle, error) {
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
	id := h.client.nextPacketId()
	var pktLen uint32 = 4 + 1 + 4 + 4 + uint32(len(h.h))
	buf := NewBuffer()
	buf.Grow(4 + pktLen)
	buf.WriteUint32(pktLen)
	buf.WriteByte(FXP_CLOSE)
	buf.WriteUint32(id)
	buf.WriteString(h.h)
	reply := h.client.Send(buf)
	replyisnil := nil == reply
	// TODO Temporary
	if replyisnil {
		return errors.New("Internal: Nil response channel.")
	}
	return parseStatusResponse(<-reply)
}
func (h *Handle) Read(offset uint64, length uint32) ([]byte, error) {
	id := h.client.nextPacketId()
	var pktLen uint32 = 4 + 1 + 4 + 4 + uint32(len(h.h)) + 8 + 4
	buf := NewBuffer()
	buf.Grow(4 + pktLen)
	buf.WriteUint32(pktLen)
	buf.WriteByte(FXP_READ)
	buf.WriteUint32(id)
	buf.WriteString(h.h)
	buf.WriteUint64(offset)
	buf.WriteUint32(length)
	reply := h.client.Send(buf)
	replyisnil := nil == reply
	// TODO Temporary
	if replyisnil {
		return nil, errors.New("Internal: Nil response channel.")
	}
	return parseDataResponse(<-reply)
}
func (h *Handle) Write(offset uint64, length uint32, data []byte) error {
	id := h.client.nextPacketId()
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
	reply := h.client.Send(buf)
	replyisnil := nil == reply
	// TODO Temporary
	if replyisnil {
		return errors.New("Internal: Nil response channel.")
	}
	return parseStatusResponse(<-reply)
}
func (c *Client) Lstat(path string) (Attrs, error) {
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
		return Attrs{}, errors.New("Internal: Nil response channel.")
	}
	return parseAttrsResponse(<-reply)
}
func (h *Handle) Fstat() (
	Attrs, error) {
	id := h.client.nextPacketId()
	var pktLen uint32 = 4 + 1 + 4 + 4 + uint32(len(h.h))
	buf := NewBuffer()
	buf.Grow(4 + pktLen)
	buf.WriteUint32(pktLen)
	buf.WriteByte(FXP_FSTAT)
	buf.WriteUint32(id)
	buf.WriteString(h.h)
	reply := h.client.Send(buf)
	replyisnil := nil == reply
	// TODO Temporary
	if replyisnil {
		return Attrs{}, errors.New("Internal: Nil response channel.")
	}
	return parseAttrsResponse(<-reply)
}
func (c *Client) Setstat(path string, flags uint32, attrs Attrs) error {
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
func (h *Handle) Fsetstat(flags uint32, attrs Attrs) error {
	id := h.client.nextPacketId()
	var pktLen uint32 = 4 + 1 + 4 + 4 + uint32(len(h.h)) + 4 + attrs.Len()
	buf := NewBuffer()
	buf.Grow(4 + pktLen)
	buf.WriteUint32(pktLen)
	buf.WriteByte(FXP_FSETSTAT)
	buf.WriteUint32(id)
	buf.WriteString(h.h)
	buf.WriteUint32(flags)
	buf.WriteAttrs(attrs)
	reply := h.client.Send(buf)
	replyisnil := nil == reply
	// TODO Temporary
	if replyisnil {
		return errors.New("Internal: Nil response channel.")
	}
	return parseStatusResponse(<-reply)
}
func (c *Client) Opendir(path string) (Handle, error) {
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
	[]Name, error) {
	id := h.client.nextPacketId()
	var pktLen uint32 = 4 + 1 + 4 + 4 + uint32(len(h.h))
	buf := NewBuffer()
	buf.Grow(4 + pktLen)
	buf.WriteUint32(pktLen)
	buf.WriteByte(FXP_READDIR)
	buf.WriteUint32(id)
	buf.WriteString(h.h)
	reply := h.client.Send(buf)
	replyisnil := nil == reply
	// TODO Temporary
	if replyisnil {
		return nil, errors.New("Internal: Nil response channel.")
	}
	return parseNameResponse(<-reply)
}
func (c *Client) Remove(path string) error {
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
func (c *Client) Mkdir(path string, flags uint32) error {
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
func (c *Client) Rmdir(path string) error {
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
func (c *Client) Realpath(path string) ([]Name, error) {
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
func (c *Client) Stat(path string) (Attrs, error) {
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
		return Attrs{}, errors.New("Internal: Nil response channel.")
	}
	return parseAttrsResponse(<-reply)
}
func (c *Client) Rename(path, newpath string) error {
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
func (c *Client) Readlink(path string) ([]Name, error) {
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
func (c *Client) Symlink(path, target string) error {
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
func (c *Client) Extended(request string, payload []byte) ([]byte, error) {
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
