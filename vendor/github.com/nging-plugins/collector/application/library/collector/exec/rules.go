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

package exec

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/webx-top/echo"

	"github.com/admpub/nging/v5/application/library/notice"

	"github.com/nging-plugins/collector/application/dbschema"
	"github.com/nging-plugins/collector/application/library/collector"
	"github.com/nging-plugins/collector/application/library/collector/sender"
)

var ErrForcedExit = errors.New(`Forced exit`)

// Rules 完整规则
type Rules struct {
	*Rule            //主页面配置以及采集规则列表
	Extra    []*Rule //扩展页面配置以及采集规则列表
	exportFn func(pageID uint, lastResult *Recv, collected echo.Store, noticeSender sender.Notice) error
	isExited func() bool
}

func NewRules() *Rules {
	return &Rules{
		Rule: &Rule{
			NgingCollectorPage: dbschema.NewNgingCollectorPage(nil),
			RuleList:           []*dbschema.NgingCollectorRule{},
		},
		Extra: []*Rule{},
	}
}

func (c *Rules) SetExportFn(exportFn func(pageID uint, lastResult *Recv, collected echo.Store, noticeSender sender.Notice) error) *Rules {
	c.exportFn = exportFn
	return c
}

func (c *Rules) SetExitedFn(exitedFn func() bool) *Rules {
	c.isExited = exitedFn
	return c
}

func (c *Rules) Collect(debug bool, noticeSender sender.Notice, progress *notice.Progress) (rs []Result, err error) {
	defer func() {
		if panicErr := recover(); panicErr != nil {
			err = fmt.Errorf(`%v`, panicErr)
			return
		}
	}()
	var fetch Fether
	timeout := int(c.Rule.NgingCollectorPage.Timeout)
	engine := c.Rule.NgingCollectorPage.Browser
	if len(engine) == 0 || engine == `default` {
		engine = `standard`
	}
	var browser collector.Browser
	browserService, ok := collector.Services.Load(engine)
	if ok {
		browser = browserService.(collector.Browser)
	} else {
		browser, ok = collector.Browsers[engine]
		if !ok {
			return nil, fmt.Errorf(`Unsupported: %s`, engine)
		}
		if err := browser.Start(echo.Store{
			`timeout`: timeout,
			`proxy`:   c.Rule.NgingCollectorPage.Proxy,
			`delay`:   c.Rule.NgingCollectorPage.Waits,
		}); err != nil {
			return nil, err
		}
		collector.Services.Store(engine, browser)
	}
	browseData := make(echo.Store)
	if len(c.Rule.Cookie) > 0 {
		rows := strings.Split(c.Rule.Cookie, "\n")
		items := make([]string, 0, len(rows))
		for _, str := range rows {
			str = strings.TrimSpace(str)
			if len(str) == 0 {
				continue
			}
			if !strings.HasSuffix(str, `;`) {
				str += `;`
			}
			items = append(items, str)
		}

		header := http.Header{}
		header.Add("Cookie", strings.Join(items, ` `))
		request := http.Request{Header: header}
		cookies := request.Cookies()
		browseData.Set(`cookie`, cookies)
	}
	if len(c.Rule.Header) > 0 {
		headers := map[string]string{}
		for _, str := range strings.Split(c.Rule.Header, "\n") {
			str = strings.TrimSpace(str)
			if len(str) == 0 {
				continue
			}
			parts := strings.SplitN(str, `:`, 2)
			if len(parts) != 2 {
				continue
			}
			parts[0] = strings.TrimSpace(parts[0])
			if len(parts[0]) == 0 {
				continue
			}
			parts[1] = strings.TrimSpace(parts[1])
			headers[parts[0]] = parts[1]
		}
		browseData.Set(`header`, headers)
	}
	fetch = func(pageURL string, charset string) ([]byte, bool, error) {
		browseData.Set(`charset`, charset)
		body, err := browser.Do(pageURL, browseData)
		return body, browser.Transcoded(), err
	}
	c.Rule.debug = debug
	c.Rule.exportFn = c.exportFn
	c.Rule.isExited = c.isExited
	// 	err = browser.Close()
	//入口页面
	topRecv := &Recv{
		Index:      -1,
		LevelIndex: -1, //子页面层级计数，用来遍历c.Extra中的元素(作为Extra切片下标)，-1表示入口页面
		//rule:       c.Rule,
		Title: ``,
		URL:   ``,
	}
	if noticeSender == nil {
		noticeSender = sender.Default
	}
	return c.Rule.Collect(
		uint64(c.NgingCollectorPage.ParentId),
		``,
		topRecv,
		fetch,
		c.Extra,
		noticeSender,
		progress,
	)
}
