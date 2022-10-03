/*

   Copyright 2016-present Wenhui Shen <www.webx.top>

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

package common

import (
	"bytes"
	"errors"
	"html/template"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"

	"github.com/coscms/forms/config"
	"github.com/webx-top/tagfast"
	"golang.org/x/sync/singleflight"
)

// Available form themes
const (
	BASE      = "base"
	BOOTSTRAP = "bootstrap3"
)

var (
	tmplDirs = map[string]string{
		BASE:      "templates",
		BOOTSTRAP: "templates",
	}
	LabelFn = func(s string) string {
		return s
	}

	//private
	cachedTemplate = make(map[string]*template.Template)
	cachedConfig   = make(map[string]*config.Config)
	lockTemplate   = new(sync.RWMutex)
	lockConfig     = new(sync.RWMutex)
	lockTmplDir    = new(sync.RWMutex)
	sg             singleflight.Group
)

const (
	PACKAGE_NAME = "github.com/coscms/forms"
)

// Input field types
const (
	BUTTON         = "button"
	CHECKBOX       = "checkbox"
	COLOR          = "color"
	DATE           = "date"
	DATETIME       = "datetime"
	DATETIME_LOCAL = "datetime-local"
	EMAIL          = "email"
	FILE           = "file"
	HIDDEN         = "hidden"
	IMAGE          = "image"
	MONTH          = "month"
	NUMBER         = "number"
	PASSWORD       = "password"
	RADIO          = "radio"
	RANGE          = "range"
	RESET          = "reset"
	SEARCH         = "search"
	SUBMIT         = "submit"
	TEL            = "tel"
	TEXT           = "text"
	TIME           = "time"
	URL            = "url"
	WEEK           = "week"
	TEXTAREA       = "textarea"
	SELECT         = "select"
	STATIC         = "static"
)

func SetTmplDir(theme, tmplDir string) {
	lockTmplDir.Lock()
	tmplDirs[theme] = tmplDir
	lockTmplDir.Unlock()
}

func TmplDir(theme string) (tmplDir string) {
	tmplDir, _ = tmplDirs[theme]
	return
}

// LookupPath creates the complete path of the desired widget template
func LookupPath(widget string) string {
	if !FileSystem.IsEmpty() {
		fp, err := FileSystem.Open(widget)
		if err != nil {
			if !errors.Is(err, fs.ErrNotExist) {
				log.Println(err.Error())
				return widget
			}
		} else {
			defer fp.Close()
			fi, err := fp.Stat()
			if err == nil && !fi.IsDir() {
				return widget
			}
		}
	}
	if !TmplExists(widget) {
		return filepath.Join(os.Getenv("GOPATH"), "src", PACKAGE_NAME, `defaults`, widget)
	}
	return widget
}

func TmplExists(tmpl string) bool {
	_, err := os.Stat(tmpl)
	return !os.IsNotExist(err)
}

func GetOrSetCachedTemplate(cachedKey string, generator func() (*template.Template, error)) (c *template.Template, err error) {
	var ok bool
	lockTemplate.RLock()
	c, ok = cachedTemplate[cachedKey]
	lockTemplate.RUnlock()
	if ok {
		return c, nil
	}
	getValue, getErr, _ := sg.Do(cachedKey, func() (interface{}, error) {
		c, err = generator()
		if err != nil {
			return nil, err
		}
		lockTemplate.Lock()
		cachedTemplate[cachedKey] = c
		lockTemplate.Unlock()
		return c, nil
	})
	if getErr != nil {
		return nil, getErr
	}
	return getValue.(*template.Template), nil
}

func ClearCachedTemplate() {
	lockTemplate.Lock()
	cachedTemplate = make(map[string]*template.Template)
	lockTemplate.Unlock()
}

func DelCachedTemplate(key string) bool {
	lockTemplate.Lock()
	defer lockTemplate.Unlock()
	if _, ok := cachedTemplate[key]; ok {
		delete(cachedTemplate, key)
		return true
	}
	return false
}

func GetOrSetCachedConfig(cachedKey string, generator func() (*config.Config, error)) (c *config.Config, err error) {
	var ok bool
	lockConfig.RLock()
	c, ok = cachedConfig[cachedKey]
	lockConfig.RUnlock()
	if ok {
		return c, nil
	}
	getValue, getErr, _ := sg.Do(cachedKey, func() (interface{}, error) {
		c, err = generator()
		if err != nil {
			return nil, err
		}
		lockConfig.Lock()
		cachedConfig[cachedKey] = c
		lockConfig.Unlock()
		return c, nil
	})
	if getErr != nil {
		return nil, getErr
	}
	return getValue.(*config.Config), nil
}

func ClearCachedConfig() {
	lockConfig.Lock()
	cachedConfig = make(map[string]*config.Config)
	lockConfig.Unlock()
}

func DelCachedConfig(key string) bool {
	lockConfig.Lock()
	defer lockConfig.Unlock()
	if _, ok := cachedConfig[key]; ok {
		delete(cachedConfig, key)
		return true
	}
	return false
}

func ParseTmpl(data interface{},
	fn_tpl template.FuncMap,
	fn_fixTpl func(tpls ...string) ([]string, error),
	tpls ...string) string {
	buf := bytes.NewBuffer(nil)
	tpf := strings.Join(tpls, `|`)
	tpl, err := GetOrSetCachedTemplate(tpf, func() (*template.Template, error) {
		c := template.New(filepath.Base(tpls[0]))
		if fn_tpl != nil {
			c.Funcs(fn_tpl)
		}
		var err error
		if fn_fixTpl != nil {
			tpls, err = fn_fixTpl(tpls...)
			if err != nil {
				return nil, err
			}
		}
		if !FileSystem.IsEmpty() {
			return c.ParseFS(FileSystem, tpls...)

		}
		return c.ParseFiles(tpls...)
	})
	if err != nil {
		return err.Error()
	}
	err = tpl.Execute(buf, data)
	if err != nil {
		return err.Error()
	}
	return buf.String()
}

func TagVal(t reflect.Type, fieldNo int, tagName string) string {
	return tagfast.Value(t, t.Field(fieldNo), tagName)
}

func Tag(t reflect.Type, f reflect.StructField, tagName string) (value string, tf tagfast.Faster) {
	return tagfast.Tag(t, f, tagName)
}
