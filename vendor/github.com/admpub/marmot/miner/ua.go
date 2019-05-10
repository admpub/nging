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
	"io/ioutil"
	"math/rand"
	"net/http"
	"path/filepath"
	"runtime"
	"strings"
)

// UserAgents Global User-Agent provide
var UserAgents = map[int]string{}
var UserAgentf http.FileSystem

// UserAgentInit User-Agent init
func UserAgentInit() {
	UserAgents = map[int]string{
		0: "Mozilla/5.0 (Macintosh; U; PPC Mac OS X; de-de) AppleWebKit/125.5.5 (KHTML, like Gecko) Safari/125.12",
		1: "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:46.0) Gecko/20100101 Firefox/46.0",
		2: "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/54.0.2840.71 Safari/537.36",
		3: "Opera/9.60 (Macintosh; Intel Mac OS X; U; en) Presto/2.1.1",
	}
	_, filename, _, _ := runtime.Caller(1)
	// this *.txt maybe not found if you exec binary, so we just fill several ua
	if UserAgentf == nil {
		UserAgentf = http.Dir(filepath.Join(filepath.Dir(filename), `config`))
	}
	f, err := UserAgentf.Open(`ua.txt`)
	if err != nil {
		return
	}
	defer f.Close()
	temp, err := ioutil.ReadAll(f)
	if err != nil {
		return
	}
	uas := strings.Split(string(temp), "\n")
	for i, ua := range uas {
		UserAgents[i] = strings.TrimSpace(strings.Replace(ua, "\r", "", -1))
	}
}

// RandomUserAgent Reback random User-Agent
func RandomUserAgent() string {
	length := len(UserAgents)
	if length == 0 {
		return ""
	}
	return UserAgents[rand.Intn(length-1)]
}
