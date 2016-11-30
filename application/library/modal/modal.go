package modal

import (
	"html/template"

	"io/ioutil"
	"os"

	"github.com/admpub/confl"
	"github.com/webx-top/echo"
)

type HTMLAttr struct {
	Attr  string
	Value interface{}
}

type Button struct {
	Attributes []HTMLAttr
	Text       string
}

type Modal struct {
	Id         string
	Custom     bool
	HeadTitle  interface{}
	Title      interface{}
	Content    interface{}
	Type       string
	ExtButtons []Button
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
				if os.IsNotExist(err) {
					var b []byte
					b, err = confl.Marshal(data)
					if err == nil {
						err = ioutil.WriteFile(v, b, os.ModePerm)
					}
				}
				return template.HTML(err.Error())
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
