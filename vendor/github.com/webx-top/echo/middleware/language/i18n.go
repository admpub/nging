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
	"fmt"
	"path/filepath"
	"strings"

	"github.com/admpub/i18n"
	"github.com/admpub/log"
	"github.com/webx-top/com"
)

var defaultInstance *I18n

type I18n struct {
	*i18n.TranslatorFactory
	Translators map[string]*i18n.Translator
	config      *Config
}

func NewI18n(c *Config) *I18n {
	f, errs := i18n.NewTranslatorFactory(c.RulesPath, c.MessagesPath, c.Fallback)
	if errs != nil && len(errs) > 0 {
		var errMsg string
		for idx, err := range errs {
			if idx > 0 {
				errMsg += "\n"
			}
			errMsg += err.Error()
		}
		if len(errMsg) > 0 {
			panic("== i18n error: " + errMsg + "\n")
		}
	}
	defaultInstance = &I18n{
		TranslatorFactory: f,
		Translators:       make(map[string]*i18n.Translator),
		config:            c,
	}
	defaultInstance.Get(c.Default)

	return defaultInstance
}

func (a *I18n) Monitor() *I18n {
	onchange := func(file string) {
		log.Info("reload language: ", file)
		defaultInstance.Reload(file)
	}
	for _, mp := range a.config.MessagesPath {
		if len(mp) == 0 {
			continue
		}
		callback := &com.MonitorEvent{
			Modify: onchange,
			Delete: onchange,
			Rename: onchange,
		}
		callback.Watch(mp, func(f string) bool {
			log.Info("changed language: ", f)
			return strings.HasSuffix(f, `.yaml`)
		})
	}
	return a
}

func (a *I18n) Get(langCode string) *i18n.Translator {
	var (
		t    *i18n.Translator
		errs []error
	)
	t, errs = a.TranslatorFactory.GetTranslator(langCode)
	if errs != nil && len(errs) > 0 {
		if a.config.Default != langCode {
			t, errs = a.TranslatorFactory.GetTranslator(a.config.Default)
		}
	}
	if errs != nil && len(errs) > 0 {
		var errMsg string
		for idx, err := range errs {
			if idx > 0 {
				errMsg += "\n"
			}
			errMsg += err.Error()
		}
		if len(errMsg) > 0 {
			panic("== i18n error: " + errMsg + "\n")
		}
	}
	a.Translators[langCode] = t
	return t
}

func (a *I18n) Reload(langCode string) {
	if strings.HasSuffix(langCode, `.yaml`) {
		langCode = strings.TrimSuffix(langCode, `.yaml`)
		langCode = filepath.Base(langCode)
	}
	a.TranslatorFactory.Reload(langCode)
	if _, ok := a.Translators[langCode]; ok {
		delete(a.Translators, langCode)
	}
}

func (a *I18n) Translate(langCode, key string, args map[string]string) string {
	t, ok := a.Translators[langCode]
	if !ok {
		t = a.Get(langCode)
	}
	translation, err := t.Translate(key, args)
	if err != nil {
		return key
	}
	return translation
}

func (a *I18n) T(langCode, key string, args ...interface{}) (t string) {
	if len(args) > 0 {
		if v, ok := args[0].(map[string]string); ok {
			t = a.Translate(langCode, key, v)
			return
		}
		t = a.Translate(langCode, key, map[string]string{})
		t = fmt.Sprintf(t, args...)
		return
	}
	t = a.Translate(langCode, key, map[string]string{})
	return
}

//T 多语言翻译
func T(langCode, key string, args ...interface{}) (t string) {
	if defaultInstance == nil {
		t = key
		if len(args) > 0 {
			t = fmt.Sprintf(t, args...)
		}
		return
	}
	t = defaultInstance.T(langCode, key, args...)
	return
}
