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
package resource

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/admpub/log"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/handler/mvc/static/minify"
	mw "github.com/webx-top/echo/middleware"
)

var (
	regexCssUrlAttr      = regexp.MustCompile(`url\(['"]?([^\)'"]+)['"]?\)`)
	regexCssImport       = regexp.MustCompile(`@import[\s]+["']([^"']+)["'][\s]*;`)
	regexCssCleanSpace   = regexp.MustCompile(`(?s)\s*(\{|\}|;|:)\s*`)
	regexCssCleanSpace2  = regexp.MustCompile(`(?s)\s{2,}`)
	regexCssCleanComment = regexp.MustCompile(`(?s)[\s]*/\*(.*?)\*/[\s]*`)
)

type urlMapInfo struct {
	AbsPath string
	Md5     string
}

func NewStatic(staticPath, rootPath string) *Static {
	if len(staticPath) > 0 && staticPath[0] != '/' {
		staticPath = `/` + staticPath
	}
	return &Static{
		StaticOptions: &mw.StaticOptions{
			Path: staticPath,
			Root: rootPath,
		},
		CombineJs:       true,
		CombineCss:      true,
		CombineSavePath: `combine`,
		Combined:        make(map[string][]string),
		Combines:        make(map[string]bool),
		urlMap:          make(map[string]*urlMapInfo),
		mutex:           &sync.Mutex{},
		logger:          log.GetLogger(`echo`),
	}
}

type Static struct {
	*mw.StaticOptions
	CombineJs       bool
	CombineCss      bool
	CombineSavePath string //合并文件保存路径，首尾均不带斜杠
	Combined        map[string][]string
	Combines        map[string]bool
	urlMap          map[string]*urlMapInfo
	mutex           *sync.Mutex
	Public          *Static
	logger          *log.Logger
	middleware      echo.MiddlewareFunc
}

// Wrapper 包装路由（作为路由时使用）
func (s *Static) Wrapper(r echo.RouteRegister) {
	r.Get(s.Path+`/*`, s)
}

// Handle 处理
func (s *Static) Handle(ctx echo.Context) error {
	file := filepath.Join(s.Root, ctx.P(0))
	if !strings.HasPrefix(file, s.Root) {
		return echo.ErrNotFound
	}
	return ctx.File(file)
}

// Middleware 中间件（作为中间件使用）
func (s *Static) Middleware() echo.MiddlewareFunc {
	if s.middleware == nil {
		s.middleware = mw.Static(s.StaticOptions)
	}
	return s.middleware
}

func (s *Static) StaticURL(staticFile string) (r string) {
	r = s.Path + "/" + staticFile
	return
}

func (s *Static) JsURL(staticFile string) (r string) {
	r = s.StaticURL("js/" + staticFile)
	return
}

func (s *Static) CssURL(staticFile string) (r string) {
	r = s.StaticURL("css/" + staticFile)
	return r
}

func (s *Static) ImgURL(staticFile string) (r string) {
	r = s.StaticURL("img/" + staticFile)
	return r
}

func (s *Static) cachedURLInfo(key string, ext string) (absPath string, fileName string) {
	if v, ok := s.urlMap[key]; ok {
		fileName = v.Md5 + "." + ext
		absPath = v.AbsPath
	} else {
		md5 := com.Md5(key)
		fileName = md5 + "." + ext
		absPath = filepath.Join(s.Root, s.CombineSavePath, fileName)
		s.urlMap[key] = &urlMapInfo{
			AbsPath: absPath,
			Md5:     md5,
		}
	}
	return
}

func (s *Static) JsTag(staticFiles ...string) template.HTML {
	var r string
	if len(staticFiles) == 1 || !s.CombineJs {
		for _, staticFile := range staticFiles {
			r += `<script type="text/javascript" src="` + s.JsURL(staticFile) + `" charset="utf-8"></script>`
		}
		return template.HTML(r)
	}
	r, combinedFile := s.cachedURLInfo(strings.Join(staticFiles, "|"), `js`)
	if s.IsCombined(r) == false || com.FileExists(r) == false {
		var content string
		for _, url := range staticFiles {
			absPath := filepath.Join(s.Root, "js", url)
			if con, err := s.genCombinedJS(absPath, url); err != nil {
				fmt.Println(err)
			} else {
				s.RecordCombined(absPath, r)
				content += con
			}
		}
		com.WriteFile(r, []byte(content))
		s.RecordCombines(r)
	}
	r = `<script type="text/javascript" src="` + s.StaticURL(path.Join(s.CombineSavePath, combinedFile)) + `" charset="utf-8"></script>`
	return template.HTML(r)
}

func (s *Static) CssTag(staticFiles ...string) template.HTML {
	var r string
	if len(staticFiles) == 1 || !s.CombineCss {
		for _, staticFile := range staticFiles {
			r += `<link rel="stylesheet" type="text/css" href="` + s.CssURL(staticFile) + `" charset="utf-8" />`
		}
		return template.HTML(r)
	}

	r, combinedFile := s.cachedURLInfo(strings.Join(staticFiles, "|"), `css`)
	if s.IsCombined(r) == false || com.FileExists(r) == false {
		var onImportFn = func(urlPath string) {
			s.RecordCombined(filepath.Join(s.Root, "css", urlPath), r)
		}
		var content string
		for _, url := range staticFiles {
			absPath := filepath.Join(s.Root, "css", url)
			if con, err := s.genCombinedCSS(absPath, url, onImportFn); err != nil {
				fmt.Println(err)
			} else {
				s.RecordCombined(absPath, r)
				content += con
			}
		}
		com.WriteFile(r, []byte(content))
		s.RecordCombines(r)
	}
	r = `<link rel="stylesheet" type="text/css" href="` + s.StaticURL(path.Join(s.CombineSavePath, combinedFile)) + `" charset="utf-8" />`
	return template.HTML(r)
}

func (s *Static) ImgTag(staticFile string, attrs ...string) template.HTML {
	var attr string
	for i, l := 0, len(attrs)-1; i < l; i++ {
		var k, v string
		k = attrs[i]
		i++
		v = attrs[i]
		attr += ` ` + k + `="` + v + `"`
	}
	r := `<img src="` + s.ImgURL(staticFile) + `"` + attr + ` />`
	return template.HTML(r)
}

func (s *Static) Register(funcMap map[string]interface{}) map[string]interface{} {
	if funcMap == nil {
		funcMap = map[string]interface{}{}
	}
	funcMap["StaticURL"] = s.StaticURL
	funcMap["JsURL"] = s.JsURL
	funcMap["CssURL"] = s.CssURL
	funcMap["ImgURL"] = s.ImgURL
	funcMap["JsTag"] = s.JsTag
	funcMap["CssTag"] = s.CssTag
	funcMap["ImgTag"] = s.ImgTag
	return funcMap
}

func (s *Static) DeleteCombined(url string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.logger.Debug(`update resource `, url)
	if val, ok := s.Combined[url]; ok {
		for _, v := range val {
			if _, has := s.Combines[v]; !has {
				s.logger.Debug(`skip resource `, url)
				continue
			}
			s.logger.Debug(`remove combines `, v)
			err := os.Remove(v)
			delete(s.Combines, v)
			if err != nil {
				s.logger.Error(err)
			}
		}
	}
}

func (s *Static) RecordCombined(fromUrl string, combineUrl string) {
	if s.Combined == nil {
		return
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if _, ok := s.Combined[fromUrl]; !ok {
		s.Combined[fromUrl] = make([]string, 0)
	}
	s.Combined[fromUrl] = append(s.Combined[fromUrl], combineUrl)
}

func (s *Static) RecordCombines(combineUrl string) {
	if s.Combines == nil {
		return
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.Combines[combineUrl] = true
}

func (s *Static) IsCombined(combineUrl string) (ok bool) {
	if s.Combines == nil {
		return
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()
	_, ok = s.Combines[combineUrl]
	return
}

func (s *Static) ClearCache() {
	for f, _ := range s.Combines {
		os.Remove(f)
	}
	s.Combined = make(map[string][]string)
	s.Combines = make(map[string]bool)
	s.urlMap = make(map[string]*urlMapInfo)
}

func (s *Static) OnUpdate(tmplDir string) func(string) {
	return func(name string) {
		if s.Public != nil {
			s.Public.ClearCache()
			if s.Public == s {
				return
			}
		}
		name = filepath.Join(tmplDir, name)
		s.DeleteCombined(name)
	}
}

func (s *Static) genCombinedJS(absPath, urlPath string) (content string, err error) {
	con, err := com.ReadFileS(absPath)
	if err != nil {
		return ``, err
	}
	content += "\n/* <from: " + urlPath + "> */\n"
	if !strings.Contains(urlPath, `/min.`) && !strings.Contains(urlPath, `.min.`) {
		b, err := minify.MinifyJS([]byte(con))
		if err != nil {
			fmt.Println(err)
		}
		con = string(b)
		con = regexCssCleanComment.ReplaceAllString(con, ``)
	}
	content += con
	return
}

func (s *Static) genCombinedCSS(absPath, urlPath string, onImportFn func(string)) (content string, err error) {
	con, err := com.ReadFileS(absPath)
	if err != nil {
		return ``, err
	}
	all := regexCssUrlAttr.FindAllStringSubmatch(con, -1)
	dir := path.Dir(s.CssURL(urlPath))
	for _, v := range all {
		res := dir
		val := v[1]
		for strings.HasPrefix(val, "../") {
			res = path.Dir(res)
			val = strings.TrimPrefix(val, "../")
		}
		con = strings.Replace(con, v[0], "url('"+res+"/"+strings.TrimLeft(val, "/")+"')", 1)
	}
	all = regexCssImport.FindAllStringSubmatch(con, -1)
	absDir := filepath.Dir(absPath)
	for _, v := range all {
		val := v[1]
		res := dir
		absRes := absDir
		for strings.HasPrefix(val, "../") {
			res = path.Dir(res)
			absRes = path.Dir(absRes)
			val = strings.TrimPrefix(val, "../")
		}
		val = strings.TrimLeft(val, "/")
		//con = strings.Replace(con, v[0], `@import "`+res+"/"+val+`";`, 1)
		if icon, err := com.ReadFileS(absRes + "/" + val); err != nil {
			fmt.Println(err)
		} else {
			if onImportFn != nil {
				onImportFn(strings.Trim(res, `/`) + "/" + val)
			}
			con = strings.Replace(con, v[0], icon, 1)
		}
	}
	content += "\n/* <from: " + urlPath + "> */\n"
	con = regexCssCleanComment.ReplaceAllString(con, ``)
	con = regexCssCleanSpace.ReplaceAllString(con, `$1`)
	con = regexCssCleanSpace2.ReplaceAllString(con, ` `)
	content += con
	return
}

// =====================
// Handle
// =====================

func (s *Static) HandleMinify(ctx echo.Context, filePathFn func(string) string) error {
	fileStr := ctx.Param(`files`)
	fileType := ctx.Param(`type`)
	files := strings.Split(fileStr, `,`)
	var (
		name    string
		content string
		mtime   time.Time
		reader  io.ReadSeeker
	)
	if len(files) < 1 || (fileType != `js` && fileType != `css`) {
		return nil
	}
	name = files[0]

	combinedSavePath, _ := s.cachedURLInfo(fileStr, fileType)
	if s.IsCombined(combinedSavePath) == false || com.FileExists(combinedSavePath) == false {
		var onImportFn = func(urlPath string) {
			s.RecordCombined(filePathFn(urlPath), combinedSavePath)
		}
		switch fileType {
		case `js`:
			for _, f := range files {
				if strings.Contains(f, `..`) {
					continue
				}
				f += `.` + fileType
				absPath := filePathFn(f)
				if con, err := s.genCombinedJS(absPath, f); err != nil {
					fmt.Println(err)
				} else {
					s.RecordCombined(absPath, combinedSavePath)
					content += con
				}
			}
		case `css`:
			for _, f := range files {
				if strings.Contains(f, `..`) {
					continue
				}
				f += `.` + fileType
				absPath := filePathFn(f)
				if con, err := s.genCombinedCSS(absPath, f, onImportFn); err != nil {
					fmt.Println(err)
				} else {
					s.RecordCombined(absPath, combinedSavePath)
					content += con
				}
			}
		}
		byteContent := []byte(content)
		com.WriteFile(combinedSavePath, byteContent)
		s.RecordCombines(combinedSavePath)
		reader = bytes.NewReader(byteContent)
		mtime = time.Now()
	} else {
		fh, err := os.Open(combinedSavePath)
		if err != nil {
			return err
		}
		defer fh.Close()
		fi, err := fh.Stat()
		if err != nil {
			return err
		}
		reader = fh
		mtime = fi.ModTime()
	}
	return ctx.ServeContent(reader, name, mtime)
}
