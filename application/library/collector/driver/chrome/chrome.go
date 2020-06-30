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

package chrome

import (
	"context"
	"time"

	"github.com/admpub/cr"
	"github.com/admpub/nging/application/library/collector"
	"github.com/chromedp/chromedp"
	"github.com/webx-top/echo"
)

func init() {
	_ = chromedp.ByQuery
	collector.Browsers[`chromedp`] = New()
}

func New() *Chrome {
	return &Chrome{
		Base: &collector.Base{},
	}
}

type Chrome struct {
	Browser *cr.Browser
	*collector.Base
}

func (s *Chrome) Start(opt echo.Store) (err error) {
	if err = s.Base.Start(opt); err != nil {
		return
	}
	chromePath := opt.String(`chromePath`)
	options := []chromedp.ExecAllocatorOption{
		chromedp.WindowSize(800, 600),
	}
	if len(chromePath) > 0 {
		options = append(options, chromedp.ExecPath(chromePath))
	}
	if len(s.Proxy) > 0 {
		options = append(options, chromedp.ProxyServer(s.Proxy))
	}
	s.Browser, err = cr.New(context.Background(), options...)
	if err != nil {
		return
	}
	if s.Timeout > 0 {
		s.Browser.SetTimeout(time.Duration(s.Timeout) * time.Second)
	}
	return nil
}

func (s *Chrome) Close() error {
	return s.Browser.Close()
}

func (s *Chrome) Name() string {
	return `chromedp`
}

func (s *Chrome) Description() string {
	return ``
}

func (s *Chrome) Do(pageURL string, data echo.Store) ([]byte, error) {
	if err := s.Browser.Navigate(pageURL); err != nil {
		return nil, err
	}
	/*
		if err := s.Browser.RunAction(chromedp.WaitReady("body", chromedp.ByQuery)); err != nil {
			return nil, err
		}
	// */
	html, err := s.Browser.GetSource()
	if err != nil {
		return nil, err
	}
	s.Sleep()
	return []byte(html), nil
}
