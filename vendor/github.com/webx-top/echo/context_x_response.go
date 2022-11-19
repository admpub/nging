package echo

import (
	"bytes"
	"encoding/xml"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
	"unicode"

	"github.com/webx-top/echo/encoding/json"
	"github.com/webx-top/echo/engine"
	"github.com/webx-top/poolx/bufferpool"
)

// Response returns *Response.
func (c *xContext) Response() engine.Response {
	return c.response
}

// Render renders a template with data and sends a text/html response with status
// code. Templates can be registered using `Echo.SetRenderer()`.
func (c *xContext) Render(name string, data interface{}, codes ...int) (err error) {
	if c.auto {
		format := c.Format()
		if render, ok := c.echo.formatRenderers[format]; ok && render != nil {
			switch v := data.(type) {
			case Data: //Skip
			case error:
				c.dataEngine.SetError(v)
			case nil:
				if c.dataEngine.GetData() == nil {
					c.dataEngine.SetData(c.Stored(), c.dataEngine.GetCode().Int())
				}
			default:
				c.dataEngine.SetData(data, c.dataEngine.GetCode().Int())
			}
			return render(c, data)
		}
	}
	c.dataEngine.SetTmplFuncs()
	if data == nil {
		data = c.dataEngine.GetData()
	}
	b, err := c.Fetch(name, data)
	if err != nil {
		return
	}
	b = bytes.TrimLeftFunc(b, unicode.IsSpace)
	c.response.Header().Set(HeaderContentType, MIMETextHTMLCharsetUTF8)
	err = c.Blob(b, codes...)
	return
}

func (c *xContext) RenderBy(name string, content func(string) ([]byte, error), data interface{}, codes ...int) (b []byte, err error) {
	c.dataEngine.SetTmplFuncs()
	if data == nil {
		data = c.dataEngine.GetData()
	}
	if c.renderer == nil {
		if c.echo.renderer == nil {
			return nil, ErrRendererNotRegistered
		}
		c.renderer = c.echo.renderer
	}
	buf := bufferpool.Get()
	defer bufferpool.Release(buf)
	if c.renderDataWrapper != nil {
		data = c.renderDataWrapper(c, data)
	}
	err = c.renderer.RenderBy(buf, name, content, data, c)
	if err != nil {
		return
	}
	b = buf.Bytes()
	return
}

// HTML sends an HTTP response with status code.
func (c *xContext) HTML(html string, codes ...int) (err error) {
	c.response.Header().Set(HeaderContentType, MIMETextHTMLCharsetUTF8)
	err = c.Blob([]byte(html), codes...)
	return
}

// String sends a string response with status code.
func (c *xContext) String(s string, codes ...int) (err error) {
	c.response.Header().Set(HeaderContentType, MIMETextPlainCharsetUTF8)
	err = c.Blob([]byte(s), codes...)
	return
}

func (c *xContext) Blob(b []byte, codes ...int) (err error) {
	if len(codes) > 0 {
		c.code = codes[0]
	}
	if c.code == 0 {
		c.code = http.StatusOK
	}
	err = c.preResponse()
	if err != nil {
		return
	}
	c.response.WriteHeader(c.code)
	_, err = c.response.Write(b)
	return
}

// JSON sends a JSON response with status code.
func (c *xContext) JSON(i interface{}, codes ...int) (err error) {
	var b []byte
	if c.echo.Debug() {
		b, err = json.MarshalIndent(i, "", "  ")
	} else {
		b, err = json.Marshal(i)
	}
	if err != nil {
		return err
	}
	return c.JSONBlob(b, codes...)
}

// JSONBlob sends a JSON blob response with status code.
func (c *xContext) JSONBlob(b []byte, codes ...int) (err error) {
	c.response.Header().Set(HeaderContentType, MIMEApplicationJSONCharsetUTF8)
	err = c.Blob(b, codes...)
	return
}

// JSONP sends a JSONP response with status code. It uses `callback` to construct
// the JSONP payload.
func (c *xContext) JSONP(callback string, i interface{}, codes ...int) (err error) {
	b, err := json.Marshal(i)
	if err != nil {
		return err
	}
	c.response.Header().Set(HeaderContentType, MIMEApplicationJavaScriptCharsetUTF8)
	b = []byte(callback + "(" + string(b) + ");")
	err = c.Blob(b, codes...)
	return
}

// XML sends an XML response with status code.
func (c *xContext) XML(i interface{}, codes ...int) (err error) {
	var b []byte
	if c.echo.Debug() {
		b, err = xml.MarshalIndent(i, "", "  ")
	} else {
		b, err = xml.Marshal(i)
	}
	if err != nil {
		return err
	}
	return c.XMLBlob(b, codes...)
}

// XMLBlob sends a XML blob response with status code.
func (c *xContext) XMLBlob(b []byte, codes ...int) (err error) {
	c.response.Header().Set(HeaderContentType, MIMEApplicationXMLCharsetUTF8)
	b = []byte(xml.Header + string(b))
	err = c.Blob(b, codes...)
	return
}

func (c *xContext) Stream(step func(w io.Writer) bool) error {
	return c.response.Stream(step)
}

func (c *xContext) SSEvent(event string, data chan interface{}) (err error) {
	hdr := c.response.Header()
	hdr.Set(HeaderContentType, MIMEEventStream)
	hdr.Set(`Cache-Control`, `no-cache`)
	hdr.Set(`Connection`, `keep-alive`)
	hdr.Set(`Transfer-Encoding`, `chunked`)
	err = c.Stream(func(w io.Writer) bool {
		recv, ok := <-data
		if !ok {
			return ok
		}
		b, _err := c.Fetch(event, recv)
		if _err != nil {
			err = _err
			return false
		}
		//c.Logger().Debugf(`SSEvent: %s`, b)
		_, _err = w.Write(b)
		if _err != nil {
			err = _err
			return false
		}
		return true
	})
	return
}

func (c *xContext) Attachment(r io.Reader, name string, modtime time.Time, inline ...bool) (err error) {
	var typ string
	if len(inline) > 0 && inline[0] {
		typ = `inline`
	} else {
		typ = `attachment`
	}
	c.response.Header().Set(HeaderContentType, ContentTypeByExtension(name))
	encodedName := URLEncode(name, true)
	c.response.Header().Set(HeaderContentDisposition, typ+"; filename="+encodedName+"; filename*=utf-8''"+encodedName)
	return c.ServeContent(r, name, modtime)
}

func (c *xContext) File(file string, fs ...http.FileSystem) (err error) {
	var f http.File
	customFS := len(fs) > 0 && fs[0] != nil
	if customFS {
		f, err = fs[0].Open(file)
	} else {
		f, err = os.Open(file)
	}
	if err != nil {
		return ErrNotFound
	}
	defer func() {
		if f != nil {
			f.Close()
		}
	}()
	fi, err := f.Stat()
	if err != nil {
		return err
	}
	if fi.IsDir() {
		f.Close()
		file = filepath.Join(file, "index.html")
		if customFS {
			f, err = fs[0].Open(file)
		} else {
			f, err = os.Open(file)
		}
		if err != nil {
			return ErrNotFound
		}
		fi, err = f.Stat()
		if err != nil {
			return err
		}
	}
	c.Response().ServeContent(f, fi.Name(), fi.ModTime())
	return nil
}

func (c *xContext) ServeContent(content io.Reader, name string, modtime time.Time) error {
	if readSeeker, ok := content.(io.ReadSeeker); ok {
		c.Response().ServeContent(readSeeker, name, modtime)
		return nil
	}
	return c.ServeCallbackContent(func(_ Context) (io.Reader, error) {
		return content, nil
	}, name, modtime)
}

func (c *xContext) ServeCallbackContent(callback func(Context) (io.Reader, error), name string, modtime time.Time) error {
	rq := c.Request()
	rs := c.Response()

	if t, err := time.Parse(http.TimeFormat, rq.Header().Get(HeaderIfModifiedSince)); err == nil && modtime.Before(t.Add(1*time.Second)) {
		rs.Header().Del(HeaderContentType)
		rs.Header().Del(HeaderContentLength)
		return c.NoContent(http.StatusNotModified)
	}
	content, err := callback(c)
	if err != nil {
		return err
	}
	if readSeeker, ok := content.(io.ReadSeeker); ok {
		c.Response().ServeContent(readSeeker, name, modtime)
		return nil
	}
	rs.Header().Set(HeaderContentType, ContentTypeByExtension(name))
	rs.Header().Set(HeaderLastModified, modtime.UTC().Format(http.TimeFormat))
	rs.WriteHeader(http.StatusOK)
	rs.KeepBody(false)
	_, err = io.Copy(rs, content)
	return err
}

// NoContent sends a response with no body and a status code.
func (c *xContext) NoContent(codes ...int) error {
	if len(codes) > 0 {
		c.code = codes[0]
	}
	if c.code == 0 {
		c.code = http.StatusOK
	}
	c.response.WriteHeader(c.code)
	return nil
}

// Redirect redirects the request with status code.
func (c *xContext) Redirect(url string, codes ...int) error {
	code := http.StatusFound
	if len(codes) > 0 {
		code = codes[0]
	}
	if code < http.StatusMultipleChoices || code > http.StatusTemporaryRedirect {
		return ErrInvalidRedirectCode
	}
	err := c.preResponse()
	if err != nil {
		return err
	}
	format := c.Format()
	if format != `html` && c.auto {
		if render, ok := c.echo.formatRenderers[format]; ok && render != nil {
			if c.dataEngine.GetData() == nil {
				c.dataEngine.SetData(c.Stored(), c.dataEngine.GetCode().Int())
			}
			c.dataEngine.SetURL(url)
			return render(c, c.dataEngine.GetData())
		}
	}
	c.response.Redirect(url, code)
	return nil
}
