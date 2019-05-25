/*
   Nging is a toolbox for webmasters
   Copyright (C) 2018-present  Wenhui Shen <swh@admpub.com>

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU Affero General Public License as published
   by the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Affero General Public License for more details.

   You should have received a copy of the GNU Affero General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/
package modal

import (
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/admpub/confl"
	"github.com/admpub/log"
	"github.com/webx-top/echo"
)

var (
	modalConfig  = map[string]Modal{}
	DefaultModal = Modal{
		ExtButtons: []Button{},
	}
	PathFixer = func(conf string) string {
		return conf
	}
	ReadConfigFile = func(file string) ([]byte, error) {
		return ioutil.ReadFile(file)
	}
	mutext = &sync.RWMutex{}
)

type HTMLAttr struct {
	Attr  string      //属性名
	Value interface{} //属性值
}

type Button struct {
	Attributes []HTMLAttr //按钮属性
	Text       string     //按钮文本
}

type Modal struct {
	Id          string   //元素id
	Custom      bool     //是否自定义整个内容区域
	HeadTitle   string   //头部标题
	Title       string   //内容标题
	Content     string   //内容
	HelpText    string   //帮助提示
	Animate     string   //动画样式class名
	Type        string   //类型：warning/primary/success/danger
	ContentType string   //内容类型：form/blackform/""
	ExtButtons  []Button //附加按钮
}

func UnmarshalFile(confile string) (Modal, error) {
	mutext.Lock()
	defer mutext.Unlock()
	confile = PathFixer(confile)
	ov, ok := modalConfig[confile]
	if ok {
		return ov, nil
	}
	b, err := ReadConfigFile(confile)
	if err == nil {
		err = confl.Unmarshal(b, &ov)
	}
	if err != nil {
		if os.IsNotExist(err) || strings.Contains(err.Error(), `cannot find the file`) {
			var b []byte
			b, err = confl.Marshal(DefaultModal)
			if err == nil {
				err = os.MkdirAll(filepath.Dir(confile), os.ModePerm)
				if err == nil {
					err = ioutil.WriteFile(confile, b, os.ModePerm)
				}
			}
		}
		if err != nil {
			return ov, err
		}
	}
	modalConfig[confile] = ov
	return ov, nil
}

func Render(ctx echo.Context, param interface{}) template.HTML {
	var data Modal
	switch v := param.(type) {
	case *Modal:
		data = *v
	case Modal:
		data = v
	case string:
		var err error
		data, err = UnmarshalFile(v)
		if err != nil {
			return template.HTML(err.Error())
		}
	}
	b, err := ctx.Fetch(`modal`, data)
	if err != nil {
		return template.HTML(err.Error())
	}
	return template.HTML(string(b))
}

func Remove(confPath string) error {
	mutext.Lock()
	defer mutext.Unlock()
	if _, ok := modalConfig[confPath]; ok {
		delete(modalConfig, confPath)
		log.Debugf(`remove: modalConfig[%s] (remains:%d)`, confPath, len(modalConfig))
	}
	return nil
}

func Clear() error {
	mutext.Lock()
	defer mutext.Unlock()
	modalConfig = map[string]Modal{}
	log.Debugf(`clear: modalConfig (remains:%d)`, len(modalConfig))
	return nil
}
