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

package xheditor

import (
	"net/url"

	uploadClient "github.com/webx-top/client/upload"
)

func init() {
	uploadClient.Register(`xheditor`, func() uploadClient.Client {
		return New()
	})
}

func New() uploadClient.Client {
	client := &XhEditor{}
	client.BaseClient = uploadClient.New(client)
	return client
}

type XhEditor struct {
	*uploadClient.BaseClient
}

func (a *XhEditor) Name() string {
	return "filedata"
}

func (a *XhEditor) Body() (file uploadClient.ReadCloserWithSize, err error) {
	file, a.Data.FileName, err = uploadClient.Receive(a.Name(), a.Context)
	if err != nil {
		return
	}
	return
}

func (a *XhEditor) Result(errMsg string) (r string) {
	var msg, publicURL string
	if a.Form("immediate") == "1" {
		publicURL = "!" + a.Data.FileURL
	} else {
		publicURL = a.Data.FileURL
	}
	switch a.Data.FileType {
	case uploadClient.TypeImage, "":
		msg = `{"url":"` + publicURL + `||||` + url.QueryEscape(a.Data.FileName) + `","localname":"` + a.Data.FileName + `","id":"` + a.Data.FileIdString() + `"}`
	case uploadClient.TypeFlash,
		uploadClient.TypeAudio, uploadClient.TypeVideo,
		"media", "file":
		fallthrough
	default:
		msg = `{"url":"` + publicURL + `","id":"` + a.Data.FileIdString() + `"}`
	}
	if len(msg) == 0 {
		msg = "{}"
	}
	r = `{"err":"` + errMsg + `","msg":` + msg + `}`
	return
}
