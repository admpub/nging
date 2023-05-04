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

type IDGenerator func(echo.Context, *Options) (string, error)

var DefaultOptions = &Options{
	EnableImage:    true,
	EnableAudio:    true,
	EnableDownload: true,
	AudioLangs:     []string{`zh`, `ru`, `en`},
	Prefix:         `/captcha`,
	CookieName:     `captchaId`,
	HeaderName:     `X-Captcha-Id`,
	IDGenerator: func(_ echo.Context, _ *Options) (string, error) {
		return captcha.New(), nil
	},
}

func New(prefix string) *Options {
	return &Options{EnableImage: true, Prefix: prefix}
}

type Options struct {
	EnableImage    bool
	EnableAudio    bool
	EnableDownload bool
	AudioLangs     []string
	Prefix         string
	CookieName     string
	HeaderName     string
	IDGenerator    IDGenerator
}

func (o *Options) SetEnableImage(enable bool) *Options {
	o.EnableImage = enable
	return o
}

func (o *Options) SetEnableAudio(enable bool) *Options {
	o.EnableAudio = enable
	return o
}

func (o *Options) SetEnableDownload(enable bool) *Options {
	o.EnableDownload = enable
	return o
}

func (o *Options) SetAudioLangs(audioLangs ...string) *Options {
	o.AudioLangs = audioLangs
	return o
}

func (o *Options) SetPrefix(prefix string) *Options {
	o.Prefix = prefix
	return o
}

func (o *Options) SetIDGenerator(h IDGenerator) *Options {
	o.IDGenerator = h
	return o
}

func (o *Options) SetCookieName(name string) *Options {
	o.CookieName = name
	return o
}

func (o *Options) SetHeaderName(name string) *Options {
	o.HeaderName = name
	return o
}

func (o Options) Wrapper(e echo.RouteRegister) {
	if o.AudioLangs == nil || len(o.AudioLangs) == 0 {
		o.AudioLangs = DefaultOptions.AudioLangs
	}
	if len(o.Prefix) == 0 {
		o.Prefix = DefaultOptions.Prefix
	}
	o.Prefix = strings.TrimRight(o.Prefix, "/")
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
	if o.IDGenerator == nil {
		o.IDGenerator = DefaultOptions.IDGenerator
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
		w := ctx.Response()
		header := w.Header()
		ids := []string{id}
		var hasCookieValue, hasHeaderValue bool
		if len(o.CookieName) > 0 {
			idByCookie := ctx.GetCookie(o.CookieName)
			if len(idByCookie) > 0 {
				ids = append(ids, idByCookie)
				hasCookieValue = true
			}
		}
		if len(o.HeaderName) > 0 {
			idByHeader := ctx.Header(o.HeaderName)
			if len(idByHeader) > 0 {
				ids = append(ids, idByHeader)
				hasHeaderValue = true
			}
		}
		if ctx.Queryx("reload").Bool() {
			var ok bool
			for _, id := range ids {
				if len(id) == 0 {
					continue
				}
				if captcha.Reload(id) {
					ok = true
					ids = []string{id}
					break
				}
			}
			if !ok && (hasCookieValue || hasHeaderValue) { // 旧的已经全部失效了，自动申请新ID
				id, err = o.IDGenerator(ctx, o)
				if err != nil {
					header.Add(`X-Captcha-ID-Error`, `generator: `+err.Error())
					return err
				}
				ids = []string{id}
				if hasCookieValue {
					ctx.SetCookie(o.CookieName, id)
				}
				if hasHeaderValue {
					header.Set(o.HeaderName, id)
				}
			}
		}
		download := o.EnableDownload && ctx.Queryx("download").Bool()
		b := bytes.NewBuffer(nil)
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
