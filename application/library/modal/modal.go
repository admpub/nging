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
package modal

import (
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/admpub/confl"
	"github.com/webx-top/echo"
)

var DefaultModal = Modal{
	ExtButtons: []Button{},
}

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

var modalConfig = map[string]Modal{}

func Render(ctx echo.Context, param interface{}) template.HTML {
	var data Modal
	if v, y := param.(*Modal); y {
		data = *v
	} else if v, y := param.(Modal); y {
		data = v
	} else if v, y := param.(string); y {
		if ov, ok := modalConfig[v]; ok {
			data = ov
		} else {
			_, err := confl.DecodeFile(v, &data)
			if err != nil {
				if os.IsNotExist(err) || strings.Contains(err.Error(), `cannot find the file`) {
					var b []byte
					data = DefaultModal
					b, err = confl.Marshal(data)
					if err == nil {
						err = os.MkdirAll(filepath.Dir(v), os.ModePerm)
						if err == nil {
							err = ioutil.WriteFile(v, b, os.ModePerm)
						}
					}
				}
				if err != nil {
					return template.HTML(err.Error())
				}
			}
			modalConfig[v] = data
		}
	}
	b, err := ctx.Fetch(`modal`, data)
	if err != nil {
		return template.HTML(err.Error())
	}
	return template.HTML(string(b))
}

func Remove(confPath string) error {
	if _, ok := modalConfig[confPath]; ok {
		delete(modalConfig, confPath)
	}
	return nil
}

func Clear() error {
	modalConfig = map[string]Modal{}
	return nil
}
