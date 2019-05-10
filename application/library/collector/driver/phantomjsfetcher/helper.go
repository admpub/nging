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
	"github.com/admpub/marmot/miner"
)

var (
	DefaultPort    = 20180
	defaultFetcher *phantomjs.Fetcher
)

func Fetch(pageURL string, jscode string, headers map[string]string) (*phantomjs.Response, error) {
	err := InitServer()
	if err != nil {
		return nil, err
	}

	fetcher := defaultFetcher

	if headers == nil {
		headers = make(map[string]string)
		if _, ok := headers["User-Agent"]; !ok {
			headers["User-Agent"] = miner.RandomUserAgent()
		}
	}
	options := *fetcher.DefaultOption
	options.Headers = headers
	//jscode := "function() {s=document.documentElement.outerHTML;document.write('<body></body>');document.body.innerText=s;}"
	jsRunAt := phantomjs.RUN_AT_DOC_END
	resp, err := fetcher.GetWithOption(pageURL, jscode, jsRunAt, &options)
	return resp, err
}

func CloseServer() error {
	if defaultFetcher == nil {
		return nil
	}
	defaultFetcher.ShutDownPhantomJSServer()
	return nil
}

func InitServer() error {
	fetcher, err := phantomjs.NewFetcher(DefaultPort, nil)
	if err != nil {
		return err
	}
	defaultFetcher = fetcher
	return nil
}
