package sftputil

import (
	"errors"
	"io"
)

// An easy way to create a session for testing.
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
