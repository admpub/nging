/*

   Copyright 2016 Wenhui Shen <www.webx.top>

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.

*/
package captcha

import (
	"bytes"
	"strconv"
	"strings"

	"github.com/webx-top/captcha"
	"github.com/webx-top/echo"
)

var DefaultOptions = &Options{
	EnableImage:    true,
	EnableAudio:    true,
	EnableDownload: true,
	AudioLangs:     []string{`zh`, `ru`, `en`},
	CookieName:     `captchaId`,
	HeaderName:     `X-Captcha-ID`,
}

type Options struct {
	EnableImage    bool
	EnableAudio    bool
	EnableDownload bool
	AudioLangs     []string
	Prefix         string
	CookieName     string
	HeaderName     string
}

func (o Options) Wrapper(e echo.RouteRegister) {
	if o.AudioLangs == nil || len(o.AudioLangs) == 0 {
		o.AudioLangs = []string{`zh`, `ru`, `en`}
	}
	if len(o.Prefix) == 0 {
		o.Prefix = `/captcha`
	} else {
		o.Prefix = strings.TrimRight(o.Prefix, "/")
	}
	e.Get(o.Prefix+"/*", Captcha(&o))
}

func Captcha(opts ...*Options) func(echo.Context) error {
	var o *Options
	if len(opts) > 0 {
		o = opts[0]
	}
	if o == nil {
		o = DefaultOptions
	}
	if len(o.CookieName) == 0 {
		o.CookieName = DefaultOptions.CookieName
	}
	if len(o.HeaderName) == 0 {
		o.HeaderName = DefaultOptions.HeaderName
	}
	return func(ctx echo.Context) (err error) {
		var id, ext string
		param := ctx.P(0)
		if p := strings.LastIndex(param, `.`); p > 0 {
			id = param[0:p]
			ext = param[p:]
		}
		if len(ext) == 0 || len(id) == 0 {
			return echo.ErrNotFound
		}
		ids := []string{id}
		if len(o.CookieName) > 0 {
			ids = append(ids, ctx.GetCookie(o.CookieName))
		}
		if len(o.HeaderName) > 0 {
			ids = append(ids, ctx.Header(o.HeaderName))
		}
		w := ctx.Response()
		header := w.Header()
		if len(ctx.Query("reload")) > 0 {
			var ok bool
			for _, id := range ids {
				if captcha.Reload(id) {
					ok = true
					ids = []string{id}
					break
				}
			}
			if !ok {
				if len(o.CookieName) > 0 {
					id = captcha.New()
					ids = []string{id}
					ctx.SetCookie(o.CookieName, id)
				} else if len(o.HeaderName) > 0 {
					id = captcha.New()
					ids = []string{id}
					header.Set(o.HeaderName, id)
				}
			}
		}
		download := o.EnableDownload && len(ctx.Query("download")) > 0
		b := bytes.NewBufferString(``)
		switch ext {
		case ".png":
			if !o.EnableImage {
				return echo.ErrNotFound
			}
			for _, id := range ids {
				if len(id) == 0 {
					continue
				}
				err = captcha.WriteImage(b, id, captcha.StdWidth, captcha.StdHeight)
				if err == nil || err != captcha.ErrNotFound {
					break
				}
			}
			if err != nil {
				if err == captcha.ErrNotFound {
					return echo.ErrNotFound
				}
				return
			}
			if download {
				header.Set(echo.HeaderContentType, "application/octet-stream")
			} else {
				header.Set(echo.HeaderContentType, "image/png")
			}
		case ".wav":
			if !o.EnableAudio {
				return echo.ErrNotFound
			}
			lang := strings.ToLower(ctx.Query("lang"))
			supported := false
			for _, supportedLang := range o.AudioLangs {
				if supportedLang == lang {
					supported = true
					break
				}
			}
			if !supported && len(o.AudioLangs) > 0 {
				lang = o.AudioLangs[0]
			}
			var au *captcha.Audio
			for _, id := range ids {
				if len(id) == 0 {
					continue
				}
				au, err = captcha.GetAudio(id, lang)
				if err == nil || err != captcha.ErrNotFound {
					break
				}
			}
			if err != nil {
				if err == captcha.ErrNotFound {
					return echo.ErrNotFound
				}
				return
			}
			length := strconv.Itoa(au.EncodedLen())
			_, err = au.WriteTo(b)
			if err != nil {
				return err
			}
			if download {
				header.Set(echo.HeaderContentType, "application/octet-stream")
			} else {
				header.Set(echo.HeaderContentType, "audio/x-wav")
			}
			header.Set("Content-Length", length)
		default:
			return nil
		}
		return ctx.Blob(b.Bytes())
	}
}
