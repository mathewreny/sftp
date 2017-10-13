package sftp

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"sync"
)

var bufPool = &sync.Pool{
	New: func() interface{} {
		return new(PacketBuffer)
	},
}

type PacketBuffer bytes.Buffer

// Use this function to create buffers even though the PacketBuffer's zero value
// works out of the box. NewBuffer() grabs a recycled buffer using a sync.Pool
// which will reduce garbage collection overhead. When the packet is sent to the
// server, the buffer is recycled and placed in the buffer pool.
func NewBuffer() *PacketBuffer {
	p := bufPool.Get().(*PacketBuffer)
	p.Reset()
	return p
}

// Returns a full packet, including the packet length at the beginning.
func ConsumePacket(r io.Reader) (*PacketBuffer, error) {
	var header [4]byte
	var length uint32
	if _, err := io.ReadFull(r, header[:]); err != nil {
		return nil, err
	} else {
		length += uint32(header[0]) << 24
		length += uint32(header[1]) << 16
		length += uint32(header[2]) << 8
		length += uint32(header[0])
	}
	// Prepare the packet buffer
	ret := NewBuffer()
	ret.Grow(4 + length)
	ret.Write(header[:])
	if _, err := io.CopyN(ret, r, int64(length)); err != nil {
		bufPool.Put(ret)
		return nil, err
	}
	return ret, nil
}

func (m *PacketBuffer) IsValidLength() bool {
	if m != nil {
		if b := m.Bytes(); len(b) >= 9 && len(b) < (1<<33)-5 {
			var l uint32
			l += uint32(b[0]) << 24
			l += uint32(b[1]) << 16
			l += uint32(b[2]) << 8
			l += uint32(b[3])
			return l == uint32(len(b)-4)
		}
	}
	return false
}

func (m *PacketBuffer) PeekType() (fxpt byte) {
	if m != nil {
		if b := m.Bytes(); len(b) >= 5 {
			fxpt = b[4]
		}
	}
	return
}

// Since every packet has an id, this function comes up a lot.
func (m *PacketBuffer) PeekId() (id uint32) {
	if m != nil {
		if b := m.Bytes(); len(b) >= 9 {
			id += uint32(b[5]) << 24
			id += uint32(b[6]) << 16
			id += uint32(b[7]) << 8
			id += uint32(b[8])
		}
	}
	return
}

func (m *PacketBuffer) Bytes() []byte {
	return (*bytes.Buffer)(m).Bytes()
}

// Note: Do not use this method for strings, instead use the WriteString method.
func (m *PacketBuffer) Write(p []byte) (int, error) {
	return (*bytes.Buffer)(m).Write(p)
}

func (m *PacketBuffer) ReadFrom(r io.Reader) (int64, error) {
	return (*PacketBuffer)(m).ReadFrom(r)
}

func (m *PacketBuffer) Read(p []byte) (int, error) {
	return (*bytes.Buffer)(m).Read(p)
}

func (m *PacketBuffer) WriteTo(w io.Writer) (int64, error) {
	return (*bytes.Buffer)(m).WriteTo(w)
}

// Make sure you save room for the ID.
func (m *PacketBuffer) Reset() {
	(*bytes.Buffer)(m).Reset()
}

func (m *PacketBuffer) Grow(l uint32) {
	(*bytes.Buffer)(m).Grow(int(l))
}

func (m *PacketBuffer) WriteByte(b byte) {
	(*bytes.Buffer)(m).WriteByte(b)
}

func (m *PacketBuffer) ReadByte() (byte, error) {
	return (*bytes.Buffer)(m).ReadByte()
}

func (m *PacketBuffer) WriteUint32(u uint32) {
	wire := [4]byte{
		byte(u >> 24),
		byte(u >> 16),
		byte(u >> 8),
		byte(u),
	}
	(*bytes.Buffer)(m).Write(wire[:])
}

func (m *PacketBuffer) ReadUint32() (uint32, error) {
	var wire [4]byte
	if _, err := (*bytes.Buffer)(m).Read(wire[:]); err != nil {
		return 0, err
	}
	u := uint32(wire[0]) << 24
	u += uint32(wire[1]) << 16
	u += uint32(wire[2]) << 8
	u += uint32(wire[3])
	return u, nil
}

func (m *PacketBuffer) WriteUint64(u uint64) {
	m.WriteUint32(uint32(u >> 32))
	m.WriteUint32(uint32(u))
}

func (m *PacketBuffer) ReadUint64() (uint64, error) {
	upper, err := m.ReadUint32()
	if err != nil {
		return 0, err
	}
	lower, err := m.ReadUint32()
	if err != nil {
		return 0, err
	}
	return uint64(lower) + (uint64(upper) << 32), nil
}

func (m *PacketBuffer) WriteString(s string) {
	// Write the length of the string first.
	m.WriteUint32(uint32(len(s)))
	// Then write the string.
	(*bytes.Buffer)(m).WriteString(s)
}

func (m *PacketBuffer) ReadString() (string, error) {
	length, err := m.ReadUint32()
	if err != nil {
		return "", err
	}
	sb := make([]byte, length)
	if n, err := m.Read(sb); err != nil {
		return "", err
	} else if uint32(n) != length {
		return "", io.ErrShortBuffer
	}
	return string(sb), nil
}

// Use the String function in all byte slice cases except for extended requests!
func (m *PacketBuffer) WriteExtendedRequestData(b []byte) {
	(*bytes.Buffer)(m).Write(b)
}

func (m *PacketBuffer) ReadExtendedRequestData() ([]byte, error) {
	return ioutil.ReadAll(m)
}

func (m *PacketBuffer) WriteExtension(exts [][2]string) {
	// Extensions do *not* use an array count prefix
	for _, e := range exts {
		m.WriteString(e[0])
		m.WriteString(e[1])
	}
}

func (m *PacketBuffer) ReadExtensions() (exts [][2]string, err error) {
	for {
		name, err := m.ReadString()
		if err == io.EOF {
			err = nil
			break
		} else if err != nil {
			break
		}
		data, err := m.ReadString()
		if err == io.EOF {
			err = io.ErrUnexpectedEOF
			break
		} else if err != nil {
			break
		}
		exts = append(exts, [2]string{name, data})
	}
	return
}

func (m *PacketBuffer) WriteNames(names []Name) {
	// Names require a count
	m.WriteUint32(uint32(len(names)))
	for _, n := range names {
		m.WriteString(n.Path)
		m.WriteString(n.Long)
		m.WriteAttrs(n.Attrs)
	}
}

func (m *PacketBuffer) ReadNames() (names []Name, err error) {
	count, err := m.ReadUint32()
	if err != nil {
		return
	}
	for i := uint32(0); i < count; i++ {
		var name Name
		name.Path, err = m.ReadString()
		if err != nil {
			return
		}
		name.Long, err = m.ReadString()
		if err != nil {
			return
		}
		name.Attrs, err = m.ReadAttrs()
		if err != nil {
			return
		}
		names = append(names, name)
	}
	return
}

func (m *PacketBuffer) WriteAttrs(a Attrs) {
	m.WriteUint32(a.Flags)
	if 0 != a.Flags&FILEXFER_ATTR_SIZE {
		m.WriteUint64(a.Size)
	}
	if 0 != a.Flags&FILEXFER_ATTR_UIDGID {
		m.WriteUint32(a.Uid)
		m.WriteUint32(a.Gid)
	}
	if 0 != a.Flags&FILEXFER_ATTR_PERMISSIONS {
		m.WriteUint32(a.Permissions)
	}
	if 0 != a.Flags&FILEXFER_ATTR_ACMODTIME {
		m.WriteUint32(a.Atime)
		m.WriteUint32(a.Mtime)
	}
	if 0 != a.Flags&FILEXFER_ATTR_EXTENDED {
		// extended needs a count.
		m.WriteUint32(uint32(len(a.Extended)))
		for _, ex := range a.Extended {
			m.WriteString(ex[0])
			m.WriteString(ex[1])
		}
	}
}

func (m *PacketBuffer) ReadAttrs() (a Attrs, err error) {
	a.Flags, err = m.ReadUint32()
	if err != nil {
		return
	}
	if 0 != a.Flags&FILEXFER_ATTR_SIZE {
		a.Size, err = m.ReadUint64()
	}
	if 0 != a.Flags&FILEXFER_ATTR_UIDGID {
		a.Uid, err = m.ReadUint32()
		if err != nil {
			return
		}
		a.Gid, err = m.ReadUint32()
		if err != nil {
			return
		}
	}
	if 0 != a.Flags&FILEXFER_ATTR_PERMISSIONS {
		a.Permissions, err = m.ReadUint32()
		if err != nil {
			return
		}
	}
	if 0 != a.Flags&FILEXFER_ATTR_ACMODTIME {
		a.Atime, err = m.ReadUint32()
		if err != nil {
			return
		}
		a.Mtime, err = m.ReadUint32()
		if err != nil {
			return
		}
	}
	var count uint32
	if 0 != a.Flags&FILEXFER_ATTR_EXTENDED {
		// extended needs a count.
		count, err = m.ReadUint32()
		if err != nil {
			return
		}
		for i := uint32(0); i < count; i++ {
			var ext [2]string
			ext[0], err = m.ReadString()
			if err != nil {
				return
			}
			ext[1], err = m.ReadString()
			if err != nil {
				return
			}
			a.Extended = append(a.Extended, ext)
		}
	}
	return
}

// Ugly validation logic condensed into a convenience function. This function
// should be used as a gateway. It checks for problems that might cause a panic
// in later uses of the buffer.
func IsValidIncomingPacketHeader(buf *PacketBuffer) error {
	if buf == nil {
		return errors.New("Buffer is nil")
	}
	if !buf.IsValidLength() {
		return errors.New("Packet length doesn't match buffer.")
	}
	id := buf.PeekId()
	if 3 > id {
		return errors.New("Packet id must be greater than two.")
	}
	switch buf.PeekType() {
	case FXP_STATUS,
		FXP_HANDLE,
		FXP_DATA,
		FXP_NAME,
		FXP_ATTRS,
		FXP_VERSION,
		FXP_EXTENDEDREPLY:
		return nil
	default:
		return errors.New("Packet is not a valid server response type.")
	}
}

// The packet header is the first 9 bytes in a packet. There are many ways this
// can be wrong. This function does not check packets for correctness based on
// any state. For instance, if an INIT packet is called twice, this function
// has no way of knowing. All the uglyness is condensed into one function.
func IsValidOutgoingPacketHeader(buf *PacketBuffer) error {
	if buf == nil {
		return errors.New("Buffer is nil.")
	}
	if !buf.IsValidLength() {
		return errors.New("Packet buffer length is invalid.")
	}
	id := buf.PeekId()
	if 3 > id {
		return errors.New("Packet id must be greater than two.")
	}
	t := buf.PeekType()
	if (t > 2 && t <= 20) || t == FXP_INIT || t == FXP_EXTENDED {
		return nil
	}
	switch t {
	case FXP_STATUS,
		FXP_HANDLE,
		FXP_DATA,
		FXP_NAME,
		FXP_ATTRS,
		FXP_VERSION,
		FXP_EXTENDEDREPLY:
		return errors.New("Packet is not an outgoing packet type.")
	default:
		return errors.New("Packet type is unknown.")
	}
}
