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
	"bytes"
	"io"
	"net/http"
	"strconv"
)

type fakeHttpResponseBody struct {
	body io.ReadSeeker
}

func (body *fakeHttpResponseBody) Read(p []byte) (n int, err error) {
	n, err = body.body.Read(p)
	if err == io.EOF {
		body.body.Seek(0, 0)
	}
	return n, err
}

func (body *fakeHttpResponseBody) Close() error {
	return nil
}

func FakeHttpResponse(status int, body string) (resp *http.Response) {
	return &http.Response{
		Status:     strconv.Itoa(status),
		StatusCode: status,
		Body:       &fakeHttpResponseBody{bytes.NewReader([]byte(body))},
		Header:     http.Header{},
	}
}
