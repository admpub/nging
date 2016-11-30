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

import (
	"encoding/json"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"

	"fmt"

	"github.com/admpub/caddyui/application/library/config"
	"github.com/admpub/caddyui/application/model"
	"github.com/webx-top/echo"
)

func ManageIndex(ctx echo.Context) error {
	return ctx.Render(`manage/index`, ctx.Flash())
}

func ManageVhostAdd(ctx echo.Context) error {
	var err error
	if ctx.IsPost() {
		m := model.NewVhost(ctx)
		m.Domain = ctx.Form(`domain`)
		m.Disabled = ctx.Form(`disabled`)
		var b []byte
		b, err = json.Marshal(ctx.Forms())
		switch {
		case err == nil:
			m.Setting = string(b)
			_, err = m.Add()
			if err != nil {
				break
			}
			fallthrough
		case 0 == 1:
			err = saveVhostData(ctx, m)
		}
	}
	return ctx.Render(`manage/vhost_edit`, err)
}

func saveVhostData(ctx echo.Context, m *model.Vhost) (err error) {
	var b []byte
	var saveFile string
	SetCaddyfileFunc(ctx)
	b, err = ctx.Fetch(`manage/caddyfile`, nil)
	if err != nil {
		return
	}
	saveFile, err = filepath.Abs(config.DefaultConfig.Sys.VhostsfileDir)
	if err != nil {
		return
	}
	saveFile = filepath.Join(saveFile, fmt.Sprint(m.Id))
	if m.Disabled == `Y` {
		err = os.Remove(saveFile)
		if os.IsNotExist(err) {
			err = nil
		}
	} else {
		err = ioutil.WriteFile(saveFile, b, os.ModePerm)
		if len(ctx.Form(`restart`)) > 0 {
			err = config.DefaultCLIConfig.CaddyRestart()
		}
	}
	return
}

func ManageVhostEdit(ctx echo.Context) error {
	id := ctx.Formx(`id`).Uint()
	if id < 1 {
		ctx.Session().AddFlash(ctx.T(`id无效`)).Save()
		return ctx.Redirect(`/manage`)
	}

	var err error
	m := model.NewVhost(ctx)
	err = m.Get(nil, `id`, id)
	if err != nil {
		ctx.Session().AddFlash(err.Error()).Save()
		return ctx.Redirect(`/manage`)
	}
	if ctx.IsPost() {
		m.Domain = ctx.Form(`domain`)
		m.Disabled = ctx.Form(`disabled`)
		var b []byte
		b, err = json.Marshal(ctx.Forms())
		switch {
		case err == nil:
			m.Setting = string(b)
			err = m.Edit(nil, `id`, id)
			if err != nil {
				break
			}
			fallthrough
		case 0 == 1:
			err = saveVhostData(ctx, m)
		}
	} else {
		var formData url.Values
		if e := json.Unmarshal([]byte(m.Setting), &formData); e == nil {
			for key, values := range formData {
				for _, v := range values {
					ctx.Request().Form().Add(key, v)
				}
			}
		}
	}
	return ctx.Render(`manage/vhost_edit`, err)
}

func ManageRestart(ctx echo.Context) error {
	if err := config.DefaultCLIConfig.CaddyRestart(); err != nil {
		return err
	}
	return ctx.String(ctx.T(`已经完成重启`))
}
