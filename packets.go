package sftp

// Automatically generated file. Do not touch.

// Used to create SFTP version 3 packets.
// PacketId is its own type to help unclutter documentation.
type PacketId uint32

func (id PacketId) Init(extended [][2]string) *Buffer {
	var pktLen uint32 = 1 + 4
	for _, ext := range extended {
		pktLen += 4 + uint32(len(ext[0]))
		pktLen += 4 + uint32(len(ext[1]))
	}

	buf := NewBuffer()
	buf.Grow(4 + pktLen)
	buf.WriteUint32(pktLen)
	buf.WriteByte(FXP_INIT)
	buf.WriteUint32(uint32(id))
	for _, ext := range extended {
		buf.WriteString(ext[0])
		buf.WriteString(ext[1])
	}
	return buf
}
func (id PacketId) Open(path string, pflags uint32, attrs Attrs) *Buffer {
	var pktLen uint32 = 1 + 4 + 4 + uint32(len(path)) + 4 + attrs.Len()
	buf := NewBuffer()
	buf.Grow(4 + pktLen)
	buf.WriteUint32(pktLen)
	buf.WriteByte(FXP_OPEN)
	buf.WriteUint32(uint32(id))
	buf.WriteString(path)
	buf.WriteUint32(pflags)
	buf.WriteAttrs(attrs)
	return buf
}
func (id PacketId) Close(handle string) *Buffer {
	var pktLen uint32 = 1 + 4 + 4 + uint32(len(handle))
	buf := NewBuffer()
	buf.Grow(4 + pktLen)
	buf.WriteUint32(pktLen)
	buf.WriteByte(FXP_CLOSE)
	buf.WriteUint32(uint32(id))
	buf.WriteString(handle)
	return buf
}
func (id PacketId) Read(handle string, offset uint64, length uint32) *Buffer {
	var pktLen uint32 = 1 + 4 + 4 + uint32(len(handle)) + 8 + 4
	buf := NewBuffer()
	buf.Grow(4 + pktLen)
	buf.WriteUint32(pktLen)
	buf.WriteByte(FXP_READ)
	buf.WriteUint32(uint32(id))
	buf.WriteString(handle)
	buf.WriteUint64(offset)
	buf.WriteUint32(length)
	return buf
}
func (id PacketId) Write(handle string, offset uint64, length uint32, data []byte) *Buffer {
	var pktLen uint32 = 1 + 4 + 4 + uint32(len(handle)) + 8 + 4 + uint32(len(data))
	buf := NewBuffer()
	buf.Grow(4 + pktLen)
	buf.WriteUint32(pktLen)
	buf.WriteByte(FXP_WRITE)
	buf.WriteUint32(uint32(id))
	buf.WriteString(handle)
	buf.WriteUint64(offset)
	buf.WriteUint32(length)
	buf.Write(data)
	return buf
}
func (id PacketId) Lstat(path string) *Buffer {
	var pktLen uint32 = 1 + 4 + 4 + uint32(len(path))
	buf := NewBuffer()
	buf.Grow(4 + pktLen)
	buf.WriteUint32(pktLen)
	buf.WriteByte(FXP_LSTAT)
	buf.WriteUint32(uint32(id))
	buf.WriteString(path)
	return buf
}
func (id PacketId) Fstat(handle string) *Buffer {
	var pktLen uint32 = 1 + 4 + 4 + uint32(len(handle))
	buf := NewBuffer()
	buf.Grow(4 + pktLen)
	buf.WriteUint32(pktLen)
	buf.WriteByte(FXP_FSTAT)
	buf.WriteUint32(uint32(id))
	buf.WriteString(handle)
	return buf
}
func (id PacketId) Setstat(path string, flags uint32, attrs Attrs) *Buffer {
	var pktLen uint32 = 1 + 4 + 4 + uint32(len(path)) + 4 + attrs.Len()
	buf := NewBuffer()
	buf.Grow(4 + pktLen)
	buf.WriteUint32(pktLen)
	buf.WriteByte(FXP_SETSTAT)
	buf.WriteUint32(uint32(id))
	buf.WriteString(path)
	buf.WriteUint32(flags)
	buf.WriteAttrs(attrs)
	return buf
}
func (id PacketId) Fsetstat(handle string, flags uint32, attrs Attrs) *Buffer {
	var pktLen uint32 = 1 + 4 + 4 + uint32(len(handle)) + 4 + attrs.Len()
	buf := NewBuffer()
	buf.Grow(4 + pktLen)
	buf.WriteUint32(pktLen)
	buf.WriteByte(FXP_FSETSTAT)
	buf.WriteUint32(uint32(id))
	buf.WriteString(handle)
	buf.WriteUint32(flags)
	buf.WriteAttrs(attrs)
	return buf
}
func (id PacketId) Opendir(path string) *Buffer {
	var pktLen uint32 = 1 + 4 + 4 + uint32(len(path))
	buf := NewBuffer()
	buf.Grow(4 + pktLen)
	buf.WriteUint32(pktLen)
	buf.WriteByte(FXP_OPENDIR)
	buf.WriteUint32(uint32(id))
	buf.WriteString(path)
	return buf
}
func (id PacketId) Readdir(handle string) *Buffer {
	var pktLen uint32 = 1 + 4 + 4 + uint32(len(handle))
	buf := NewBuffer()
	buf.Grow(4 + pktLen)
	buf.WriteUint32(pktLen)
	buf.WriteByte(FXP_READDIR)
	buf.WriteUint32(uint32(id))
	buf.WriteString(handle)
	return buf
}
func (id PacketId) Remove(path string) *Buffer {
	var pktLen uint32 = 1 + 4 + 4 + uint32(len(path))
	buf := NewBuffer()
	buf.Grow(4 + pktLen)
	buf.WriteUint32(pktLen)
	buf.WriteByte(FXP_REMOVE)
	buf.WriteUint32(uint32(id))
	buf.WriteString(path)
	return buf
}
func (id PacketId) Mkdir(path string, flags uint32) *Buffer {
	var pktLen uint32 = 1 + 4 + 4 + uint32(len(path)) + 4
	buf := NewBuffer()
	buf.Grow(4 + pktLen)
	buf.WriteUint32(pktLen)
	buf.WriteByte(FXP_MKDIR)
	buf.WriteUint32(uint32(id))
	buf.WriteString(path)
	buf.WriteUint32(flags)
	return buf
}
func (id PacketId) Rmdir(path string) *Buffer {
	var pktLen uint32 = 1 + 4 + 4 + uint32(len(path))
	buf := NewBuffer()
	buf.Grow(4 + pktLen)
	buf.WriteUint32(pktLen)
	buf.WriteByte(FXP_RMDIR)
	buf.WriteUint32(uint32(id))
	buf.WriteString(path)
	return buf
}
func (id PacketId) Realpath(path string) *Buffer {
	var pktLen uint32 = 1 + 4 + 4 + uint32(len(path))
	buf := NewBuffer()
	buf.Grow(4 + pktLen)
	buf.WriteUint32(pktLen)
	buf.WriteByte(FXP_REALPATH)
	buf.WriteUint32(uint32(id))
	buf.WriteString(path)
	return buf
}
func (id PacketId) Stat(path string) *Buffer {
	var pktLen uint32 = 1 + 4 + 4 + uint32(len(path))
	buf := NewBuffer()
	buf.Grow(4 + pktLen)
	buf.WriteUint32(pktLen)
	buf.WriteByte(FXP_STAT)
	buf.WriteUint32(uint32(id))
	buf.WriteString(path)
	return buf
}
func (id PacketId) Rename(path, newpath string) *Buffer {
	var pktLen uint32 = 1 + 4 + 4 + uint32(len(path)) + 4 + uint32(len(newpath))
	buf := NewBuffer()
	buf.Grow(4 + pktLen)
	buf.WriteUint32(pktLen)
	buf.WriteByte(FXP_RENAME)
	buf.WriteUint32(uint32(id))
	buf.WriteString(path)
	buf.WriteString(newpath)
	return buf
}
func (id PacketId) Readlink(path string) *Buffer {
	var pktLen uint32 = 1 + 4 + 4 + uint32(len(path))
	buf := NewBuffer()
	buf.Grow(4 + pktLen)
	buf.WriteUint32(pktLen)
	buf.WriteByte(FXP_READLINK)
	buf.WriteUint32(uint32(id))
	buf.WriteString(path)
	return buf
}
func (id PacketId) Symlink(path, target string) *Buffer {
	var pktLen uint32 = 1 + 4 + 4 + uint32(len(path)) + 4 + uint32(len(target))
	buf := NewBuffer()
	buf.Grow(4 + pktLen)
	buf.WriteUint32(pktLen)
	buf.WriteByte(FXP_SYMLINK)
	buf.WriteUint32(uint32(id))
	buf.WriteString(path)
	buf.WriteString(target)
	return buf
}
func (id PacketId) Extended(request string, payload []byte) *Buffer {
	var pktLen uint32 = 1 + 4 + 4 + uint32(len(request)) + uint32(len(payload))
	buf := NewBuffer()
	buf.Grow(4 + pktLen)
	buf.WriteUint32(pktLen)
	buf.WriteByte(FXP_EXTENDED)
	buf.WriteUint32(uint32(id))
	buf.WriteString(request)
	buf.Write(payload)
	return buf
}
