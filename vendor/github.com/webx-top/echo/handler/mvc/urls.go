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
package mvc

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/webx-top/com"
	"github.com/webx-top/echo"
)

func NewURLs(project string, mvc *Application) *URLs {
	return &URLs{
		projectPath: `github.com/webx-top/` + project,
		extensions:  map[string]map[string]int{},
		Application: mvc,
	}
}

type URLs struct {
	projectPath  string
	extensions   map[string]map[string]int
	*Application `json:"-" xml:"-"`
}

func (a *URLs) SetProjectPath(projectPath string) {
	a.projectPath = strings.TrimSuffix(projectPath, `/`)
}

func (a *URLs) urlRecovery(s string) string {
	if a.URLRecovery != nil {
		return a.URLRecovery(s)
	}
	return s
}

func (a *URLs) Build(mdl string, ctl string, act string, params ...interface{}) (r string) {
	module, ok := a.Application.moduleNames[mdl]
	if !ok {
		return
	}
	pkg := a.projectPath + `/app/` + module.Name + `/controller`
	key := ``
	if len(ctl) == 0 {
		key = pkg + `.` + a.urlRecovery(act)
	} else {
		key = pkg + `.(*` + a.urlRecovery(ctl) + `).` + a.urlRecovery(act) + `-fm`
	}
	r = module.Router().URL(key, params...)
	if len(module.Domain) > 0 {
		scheme := `http`
		if a.Application.SessionOptions.Secure {
			scheme = `https`
		}
		r = scheme + `://` + module.Domain + r
	}
	return
}

func (a *URLs) BuildFromPath(ppath string, args ...map[string]interface{}) (r string) {
	var mdl, ctl, act string
	uris := strings.SplitN(ppath, "?", 2)
	ret := strings.SplitN(uris[0], `/`, 3)
	switch len(ret) {
	case 3:
		act = ret[2]
		ctl = ret[1]
		mdl = ret[0]
	case 2:
		act = ret[1]
		mdl = ret[0]
	default:
		return
	}
	module, ok := a.Application.moduleNames[mdl]
	if !ok {
		return
	}
	pkg := a.projectPath + `/app/` + module.Name + `/controller`
	key := ``
	if len(ctl) == 0 {
		key = pkg + `.` + a.urlRecovery(act)
	} else {
		key = pkg + `.(*` + a.urlRecovery(ctl) + `).` + a.urlRecovery(act) + `-fm`
	}
	var params url.Values
	if len(uris) > 1 {
		params, _ = url.ParseQuery(uris[1])
	}
	if len(args) > 0 {
		for k, v := range args[0] {
			params.Set(k, fmt.Sprintf("%v", v))
		}
	}
	r = module.Router().URL(key, params)
	if len(module.Domain) > 0 {
		scheme := `http`
		if a.Application.SessionOptions.Secure {
			scheme = `https`
		}
		r = scheme + `://` + module.Domain + r
	}
	return
}

func (a *URLs) Set(h interface{}) (pkg string, ctl string, act string) {
	key := echo.HandlerName(h)
	if _, ok := a.extensions[key]; !ok {
		a.extensions[key] = map[string]int{}
	}
	pkg, ctl, act = com.ParseFuncName(key)
	return
}

func (a *URLs) SetExtensions(key string, exts []string) *URLs {
	e := map[string]int{}
	for key, val := range exts {
		e[val] = key
	}
	a.extensions[key] = e
	return a
}

func (a *URLs) AllowFormat(key string, ext string) (ok bool) {
	if ex, y := a.extensions[key]; !y || ex == nil || len(ex) < 1 {
		ok = true
	} else {
		_, ok = ex[ext]
	}
	return
}
