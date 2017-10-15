package sftp

import (
	"errors"
	"fmt"
	"github.com/mathewreny/sftp/sftputil"
	"io"
	"sync"
	"testing"
	"time"
)

type versionWriter struct {
	sync.Once
	b    []byte
	done chan struct{}
}

func (w *versionWriter) Write(p []byte) (n int, err error) {
	n = len(p)
	w.b = append(w.b, p...)
	if len(w.b) >= 9 {
		go w.Do(func() {
			//<-time.After(1 * time.Second)
			close(w.done)
			<-time.After(1 * time.Second)
		})
		if w.b[4] != FXP_INIT {
			err = errors.New("Did not send FXP_INIT packet type.")
			return
		} else if 3 != w.b[8] && 0 != w.b[5]+w.b[6]+w.b[7] {
			err = errors.New("Did not send SFTP version 3 identifier.")
			return
		}
	}
	return
}
func (w *versionWriter) Close() error {
	select {
	case <-w.done:
	default:
		panic("done has not been closed.")
	}
	return nil
}

func newVersionTestSession() Session {
	w := &versionWriter{done: make(chan struct{})}
	r, pw := io.Pipe()
	go func() {
		<-w.done
		b := NewBuffer()
		b.WriteUint32(5)
		b.WriteByte(FXP_VERSION)
		b.WriteUint32(3)
		pw.Write(b.Bytes())
	}()
	return &sftputil.TestSession{R: r, WC: w}
}

func TestVersion(t *testing.T) {
	s := newVersionTestSession()
	c, err := NewClient(s)
	if err != nil {
		t.Fatal("Client couldn't start. NewClient failed.")
	}
	_, err = c.Init(nil)
	if err != nil {
		fmt.Println(err.Error())
		t.Fatal("Client couldn't Init.")
	}
}

func TestVersion2(t *testing.T) {
	packet := NewBuffer()
	packet.WriteUint32(5)
	packet.WriteByte(FXP_VERSION)
	packet.WriteUint32(3)

}
