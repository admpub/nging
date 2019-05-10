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

package standard

import (
	"net/url"
	"time"

	"github.com/admpub/marmot/miner"
	"github.com/admpub/nging/application/library/collector"
	"github.com/webx-top/echo"
)

func init() {
	collector.Browsers[`standard`] = New()
}

func New() *Standard {
	return &Standard{
		Worker: miner.DefaultWorker,
		Base:   &collector.Base{},
	}
}

type Standard struct {
	Worker *miner.Worker
	*collector.Base
}

func (s *Standard) Start(opt echo.Store) (err error) {
	if err = s.Base.Start(opt); err != nil {
		return
	}
	var proxyCfg interface{}
	if len(s.Proxy) > 0 {
		proxyCfg = s.Proxy
	}
	if s.Timeout > 0 {
		s.Worker, err = miner.New(proxyCfg, time.Duration(s.Timeout)*time.Second)
	} else {
		s.Worker, err = miner.New(proxyCfg)
	}
	if err != nil {
		return
	}
	s.Worker.SetUserAgent(miner.RandomUserAgent())
	return nil
}

func (s *Standard) Close() error {
	s.Worker.ClearAll()
	return nil
}

func (s *Standard) Name() string {
	return `standard`
}

func (s *Standard) Description() string {
	return ``
}

func (s *Standard) Do(pageURL string, data echo.Store) ([]byte, error) {
	charset := data.String(`charset`)
	method := data.String(`method`, miner.GET)
	s.Worker.
		SetURL(pageURL).
		SetMethod(method).
		SetDetectCharset(true).
		SetResponseCharset(charset)

	if formData, ok := data.Get(`formData`).(url.Values); ok {
		s.Worker.SetForm(formData)
	}

	sleepSeconds := s.SleepSeconds()
	if sleepSeconds > 0 {
		s.Worker.SetWaitTime(sleepSeconds)
	}
	return s.Worker.Go()
}
