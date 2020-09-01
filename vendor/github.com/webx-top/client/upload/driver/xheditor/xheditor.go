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
	"github.com/webx-top/echo"
)

func init() {
	uploadClient.Register(`xheditor`, func() uploadClient.Client {
		return New()
	})
}

var FormField = `filedata`

func New() uploadClient.Client {
	client := &XhEditor{}
	client.BaseClient = uploadClient.New(client, FormField)
	return client
}

type XhEditor struct {
	*uploadClient.BaseClient
}

func (a *XhEditor) BuildResult() {
	var publicURL string
	if a.Form("immediate") == "1" {
		publicURL = "!" + a.Data.FileURL
	} else {
		publicURL = a.Data.FileURL
	}
	data := echo.H{
		`id`: a.Data.FileIdString(),
	}
	switch a.Data.FileType {
	case uploadClient.TypeImage, "":
		data[`url`] = publicURL + `||||` + url.QueryEscape(a.Data.FileName)
		data[`localname`] = a.Data.FileName
	case uploadClient.TypeFlash,
		uploadClient.TypeAudio, uploadClient.TypeVideo,
		"media", "file":
		fallthrough
	default:
		data[`url`] = publicURL
	}
	var errMsg string
	if a.GetError() != nil {
		errMsg = a.Error()
	}
	a.RespData = echo.H{
		`err`: errMsg,
		`msg`: data,
	}
}
