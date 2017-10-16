package sftp

// Automatically generated file. Do not touch.

// Extended is an optinal list of `name` and `data` that you want to support.
func (c *Client) Init(extended [][2]string) ([][2]string, error) {
	id := NextId(c)
	buf := id.BufferInit(extended)
	reply := <-c.send(buf)
	// TODO Temporary
	return parseVersionResponse(reply)
}

// Open a file which is represented by a `Handle`.
func (c *Client) Open(path string, pflags uint32, attrs Attrs) (Handle, error) {
	id := NextId(c)
	buf := id.BufferOpen(path, pflags, attrs)
	reply := <-c.send(buf)
	// TODO Temporary
	return parseHandleResponse(reply, c)
}
func (h *Handle) Close() error {
	id := NextId(h.client)
	buf := id.BufferClose(h.h)
	reply := <-h.client.send(buf)
	// TODO Temporary
	return parseStatusResponse(reply)
}
func (h *Handle) Read(offset uint64, length uint32) ([]byte, error) {
	id := NextId(h.client)
	buf := id.BufferRead(h.h, offset, length)
	reply := <-h.client.send(buf)
	// TODO Temporary
	return parseDataResponse(reply)
}
func (h *Handle) Write(offset uint64, length uint32, data []byte) error {
	id := NextId(h.client)
	buf := id.BufferWrite(h.h, offset, length, data)
	reply := <-h.client.send(buf)
	// TODO Temporary
	return parseStatusResponse(reply)
}
func (c *Client) Lstat(path string) (Attrs, error) {
	id := NextId(c)
	buf := id.BufferLstat(path)
	reply := <-c.send(buf)
	// TODO Temporary
	return parseAttrsResponse(reply)
}
func (h *Handle) Fstat() (Attrs, error) {
	id := NextId(h.client)
	buf := id.BufferFstat(h.h)
	reply := <-h.client.send(buf)
	// TODO Temporary
	return parseAttrsResponse(reply)
}
func (c *Client) Setstat(path string, flags uint32, attrs Attrs) error {
	id := NextId(c)
	buf := id.BufferSetstat(path, flags, attrs)
	reply := <-c.send(buf)
	// TODO Temporary
	return parseStatusResponse(reply)
}
func (h *Handle) Fsetstat(flags uint32, attrs Attrs) error {
	id := NextId(h.client)
	buf := id.BufferFsetstat(h.h, flags, attrs)
	reply := <-h.client.send(buf)
	// TODO Temporary
	return parseStatusResponse(reply)
}
func (c *Client) Opendir(path string) (Handle, error) {
	id := NextId(c)
	buf := id.BufferOpendir(path)
	reply := <-c.send(buf)
	// TODO Temporary
	return parseHandleResponse(reply, c)
}
func (h *Handle) Readdir() ([]Name, error) {
	id := NextId(h.client)
	buf := id.BufferReaddir(h.h)
	reply := <-h.client.send(buf)
	// TODO Temporary
	return parseNameResponse(reply)
}
func (c *Client) Remove(path string) error {
	id := NextId(c)
	buf := id.BufferRemove(path)
	reply := <-c.send(buf)
	// TODO Temporary
	return parseStatusResponse(reply)
}
func (c *Client) Mkdir(path string, flags uint32) error {
	id := NextId(c)
	buf := id.BufferMkdir(path, flags)
	reply := <-c.send(buf)
	// TODO Temporary
	return parseStatusResponse(reply)
}
func (c *Client) Rmdir(path string) error {
	id := NextId(c)
	buf := id.BufferRmdir(path)
	reply := <-c.send(buf)
	// TODO Temporary
	return parseStatusResponse(reply)
}
func (c *Client) Realpath(path string) ([]Name, error) {
	id := NextId(c)
	buf := id.BufferRealpath(path)
	reply := <-c.send(buf)
	// TODO Temporary
	return parseNameResponse(reply)
}
func (c *Client) Stat(path string) (Attrs, error) {
	id := NextId(c)
	buf := id.BufferStat(path)
	reply := <-c.send(buf)
	// TODO Temporary
	return parseAttrsResponse(reply)
}
func (c *Client) Rename(path, newpath string) error {
	id := NextId(c)
	buf := id.BufferRename(path, newpath)
	reply := <-c.send(buf)
	// TODO Temporary
	return parseStatusResponse(reply)
}
func (c *Client) Readlink(path string) ([]Name, error) {
	id := NextId(c)
	buf := id.BufferReadlink(path)
	reply := <-c.send(buf)
	// TODO Temporary
	return parseNameResponse(reply)
}
func (c *Client) Symlink(path, target string) error {
	id := NextId(c)
	buf := id.BufferSymlink(path, target)
	reply := <-c.send(buf)
	// TODO Temporary
	return parseStatusResponse(reply)
}
func (c *Client) Extended(request string, payload []byte) ([]byte, error) {
	id := NextId(c)
	buf := id.BufferExtended(request, payload)
	reply := <-c.send(buf)
	// TODO Temporary
	return parseExtendedReplyResponse(reply)
}
