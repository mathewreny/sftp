package sftputil

import (
	"errors"
	"io"
	"sync"
)

// An easy way to create a session.
// Just supply an io.Reader and io.WriteCloser.
type TestSession struct {
	R  io.Reader
	WC io.WriteCloser
}

func (ts *TestSession) Close() error {
	return ts.WC.Close()
}

func (ts *TestSession) StdinPipe() (io.WriteCloser, error) {
	return ts.WC, nil
}

func (ts *TestSession) StdoutPipe() (io.Reader, error) {
	return ts.R, nil
}

func (ts *TestSession) RequestSubsystem(subsystem string) error {
	if subsystem != "sftp" {
		return errors.New("Improper subsystem")
	} else {
		return nil
	}
}

// Will reply with the provided Packet.
type PacketSession struct {
	sync.Once
	rc       io.ReadCloser
	Response io.Reader
	// Packet that should be received. The ID is not checked.
	// If empty, the prototype packet is not checked.
	Prototype []byte
	proto     struct {
		sync.Mutex
		b          []byte
		pastheader bool
		unread     []byte
	}
	ready chan struct{}
}

func (ps *PacketSession) init() {
	var pw io.Writer
	ps.rc, pw = io.Pipe()
	ps.ready = make(chan struct{})
	ps.proto.b = append([]byte{}, ps.Prototype...)
	copy(ps.proto.b[5:9], []byte{0, 0, 0, 0}) // zero out the ID.
	go func() {
		<-ps.ready
		io.Copy(pw, ps.Response)
	}()
}

func (ps *PacketSession) Close() error {
	ps.Do(ps.init)
	return ps.rc.Close()
}

func (ps *PacketSession) StdinPipe() (io.WriteCloser, error) {
	ps.Do(ps.init)
	return ps, nil
}

func (ps *PacketSession) StdoutPipe() (io.Reader, error) {
	ps.Do(ps.init)
	return ps.rc, nil
}

func (ps *PacketSession) RequestSubsystem(subsystem string) error {
	if subsystem != "sftp" {
		return errors.New("Subsystem not requested.")
	}
	return nil
}

// Don't directly call this when testing.
func (ps *PacketSession) Write(p []byte) (n int, err error) {
	ps.proto.Lock()
	defer ps.proto.Unlock()
	n = len(p)
	if len(ps.proto.b) < n+len(ps.proto.unread) {
		err = errors.New("Prototype does not match outgoing packet. Prototype too small.")
		return
	}

	if ps.proto.pastheader {
		if string(ps.proto.b[:n]) != string(p) {
			err = errors.New("Prototype does not match outgoing packet.")
		} else {
			ps.proto.b = ps.proto.b[n:]
		}
	} else {
		ps.proto.unread = append(ps.proto.unread, p...)
		urlen := len(ps.proto.unread)
		if urlen >= 9 {
			ps.proto.pastheader = true
			copy(ps.proto.unread[5:9], []byte{0, 0, 0, 0})
			if string(ps.proto.b[:urlen]) != string(ps.proto.unread) {
				err = errors.New("Prototype does not match outgoing packet.")
			}
			ps.proto.b = ps.proto.b[urlen:]
		}
	}
	return
}
