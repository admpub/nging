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

package webdriver

import (
	"github.com/admpub/log"
	"github.com/tebeka/selenium"
	"github.com/webx-top/echo"

	"github.com/nging-plugins/collector/application/library/collector"
)

func init() {
	collector.Browsers[`webdriver`] = New()
}

func New() *WebDriver {
	return &WebDriver{
		Base: &collector.Base{},
	}
}

type WebDriver struct {
	client selenium.WebDriver
	server *selenium.Service
	*collector.Base
}

func (s *WebDriver) Start(opt echo.Store) (err error) {
	if err = s.Base.Start(opt); err != nil {
		return
	}
	_, err = StartServer(s.Base)
	if err != nil {
		log.Error(err.Error())
	}
	s.client, err = InitClient(s.Base)
	return
}

func (s *WebDriver) Close() error {
	err := s.client.Quit()
	return err
}

func (s *WebDriver) Name() string {
	return `webdriver`
}

func (s *WebDriver) Description() string {
	return ``
}

func (s *WebDriver) Do(pageURL string, data echo.Store) ([]byte, error) {
	p := Page{client: s.client}
	err := p.client.Get(pageURL)
	if err != nil {
		return nil, err
	}
	var html string
	html, err = p.client.PageSource()
	if err != nil {
		return nil, err
	}
	s.Sleep()
	return []byte(html), err
}
