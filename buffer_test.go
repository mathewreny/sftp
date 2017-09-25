package sftp

import "testing"

func TestRW(t *testing.T) {
	buf := NewBuffer()
	var hi [5]byte
	if _, err := buf.Write([]byte("Hello")); err != nil {
		t.Fatal(err)
	} else if _, err = buf.Read(hi[:]); err != nil {
		t.Fatal(err)
	} else if string(hi[:]) != "Hello" {
		t.Fatal("Got incorrect string")
	}
}
