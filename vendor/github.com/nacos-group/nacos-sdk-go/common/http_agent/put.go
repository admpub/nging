/*
 * Copyright 1999-2020 Alibaba Group Holding Ltd.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package http_agent

import (
	"log"
	"net/http"
	"strings"
	"time"
)

func put(path string, header http.Header, timeoutMs uint64, params map[string]string) (response *http.Response, err error) {
	client := http.Client{}
	client.Timeout = time.Millisecond * time.Duration(timeoutMs)
	var body string
	for key, value := range params {
		if len(value) > 0 {
			body += key + "=" + value + "&"
		}
	}
	if strings.HasSuffix(body, "&") {
		body = body[:len(body)-1]
	}
	request, errNew := http.NewRequest(http.MethodPut, path, strings.NewReader(body))
	if errNew != nil {
		err = errNew
		return
	}
	request.Header = header
	resp, errDo := client.Do(request)
	if errDo != nil {
		log.Println(errDo)
		err = errDo
	} else {
		response = resp
	}
	return
}
