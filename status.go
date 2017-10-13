package sftp

import (
	"errors"
	"io"
)

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

type Status struct {
	Id      uint32
	Code    uint32
	Message string
	LangTag string
}

func (s Status) Error() string {
	return "Status " + StatusMessages[s.Code] + ": " + s.Message + "."
}

// Returns true if buf is a status packet and the status code is not STATUS_OK.
//
// This function does not modify buf. It should be called before NewStatus when
// the client expects a status response from the server, and doesn't care about
// the message if the status is "ok".
func IsStatusAndNotOk(buf *Buffer) bool {
	b := buf.Bytes()
	if 13 > len(b) {
		return false
	}
	isStatus := b[4] == FXP_STATUS
	isOk := 0 == b[9]|b[10]|b[11]|b[12]
	return isStatus && !isOk
}

func ParseStatus(r io.Reader) (s Status, err error) {
	buf, ok := r.(*Buffer)
	if !ok {
		buf = NewBuffer()
		_, err = CopyPacket(buf, r)
		if err != nil {
			return
		}
	}
	s, ok = parseStatus(buf)
	if !ok {
		err = errors.New("Reader did not provide a status packet.")
	}
	return
}

func parseStatus(buf *Buffer) (s Status, valid bool) {
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

func (s Status) Buffer() *Buffer {
	buf := NewBuffer()
	buf.WriteUint32(21 + uint32(len(s.Message)+len(s.LangTag)))
	buf.WriteByte(FXP_STATUS)
	buf.WriteUint32(s.Id)
	buf.WriteUint32(s.Code)
	buf.WriteString(s.Message)
	buf.WriteString(s.LangTag)
	return buf
}

// The buf parameter is optional. When provided, the buffer is recycled and the
// status packet's id is set to the buffer's id. When buf is nil, a buffer is
// grabbed from the global pool. Custom status codes are allowed. They should be
// present in the "UserDefinedStatusCodes" package variable.
func NewStatus(code uint32) Status {
	return Status{Code: code, Message: StatusMessages[code], LangTag: "en"}
}
