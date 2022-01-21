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

package dropzone

import (
	"net/http"

	uploadClient "github.com/webx-top/client/upload"
	"github.com/webx-top/echo"
)

func init() {
	uploadClient.Register(`dropzone`, func() uploadClient.Client {
		return New()
	})
}

var FormField = `file`

func New() uploadClient.Client {
	client := &Dropzone{}
	client.BaseClient = uploadClient.New(client, FormField)
	client.BaseClient.SetFieldMapping(MappingChunkInfo)
	return client
}

type Dropzone struct {
	*uploadClient.BaseClient
}

func (a *Dropzone) BuildResult() uploadClient.Client {
	if a.GetError() == nil {
		a.RespData = echo.H{
			`result`: echo.H{
				`url`: a.Data.FileURL,
				`id`:  a.Data.FileIDString(),
			},
			`error`: nil,
		}
	} else {
		a.Code = http.StatusInternalServerError
		a.ContentType = `string`
		a.RespData = a.ErrorString()
	}
	return a
}

var MappingChunkInfo = map[string]string{
	`fileUUID`:        `dzuuid`,
	`chunkIndex`:      `dzchunkindex`,
	`fileTotalBytes`:  `dztotalfilesize`,
	`fileChunkBytes`:  `dzchunksize`,
	`fileTotalChunks`: `dztotalchunkcount`,
	//`<unsuppored>`: `dzchunkbyteoffset`,
}
