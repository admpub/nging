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

package markdown

import (
	"io"

	"net/url"
	"time"

	uploadClient "github.com/webx-top/client/upload"
)

func init() {
	uploadClient.Register(`markdown`, func() uploadClient.Client {
		return New()
	})
}

func New() uploadClient.Client {
	client := &Markdown{}
	client.BaseClient = uploadClient.New(client)
	return client
}

type Markdown struct {
	*uploadClient.BaseClient
}

func (a *Markdown) Name() string {
	return "editormd-image-file"
}

func (a *Markdown) Body() (file io.ReadCloser, err error) {
	file, a.Data.FileName, err = uploadClient.Receive(a.Name(), a.Context)
	if err != nil {
		return
	}
	return
}

func (a *Markdown) Result(errMsg string) (r string) {
	succed := "0" // 0 表示上传失败，1 表示上传成功
	if len(errMsg) > 0 {
		succed = "1"
	}
	callback := a.Form(`callback`)
	dialogID := a.Form(`dialog_id`)
	if len(callback) > 0 && len(dialogID) > 0 {
		//跨域上传返回操作
		nextURL := callback + "?dialog_id=" + dialogID + "&temp=" + time.Now().String() + "&success=" + succed + "&message=" + url.QueryEscape(errMsg) + "&url=" + a.Data.FileURL
		a.Redirect(nextURL)
	} else {
		r = `{"success":` + succed + `,"message":"` + errMsg + `","url":"` + a.Data.FileURL + `","id":"` + a.Data.FileIdString() + `"}`
	}
	return
}
