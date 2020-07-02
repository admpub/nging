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
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/admpub/log"
	"github.com/webx-top/chardet"
	"golang.org/x/net/html/charset"
)

// Wait some secord
func Wait(waittime int) {
	if waittime <= 0 {
		return
	}
	// debug
	log.Debugf("Wait %d Second.", waittime)
	time.Sleep(time.Duration(waittime) * time.Second)
}

// CopyM Header map[string][]string ,can use to copy a http header, so that they are not effect each other
func CopyM(h http.Header) http.Header {
	if h == nil || len(h) == 0 {
		return h
	}
	h2 := make(http.Header, len(h))
	for k, vv := range h {
		vv2 := make([]string, len(vv))
		copy(vv2, vv)
		h2[k] = vv2
	}
	return h2
}

//TooShortSizes if a file size small than sizes(KB) ,it will be throw a error
func TooShortSizes(data []byte, sizes float64) error {
	if float64(len(data))/1000 < sizes {
		return fmt.Errorf("FileSize:%d bytes,%d kb < %f kb dead too sort", len(data), len(data)/1000, sizes)
	}
	return nil
}

// OutputMaps Just debug a map
func OutputMaps(info string, args map[string][]string) {
	s := "\n"
	for k, v := range args {
		s = s + fmt.Sprintf("%-25s| %-6s\n", k, strings.Join(v, "||"))
	}
	Debugf("[GoWorker] %s", s)
}

func FixCharset(body []byte, extra interface{}, detectCharset bool, defaultEncoding ...string) ([]byte, error) {
	if len(body) == 0 {
		return body, nil
	}

	if len(defaultEncoding) > 0 && len(defaultEncoding[0]) > 0 {
		return encodeBytes(body, "text/plain; charset="+defaultEncoding[0])
	}

	var contentType string
	if resp, ok := extra.(*http.Response); ok {
		contentType = resp.Header.Get(`Content-Type`)
	} else {
		contentType = fmt.Sprint(extra)
	}

	contentType = strings.ToLower(contentType)
	if strings.Contains(contentType, "image/") ||
		strings.Contains(contentType, "video/") ||
		strings.Contains(contentType, "audio/") ||
		strings.Contains(contentType, "font/") {
		// These MIME types should not have textual data.

		return body, nil
	}

	if !strings.Contains(contentType, "charset") {
		if !detectCharset {
			return body, nil
		}
		d := chardet.NewTextDetector()
		r, err := d.DetectBest(body)
		if err != nil {
			return body, err
		}
		contentType = "text/plain; charset=" + r.Charset
	}
	if strings.Contains(contentType, "utf-8") || strings.Contains(contentType, "utf8") {
		return body, nil
	}
	return encodeBytes(body, contentType)
}

func encodeBytes(b []byte, contentType string) ([]byte, error) {
	r, err := charset.NewReader(bytes.NewReader(b), contentType)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(r)
}
