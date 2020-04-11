//
// 	Copyright 2017 by marmot author: gdccmcm14@live.com.
// 	Licensed under the Apache License, Version 2.0 (the "License");
// 	you may not use this file except in compliance with the License.
// 	You may obtain a copy of the License at
// 		http://www.apache.org/licenses/LICENSE-2.0
// 	Unless required by applicable law or agreed to in writing, software
// 	distributed under the License is distributed on an "AS IS" BASIS,
// 	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// 	See the License for the specific language governing permissions and
// 	limitations under the License
//

package miner

import (
	"context"
	"net/http"
	"net/url"

	"github.com/admpub/pester"
)

// DefaultWorker Global Worker
var DefaultWorker *Worker

func init() {
	UserAgentInit()

	// New a Worker
	worker := new(Worker)
	worker.Header = http.Header{}
	worker.Data = url.Values{}
	worker.BinaryData = []byte{}
	worker.Client = Client

	// Global Worker!
	DefaultWorker = worker

}

// ToString This make effect only your Worker exec serial! Attention!
// Change Your Raw data To string
func ToString() string {
	return DefaultWorker.ToString()
}

// JSONToString This make effect only your Worker exec serial! Attention!
// Change Your JSON'like Raw data to string
func JSONToString() (string, error) {
	return DefaultWorker.JSONToString()
}

func Get() (body []byte, e error) {
	return DefaultWorker.Get()
}

func Delete() (body []byte, e error) {
	return DefaultWorker.Delete()
}

func Go() (body []byte, e error) {
	return DefaultWorker.Go()
}

func GoByMethod(method string) (body []byte, e error) {
	return DefaultWorker.SetMethod(method).Go()
}

func OtherGo(method, contenttype string) (body []byte, e error) {
	return DefaultWorker.OtherGo(method, contenttype)
}

func Post() (body []byte, e error) {
	return DefaultWorker.Post()
}

func PostJSON() (body []byte, e error) {
	return DefaultWorker.PostJSON()
}

func PostFile() (body []byte, e error) {
	return DefaultWorker.PostFile()
}

func PostXML() (body []byte, e error) {
	return DefaultWorker.PostXML()
}

func Put() (body []byte, e error) {
	return DefaultWorker.Put()
}
func PutJSON() (body []byte, e error) {
	return DefaultWorker.PutJSON()
}

func PutFile() (body []byte, e error) {
	return DefaultWorker.PutFile()
}

func PutXML() (body []byte, e error) {
	return DefaultWorker.PutXML()
}

func SetDetectCharset(on bool) *Worker {
	return DefaultWorker.SetDetectCharset(on)
}

func SetHeaderParm(k, v string) *Worker {
	return DefaultWorker.SetHeaderParm(k, v)
}

func SetMaxRetries(maxRetries int) *Worker {
	return DefaultWorker.SetMaxRetries(maxRetries)
}

func SetPesterOptions(options ...pester.ApplyOptions) *Worker {
	return DefaultWorker.SetPesterOptions(options...)
}

func SetResponseCharset(charset string) *Worker {
	return DefaultWorker.SetResponseCharset(charset)
}

func SetCookie(v string) *Worker {
	return DefaultWorker.SetCookie(v)
}

func SetCookieByFile(file string) (*Worker, error) {
	return DefaultWorker.SetCookieByFile(file)
}

func SetUserAgent(ua string) *Worker {
	return DefaultWorker.SetUserAgent(ua)
}
func SetRefer(refer string) *Worker {
	return DefaultWorker.SetRefer(refer)
}

func SetURL(url string) *Worker {
	return DefaultWorker.SetURL(url)
}

func SetFileInfo(fileName, fileFormName string) *Worker {
	return DefaultWorker.SetFileInfo(fileName, fileFormName)
}

func SetMethod(method string) *Worker {
	return DefaultWorker.SetMethod(method)
}

func SetWaitTime(num int) *Worker {
	return DefaultWorker.SetWaitTime(num)
}

func SetBinary(data []byte) *Worker {
	return DefaultWorker.SetBinary(data)
}

func SetForm(form url.Values) *Worker {
	return DefaultWorker.SetForm(form)
}

func SetFormParm(k, v string) *Worker {
	return DefaultWorker.SetFormParm(k, v)
}

func SetContext(ctx context.Context) *Worker {
	return DefaultWorker.SetContext(ctx)
}

func SetBeforeAction(fc func(context.Context, *Worker)) *Worker {
	return DefaultWorker.SetBeforeAction(fc)
}

func SetAfterAction(fc func(context.Context, *Worker)) *Worker {
	return DefaultWorker.SetAfterAction(fc)
}

func Clear() *Worker {
	return DefaultWorker.Clear()
}

func ClearAll() *Worker {
	return DefaultWorker.ClearAll()
}

func ClearCookie() *Worker {
	return DefaultWorker.ClearCookie()
}

func GetCookies() []*http.Cookie {
	return DefaultWorker.GetCookies()
}
