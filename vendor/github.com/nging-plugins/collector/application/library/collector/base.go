/*
   Nging is a toolbox for webmasters
   Copyright (C) 2018-present  Wenhui Shen <swh@admpub.com>

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU Affero General Public License as published
   by the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Affero General Public License for more details.

   You should have received a copy of the GNU Affero General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package collector

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/webx-top/com"
	"github.com/webx-top/echo"
)

type Base struct {
	Proxy            string
	Timeout          int
	RandomDelay      string
	minWait, maxWait int
	Cookies          []*http.Cookie
	CookieString     string
	Header           map[string]string
	UserAgent        string
}

func (s *Base) Start(opt echo.Store) (err error) {
	s.Proxy = opt.String(`proxy`)
	s.Timeout = opt.Int(`timeout`)
	s.RandomDelay = opt.String(`delay`)
	s.UserAgent = opt.String(`userAgent`)
	if len(s.RandomDelay) > 0 {
		waits := strings.SplitN(s.RandomDelay, `-`, 2)
		switch len(waits) {
		case 2:
			s.maxWait, _ = strconv.Atoi(waits[1])
			fallthrough
		case 1:
			s.minWait, _ = strconv.Atoi(waits[0])
		}
	}
	if hd, ok := opt.Get(`header`).(map[string]string); ok {
		s.Header = hd
	}
	switch cookieData := opt.Get(`cookie`).(type) {
	case string:
		if len(cookieData) > 0 {
			s.CookieString = cookieData
		}
	case *http.Cookie:
		if cookieData.Valid() == nil {
			s.Cookies = []*http.Cookie{cookieData}
		}
	case []*http.Cookie:
		s.Cookies = cookieData
	}
	return nil
}

func (s *Base) Close() error {
	return nil
}

func (s *Base) Name() string {
	return `undefined`
}

func (s *Base) Description() string {
	return ``
}

func (s *Base) Transcoded() bool {
	return true
}

func (s *Base) Do(pageURL string, data echo.Store) ([]byte, error) {
	return nil, nil
}

func (s *Base) SleepSeconds() int {
	if s.minWait > 0 || s.maxWait > 0 {
		return com.RandRangeInt(s.minWait, s.maxWait)
	}
	return 0
}

func (s *Base) Sleep() Browser {
	if s.minWait > 0 || s.maxWait > 0 {
		delay := com.RandRangeInt(s.minWait, s.maxWait)
		if delay > 0 {
			time.Sleep(time.Duration(delay) * time.Second)
		}
	}
	return s
}
