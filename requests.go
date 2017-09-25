package sftp

// Automatically generated file. Do not touch.

func (c *Conn) Open(path string, pflags uint32, attrs FxpAttrs) (handle chan string, status chan error) {
	id := c.generatePacketId()
	var pktLen uint32 = 4 + 1 + 4 + 4 + uint32(len(path)) + 4 + attrs.Len()
	buf := NewBuffer()
	buf.Grow(4 + pktLen)
	buf.WriteUint32(pktLen)
	buf.WriteByte(FXP_OPEN)
	buf.WriteUint32(id)
	buf.WriteString(path)
	buf.WriteUint32(pflags)
	buf.WriteAttrs(attrs)
	handle = c.handleResponse(id)
	status = c.statusResponse(id)
	c.send(id, buf)
	return
}
func (c *Conn) Close(handle string) (status chan error) {
	id := c.generatePacketId()
	var pktLen uint32 = 4 + 1 + 4 + 4 + uint32(len(handle))
	buf := NewBuffer()
	buf.Grow(4 + pktLen)
	buf.WriteUint32(pktLen)
	buf.WriteByte(FXP_CLOSE)
	buf.WriteUint32(id)
	buf.WriteString(handle)
	status = c.statusResponse(id)
	c.send(id, buf)
	return
}
func (c *Conn) Read(handle string, offset uint64, length uint32) (response chan []byte, status chan error) {
	id := c.generatePacketId()
	var pktLen uint32 = 4 + 1 + 4 + 4 + uint32(len(handle)) + 8 + 4
	buf := NewBuffer()
	buf.Grow(4 + pktLen)
	buf.WriteUint32(pktLen)
	buf.WriteByte(FXP_READ)
	buf.WriteUint32(id)
	buf.WriteString(handle)
	buf.WriteUint64(offset)
	buf.WriteUint32(length)
	response = c.dataResponse(id)
	status = c.statusResponse(id)
	c.send(id, buf)
	return
}
func (c *Conn) Write(handle string, offset uint64, length uint32, data []byte) (status chan error) {
	id := c.generatePacketId()
	var pktLen uint32 = 4 + 1 + 4 + 4 + uint32(len(handle)) + 8 + 4 + uint32(len(data))
	buf := NewBuffer()
	buf.Grow(4 + pktLen)
	buf.WriteUint32(pktLen)
	buf.WriteByte(FXP_WRITE)
	buf.WriteUint32(id)
	buf.WriteString(handle)
	buf.WriteUint64(offset)
	buf.WriteUint32(length)
	buf.Write(data)
	status = c.statusResponse(id)
	c.send(id, buf)
	return
}
func (c *Conn) Lstat(path string) (response chan FxpAttrs, status chan error) {
	id := c.generatePacketId()
	var pktLen uint32 = 4 + 1 + 4 + 4 + uint32(len(path))
	buf := NewBuffer()
	buf.Grow(4 + pktLen)
	buf.WriteUint32(pktLen)
	buf.WriteByte(FXP_LSTAT)
	buf.WriteUint32(id)
	buf.WriteString(path)
	response = c.attrsResponse(id)
	status = c.statusResponse(id)
	c.send(id, buf)
	return
}
func (c *Conn) Fstat(handle string) (response chan FxpAttrs, status chan error) {
	id := c.generatePacketId()
	var pktLen uint32 = 4 + 1 + 4 + 4 + uint32(len(handle))
	buf := NewBuffer()
	buf.Grow(4 + pktLen)
	buf.WriteUint32(pktLen)
	buf.WriteByte(FXP_FSTAT)
	buf.WriteUint32(id)
	buf.WriteString(handle)
	response = c.attrsResponse(id)
	status = c.statusResponse(id)
	c.send(id, buf)
	return
}
func (c *Conn) Setstat(path string, flags uint32, attrs FxpAttrs) (status chan error) {
	id := c.generatePacketId()
	var pktLen uint32 = 4 + 1 + 4 + 4 + uint32(len(path)) + 4 + attrs.Len()
	buf := NewBuffer()
	buf.Grow(4 + pktLen)
	buf.WriteUint32(pktLen)
	buf.WriteByte(FXP_SETSTAT)
	buf.WriteUint32(id)
	buf.WriteString(path)
	buf.WriteUint32(flags)
	buf.WriteAttrs(attrs)
	status = c.statusResponse(id)
	c.send(id, buf)
	return
}
func (c *Conn) Fsetstat(handle string, flags uint32, attrs FxpAttrs) (status chan error) {
	id := c.generatePacketId()
	var pktLen uint32 = 4 + 1 + 4 + 4 + uint32(len(handle)) + 4 + attrs.Len()
	buf := NewBuffer()
	buf.Grow(4 + pktLen)
	buf.WriteUint32(pktLen)
	buf.WriteByte(FXP_FSETSTAT)
	buf.WriteUint32(id)
	buf.WriteString(handle)
	buf.WriteUint32(flags)
	buf.WriteAttrs(attrs)
	status = c.statusResponse(id)
	c.send(id, buf)
	return
}
func (c *Conn) Opendir(path string) (handle chan string, status chan error) {
	id := c.generatePacketId()
	var pktLen uint32 = 4 + 1 + 4 + 4 + uint32(len(path))
	buf := NewBuffer()
	buf.Grow(4 + pktLen)
	buf.WriteUint32(pktLen)
	buf.WriteByte(FXP_OPENDIR)
	buf.WriteUint32(id)
	buf.WriteString(path)
	handle = c.handleResponse(id)
	status = c.statusResponse(id)
	c.send(id, buf)
	return
}
func (c *Conn) Readdir(handle string) (response chan []FxpName, status chan error) {
	id := c.generatePacketId()
	var pktLen uint32 = 4 + 1 + 4 + 4 + uint32(len(handle))
	buf := NewBuffer()
	buf.Grow(4 + pktLen)
	buf.WriteUint32(pktLen)
	buf.WriteByte(FXP_READDIR)
	buf.WriteUint32(id)
	buf.WriteString(handle)
	response = c.nameResponse(id)
	status = c.statusResponse(id)
	c.send(id, buf)
	return
}
func (c *Conn) Remove(path string) (status chan error) {
	id := c.generatePacketId()
	var pktLen uint32 = 4 + 1 + 4 + 4 + uint32(len(path))
	buf := NewBuffer()
	buf.Grow(4 + pktLen)
	buf.WriteUint32(pktLen)
	buf.WriteByte(FXP_REMOVE)
	buf.WriteUint32(id)
	buf.WriteString(path)
	status = c.statusResponse(id)
	c.send(id, buf)
	return
}
func (c *Conn) Mkdir(path string, flags uint32) (status chan error) {
	id := c.generatePacketId()
	var pktLen uint32 = 4 + 1 + 4 + 4 + uint32(len(path)) + 4
	buf := NewBuffer()
	buf.Grow(4 + pktLen)
	buf.WriteUint32(pktLen)
	buf.WriteByte(FXP_MKDIR)
	buf.WriteUint32(id)
	buf.WriteString(path)
	buf.WriteUint32(flags)
	status = c.statusResponse(id)
	c.send(id, buf)
	return
}
func (c *Conn) Rmdir(path string) (status chan error) {
	id := c.generatePacketId()
	var pktLen uint32 = 4 + 1 + 4 + 4 + uint32(len(path))
	buf := NewBuffer()
	buf.Grow(4 + pktLen)
	buf.WriteUint32(pktLen)
	buf.WriteByte(FXP_RMDIR)
	buf.WriteUint32(id)
	buf.WriteString(path)
	status = c.statusResponse(id)
	c.send(id, buf)
	return
}
func (c *Conn) Realpath(path string) (response chan []FxpName, status chan error) {
	id := c.generatePacketId()
	var pktLen uint32 = 4 + 1 + 4 + 4 + uint32(len(path))
	buf := NewBuffer()
	buf.Grow(4 + pktLen)
	buf.WriteUint32(pktLen)
	buf.WriteByte(FXP_REALPATH)
	buf.WriteUint32(id)
	buf.WriteString(path)
	response = c.nameResponse(id)
	status = c.statusResponse(id)
	c.send(id, buf)
	return
}
func (c *Conn) Stat(path string) (response chan FxpAttrs, status chan error) {
	id := c.generatePacketId()
	var pktLen uint32 = 4 + 1 + 4 + 4 + uint32(len(path))
	buf := NewBuffer()
	buf.Grow(4 + pktLen)
	buf.WriteUint32(pktLen)
	buf.WriteByte(FXP_STAT)
	buf.WriteUint32(id)
	buf.WriteString(path)
	response = c.attrsResponse(id)
	status = c.statusResponse(id)
	c.send(id, buf)
	return
}
func (c *Conn) Rename(path, newpath string) (status chan error) {
	id := c.generatePacketId()
	var pktLen uint32 = 4 + 1 + 4 + 4 + uint32(len(path)) + 4 + uint32(len(newpath))
	buf := NewBuffer()
	buf.Grow(4 + pktLen)
	buf.WriteUint32(pktLen)
	buf.WriteByte(FXP_RENAME)
	buf.WriteUint32(id)
	buf.WriteString(path)
	buf.WriteString(newpath)
	status = c.statusResponse(id)
	c.send(id, buf)
	return
}
func (c *Conn) Readlink(path string) (response chan []FxpName, status chan error) {
	id := c.generatePacketId()
	var pktLen uint32 = 4 + 1 + 4 + 4 + uint32(len(path))
	buf := NewBuffer()
	buf.Grow(4 + pktLen)
	buf.WriteUint32(pktLen)
	buf.WriteByte(FXP_READLINK)
	buf.WriteUint32(id)
	buf.WriteString(path)
	response = c.nameResponse(id)
	status = c.statusResponse(id)
	c.send(id, buf)
	return
}
func (c *Conn) Symlink(path, target string) (status chan error) {
	id := c.generatePacketId()
	var pktLen uint32 = 4 + 1 + 4 + 4 + uint32(len(path)) + 4 + uint32(len(target))
	buf := NewBuffer()
	buf.Grow(4 + pktLen)
	buf.WriteUint32(pktLen)
	buf.WriteByte(FXP_SYMLINK)
	buf.WriteUint32(id)
	buf.WriteString(path)
	buf.WriteString(target)
	status = c.statusResponse(id)
	c.send(id, buf)
	return
}
func (c *Conn) Extended(request string, payload []byte) (response chan []byte, status chan error) {
	id := c.generatePacketId()
	var pktLen uint32 = 4 + 1 + 4 + 4 + uint32(len(request)) + uint32(len(payload))
	buf := NewBuffer()
	buf.Grow(4 + pktLen)
	buf.WriteUint32(pktLen)
	buf.WriteByte(FXP_EXTENDED)
	buf.WriteUint32(id)
	buf.WriteString(request)
	buf.Write(payload)
	response = c.extendedReplyResponse(id)
	status = c.statusResponse(id)
	c.send(id, buf)
	return
}
