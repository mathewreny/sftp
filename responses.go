package sftp

import (
	"errors"
)

func parseVersionResponse(buf *Buffer) (extended [][2]string, err error) {
	var length, version uint32
	if buf == nil {
		return nil, errors.New("Nil buffer provided.")
	} else if IsStatusAndNotOk(buf) {
		err, _ = ParseStatus(buf)
		return
	}
	if length, err = buf.ReadUint32(); err != nil {
		return nil, err
	} else if length != uint32(len(buf.Bytes())) {
		return nil, errors.New("Buf has invalid length.")
	}
	var t byte
	if t, err = buf.ReadByte(); err != nil {
		return nil, err
	} else if t != FXP_VERSION {
		return nil, errors.New("Invalid packet response type.")
	}
	if version, err = buf.ReadUint32(); err != nil {
		return nil, err
	} else if version != 3 {
		return nil, errors.New("Invalid SFTP version.")
	}
	for 0 != len(buf.Bytes()) {
		var temp [2]string
		if temp[0], err = buf.ReadString(); err != nil {
			return
		}
		if temp[1], err = buf.ReadString(); err != nil {
			return
		}
		extended = append(extended, temp)
	}
	return
}

func parseHandleResponse(buf *Buffer, c *Client) (h Handle, err error) {
	if buf == nil {
		err = errors.New("Nil buffer provided.")
		return
	}
	if IsStatusAndNotOk(buf) {
		err, _ = ParseStatus(buf)
		return
	}
	var length uint32
	length, err = buf.ReadUint32()
	if err != nil {
		return
	} else if length != uint32(len(buf.Bytes())) {
		err = errors.New("Buf has invalid length.")
		return
	}
	var ptype byte
	ptype, err = buf.ReadByte()
	if err != nil {
		return
	} else if ptype != FXP_HANDLE {
		err = errors.New("Invalid packet response type.")
		return
	}
	var str string
	str, err = buf.ReadString()
	if err != nil {
		return
	}
	h.h = str
	h.client = c
	return
}

func parseStatusResponse(buf *Buffer) error {
	if buf == nil {
		return errors.New("Nil buffer provided.")
	}
	if IsStatusAndNotOk(buf) {
		err, _ := ParseStatus(buf)
		return err
	}
	return nil
}

func parseDataResponse(buf *Buffer) ([]byte, error) {
	if buf == nil {
		return nil, errors.New("Nil buffer provided.")
	}
	if IsStatusAndNotOk(buf) {
		s, _ := ParseStatus(buf)
		return nil, s
	}
	if length, err := buf.ReadUint32(); err != nil {
		return nil, err
	} else if length != uint32(len(buf.Bytes())) {
		return nil, errors.New("Buf has invalid length.")
	}
	if t, err := buf.ReadByte(); err != nil {
		return nil, err
	} else if t != FXP_DATA {
		return nil, errors.New("Invalid packet response type.")
	}
	if _, err := buf.ReadUint32(); err != nil {
		return nil, err
	}
	strdata, err := buf.ReadString()
	if err != nil {
		return nil, err
	}
	data := make([]byte, len(strdata))
	copy(data, strdata)
	return data, nil
}

func parseNameResponse(buf *Buffer) (names []Name, err error) {
	if buf == nil {
		err = errors.New("Nil buffer provided.")
		return
	}
	if IsStatusAndNotOk(buf) {
		err, _ = ParseStatus(buf)
		return
	}
	var length uint32
	if length, err = buf.ReadUint32(); err != nil {
		return
	} else if length != uint32(len(buf.Bytes())) {
		err = errors.New("Buf has invalid length.")
		return
	}
	var t byte
	if t, err = buf.ReadByte(); err != nil {
		return
	} else if t != FXP_NAME {
		err = errors.New("Invalid packet response type.")
		return
	}
	if _, err = buf.ReadUint32(); err != nil { // throw away id
		return
	}
	names, err = buf.ReadNames()
	return
}

func parseAttrsResponse(buf *Buffer) (a Attrs, err error) {
	if buf == nil {
		err = errors.New("Nil buffer provided.")
		return
	}
	if IsStatusAndNotOk(buf) {
		err, _ = ParseStatus(buf)
		return
	}
	var length uint32
	if length, err = buf.ReadUint32(); err != nil {
		return
	} else if length != uint32(len(buf.Bytes())) {
		err = errors.New("Buf has invalid length.")
		return
	}
	var t byte
	if t, err = buf.ReadByte(); err != nil {
		return
	} else if t != FXP_ATTRS {
		err = errors.New("Invalid packet response type.")
		return
	}
	if _, err = buf.ReadUint32(); err != nil {
		return
	}
	a, err = buf.ReadAttrs()
	return
}

func parseExtendedReplyResponse(buf *Buffer) (data []byte, err error) {
	if buf == nil {
		err = errors.New("Nil buffer provided.")
		return
	}
	if IsStatusAndNotOk(buf) {
		err, _ = ParseStatus(buf)
		return
	}
	var length uint32
	if length, err = buf.ReadUint32(); err != nil {
		return
	} else if length != uint32(len(buf.Bytes())) {
		err = errors.New("Buf has invalid length.")
		return
	}
	var t byte
	if t, err = buf.ReadByte(); err != nil {
		return
	} else if t != FXP_EXTENDEDREPLY {
		err = errors.New("Invalid packet response type.")
		return
	}
	// copy the raw data to facillitate garbage collection.
	bufdata := buf.Bytes()
	data = make([]byte, len(bufdata))
	copy(data, bufdata)
	return
}
