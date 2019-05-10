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
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/admpub/log"
)

// NewWorker New a worker, if ipstring is a proxy address, New a proxy client.
// Proxy address such as:
// 		http://[user]:[password@]ip:port, [] stand it can choose or not. case: socks5://127.0.0.1:1080
func NewWorker(ipstring interface{}, timeout ...time.Duration) (*Worker, error) {
	worker := new(Worker)
	worker.Header = http.Header{}
	worker.Data = url.Values{}
	worker.BinaryData = []byte{}
	if ipstring != nil {
		client, err := NewProxyClient(strings.ToLower(ipstring.(string)))
		worker.Client = client
		worker.Ipstring = ipstring.(string)
		return worker, err
	}
	client, err := NewClient(timeout...)
	worker.Client = client
	worker.Ipstring = "localhost"
	return worker, err
}

// New Alias Name for NewWorker
func New(ipstring interface{}, timeout ...time.Duration) (*Worker, error) {
	return NewWorker(ipstring, timeout...)
}

// NewWorkerByClient New Worker by Your Client
func NewWorkerByClient(client *http.Client) *Worker {
	worker := new(Worker)
	worker.Header = http.Header{}
	worker.Data = url.Values{}
	worker.BinaryData = []byte{}

	// API must can set timeout
	if DefaultTimeOut != 0 && client.Timeout < 1 {
		client.Timeout = time.Second * time.Duration(DefaultTimeOut)
	}
	worker.Client = client
	return worker
}

// NewAPI New API Worker, No Cookie Keep.
func NewAPI() *Worker {
	return NewWorkerByClient(NoCookieClient)
}

// Go Auto decide which method, Default Get.
func (worker *Worker) Go() (body []byte, e error) {
	switch strings.ToUpper(worker.Method) {
	case POST:
		return worker.Post()
	case POSTJSON:
		return worker.PostJSON()
	case POSTXML:
		return worker.PostXML()
	case POSTFILE:
		return worker.PostFile()
	case PUT:
		return worker.Put()
	case PUTJSON:
		return worker.PutJSON()
	case PUTXML:
		return worker.PutXML()
	case PUTFILE:
		return worker.PutFile()
	case DELETE:
		return worker.Delete()
	case OTHER:
		return []byte(""), errors.New("please use method OtherGo(method, contentType string) or OtherGoBinary(method, contentType string)")
	default:
		return worker.Get()
	}
}

func (worker *Worker) GoByMethod(method string) (body []byte, e error) {
	return worker.SetMethod(method).Go()
}

// ToString This make effect only your worker exec serial! Attention!
// Change Your Raw data To string
func (worker *Worker) ToString() string {
	if worker.Raw == nil {
		return ""
	}
	return string(worker.Raw)
}

// JSONToString This make effect only your worker exec serial! Attention!
// Change Your JSON'like Raw data to string
func (worker *Worker) JSONToString() (string, error) {
	if worker.Raw == nil {
		return "", nil
	}
	temp, err := json.Marshal(worker.Raw)
	if err != nil {
		return "", err
	}
	return string(temp), nil
}

// Main method I make!
func (worker *Worker) sent(method, contentType string, binary bool) (body []byte, e error) {
	// Lock it for save
	worker.mux.Lock()
	defer worker.mux.Unlock()

	// Before FAction we can change or add something before Go()
	if worker.BeforeAction != nil {
		worker.BeforeAction(worker.Ctx, worker)
	}

	// Wait if must
	if worker.Wait > 0 {
		Wait(worker.Wait)
	}

	// For debug
	Debugf("[GoWorker] %s %s", method, worker.Url)

	// New a Request
	var (
		request = &http.Request{}
		err     error
	)

	// If binary parm value is true and BinaryData is not empty
	// suit for POSTJSON(), PostFile()
	if len(worker.BinaryData) != 0 && binary {
		contentReader := bytes.NewReader(worker.BinaryData)
		request, err = http.NewRequest(method, worker.Url, contentReader)
	} else if len(worker.Data) != 0 { // such POST() from table form
		contentReader := strings.NewReader(worker.Data.Encode())
		request, err = http.NewRequest(method, worker.Url, contentReader)
	} else {
		request, err = http.NewRequest(method, worker.Url, nil)
	}
	if err != nil {
		return nil, err
	}

	// Close avoid EOF
	// For client requests, setting this field prevents re-use of
	// TCP connections between requests to the same hosts, as if
	// Transport.DisableKeepAlives were set.
	// todo
	// maybe you want long connection
	//request.Close = true

	// Clone Header, I add some HTTP header!
	request.Header = CloneHeader(worker.Header)

	// In fact content type must not empty
	if contentType != "" {
		request.Header.Set("Content-Type", contentType)
	}
	worker.Request = request

	// Debug for RequestHeader
	OutputMaps("Request header", request.Header)

	// Tolerate abnormal way to create a Worker
	if worker.Client == nil {
		worker.Client = Client
	}

	// Do it
	response, err := worker.Client.Do(request)

	// Close it attention response may be nil
	if response != nil {
		//response.Close = true
		defer response.Body.Close()
	}

	if err != nil {
		// I count Error time
		worker.Errortimes++
		return nil, err
	}

	// Debug
	OutputMaps("Response header", response.Header)
	Debugf("[GoWorker] %v %s", response.Proto, response.Status)

	// Read output
	body, e = ioutil.ReadAll(response.Body)
	if e == nil {
		body, e = FixCharset(body, response, worker.DetectCharset, worker.ResponseCharset)
	}
	if e != nil {
		log.Error(e)
	}
	worker.Raw = body

	worker.Statuscode = response.StatusCode
	worker.Preurl = worker.Url
	worker.Response = response
	worker.Fetchtimes++

	// After action
	if worker.AfterAction != nil {
		worker.AfterAction(worker.Ctx, worker)
		body = worker.Raw
	}
	return
}

// Get method
func (worker *Worker) Get() (body []byte, e error) {
	worker.Clear()
	return worker.sent(GET, "", false)
}

func (worker *Worker) Delete() (body []byte, e error) {
	worker.Clear()
	return worker.sent(DELETE, "", false)
}

// Post Almost include bellow:
/*
	"application/x-www-form-urlencoded"
	"application/json"
	"text/xml"
	"multipart/form-data"
*/
func (worker *Worker) Post() (body []byte, e error) {
	return worker.sent(POST, HTTPFORMContentType, false)
}

func (worker *Worker) PostJSON() (body []byte, e error) {
	return worker.sent(POST, HTTPJSONContentType, true)
}

func (worker *Worker) PostXML() (body []byte, e error) {
	return worker.sent(POST, HTTPXMLContentType, true)
}

func (worker *Worker) PostFile() (body []byte, e error) {
	return worker.sentFile(POST)

}

func (worker *Worker) sentFile(method string) ([]byte, error) {
	if worker.FileName == "" || worker.FileFormName == "" {
		return nil, errors.New("fileName or fileFormName must not empty")
	}
	if len(worker.BinaryData) == 0 {
		return nil, errors.New("BinaryData must not empty")
	}

	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	fileWriter, err := bodyWriter.CreateFormFile(worker.FileFormName, worker.FileName)
	if err != nil {
		return nil, err
	}

	fileWriter.Write(worker.BinaryData)
	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	worker.SetBinary(bodyBuf.Bytes())

	return worker.sent(method, contentType, true)
}

// Put .
func (worker *Worker) Put() (body []byte, e error) {
	return worker.sent(PUT, HTTPFORMContentType, false)
}

func (worker *Worker) PutJSON() (body []byte, e error) {
	return worker.sent(PUT, HTTPJSONContentType, true)
}

func (worker *Worker) PutXML() (body []byte, e error) {
	return worker.sent(PUT, HTTPXMLContentType, true)
}

func (worker *Worker) PutFile() (body []byte, e error) {
	return worker.sentFile(PUT)

}

/*
OtherGo Other Method

     Method         = "OPTIONS"                ; Section 9.2
                    | "GET"                    ; Section 9.3
                    | "HEAD"                   ; Section 9.4
                    | "POST"                   ; Section 9.5
                    | "PUT"                    ; Section 9.6
                    | "DELETE"                 ; Section 9.7
                    | "TRACE"                  ; Section 9.8
                    | "CONNECT"                ; Section 9.9
                    | extension-method
   extension-method = token
     token          = 1*<any CHAR except CTLs or separators>


Content Type

	"application/x-www-form-urlencoded"
	"application/json"
	"text/xml"
	"multipart/form-data"
*/
func (worker *Worker) OtherGo(method, contentType string) (body []byte, e error) {
	return worker.sent(method, contentType, false)
}

func (worker *Worker) OtherGoBinary(method, contentType string) (body []byte, e error) {
	return worker.sent(method, contentType, true)
}
