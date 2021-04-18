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

import "net/http"

//go:generate mockgen -destination ../../mock/mock_http_agent_interface.go -package mock -source=./http_agent_interface.go

type IHttpAgent interface {
	Get(path string, header http.Header, timeoutMs uint64, params map[string]string) (response *http.Response, err error)
	Post(path string, header http.Header, timeoutMs uint64, params map[string]string) (response *http.Response, err error)
	Delete(path string, header http.Header, timeoutMs uint64, params map[string]string) (response *http.Response, err error)
	Put(path string, header http.Header, timeoutMs uint64, params map[string]string) (response *http.Response, err error)
	RequestOnlyResult(method string, path string, header http.Header, timeoutMs uint64, params map[string]string) string
	Request(method string, path string, header http.Header, timeoutMs uint64, params map[string]string) (response *http.Response, err error)
}
