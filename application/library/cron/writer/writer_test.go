package writer

import (
	"testing"
	"bufio"
	"io"
	"bytes"
	"strings"
	"unicode/utf8"

	"github.com/webx-top/echo/testing/test"
)

var expectedCn = `中国中国中国中国中国中国中国中国中
......
世界世界世界世界世界世界世界世界`

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
	cn := []byte(`中1`)
	test.Eq(t, true, utf8.RuneStart(cn[0]))
	test.Eq(t, false, utf8.RuneStart(cn[1]))
	test.Eq(t, false, utf8.RuneStart(cn[2]))
	test.Eq(t, true, utf8.RuneStart(cn[3]))


	w = New(100)
	s = strings.Repeat("中国",30)+`CCCCCCCCCCCCCCCCCCCCCCCCCC`+strings.Repeat("你好",1000)+strings.Repeat("世界",50)
	b = bytes.NewBufferString(s)
	r = bufio.NewReader(b)
	_, err = io.Copy(w, r)
	if err != nil {
		t.Fatal(err)
	}
	bz = w.Bytes()
	test.Eq(t, 8, bytes.Count(bz,[]byte(`中国`)))
	test.Eq(t, 8, bytes.Count(bz,[]byte(`世界`)))
	test.Eq(t, true, utf8.RuneStart(bz[48]))
	test.Eq(t, expectedCn, string(bz))
	w.Write([]byte(`团结团结团结团结`))
	bz = w.Bytes()
	test.Eq(t, 4, bytes.Count(bz,[]byte(`团结`)))
}