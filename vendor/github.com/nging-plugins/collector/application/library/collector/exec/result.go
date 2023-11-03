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
	"encoding/gob"
	"encoding/json"
	"net/http"
	"path/filepath"
	"time"

	"github.com/admpub/color"
	"github.com/admpub/marmot/miner"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"

	"github.com/admpub/nging/v5/application/library/common"

	"github.com/nging-plugins/collector/application/dbschema"
)

func init() {
	gob.Register([]*dbschema.NgingCollectorRule{})
	_ = color.Red
	miner.UserAgentf = http.Dir(filepath.Join(echo.Wd(), `config`))
}

var (
	//{(1-2,2-3:2)} 分为两个部分，用":"分隔，":"左边部分定义页码范围，右边定义步进
	pagingFlagLeft  = "{("
	pagingFlagRight = ")}"
	emptyRecv       = &Recv{IsEmpty: true}
)

type (
	Fether func(pageURL string, charset string) (body []byte, transcoded bool, err error)
	Result struct { //收集测试结果
		Title     string
		URL       string
		Result    interface{}
		Type      string //map/slice
		StartTime time.Time
		EndTime   time.Time
		Elapsed   time.Duration
	}
	Recv struct { //接收结果
		Index      int
		IsEmpty    bool        //是否为空结果
		LevelIndex int         //层级索引
		URLIndex   int         //网址列表索引
		Result     interface{} //采集结果数据
		//rule       *Rule       //页面规则
		Title  string //页面标题
		URL    string //网址
		parent *Recv
	}
)

func (c *Result) IsSlice() bool {
	return c.Type == `array`
}

func (c *Result) IsMap() bool {
	return c.Type == `map`
}

func (c *Result) ElapsedString(lang string) string {
	duration := com.ParseDuration(c.Elapsed, lang)
	return duration.String()
}

// ================== Recv ====================

func (c *Recv) String() string {
	b, _ := json.MarshalIndent(c, ``, `  `)
	return string(b)
}

func (c *Recv) ResultByIndex(index int) interface{} {
	if res, ok := c.Result.([]interface{}); ok {
		if len(res) > index {
			return res[index]
		}
	}
	return echo.H{}
}

// func (c *Recv) Rule() *Rule {
// 	return c.rule
// }

func (c *Recv) ParentItem() interface{} {
	return c.Parent().ResultByIndex(c.URLIndex)
}

func (c *Recv) ParentResult() interface{} {
	return c.Parent().Result
}

func (c *Recv) ParentsResult(lasts ...int) interface{} {
	return c.Parents(lasts...).Result
}

func (c *Recv) ParentsItem(lasts ...int) interface{} {
	return c.Parents(lasts...).ResultByIndex(c.URLIndex)
}

func (c *Recv) Parent() *Recv {
	if c.parent != nil {
		return c.parent
	}
	return emptyRecv
}

func (c *Recv) Parents(lasts ...int) *Recv {
	var last int
	if len(lasts) > 0 {
		last = lasts[0]
	}
	if last < 1 {
		return c
	}
	r := c
	for i := 1; i <= last; i++ {
		if r.parent == nil {
			return emptyRecv
		}
		r = r.parent
	}
	if r == nil {
		return emptyRecv
	}
	return r
}

func (_ *Recv) UniqueID() (string, error) {
	return common.UniqueID()
}
