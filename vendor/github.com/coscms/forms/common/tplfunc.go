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

//Package common This package provides basic constants used by forms packages.
package common

import (
	"html/template"
	"io/fs"
	"io/ioutil"
	"path/filepath"
)

var TplFuncs = func() template.FuncMap {
	return template.FuncMap{}
}

func ParseFiles(files ...string) (*template.Template, error) {
	if !FileSystem.IsEmpty() {
		return ParseFS(FileSystem, files...)
	}
	name := filepath.Base(files[0])
	b, err := ioutil.ReadFile(files[0])
	if err != nil {
		return nil, err
	}
	tmpl := template.New(name)
	tmpl.Funcs(TplFuncs())
	tmpl = template.Must(tmpl.Parse(string(b)))
	if len(files) > 1 {
		tmpl, err = tmpl.ParseFiles(files[1:]...)
	}
	return tmpl, err
}

func ParseFS(fs fs.FS, files ...string) (*template.Template, error) {
	name := filepath.Base(files[0])
	tmpl := template.New(name)
	tmpl.Funcs(TplFuncs())
	fp, err := fs.Open(files[0])
	if err != nil {
		return tmpl, err
	}
	b, err := ioutil.ReadAll(fp)
	fp.Close()
	if err != nil {
		return tmpl, err
	}
	tmpl = template.Must(tmpl.Parse(string(b)))
	if len(files) > 1 {
		tmpl, err = tmpl.ParseFS(fs, files[1:]...)
	}
	return tmpl, err
}
