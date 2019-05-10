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

	"github.com/admpub/nging/application/dbschema"
	"github.com/admpub/nging/application/library/collector"
	"github.com/admpub/nging/application/library/collector/sender"
	"github.com/admpub/nging/application/library/notice"
	"github.com/webx-top/echo"
)

var ErrForcedExit = errors.New(`Forced exit`)

// Rules 完整规则
type Rules struct {
	*Rule            //主页面规则
	Extra    []*Rule //扩展页面规则
	exportFn func(pageID uint, lastResult *Recv, collected echo.Store, noticeSender sender.Notice) error
	isExited func() bool
}

func NewRules() *Rules {
	return &Rules{
		Rule: &Rule{
			CollectorPage: &dbschema.CollectorPage{},
			RuleList:      []*dbschema.CollectorRule{},
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

func (c *Rules) Collect(debug bool, noticeSender sender.Notice, progress *notice.Progress) ([]Result, error) {
	var fetch Fether
	timeout := int(c.Rule.CollectorPage.Timeout)
	engine := c.Rule.CollectorPage.Browser
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
			`proxy`:   c.Rule.CollectorPage.Proxy,
			`delay`:   c.Rule.CollectorPage.Waits,
		}); err != nil {
			return nil, err
		}
		collector.Services.Store(engine, browser)
	}
	browseData := make(echo.Store)
	fetch = func(pageURL string, charset string) ([]byte, bool, error) {
		browseData.Set(`charset`, charset)
		body, err := browser.Do(pageURL, browseData)
		return body, browser.Transcoded(), err
	}
	c.Rule.debug = debug
	c.Rule.exportFn = c.exportFn
	c.Rule.isExited = c.isExited
	// 	err = browser.Close()
	index := -1 //子页面层级计数，用来遍历c.Extra中的元素，-1表示入口页面
	//入口页面
	c.Rule.result = &Recv{
		index: -1,
		rule:  c.Rule,
		title: ``,
		url:   ``,
	}
	if noticeSender == nil {
		noticeSender = sender.Default
	}
	return c.Rule.Collect(
		uint64(c.CollectorPage.ParentId),
		fetch,
		index,
		c.Extra,
		noticeSender,
		progress,
	)
}
