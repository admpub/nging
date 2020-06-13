/*
   Nging is a toolbox for webmasters
   Copyright (C) 2018-present  Wenhui Shen <swh@admpub.com>

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU Affero General Public License as published
   by the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Affero General Public License for more details.

   You should have received a copy of the GNU Affero General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package writer

import (
	"bytes"
	"io"
	"strings"
	"unicode/utf8"
)

var (
	dot6str   = "\n" + `......` + "\n"
	dot6bytes = []byte(dot6str)

	//NotRecordPrefixFlag 不记录日志的前缀标识
	NotRecordPrefixFlag = `--/ignore/--`
)

type OutputWriter interface {
	io.Writer
	String() string
	Bytes() []byte
}

func New(max uint64) *cmdRec {
	return &cmdRec{
		buf:  new(bytes.Buffer),
		max:  max / 2,
		last: []byte{},
	}
}

type cmdRec struct {
	buf    *bytes.Buffer
	max    uint64
	start  uint64
	end    uint64
	last   []byte
	ignore bool
}

func GetRuneStartIndex(end int, p []byte) int {
	n := len(p)
	for ; end < n; end++ {
		if utf8.RuneStart(p[end]) {
			break
		}
	}
	return end
}

func (c *cmdRec) Write(p []byte) (n int, err error) {
	if c.ignore {
		n = len(p)
		return
	}
	if c.start == 0 && strings.HasPrefix(string(p), NotRecordPrefixFlag) {
		c.ignore = true
		n = len(p)
		return
	}
	n = len(p)
	size := uint64(n)
	if c.start < c.max {
		remain := c.max - c.start
		if remain < uint64(n) {
			end := int(remain)
			end = GetRuneStartIndex(end, p)
			rp := p[end:]
			p = p[:end]
			var actualN int
			actualN, err = c.buf.Write(p)
			c.start += uint64(actualN)
			p = rp
			size = uint64(len(p))
		} else {
			n, err = c.buf.Write(p)
			c.start += uint64(n)
			return
		}
	}
	if c.end >= c.max {
		if c.max > size {
			end := int(c.max - size)
			end = GetRuneStartIndex(end, c.last)
			c.last = append(c.last[0:end], p...)
		} else if c.max == size {
			c.last = p
		} else {
			end := int(size - c.max)
			end = GetRuneStartIndex(end, p)
			c.last = p[end:]
		}
		return
	}
	remain := c.max - c.end
	if remain < size {
		if c.max > size {
			end := int(c.max - size)
			end = GetRuneStartIndex(end, c.last)
			c.last = append(c.last[0:end], p...)
		}else if c.max == size {
			c.last = p
		} else {
			end := int(size - c.max)
			end = GetRuneStartIndex(end, p)
			c.last = p[end:]
		}
		c.end = uint64(len(c.last))
		return
	}
	c.end += size
	c.last = append(c.last, p...)
	return
}

// String returns the contents of the unread portion of the buffer
// as a string. If the Buffer is a nil pointer, it returns "<nil>".
func (c *cmdRec) String() string {
	if c.buf == nil {
		// Special case, useful in debugging.
		return string(c.last)
	}
	s := c.buf.String()
	if len(s) > 0 && len(c.last) > 0 {
		s += dot6str + string(c.last)
	}
	return s
}

func (c *cmdRec) Bytes() []byte {
	if c.buf == nil {
		// Special case, useful in debugging.
		return c.last
	}
	b := c.buf.Bytes()
	if len(b) > 0 && len(c.last) > 0 {
		b = append(b, dot6bytes...)
	}
	return append(b, c.last...)
}
