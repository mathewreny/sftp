package sftp

var (
	StatusMessages = make(map[uint32]string)
	statusTypes    = []string{
		"OK", "EOF", "No Such File", "Permission Denied", "Failure",
		"Bad Message", "No Connection", "Connection Lost", "Op Unsupported",
	}
)

func init() {
	for i, msg := range statusTypes {
		StatusMessages[uint32(i)] = msg
	}
}

type FxpStatus struct {
	Id      uint32
	Code    uint32
	Message string
	LangTag string
}

func (s FxpStatus) Error() string {
	return "Status " + StatusMessages[s.Code] + ": " + s.Message + "."
}

// Returns true if buf is a status packet and the status code is not STATUS_OK.
//
// This function does not modify buf. It should be called before NewStatus when
// the client expects a status response from the server, and doesn't care about
// the message if the status is "ok".
func IsStatusAndNotOk(buf *PacketBuffer) bool {
	b := buf.Bytes()
	if 13 > len(b) {
		return false
	}
	isStatus := b[4] == FXP_STATUS
	isOk := 0 == b[9]|b[10]|b[11]|b[12]
	return isStatus && !isOk
}

// This constructor will consume the provided packet buffer. Do not use it after
// calling NewStatus.
func NewStatus(buf *PacketBuffer) (s FxpStatus, valid bool) {
	buf.ReadUint32() // Throw away length
	if t, _ := buf.ReadByte(); t == FXP_STATUS {
		s.Id, _ = buf.ReadUint32()
		s.Code, _ = buf.ReadUint32()
		s.Message, _ = buf.ReadString()
		s.LangTag, _ = buf.ReadString()
		valid = true
	}
	bufPool.Put(buf)
	return
}

// The buf parameter is optional. When provided, the buffer is recycled and the
// status packet's id is set to the buffer's id. When buf is nil, a buffer is
// grabbed from the global pool. Custom status codes are allowed. They should be
// present in the "UserDefinedStatusCodes" package variable.
func BufferStatusCode(buf *PacketBuffer, code uint32) *PacketBuffer {
	return BufferStatus(buf, code, 0, StatusMessages[code], "en")
}

func BufferStatus(b *PacketBuffer, c, id uint32, m, lang string) *PacketBuffer {
	if b == nil {
		b = bufPool.Get().(*PacketBuffer)
	} else if id == 0 && 9 <= len(b.Bytes()) {
		b.ReadUint32() // throw away old length
		b.ReadByte()   // throw away old type
		id, _ = b.ReadUint32()
	}
	var length = uint32(1 + 4 + 4 + 4 + len(m) + 4 + len(lang))
	b.Reset()
	b.Grow(4 + length)
	b.WriteUint32(length)
	b.WriteByte(FXP_STATUS)
	b.WriteUint32(id)
	b.WriteUint32(c)
	b.WriteString(m)
	b.WriteString(lang)
	return b
}
