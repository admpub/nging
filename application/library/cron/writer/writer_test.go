package writer

import (
	"testing"
	"bufio"
	"io"
	"bytes"
	"strings"

	"github.com/webx-top/echo/testing/test"
)

func TestWrite(t *testing.T) {
	w := New(100)
	s := strings.Repeat("A",50)+`CCCCCCCCCCCCCCCCCCCCCCCCCC`+strings.Repeat("B",5000)+strings.Repeat("D",50)
	b := bytes.NewBufferString(s)
	r := bufio.NewReader(b)
	_, err := io.Copy(w, r)
	if err != nil {
		t.Fatal(err)
	}
	bz := w.Bytes()
	test.Eq(t, 50, bytes.Count(bz,[]byte(`A`)))
	test.Eq(t, 50, bytes.Count(bz,[]byte(`D`)))

	w2 := New(100)
	sb := []byte(s)
	_,err = w2.Write(sb[0:1000])
	if err != nil {
		t.Fatal(err)
	}
	_,err = w2.Write(sb[1000:2000])
	if err != nil {
		t.Fatal(err)
	}
	_,err = w2.Write(sb[2000:3000])
	if err != nil {
		t.Fatal(err)
	}
	_,err = w2.Write(sb[3000:])
	if err != nil {
		t.Fatal(err)
	}
	bz = w2.Bytes()
	test.Eq(t, 50, bytes.Count(bz,[]byte(`A`)))
	test.Eq(t, 50, bytes.Count(bz,[]byte(`D`)))
}