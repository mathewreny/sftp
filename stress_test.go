package sftp

import (
	//	"bufio"
	"fmt"
	"github.com/mathewreny/sftp/sftputil"
	"io"
	//	"os"
	"sync"
	"testing"
	//	"time"
)

func TestFastStress(t *testing.T) {
	var pr io.ReadCloser
	var pw io.WriteCloser
	s := new(sftputil.TestSession)
	s.R, pw = io.Pipe()
	pr, s.WC = io.Pipe()
	c, err := NewClient(s)
	if err != nil {
		fmt.Println(err)
		t.Fatal("Client didn't start.")
	}

	//	pkts := make(chan *Buffer)
	go func() {
		for {
			buf := NewBuffer()
			CopyPacket(buf, pr)
			fmt.Printf("--RAW: %x\n", buf.Bytes())
			fmt.Println("------SENT PACKET", buf.PeekId())
			buf.ReadUint32()
			buf.ReadByte()
			id, _ := buf.ReadUint32()
			a := Attrs{}
			buf.Reset()
			buf.WriteUint32(5 + a.Len())
			buf.WriteByte(FXP_ATTRS)
			buf.WriteUint32(id)
			buf.WriteAttrs(a)
			io.Copy(pw, buf)
		}
	}()
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			for sends := 0; sends < 30; sends++ {
				if _, err := c.Stat(""); err != nil {
					t.Fatal(err)
				}
				fmt.Println("STAT.")
			}
			wg.Done()
		}()
	}
	wg.Wait()
	pw.Close()
	fmt.Println("DONE WITH EVERYTHIGN")
}
