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

func (a *Webuploader) Result() (r string) {
	cid := a.Form("id")
	if len(cid) == 0 {
		form := a.Request().MultipartForm()
		if form != nil && form.Value != nil {
			if v, ok := form.Value["id"]; ok && len(v) > 0 {
				cid = v[0]
			}
		}
	}
	if a.GetError() == nil {
		r = `{"jsonrpc":"2.0","result":{"url":"` + a.Data.FileURL + `","id":"` + a.Data.FileIdString() + `","containerid":"` + cid + `"},"error":null}`
		return
	}
	code := "100"
	r = `{"jsonrpc":"2.0","result":{"url":"` + a.Data.FileURL + `","id":"` + a.Data.FileIdString() + `","containerid":"` + cid + `"},"error":{"code":"` + code + `","message":"` + a.Error() + `"}}`

	return
}
