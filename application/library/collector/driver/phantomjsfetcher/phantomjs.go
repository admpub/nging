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

package phantomjsfetcher

import (
	phantomjs "github.com/admpub/go-phantomjs-fetcher"
	"github.com/admpub/nging/application/library/collector"
	"github.com/webx-top/echo"
)

func init() {
	collector.Browsers[`phantomjs`] = New()
}

func New() *PhantomJS {
	return &PhantomJS{
		Base: &collector.Base{},
	}
}

type PhantomJS struct {
	Fetcher *phantomjs.Fetcher
	*collector.Base
}

func (s *PhantomJS) Start(opt echo.Store) (err error) {
	if err = s.Base.Start(opt); err != nil {
		return
	}
	InitServer()
	return
}

func (s *PhantomJS) Close() error {
	return CloseServer()
}

func (s *PhantomJS) Name() string {
	return `phantomjs`
}

func (s *PhantomJS) Description() string {
	return ``
}

func (s *PhantomJS) Do(pageURL string, data echo.Store) ([]byte, error) {
	jscode := data.String(`jscode`)
	headers, _ := data.Get(`headers`).(map[string]string)
	resp, err := Fetch(pageURL, jscode, headers)
	if err != nil {
		return nil, err
	}
	s.Sleep()
	return []byte(resp.Content), nil
}
