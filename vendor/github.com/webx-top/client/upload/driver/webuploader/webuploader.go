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

package webuploader

import (
	uploadClient "github.com/webx-top/client/upload"
	"github.com/webx-top/echo"
)

func init() {
	uploadClient.Register(`webuploader`, func() uploadClient.Client {
		return New()
	})
}

var FormField = `file`

func New() uploadClient.Client {
	client := &Webuploader{}
	client.BaseClient = uploadClient.New(client, FormField)
	return client
}

type Webuploader struct {
	*uploadClient.BaseClient
}

func (a *Webuploader) BuildResult() uploadClient.Client {
	cid := a.Form("id")
	if len(cid) == 0 {
		form, err := a.Request().MultipartForm()
		if err != nil {
			a.SetError(err)
		}
		if form != nil && form.Value != nil {
			if v, ok := form.Value["id"]; ok && len(v) > 0 {
				cid = v[0]
			}
		}
	}
	data := echo.H{
		`jsonrpc`: `2.0`,
		`result`: echo.H{
			`url`:         a.Data.FileURL,
			`id`:          a.Data.FileIDString(),
			`containerid`: cid,
		},
		`error`: nil,
	}
	if a.GetError() != nil {
		data[`error`] = echo.H{
			`code`:    100,
			`message`: a.ErrorString(),
		}
	}
	a.RespData = data
	return a
}
