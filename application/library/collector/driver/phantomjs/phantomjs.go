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

package phantomjs

import (
	"github.com/admpub/nging/application/library/collector"
	pjs "github.com/admpub/phantomjs"
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
	Process *pjs.Process
	*collector.Base
}

func (s *PhantomJS) Start(opt echo.Store) (err error) {
	if err = s.Base.Start(opt); err != nil {
		return
	}
	if len(s.Proxy) > 0 {
		s.Process, err = pjs.NewWithProxy(s.Proxy)
		if err != nil {
			return
		}
	} else {
		s.Process = pjs.NewProcess()
	}
	err = s.Process.Open()
	return
}

func (s *PhantomJS) Close() error {
	return s.Process.Close()
}

func (s *PhantomJS) Name() string {
	return `phantomjs`
}

func (s *PhantomJS) Description() string {
	return ``
}

func (s *PhantomJS) Do(pageURL string, data echo.Store) ([]byte, error) {
	page, err := s.Process.CreateWebPage()
	if err != nil {
		return nil, err
	}
	defer page.Close()
	// Open a URL.
	err = page.Open(pageURL)
	if err != nil {
		return nil, err
	}
	var html string
	html, err = page.Content()
	s.Sleep()
	return []byte(html), nil
}
