package sftp

import (
	"github.com/mathewreny/sftp/sftputil"
	"io"
	"log"
	"math/rand"
	"sync"
	"testing"
	"time"
)

const verylongpath = "hguaywgflaywgelfiauwglfiugwilyug4lfi3ygfihaglfa3iglahgfliay3gfp4ia3ubfpciua3byfpc9ab3pftybv3iuayt3lvr7capv3ry8bac4yiur4iby4acbityr4pi7ta3pr4937t4bcpr397t4abcpra973tpr4a7tprait3orliyuag3br4uya3gor8a3t4bora3874torcva3i7rv937trvoi7atyro9374rtyoa837rtoaivuygfoiauygfoa8374tgr8bc7ao874or74rtoai4v734otva74rtvoa73taorc87taorciytgrbgorcaibgr3boca2t3rvkiuy23rtkvaiuyfgobq8ytgo8q7t2o87rtaowiuetyrlaw4itucalwiu4ytrvali7t3roalirtfor8catoraliytgfweliygvaliuywvcialgbgzjzgyflq7t"

func TestConcurrentEgressSequentialIngressStress(t *testing.T) {
	var pr io.ReadCloser
	var pw io.WriteCloser
	s := new(sftputil.TestSession)
	s.R, pw = io.Pipe()
	pr, s.WC = io.Pipe()
	c, err := NewClient(s)
	if err != nil {
		t.Fatal("Client didn't start.", err)
	}

	//	pkts := make(chan *Buffer)
	go func() {
		for {
			buf := NewBuffer()
			if _, err := CopyPacket(buf, pr); err != nil {
				break
			}
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
	numgo := 50
	numsends := 1000
	log.Println("Starting", numgo, "goroutines sending", numsends, "requests each.")
	var wg sync.WaitGroup
	for i := 0; i < numgo; i++ {
		wg.Add(1)
		gx := i
		go func() {
			for sends := 0; sends < numsends; sends++ {
				if _, err := c.Stat(verylongpath); err != nil {
					t.Fatal(err)
				}
			}
			log.Println("Goroutine", gx, " sent/received all", numsends, "requests.")
			wg.Done()
		}()
	}
	wg.Wait()
	pw.Close()
	log.Println("Finished with fast sequential (FIFO) sending test.")
	log.Println()
}

func TestConcurrentEgressRandomDelayedIngressStress(t *testing.T) {
	var pr io.ReadCloser
	var pw io.WriteCloser
	s := new(sftputil.TestSession)
	s.R, pw = io.Pipe()
	pr, s.WC = io.Pipe()
	c, err := NewClient(s)
	if err != nil {
		t.Fatal("Client didn't start.", err)
	}
	const numgo, numsends = 500, 100
	log.Println("STARTING", numgo, "GOROUTINES SENDING", numsends, "REQUESTS EACH")
	const timemin, timedelta = 25, 75
	log.Println("RESPONSES RANDOMLY TAKE", timemin, "TO", timemin+timedelta, "MILLISECONDS.")

	pkts := make(chan *Buffer, 30)
	go func() {
		for {
			buf := NewBuffer()
			if _, err := CopyPacket(buf, pr); err != nil {
				break
			}
			go func() {
				mills := timemin + (time.Duration(rand.Int()) % timedelta)
				<-time.After(mills * time.Millisecond)
				pkts <- buf
			}()
		}
		//close(pkts)
	}()
	write := make(chan *Buffer, 30)
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

			if 0 == id%1000 {
				log.Println("Reached packet", id)
			}

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
	for i := 0; i < numgo; i++ {
		wg.Add(1)
		go func() {
			for sends := 0; sends < numsends; sends++ {
				if _, err := c.Stat(verylongpath); err != nil {
					t.Fatal(err)
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
	pw.Close()
	log.Println("DONE WITH OUT OF ORDER TEST")
	log.Println()
}

func TestConcurrentEgressRandomDelayedIngressBigQueueStress(t *testing.T) {
	var pr io.ReadCloser
	var pw io.WriteCloser
	s := new(sftputil.TestSession)
	s.R, pw = io.Pipe()
	pr, s.WC = io.Pipe()
	c, err := NewClient(s)
	if err != nil {
		t.Fatal("Client didn't start.", err)
	}

	const numgo, numsends = 5000, 10
	log.Println("STARTING", numgo, "GOROUTINES SENDING", numsends, "REQUESTS EACH")
	const timemin, timedelta = 100, 400
	log.Println("RESPONSES RANDOMLY TAKE", timemin, "TO", timemin+timedelta, "MILLISECONDS.")
	pkts := make(chan *Buffer, 30)
	go func() {
		for {
			buf := NewBuffer()
			if _, err := CopyPacket(buf, pr); err != nil {
				break
			}
			go func() {
				mills := timemin + time.Duration(rand.Int()%timedelta)
				<-time.After(mills * time.Millisecond)
				pkts <- buf
			}()
		}
		//close(pkts)
	}()
	write := make(chan *Buffer, 30)
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

			if 0 == id%1000 {
				log.Println("Reached packet", id)
			}

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
	for i := 0; i < numgo; i++ {
		wg.Add(1)
		go func() {
			for sends := 0; sends < numsends; sends++ {
				if _, err := c.Stat(verylongpath); err != nil {
					t.Fatal(err)
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
	pw.Close()
	log.Println("DONE WITH LARGE QUEUE OUT OF ORDER TEST")
	log.Println()
}
