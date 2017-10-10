package sftp

import (
	"testing"
)

var tc struct {
	Pipe struct {
		Out io.ReadCloser
		In  io.WriteCloser
	}
	Conn *Conn
}

func tcReset() {
	if tc.Conn != nil {
		tc.Conn.Close()
	}
	if tc.Pipe.Out != nil {
		tc.Pipe.Out.Close()
	}
	if tc.Pipe.In != nil {
		tc.Pipe.In.Close()
	}
	tc.Conn = new(Conn)
	tc.Conn.ingress = make(chan *PacketBuffer, 20)
	tc.Conn.response.chans = make(map[uint32]chan<- *PacketBuffer)
	tc.Pipe.Out, tc.Conn.w = io.Pipe()
	tc.Conn.r, tc.Pipe.In = io.Pipe()
}

func init() {
	tcReset()
}

// TODO finish these tests. Move to TDD.
