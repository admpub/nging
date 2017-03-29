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
}

type Options struct {
	EnableImage    bool
	EnableAudio    bool
	EnableDownload bool
	AudioLangs     []string
	Prefix         string
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
		if len(ctx.Query("reload")) > 0 {
			captcha.Reload(id)
		}
		w := ctx.Response()
		header := w.Header()
		download := o.EnableDownload && len(ctx.Query("download")) > 0
		b := bytes.NewBufferString(``)
		switch ext {
		case ".png":
			if !o.EnableImage {
				return echo.ErrNotFound
			}
			if download {
				header.Set(echo.HeaderContentType, "application/octet-stream")
			} else {
				header.Set(echo.HeaderContentType, "image/png")
			}
			err = captcha.WriteImage(b, id, captcha.StdWidth, captcha.StdHeight)
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
			au, err = captcha.GetAudio(id, lang)
			if err != nil {
				return
			}
			if download {
				header.Set(echo.HeaderContentType, "application/octet-stream")
			} else {
				header.Set(echo.HeaderContentType, "audio/x-wav")
			}
			header.Set("Content-Length", strconv.Itoa(au.EncodedLen()))
			_, err = au.WriteTo(b)
		default:
			return nil
		}
		if err != nil {
			return
		}
		return ctx.Blob(b.Bytes())
	}
}
