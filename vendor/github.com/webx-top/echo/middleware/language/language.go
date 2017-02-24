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
package language

import (
	"regexp"
	"strings"

	"github.com/webx-top/echo"
	"github.com/webx-top/echo/engine"
)

var (
	LangVarName        = `lang`
	DefaultLang        = `zh-cn`
	headerAcceptRemove = regexp.MustCompile(`;q=[0-9.]+`)
)

func New(c ...*Config) *Language {
	lang := &Language{
		List:    make(map[string]bool),
		Index:   make([]string, 0),
		Default: DefaultLang,
	}
	if len(c) > 0 {
		lang.Init(c[0])
	}
	return lang
}

type Language struct {
	List    map[string]bool //语种列表
	Index   []string        //索引
	Default string          //默认语种
	I18n    *I18n
}

func (a *Language) Init(c *Config) {
	if c.AllList != nil {
		for _, lang := range c.AllList {
			a.Set(lang, true, lang == c.Default)
		}
	} else {
		a.Set(c.Default, true, true)
		if c.Default != `en` {
			a.Set(`en`, true)
		}
	}
	a.I18n = NewI18n(c)
	if c.Reload {
		a.I18n.Monitor()
	}
}

func (a *Language) Set(lang string, on bool, args ...bool) *Language {
	if a.List == nil {
		a.List = make(map[string]bool)
	}
	if _, ok := a.List[lang]; !ok {
		a.Index = append(a.Index, lang)
	}
	a.List[lang] = on
	if on && len(args) > 0 && args[0] {
		a.Default = lang
	}
	return a
}

func (a *Language) DetectURI(r engine.Request) string {
	p := strings.TrimPrefix(r.URL().Path(), `/`)
	s := strings.Index(p, `/`)
	var lang string
	if s != -1 {
		lang = p[0:s]
	} else {
		lang = p
	}
	if len(lang) > 0 {
		if on, ok := a.List[lang]; ok {
			r.URL().SetPath(strings.TrimPrefix(p, lang))
			if !on {
				lang = ""
			}
		} else {
			lang = ""
		}
	}
	return lang
}

func (a *Language) Valid(lang string) bool {
	if len(lang) > 0 {
		if on, ok := a.List[lang]; ok {
			return on
		}
	}
	return false
}

func (a *Language) DetectHeader(r engine.Request) string {
	al := r.Header().Get(`Accept-Language`)
	al = headerAcceptRemove.ReplaceAllString(al, ``)
	lg := strings.SplitN(al, `,`, 5)
	for _, lang := range lg {
		lang = strings.ToLower(lang)
		if a.Valid(lang) {
			return lang
		}
	}
	return a.Default
}

func (a *Language) Middleware() echo.MiddlewareFunc {
	return echo.MiddlewareFunc(func(h echo.Handler) echo.Handler {
		return echo.HandlerFunc(func(c echo.Context) error {
			lang := c.Query(LangVarName)
			var hasCookie bool
			if !a.Valid(lang) {
				lang = a.DetectURI(c.Request())
				if !a.Valid(lang) {
					lang = c.GetCookie(LangVarName)
					if !a.Valid(lang) {
						lang = a.DetectHeader(c.Request())
					} else {
						hasCookie = true
					}
				}
			}
			if !hasCookie {
				c.SetCookie(LangVarName, lang)
			}
			c.SetTranslator(NewTranslate(lang, a.I18n))
			return h.Handle(c)
		})
	})
}
