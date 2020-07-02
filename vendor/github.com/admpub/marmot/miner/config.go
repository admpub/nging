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
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/admpub/pester"
)

// Worker is the main object to sent http request and return result of response
type Worker struct {
	// In order fast chain func call I put the basic config below
	URL          string         // Which url we want
	Method       string         // Get/Post method
	Header       http.Header    // Http header
	Data         url.Values     // Sent by form data
	FileName     string         // FileName which sent to remote
	FileFormName string         // File Form Name which sent to remote
	BinaryData   []byte         // Sent by binary data, can together with File
	Wait         int            // Wait Time
	Client       *http.Client   // Our Client
	Request      *http.Request  // Debug
	Response     *http.Response // Debug
	Raw          []byte         // Raw data we get

	// The name below is not so good but has already been used in many project, so bear it.
	PreURL     string // Pre url
	StatusCode int    // the last url response code, such as 404
	FetchTimes int    // Url fetch number times
	ErrorTimes int    // Url fetch error times
	IP         string // worker proxy ip, just for user to record their proxy ip, default: localhost

	// AOP like Java
	Ctx           context.Context
	BeforeAction  func(context.Context, *Worker)
	AfterAction   func(context.Context, *Worker)
	PesterOptions []pester.ApplyOptions

	// ResponseCharset is the character encoding of the response body.
	// Leave it blank to allow automatic character encoding of the response body.
	ResponseCharset string

	// DetectCharset can enable character encoding detection for non-utf8 response bodies
	// without explicit charset declaration. This feature uses https://github.com/webx-top/chardet
	DetectCharset bool

	MaxRetries int

	mux sync.RWMutex // Lock, execute concurrently please use worker Pool!
}

// SetHeader Java Bean Chain pattern
func (worker *Worker) SetHeader(header http.Header) *Worker {
	worker.Header = header
	return worker
}

// SetHeader Default Worker SetHeader!
func SetHeader(header http.Header) *Worker {
	return DefaultWorker.SetHeader(header)
}

func (worker *Worker) SetHeaderParm(k, v string) *Worker {
	worker.Header.Set(k, v)
	return worker
}

func (worker *Worker) SetDetectCharset(on bool) *Worker {
	worker.DetectCharset = on
	return worker
}

func (worker *Worker) SetMaxRetries(maxRetries int) *Worker {
	worker.MaxRetries = maxRetries
	return worker
}

func (worker *Worker) SetPesterOptions(options ...pester.ApplyOptions) *Worker {
	worker.PesterOptions = options
	return worker
}

func (worker *Worker) SetResponseCharset(charset string) *Worker {
	worker.ResponseCharset = charset
	return worker
}

func (worker *Worker) SetCookie(v string) *Worker {
	worker.SetHeaderParm("Cookie", v)
	return worker
}

// SetCookieByFile Set Cookie by file.
func (worker *Worker) SetCookieByFile(file string) (*Worker, error) {
	haha, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	cookie := string(haha)
	cookie = strings.Replace(cookie, " ", "", -1)
	cookie = strings.Replace(cookie, "\n", "", -1)
	cookie = strings.Replace(cookie, "\r", "", -1)
	sconfig := worker.SetCookie(cookie)
	return sconfig, nil
}

func (worker *Worker) SetUserAgent(ua string) *Worker {
	worker.Header.Set("User-Agent", ua)
	return worker
}

func (worker *Worker) SetRefer(refer string) *Worker {
	worker.Header.Set("Referer", refer)
	return worker
}

func (worker *Worker) SetHost(host string) *Worker {
	worker.Header.Set("Host", host)
	return worker
}

// SetURL at the same time SetHost
func (worker *Worker) SetURL(url string) *Worker {
	worker.URL = url
	//https://www.zhihu.com/people/
	temp := strings.Split(url, "//")
	if len(temp) >= 2 {
		worker.SetHost(strings.Split(temp[1], "/")[0])
	}
	return worker
}

func (worker *Worker) SetFileInfo(fileName, fileFormName string) *Worker {
	worker.FileName = fileName
	worker.FileFormName = fileFormName
	return worker
}

func (worker *Worker) SetMethod(method string) *Worker {
	temp := GET
	v := strings.ToUpper(method)
	switch v {
	case GET, POST, POSTFILE, POSTJSON, POSTXML, PUT, PUTFILE, PUTJSON, PUTXML, DELETE:
		temp = v
	default:
		temp = OTHER
	}
	worker.Method = temp
	return worker
}

func (worker *Worker) SetWaitTime(num int) *Worker {
	if num <= 0 {
		num = 1
	}
	worker.Wait = num
	return worker
}

func (worker *Worker) SetBinary(data []byte) *Worker {
	worker.BinaryData = data
	return worker
}

func (worker *Worker) SetForm(form url.Values) *Worker {
	worker.Data = form
	return worker
}

func (worker *Worker) SetFormParm(k, v string) *Worker {
	worker.Data.Set(k, v)
	return worker
}

// Set Context so Action can soft
func (worker *Worker) SetContext(ctx context.Context) *Worker {
	worker.Ctx = ctx
	return worker
}

func (worker *Worker) SetBeforeAction(fc func(context.Context, *Worker)) *Worker {
	worker.BeforeAction = fc
	return worker
}

func (worker *Worker) SetAfterAction(fc func(context.Context, *Worker)) *Worker {
	worker.AfterAction = fc
	return worker
}

// Clear data we sent
func (worker *Worker) Clear() *Worker {
	worker.Data = url.Values{}
	worker.BinaryData = []byte{}
	return worker
}

// All clear include header
func (worker *Worker) ClearAll() *Worker {
	worker.Header = http.Header{}
	worker.Data = url.Values{}
	worker.BinaryData = []byte{}
	return worker
}

// ClearCookie Clear Cookie
func (worker *Worker) ClearCookie() *Worker {
	worker.Header.Del("Cookie")
	return worker
}

// GetCookies Get Cookies
func (worker *Worker) GetCookies() []*http.Cookie {
	if worker.Response != nil {
		return worker.Response.Cookies()
	} else {
		return []*http.Cookie{}
	}
}
