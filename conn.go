package sftp

import (
	"io"
)

// Interface composed of functions used from golang.org/x/crypto/ssh#Session
type SshSession interface {
	Close() error
	RequestSubsystem(subsystem string) error
	StdinPipe() (io.WriteCloser, error)
	StdoutPipe() (io.Reader, error)
}

type Conn struct {
	w       io.Writer
	r       io.Reader
	closers []io.Closer
}

func (c *Conn) Write(p []byte) (int, error) { return c.w.Write(p) }
func (c *Conn) Read(p []byte) (int, error)  { return c.r.Read(p) }

func (c *Conn) Close() (err error) {
	for _, x := range c.closers {
		if err != nil {
			x.Close()
		} else {
			err = x.Close()
		}
	}
	return
}

func NewConn(s SshSession) (*Conn, error) {
	err := s.RequestSubsystem("sftp")
	if err != nil {
		return nil, err
	}
	wc, err := s.StdinPipe()
	if err != nil {
		return nil, err
	}
	r, err := s.StdoutPipe()
	if err != nil {
		return nil, err
	}
	return &Conn{r: r, w: wc, closers: []io.Closer{wc, s}}, nil
}
