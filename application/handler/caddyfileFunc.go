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
package handler

import "github.com/webx-top/echo"

func addonAttr(ctx echo.Context) {
	ctx.SetFunc(`AddonAttr`, func(addon string, item string) string {
		if len(addon) > 0 {
			addon += `_`
		}
		k := addon + item
		v := ctx.Form(k)
		if len(v) == 0 {
			return ``
		}
		return item + `   ` + v
	})
}

func iteratorKV(ctx echo.Context) {
	ctx.SetFunc(`IteratorKV`, func(addon string, item string, prefix string) string {
		if len(addon) > 0 {
			addon += `_`
		}
		k := addon + item + `_k`
		keys := ctx.FormValues(k)

		k = addon + item + `_v`
		values := ctx.FormValues(k)

		r := ``
		l := len(values)
		t := ``
		for i, v := range keys {
			if i < l {
				r += t + prefix + v + `   ` + values[i]
				t = "\n"
			}
		}
		return r
	})
}

func SetCaddyfileFunc(ctx echo.Context) {
	addonAttr(ctx)
	iteratorKV(ctx)
}
