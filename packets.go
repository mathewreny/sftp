package sftp

// Automatically generated file. Do not touch.

type PacketId uint32

// The first packet sent over an SFTP session. Extended is an optinal list of `name` and `data`
// that you want to support.
func (id PacketId) BufferInit(extended [][2]string) *Buffer {
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

// The pflags variable is a bit mask. For more information, see this package's const declarations
// beginning with the "FXF" prefix.
func (id PacketId) BufferOpen(path string, pflags uint32, attrs Attrs) *Buffer {
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
func (id PacketId) BufferClose(handle string) *Buffer {
	var pktLen uint32 = 1 + 4 + 4 + uint32(len(handle))
	buf := NewBuffer()
	buf.Grow(4 + pktLen)
	buf.WriteUint32(pktLen)
	buf.WriteByte(FXP_CLOSE)
	buf.WriteUint32(uint32(id))
	buf.WriteString(handle)
	return buf
}
func (id PacketId) BufferRead(handle string, offset uint64, length uint32) *Buffer {
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
func (id PacketId) BufferWrite(handle string, offset uint64, length uint32, data []byte) *Buffer {
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
func (id PacketId) BufferLstat(path string) *Buffer {
	var pktLen uint32 = 1 + 4 + 4 + uint32(len(path))
	buf := NewBuffer()
	buf.Grow(4 + pktLen)
	buf.WriteUint32(pktLen)
	buf.WriteByte(FXP_LSTAT)
	buf.WriteUint32(uint32(id))
	buf.WriteString(path)
	return buf
}
func (id PacketId) BufferFstat(handle string) *Buffer {
	var pktLen uint32 = 1 + 4 + 4 + uint32(len(handle))
	buf := NewBuffer()
	buf.Grow(4 + pktLen)
	buf.WriteUint32(pktLen)
	buf.WriteByte(FXP_FSTAT)
	buf.WriteUint32(uint32(id))
	buf.WriteString(handle)
	return buf
}
func (id PacketId) BufferSetstat(path string, flags uint32, attrs Attrs) *Buffer {
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
func (id PacketId) BufferFsetstat(handle string, flags uint32, attrs Attrs) *Buffer {
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
func (id PacketId) BufferOpendir(path string) *Buffer {
	var pktLen uint32 = 1 + 4 + 4 + uint32(len(path))
	buf := NewBuffer()
	buf.Grow(4 + pktLen)
	buf.WriteUint32(pktLen)
	buf.WriteByte(FXP_OPENDIR)
	buf.WriteUint32(uint32(id))
	buf.WriteString(path)
	return buf
}
func (id PacketId) BufferReaddir(handle string) *Buffer {
	var pktLen uint32 = 1 + 4 + 4 + uint32(len(handle))
	buf := NewBuffer()
	buf.Grow(4 + pktLen)
	buf.WriteUint32(pktLen)
	buf.WriteByte(FXP_READDIR)
	buf.WriteUint32(uint32(id))
	buf.WriteString(handle)
	return buf
}
func (id PacketId) BufferRemove(path string) *Buffer {
	var pktLen uint32 = 1 + 4 + 4 + uint32(len(path))
	buf := NewBuffer()
	buf.Grow(4 + pktLen)
	buf.WriteUint32(pktLen)
	buf.WriteByte(FXP_REMOVE)
	buf.WriteUint32(uint32(id))
	buf.WriteString(path)
	return buf
}
func (id PacketId) BufferMkdir(path string, flags uint32) *Buffer {
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
func (id PacketId) BufferRmdir(path string) *Buffer {
	var pktLen uint32 = 1 + 4 + 4 + uint32(len(path))
	buf := NewBuffer()
	buf.Grow(4 + pktLen)
	buf.WriteUint32(pktLen)
	buf.WriteByte(FXP_RMDIR)
	buf.WriteUint32(uint32(id))
	buf.WriteString(path)
	return buf
}
func (id PacketId) BufferRealpath(path string) *Buffer {
	var pktLen uint32 = 1 + 4 + 4 + uint32(len(path))
	buf := NewBuffer()
	buf.Grow(4 + pktLen)
	buf.WriteUint32(pktLen)
	buf.WriteByte(FXP_REALPATH)
	buf.WriteUint32(uint32(id))
	buf.WriteString(path)
	return buf
}
func (id PacketId) BufferStat(path string) *Buffer {
	var pktLen uint32 = 1 + 4 + 4 + uint32(len(path))
	buf := NewBuffer()
	buf.Grow(4 + pktLen)
	buf.WriteUint32(pktLen)
	buf.WriteByte(FXP_STAT)
	buf.WriteUint32(uint32(id))
	buf.WriteString(path)
	return buf
}
func (id PacketId) BufferRename(path, newpath string) *Buffer {
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
func (id PacketId) BufferReadlink(path string) *Buffer {
	var pktLen uint32 = 1 + 4 + 4 + uint32(len(path))
	buf := NewBuffer()
	buf.Grow(4 + pktLen)
	buf.WriteUint32(pktLen)
	buf.WriteByte(FXP_READLINK)
	buf.WriteUint32(uint32(id))
	buf.WriteString(path)
	return buf
}
func (id PacketId) BufferSymlink(path, target string) *Buffer {
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

// The payload's length must be less than or equal to `max(uint32)-17-len(request)`. Ditto for the
// payload request.
func (id PacketId) BufferExtended(request string, payload []byte) *Buffer {
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

// The version packet is a response to the Init packet. Extended is an optinal slice of extensions
// that the server supports. The array's first index is `name` and the second is `data`.
func (id PacketId) BufferVersion(extended [][2]string) *Buffer {
	var pktLen uint32 = 1 + 4
	for _, ext := range extended {
		pktLen += 4 + uint32(len(ext[0]))
		pktLen += 4 + uint32(len(ext[1]))
	}

	buf := NewBuffer()
	buf.Grow(4 + pktLen)
	buf.WriteUint32(pktLen)
	buf.WriteByte(FXP_VERSION)
	buf.WriteUint32(uint32(id))
	for _, ext := range extended {
		buf.WriteString(ext[0])
		buf.WriteString(ext[1])
	}
	return buf
}

// See the Status type for an alternate way to create this packet.
func (id PacketId) BufferStatus(code uint32, message, lang string) *Buffer {
	var pktLen uint32 = 1 + 4 + 4 + 4 + uint32(len(message)) + 4 + uint32(len(lang))
	buf := NewBuffer()
	buf.Grow(4 + pktLen)
	buf.WriteUint32(pktLen)
	buf.WriteByte(FXP_STATUS)
	buf.WriteUint32(uint32(id))
	buf.WriteUint32(code)
	buf.WriteString(message)
	buf.WriteString(lang)
	return buf
}

// See the Handle type for more information. The handle string's length MUST NOT exceed 256 bytes
// according to the protocol. Callers are responsible for enforcing this limit.
func (id PacketId) BufferHandle(handle string) *Buffer {
	var pktLen uint32 = 1 + 4 + 4 + uint32(len(handle))
	buf := NewBuffer()
	buf.Grow(4 + pktLen)
	buf.WriteUint32(pktLen)
	buf.WriteByte(FXP_HANDLE)
	buf.WriteUint32(uint32(id))
	buf.WriteString(handle)
	return buf
}

// Data is returned from sftp read requests. Obviously the provided length must match len(data).
// The Length must be less than or equal to `max(uint32) - 13`. Callers are responsible for
// enforcing this limit.
func (id PacketId) BufferData(length uint32, data []byte) *Buffer {
	var pktLen uint32 = 1 + 4 + 4 + uint32(len(data))
	buf := NewBuffer()
	buf.Grow(4 + pktLen)
	buf.WriteUint32(pktLen)
	buf.WriteByte(FXP_DATA)
	buf.WriteUint32(uint32(id))
	buf.WriteUint32(length)
	buf.Write(data)
	return buf
}

// See the Name type for more information.
func (id PacketId) BufferName(names []Name) *Buffer {
	var pktLen uint32 = 1 + 4
	buf := NewBuffer()
	buf.Grow(4 + pktLen)
	buf.WriteUint32(pktLen)
	buf.WriteByte(FXP_NAME)
	buf.WriteUint32(uint32(id))
	return buf
}

// See the Attrs type for more information.
func (id PacketId) BufferAttrs(attrs Attrs) *Buffer {
	var pktLen uint32 = 1 + 4 + attrs.Len()
	buf := NewBuffer()
	buf.Grow(4 + pktLen)
	buf.WriteUint32(pktLen)
	buf.WriteByte(FXP_ATTRS)
	buf.WriteUint32(uint32(id))
	buf.WriteAttrs(attrs)
	return buf
}

// The extended reply allows servers to extend the SFTP version 3 protocol. The data's length must
// be less than or equal to `max(uint32) - 9`. Callers are responsible for enforcing this limit.
func (id PacketId) BufferExtendedReply(data []byte) *Buffer {
	var pktLen uint32 = 1 + 4 + uint32(len(data))
	buf := NewBuffer()
	buf.Grow(4 + pktLen)
	buf.WriteUint32(pktLen)
	buf.WriteByte(FXP_EXTENDEDREPLY)
	buf.WriteUint32(uint32(id))
	buf.Write(data)
	return buf
}
