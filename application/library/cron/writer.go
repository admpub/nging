package cron

import "bytes"

var dot6bytes = []byte(` ...... `)

func NewCmdRec(max uint64) *cmdRec {
	return &cmdRec{
		buf:  new(bytes.Buffer),
		max:  max / 2,
		last: []byte{},
	}
}

type cmdRec struct {
	buf   *bytes.Buffer
	max   uint64
	start uint64
	end   uint64
	last  []byte
}

func (c *cmdRec) Write(p []byte) (n int, err error) {
	if c.start < c.max {
		n, err = c.buf.Write(p)
		c.start += uint64(n)
		return
	}
	size := uint64(len(p))
	if c.end > c.max {
		if c.max > size {
			c.last = append(c.last[0:c.max-size], p...)
		} else if c.max == size {
			c.last = p
		} else {
			start := size - c.max
			c.last = p[start:]
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
		s += ` ...... ` + string(c.last)
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
