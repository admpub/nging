package telnet

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
	"sync/atomic"
	"time"
	"unicode"
)

const (
	CR = byte('\r')
	LF = byte('\n')
)

const (
	// SE                  240    End of subnegotiation parameters.
	cmdSE = 240
	// NOP                 241    No operation.
	cmdNOP = 241
	// Data Mark           242    The data stream portion of a Synch.
	//                            This should always be accompanied
	//                            by a TCP Urgent notification.
	cmdData = 242

	// Break               243    NVT character BRK.
	cmdBreak = 243
	// Interrupt Process   244    The function IP.
	cmdIP = 244
	// Abort output        245    The function AO.
	cmdAO = 245
	// Are You There       246    The function AYT.
	cmdAYT = 246
	// Erase character     247    The function EC.
	cmdEC = 247
	// Erase Line          248    The function EL.
	cmdEL = 248
	// Go ahead            249    The GA signal.
	cmdGA = 249
	// SB                  250    Indicates that what follows is
	//                            subnegotiation of the indicated
	//                            option.
	cmdSB = 250 // FA

	// WILL (option code)  251    Indicates the desire to begin
	//                            performing, or confirmation that
	//                            you are now performing, the
	//                            indicated option.
	cmdWill = 251 // FB
	// WON'T (option code) 252    Indicates the refusal to perform,
	//                            or continue performing, the
	//                            indicated option.
	cmdWont = 252 // FC
	// DO (option code)    253    Indicates the request that the
	//                            other party perform, or
	//                            confirmation that you are expecting
	//                            the other party to perform, the
	//                            indicated option.
	cmdDo = 253 // FD
	// DON'T (option code) 254    Indicates the demand that the
	//                            other party stop performing,
	//                            or confirmation that you are no
	//                            longer expecting the other party
	//                            to perform, the indicated option.
	cmdDont = 254 // FE

	// IAC                 255    Data Byte 255.
	cmdIAC = 255 //FF

)

const (

	// 1(0x01)    回显(echo)
	optEcho = 1
	// 3(0x03)    抑制继续进行(传送一次一个字符方式可以选择这个选项)
	optSuppressGoAhead = 3
	// 24(0x18)   终端类型
	optWndType = 24
	// 31(0x1F)   窗口大小
	optWndSize = 31
	// 32(0x20)   终端速率
	optRate = 32

// 33(0x21)   远程流量控制
// 34(0x22)   行方式
// 36(0x24)   环境变量
)

// Conn implements net.Conn interface for Telnet protocol plus some set of
// Telnet specific methods.
type Conn struct {
	net.Conn
	r      *bufio.Reader
	ticker *time.Ticker

	is_closed     int32
	unixWriteMode bool

	cliSuppressGoAhead bool
	cliEcho            bool

	rows, columns byte
}

func NewConn(conn net.Conn) (*Conn, error) {
	return NewConnWithRead(conn, conn)
}

func NewConnWithRead(conn net.Conn, rd io.Reader) (*Conn, error) {
	c := Conn{
		Conn: conn,
		r:    bufio.NewReaderSize(rd, 256),
	}
	c.is_closed = 0
	c.ticker = time.NewTicker(1 * time.Second)
	go c.run()
	return &c, nil
}

func Dial(network, addr string) (*Conn, error) {
	conn, err := net.Dial(network, addr)
	if err != nil {
		return nil, err
	}
	return NewConn(conn)
}

func DialTimeout(network, addr string, timeout time.Duration) (*Conn, error) {
	conn, err := net.DialTimeout(network, addr, timeout)
	if err != nil {
		return nil, err
	}
	return NewConn(conn)
}

func (c *Conn) run() {
	defer func() {
		if o := recover(); nil != o {
			fmt.Println(o)
		}
		c.ticker.Stop()
	}()

	for 0 == atomic.LoadInt32(&c.is_closed) {
		_, ok := <-c.ticker.C
		if !ok {
			break
		}
		if e := c.noop(); nil != e {
			break
		}
	}
}

func (c *Conn) noop() error {
	_, e := c.Write([]byte{cmdNOP})
	return e
}

func (c *Conn) Close() error {
	if !atomic.CompareAndSwapInt32(&c.is_closed, 0, 1) {
		return nil
	}

	return c.Conn.Close()
}

// SetUnixWriteMode sets flag that applies only to the Write method.
// If set, Write converts any '\n' (LF) to '\r\n' (CR LF).
func (c *Conn) SetUnixWriteMode(uwm bool) {
	c.unixWriteMode = uwm
}

func (c *Conn) do(option byte) error {
	//log.Println("do:", option)
	_, err := c.Conn.Write([]byte{cmdIAC, cmdDo, option})
	return err
}

func (c *Conn) dont(option byte) error {
	//log.Println("dont:", option)
	_, err := c.Conn.Write([]byte{cmdIAC, cmdDont, option})
	return err
}

func (c *Conn) will(option byte) error {
	//log.Println("will:", option)
	_, err := c.Conn.Write([]byte{cmdIAC, cmdWill, option})
	return err
}

func (c *Conn) wont(option byte) error {
	//log.Println("wont:", option)
	_, err := c.Conn.Write([]byte{cmdIAC, cmdWont, option})
	return err
}

func (c *Conn) SetWindowSize(rows, columns byte) error {
	c.rows = rows
	c.columns = columns
	return c.will(optWndSize)
}

func (c *Conn) cmd(cmd byte) error {
	switch cmd {
	case cmdGA:
		return nil
	case cmdDo, cmdDont, cmdWill, cmdWont:
	case cmdSB:

		var data []byte
		o, err := c.r.ReadByte()
		if err != nil {
			return err
		}
		for {
			char, err := c.r.ReadByte()
			if err != nil {
				return errors.New("read IAC SE of IAC SB '" + strconv.FormatInt(int64(char), 10) + "' fail, " + err.Error())
			}
			if char != cmdIAC {
				data = append(data, char)
				continue
			}

			char, err = c.r.ReadByte()
			if err != nil {
				return errors.New("read IAC SE of IAC SB '" + strconv.FormatInt(int64(char), 10) + "' fail, " + err.Error())
			}

			if char == cmdSE {
				break
			}

			data = append(data, cmdIAC)
			data = append(data, char)
		}

		switch o {
		case optWndType:
			if len(data) == 1 && data[0] == 1 {
				// IAC SB TERMINAL-TYPE IS xterm IAC SE
				_, err = c.Conn.Write([]byte{cmdIAC, cmdSB, optWndType, 0, 'x', 't', 'e', 'r', 'm', cmdIAC, cmdSE})
				return err
			}
		}
		return nil
	default:
		return nil //fmt.Errorf("unknwn command: %d", cmd)
	}
	// Read an option
	o, err := c.r.ReadByte()
	if err != nil {
		return err
	}
	//log.Println("received cmd:", cmd, o)
	switch o {
	case optEcho:
		// Accept any echo configuration.
		switch cmd {
		case cmdDo:
			if !c.cliEcho {
				c.cliEcho = true
				err = c.will(o)
			}
		case cmdDont:
			if c.cliEcho {
				c.cliEcho = false
				err = c.wont(o)
			}
		case cmdWill:
			err = c.do(o)
		case cmdWont:
			err = c.dont(o)
		}
	case optSuppressGoAhead:
		// We don't use GA so can allways accept every configuration
		switch cmd {
		case cmdDo:
			if !c.cliSuppressGoAhead {
				c.cliSuppressGoAhead = true
				err = c.will(o)
			}
		case cmdDont:
			if c.cliSuppressGoAhead {
				c.cliSuppressGoAhead = false
				err = c.wont(o)
			}
		case cmdWill:
			err = c.do(o)
		case cmdWont:
			err = c.dont(o)

		}
	case optWndSize:
		if cmd == cmdDo {
			_, err = c.Conn.Write([]byte{cmdIAC, cmdSB, optWndSize, 0, c.columns, 0, c.rows, cmdIAC, cmdSE})
		}
	case optWndType:
		// Accept any echo configuration.
		switch cmd {
		case cmdDo:
			err = c.will(o)
		case cmdDont:
		case cmdWill, cmdWont:
			err = c.dont(o)
		}
	default:
		// Deny any other option
		switch cmd {
		case cmdDo:
			err = c.wont(o)
		case cmdDont:
		// nop
		case cmdWill, cmdWont:
			err = c.dont(o)
		}
	}
	return err
}

func (c *Conn) tryReadByte() (b byte, retry bool, err error) {
	b, err = c.r.ReadByte()
	if err != nil || b != cmdIAC {
		return
	}
	b, err = c.r.ReadByte()
	if err != nil {
		return
	}
	if b != cmdIAC {
		err = c.cmd(b)
		if err != nil {
			return
		}
		retry = true
	}
	return
}

// SetEcho tries to enable/disable echo on server side. Typically telnet
// servers doesn't support this.
func (c *Conn) SetEcho(echo bool) error {
	if echo {
		return c.do(optEcho)
	}
	return c.dont(optEcho)
}

// ReadByte works like bufio.ReadByte
func (c *Conn) ReadByte() (b byte, err error) {
	retry := true
	for retry && err == nil {
		b, retry, err = c.tryReadByte()
	}
	return
}

// ReadRune works like bufio.ReadRune
func (c *Conn) ReadRune() (r rune, size int, err error) {
loop:
	r, size, err = c.r.ReadRune()
	if err != nil {
		return
	}
	if r != unicode.ReplacementChar || size != 1 {
		// Properly readed rune
		return
	}
	// Bad rune
	err = c.r.UnreadRune()
	if err != nil {
		return
	}
	// Read telnet command or escaped IAC
	_, retry, err := c.tryReadByte()
	if err != nil {
		return
	}
	if retry {
		// This bad rune was a begining of telnet command. Try read next rune.
		goto loop
	}
	// Return escaped IAC as unicode.ReplacementChar
	return
}

// Read is for implement an io.Reader interface
func (c *Conn) Read(buf []byte) (int, error) {
	var n int
	for n < len(buf) {
		b, err := c.ReadByte()
		if err != nil {
			return n, err
		}
		//log.Printf("char: %d %q", b, b)
		buf[n] = b
		n++
		if c.r.Buffered() == 0 {
			// Try don't block if can return some data
			break
		}
	}
	return n, nil
}

// ReadBytes works like bufio.ReadBytes
func (c *Conn) ReadBytes(delim byte) ([]byte, error) {
	var line []byte
	for {
		b, err := c.ReadByte()
		if err != nil {
			return nil, err
		}
		line = append(line, b)
		if b == delim {
			break
		}
	}
	return line, nil
}

// SkipBytes works like ReadBytes but skips all read data.
func (c *Conn) SkipBytes(delim byte) error {
	for {
		b, err := c.ReadByte()
		if err != nil {
			return err
		}
		if b == delim {
			break
		}
	}
	return nil
}

// ReadString works like bufio.ReadString
func (c *Conn) ReadString(delim byte) (string, error) {
	bytes, err := c.ReadBytes(delim)
	return string(bytes), err
}

func (c *Conn) readUntil(buf *bytes.Buffer, delims [][]byte) (int, error) {
	if len(delims) == 0 {
		return 0, nil
	}
	p := make([][]byte, len(delims))
	for i, s := range delims {
		if len(s) == 0 {
			return 0, nil
		}
		p[i] = s
	}

	for {
		b, err := c.ReadByte()
		if err != nil {
			return 0, err
		}

		if nil != buf {
			buf.WriteByte(b)
		}
		for i, s := range p {
			if s[0] == b {
				if len(s) == 1 {
					return i, nil
				}
				p[i] = s[1:]
			} else {
				p[i] = delims[i]
			}
		}
	}
	panic(nil)
}

// ReadUntilIndex reads from connection until one of delimiters occurs. Returns
// read data and an index of delimiter or error.
// func (c *Conn) ReadUntilIndex(delims ...[]byte) ([]byte, int, error) {
// 	return c.readUntil(true, delims...)
// }

// // ReadUntil works like ReadUntilIndex but don't return a delimiter index.
// func (c *Conn) ReadUntil(delims ...string) ([]byte, error) {
// 	d, _, err := c.readUntil(true, delims...)
// 	return d, err
// }

// SkipUntilIndex works like ReadUntilIndex but skips all read data.
// func (c *Conn) SkipUntilIndex(delims ...[]byte) (int, error) {
// 	_, i, err := c.readUntil(false, delims...)
// 	return i, err
// }

// SkipUntil works like ReadUntil but skips all read data.
// func (c *Conn) SkipUntil(delims ...[]byte) error {
// 	_, _, err := c.readUntil(false, delims...)
// 	return err
// }

var (
	noop_request  = []byte("echo abcdefghigklmn")
	noop_response = []byte("abcdefghigklmn")
)

func (c *Conn) Noop() error {
	if e := c.Sendln(nil, 1*time.Second, noop_request); nil != e {
		return e
	}
	// if _, e := c.Expect(nil, 1*time.Second, noop_response); nil != e {
	// 	return e
	// }
	if _, e := c.Expect(nil, 1*time.Second, [][]byte{[]byte("$"), []byte("#"), []byte(">")}); nil != e {
		return e
	}
	return nil
}

func (c *Conn) Expect(buf *bytes.Buffer, timeout time.Duration, delims [][]byte) (int, error) {
	if e := c.SetReadDeadline(time.Now().Add(timeout)); nil != e {
		return 0, e
	}
	return c.readUntil(buf, delims)
}

// func (c *Conn) ReadUntil(buf *bytes.Buffer, timeout time.Duration, d ...[]byte) (int, error) {
// 	if e := c.SetReadDeadline(time.Now().Add(timeout)); nil != e {
// 		return nil, 0, e
// 	}
// 	return c.readUntil(true, d...)
// }

func (c *Conn) Sendln(buf *bytes.Buffer, timeout time.Duration, s []byte) error {
	if e := c.SetWriteDeadline(time.Now().Add(timeout)); nil != e {
		return e
	}

	copy_buffer := s
	if !bytes.HasSuffix(s, []byte("\n")) {
		copy_buffer = make([]byte, len(s)+1)
		copy(copy_buffer, s)
		copy_buffer[len(s)] = '\n'
	}

	if nil != buf {
		buf.Write(copy_buffer)
	}
	_, err := c.Write(copy_buffer)
	return err
}

func (c *Conn) Send(buf *bytes.Buffer, timeout time.Duration, s []byte) error {
	if e := c.SetWriteDeadline(time.Now().Add(timeout)); nil != e {
		return e
	}

	if nil != buf {
		buf.Write(s)
	}
	_, err := c.Write(s)
	return err
}

// func (c *Conn) Signal(sig int) error {
// 	return nil
// }

// Write is for implement an io.Writer interface
func (c *Conn) Write(buf []byte) (int, error) {
	search := "\xff"
	if c.unixWriteMode {
		search = "\xff\n"
	}
	var (
		n   int
		err error
	)
	for len(buf) > 0 {
		var k int
		i := bytes.IndexAny(buf, search)
		if i == -1 {
			k, err = c.Conn.Write(buf)
			n += k
			break
		}
		k, err = c.Conn.Write(buf[:i])
		n += k
		if err != nil {
			break
		}
		switch buf[i] {
		case LF:
			k, err = c.Conn.Write([]byte{CR, LF})
		case cmdIAC:
			k, err = c.Conn.Write([]byte{cmdIAC, cmdIAC})
		}
		n += k
		if err != nil {
			break
		}
		buf = buf[i+1:]
	}
	return n, err
}
