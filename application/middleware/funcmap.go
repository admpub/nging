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
package middleware

import (
	"fmt"
	"html/template"
	"time"

	"strconv"

	"github.com/admpub/nging/application/library/config"
	"github.com/admpub/nging/application/library/errors"
	"github.com/admpub/nging/application/library/modal"
	"github.com/webx-top/echo"
)

func FuncMap() echo.MiddlewareFunc {
	return func(h echo.Handler) echo.Handler {
		return echo.HandlerFunc(func(c echo.Context) error {
			c.SetFunc(`Now`, time.Now)
			c.SetFunc(`HasString`, hasString)
			c.SetFunc(`Date`, date)
			c.SetFunc(`Modal`, func(data interface{}) template.HTML {
				return modal.Render(c, data)
			})
			c.SetFunc(`IsMessage`, errors.IsMessage)
			c.SetFunc(`Languages`, func() []string {
				return config.DefaultConfig.Language.AllList
			})
			c.SetFunc(`IsError`, errors.IsError)
			c.SetFunc(`IsOk`, errors.IsOk)
			c.SetFunc(`Message`, errors.Message)
			c.SetFunc(`Ok`, errors.Ok)
			c.SetFunc(`IndexStrSlice`, indexStrSlice)
			c.SetFunc(`Version`, config.Version)
			return h.Handle(c)
		})
	}
}

func indexStrSlice(slice []string, index int) string {
	if slice == nil {
		return ``
	}
	if index >= len(slice) {
		return ``
	}
	return slice[index]
}

func hasString(slice []string, str string) bool {
	if slice == nil {
		return false
	}
	for _, v := range slice {
		if v == str {
			return true
		}
	}
	return false
}

func date(timestamp interface{}) time.Time {
	if v, y := timestamp.(int64); y {
		return time.Unix(v, 0)
	}
	if v, y := timestamp.(uint); y {
		return time.Unix(int64(v), 0)
	}
	v, _ := strconv.ParseInt(fmt.Sprint(timestamp), 10, 64)
	return time.Unix(v, 0)
}
