package sftp

import (
	//	"bufio"
	"fmt"
	"github.com/mathewreny/sftp/sftputil"
	"io"
	//	"os"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestConcurrentEgressSequentialIngressStress(t *testing.T) {
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
			if _, err := CopyPacket(buf, pr); err != nil {
				break
			}
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
			bufPool.Put(buf)
		}
	}()
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			for sends := 0; sends < 100; sends++ {
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

func TestConcurrentEgressRandomDelayedIngressStress(t *testing.T) {
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

	pkts := make(chan *Buffer, 30)
	go func() {
		for {
			buf := NewBuffer()
			if _, err := CopyPacket(buf, pr); err != nil {
				break
			}
			go func() {
				<-time.After(25 + time.Duration(rand.Int()%75)*time.Millisecond)
				pkts <- buf
			}()
		}
		//close(pkts)
	}()
	write := make(chan *Buffer)
	go func() {
		bufq := newBufqueue()
		for {
			if bufq.Peek() == nil {
				bufq.Push(<-pkts)
			} else {
				select {
				case buf := <-pkts:
					bufq.Push(buf)
				case write <- bufq.Peek():
					bufq.Pop()
				}
			}
		}
	}()
	go func() {
		for buf := range write {
			buf.ReadUint32()
			buf.ReadByte()
			id, _ := buf.ReadUint32()
			fmt.Print(",", id)

			a := Attrs{}
			buf.Reset()
			buf.WriteUint32(5 + a.Len())
			buf.WriteByte(FXP_ATTRS)
			buf.WriteUint32(id)
			buf.WriteAttrs(a)
			io.Copy(pw, buf)
			bufPool.Put(buf)
		}
	}()
	var wg sync.WaitGroup
	for i := 0; i < 500; i++ {
		wg.Add(1)
		go func() {
			for sends := 0; sends < 100; sends++ {
				if _, err := c.Stat(""); err != nil {
					t.Fatal(err)
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
	pw.Close()
	fmt.Println("DONE WITH EVERYTHIGN")
}
