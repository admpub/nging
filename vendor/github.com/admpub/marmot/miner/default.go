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
	"net/http"
	"net/url"
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
