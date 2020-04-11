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

import "net/http"

const (
	// WaitTime Default wait time
	WaitTime = 5

	// HTTP method

	GET      = "GET"
	POST     = "POST"
	POSTJSON = "POSTJSON"
	POSTXML  = "POSTXML"
	POSTFILE = "POSTFILE"
	PUT      = "PUT"
	PUTJSON  = "PUTJSON"
	PUTXML   = "PUTXML"
	PUTFILE  = "PUTFILE"
	DELETE   = "DELETE"
	OTHER    = "OTHER" // this stand for you can use other method this lib not own.

	// HTTP content type

	HTTPFORMContentType = "application/x-www-form-urlencoded"
	HTTPJSONContentType = "application/json"
	HTTPXMLContentType  = "text/xml"
	HTTPFILEContentType = "multipart/form-data"
)

var (
	// Browser User-Agent, Our default Http ua header!
	defaultUserAgent = "Marmot"

	DefaultHeader = map[string][]string{
		"User-Agent": {
			defaultUserAgent,
		},
	}

	// DefaultTimeOut http get and post No timeout
	DefaultTimeOut = 30
)

// SetGlobalTimeout Set global timeout, it can only by this way!
func SetGlobalTimeout(num int) {
	DefaultTimeOut = num
}

// MergeCookie Merge Cookie, not use
func MergeCookie(before []*http.Cookie, after []*http.Cookie) []*http.Cookie {
	cs := make(map[string]*http.Cookie)

	for _, b := range before {
		cs[b.Name] = b
	}

	for _, a := range after {
		if a.Value != "" {
			cs[a.Name] = a
		}
	}

	res := make([]*http.Cookie, 0, len(cs))

	for _, q := range cs {
		res = append(res, q)
	}

	return res

}

// CloneHeader Clone a header, If not exist Ua, Set our Ua!
func CloneHeader(h map[string][]string) map[string][]string {
	if h == nil || len(h) == 0 {
		h = DefaultHeader
		return h
		//return map[string][]string{}
	}

	if len(h["User-Agent"]) == 0 {
		h["User-Agent"] = []string{defaultUserAgent}
	}
	return CopyM(h)
}
