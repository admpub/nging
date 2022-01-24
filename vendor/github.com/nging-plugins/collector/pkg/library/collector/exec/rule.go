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
	"bytes"
	"fmt"
	"html/template"
	"io"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	// import collector browser driver

	download "github.com/admpub/go-download/v2"
	"github.com/admpub/gopiper"
	"github.com/admpub/log"
	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/engine"
	"github.com/webx-top/echo/middleware/tplfunc"

	"github.com/admpub/nging/v4/application/library/charset"
	"github.com/admpub/nging/v4/application/library/common"
	"github.com/admpub/nging/v4/application/library/notice"

	"github.com/nging-plugins/collector/pkg/dbschema"
	_ "github.com/nging-plugins/collector/pkg/library/collector/driver/chrome"
	_ "github.com/nging-plugins/collector/pkg/library/collector/driver/webdriver"
	"github.com/nging-plugins/collector/pkg/library/collector/sender"

	//_ "github.com/nging-plugins/collector/pkg/library/collector/driver/phantomjs" //高CPU占用
	//_ "github.com/nging-plugins/collector/pkg/library/collector/driver/phantomjsfetcher" //高CPU占用
	_ "github.com/nging-plugins/collector/pkg/library/collector/driver/standard"
)

var RegexpTitle = regexp.MustCompile(`(?i)<title[\s]*>([^<]+)</title[\s]*>`)

// Rule 页面规则
type Rule struct {
	*dbschema.NgingCollectorPage
	RuleList []*dbschema.NgingCollectorRule
	debug    bool
	exportFn func(pageID uint, result *Recv, collected echo.Store, noticeSender sender.Notice) error
	isExited func() bool
	result   *Recv
}

func (c *Rule) IsExited() bool {
	if c.isExited == nil {
		return false
	}
	return c.isExited()
}

func (c *Rule) Result() *Recv {
	return c.result
}

func (c *Rule) ParseTmplContent(tmplContent string) (string, error) {
	if len(tmplContent) == 0 || strings.Index(tmplContent, `{{`) < 0 {
		return tmplContent, nil
	}
	md5 := com.Md5(tmplContent)
	t := template.New(md5).Funcs(tplfunc.TplFuncMap)
	_, err := t.Parse(tmplContent)
	if err != nil {
		return ``, err
	}
	buf := bytes.NewBuffer(nil)
	err = t.Execute(buf, c.result)
	if err != nil {
		return ``, err
	}
	return strings.TrimSpace(buf.String()), err
}

func (c *Rule) Collect(parentID uint64,
	fetch Fether, index int, extra []*Rule,
	noticeSender sender.Notice,
	progress *notice.Progress) ([]Result, error) {
	if c.IsExited() {
		return nil, ErrForcedExit
	}
	enterURL, err := c.ParseTmplContent(c.NgingCollectorPage.EnterUrl)
	if err != nil {
		return nil, err
	}
	if sendErr := noticeSender(`<collector> URL Template Result: `+enterURL, 1, progress); sendErr != nil {
		return nil, sendErr
	}
	if len(enterURL) == 0 {
		return nil, nil
	}
	index++
	var (
		urlList []string
		result  []Result
	)
	for _, v := range strings.Split(enterURL, "\n") {
		if c.IsExited() {
			return result, ErrForcedExit
		}
		v = strings.TrimSpace(v)
		if len(v) == 0 {
			continue
		}
		position := strings.Index(v, pagingFlagLeft)
		if position == -1 {
			urlList = append(urlList, v)
			continue
		}
		prefix := v[0:position]
		remVal := v[position+len(pagingFlagLeft):]
		position = strings.Index(remVal, pagingFlagRight)
		if position == -1 {
			urlList = append(urlList, v)
			continue
		}
		expr := remVal[0:position]
		suffix := remVal[position+len(pagingFlagRight):]
		com.SeekRangeNumbers(expr, func(page int) bool {
			urlList = append(urlList, prefix+strconv.Itoa(page)+suffix)
			return true
		})
	}
	if progress != nil && index == 1 {
		progress.Finish = 0
		progress.Total = int64(len(urlList))
	}
	// collection 的类型有两种可能：[]interface{} / map[string]interface{}
	var collection interface{}
	historyMdl := dbschema.NewNgingCollectorHistory(c.Context())
	for pageKey, pageURL := range urlList {
		if c.IsExited() {
			return result, ErrForcedExit
		}
		urlMD5 := com.Md5(pageURL)
		var (
			historyID  uint64
			content    []byte
			encoded    []byte
			transcoded bool
			ruleMd5    string
			contentMd5 string
		)
		c.result.url = pageURL
		c.result.result = collection
		if !c.debug { //非测试模式才保存到数据库
			err = historyMdl.Get(nil, `url_md5`, urlMD5)
			if err != nil {
				if err != db.ErrNoMoreRows {
					return nil, err
				}
				//不存在记录
			} else if historyMdl.Id > 0 {
				if c.NgingCollectorPage.DuplicateRule == `url` {
					continue
				}
				historyID = historyMdl.Id
			}
			encoded, err = com.GobEncode(c.RuleList)
			if err != nil {
				return nil, err
			}
			ruleMd5 = com.ByteMd5(encoded)
			if historyID > 0 && historyMdl.RuleMd5 == ruleMd5 && c.NgingCollectorPage.DuplicateRule == `rule` { //规则没有更改过的情况下，如果已经采集过则跳过
				continue
			}
		}
		if sendErr := noticeSender(`<collector> Fetching URL: `+pageURL, 1, progress); sendErr != nil {
			return nil, sendErr
		}
		startTime := time.Now()
		content, transcoded, err = fetch(pageURL, c.NgingCollectorPage.Charset)
		if err != nil {
			if err == io.EOF {
				log.Error(err)
				continue
			}
			return nil, err
		}
		if !c.debug { //非测试模式才保存到数据库
			contentMd5 = com.ByteMd5(content)
			if historyID > 0 && historyMdl.Content == contentMd5 && c.NgingCollectorPage.DuplicateRule == `content` { //规则没有更改过的情况下，如果已经采集过则跳过
				continue
			}
			historyMdl.Reset()
			historyMdl.RuleMd5 = ruleMd5
			historyMdl.Created = uint(time.Now().Unix())
			historyMdl.Url = pageURL
			historyMdl.UrlMd5 = urlMD5
			historyMdl.PageId = c.NgingCollectorPage.Id
			historyMdl.PageParentId = c.NgingCollectorPage.ParentId
			historyMdl.PageRootId = c.NgingCollectorPage.RootId
			historyMdl.ParentId = parentID
			historyMdl.HasChild = c.NgingCollectorPage.HasChild
			historyMdl.Content = contentMd5
		}
		if !transcoded {
			if len(c.NgingCollectorPage.Charset) < 1 {
				c.NgingCollectorPage.Charset = `utf-8`
			}
			// 字符集转码
			if strings.ToLower(c.NgingCollectorPage.Charset) != `utf-8` {
				content, err = charset.Convert(c.NgingCollectorPage.Charset, `utf-8`, content)
				if err != nil {
					return nil, err
				}
			}
		}
		subItems := []gopiper.PipeItem{}
		for _, rule := range c.RuleList {
			subItem := gopiper.PipeItem{
				Name:     rule.Name,
				Type:     rule.Type,
				Selector: rule.Rule,
				Filter:   rule.Filter,
			}
			subItem.Selector, err = c.ParseTmplContent(subItem.Selector)
			if err != nil {
				return nil, err
			}
			subItem.Filter, err = c.ParseTmplContent(subItem.Filter)
			if err != nil {
				return nil, err
			}
			subItems = append(subItems, subItem)
		}
		var pipe gopiper.PipeItem
		if c.NgingCollectorPage.Type == `list` {
			child := gopiper.PipeItem{
				Type:    `map`,
				SubItem: subItems,
			}
			pipe = gopiper.PipeItem{
				Type:     `array`,
				Selector: c.NgingCollectorPage.ScopeRule,
				SubItem:  []gopiper.PipeItem{child},
			}
			pipe.Filter, err = c.ParseTmplContent(pipe.Filter)
			if err != nil {
				return nil, err
			}
			pipe.Selector, err = c.ParseTmplContent(pipe.Selector)
			if err != nil {
				return nil, err
			}
		} else {
			pipe = gopiper.PipeItem{
				Type:    `map`,
				SubItem: subItems,
			}
		}
		pipe.SetFetcher(func(pURL string) (body []byte, err error) {
			pURL = com.AbsURL(pageURL, pURL)
			body, transcoded, err := fetch(pURL, c.NgingCollectorPage.Charset)
			if !transcoded {
				if len(c.NgingCollectorPage.Charset) < 1 {
					c.NgingCollectorPage.Charset = `utf-8`
				}
				// 字符集转码
				if strings.ToLower(c.NgingCollectorPage.Charset) != `utf-8` {
					body, err = charset.Convert(c.NgingCollectorPage.Charset, `utf-8`, body)
					if err != nil {
						return
					}
				}
			}
			return
		})
		pipe.SetStorer(func(fileURL string, savePath string, fetched bool) (newPath string, err error) {
			newPath = savePath
			newPath, err = c.ParseTmplContent(newPath)
			if err != nil {
				return newPath, err
			}
			if strings.HasSuffix(newPath, `/`) || strings.HasSuffix(newPath, `\`) {
				newPath += path.Base(fileURL)
			}
			saveTo := filepath.Join(echo.Wd(), newPath)
			saveDir := filepath.Dir(saveTo)
			err = com.MkdirAll(saveDir, os.ModePerm)
			if err != nil {
				return newPath, err
			}
			if sendErr := noticeSender(`download file: `+fileURL+` => `+saveTo, 1, progress); sendErr != nil {
				return newPath, sendErr
			}
			if fetched {
				err = os.WriteFile(saveTo, []byte(fileURL), os.ModePerm)
				return
			}
			fileURL = com.AbsURL(pageURL, fileURL)
			_, err = download.Download(fileURL, saveTo, nil)
			return
		})
		collection, err = pipe.PipeBytes(content, c.NgingCollectorPage.ContentType)
		if err != nil {
			if err != gopiper.ErrInvalidContent { //跳过无效内容
				if sendErr := noticeSender(err.Error(), 0, progress); sendErr != nil {
					return result, sendErr
				}
			}
			if c.debug {
				return result, err
			}
			err = nil
			continue
		}
		if collection == nil {
			continue
		}

		// 自动获取页面标题
		var pageTitle string
		switch c.NgingCollectorPage.ContentType {
		case `html`:
			find := RegexpTitle.FindAllStringSubmatch(engine.Bytes2str(content), 1)
			if len(find) > 0 && len(find[0]) > 1 {
				pageTitle = strings.TrimSpace(find[0][1])
			}
			if len(pageTitle) == 0 {
				if mp, ok := collection.(map[string]interface{}); ok {
					if tt, ok := mp[`title`]; ok {
						pageTitle = com.Str(tt)
					}
				}
			}
		case `json`:
			fallthrough
		case `text`:
			if mp, ok := collection.(map[string]interface{}); ok {
				if tt, ok := mp[`title`]; ok {
					pageTitle = com.Str(tt)
				}
			}
		}
		c.result.title = pageTitle
		c.result.result = collection
		// 记录第一个网址数据
		if pageKey == 0 {
			endTime := time.Now()
			result = append(result, Result{
				Title:     pageTitle,
				URL:       pageURL,
				Result:    collection,
				Type:      pipe.Type,
				StartTime: startTime,
				EndTime:   endTime,
				Elapsed:   endTime.Sub(startTime),
			})
		}
		encoded, err = com.JSONEncode(collection)
		if err != nil {
			return result, err
		}
		//historyMdl.Data = string(encoded)
		err = common.WriteCache(`colloctor`, urlMD5+`.json`, encoded)
		if err != nil {
			return result, err
		}

		if !c.debug { //非测试模式才保存到数据库
			historyMdl.Title = pageTitle
			if historyID > 0 {
				err = historyMdl.Update(nil, `id`, historyID)
			} else {
				_, err = historyMdl.Insert()
				if err == nil {
					historyID = historyMdl.Id
				}
			}
			if err != nil {
				return result, err
			}
		}
		//msgbox.Table(ctx.T(`Result`), collection, 200)
		//color.Red(`(%d) `+pageURL, index)
		if len(extra) > index {
			lastRuleForm := c
			for _, pageRuleForm := range extra[index:] {
				if c.IsExited() {
					return result, ErrForcedExit
				}
				pageRuleForm.result = &Recv{
					index:  pageKey,
					rule:   pageRuleForm,
					parent: lastRuleForm.result,
				}
				pageRuleForm.debug = c.debug
				pageRuleForm.exportFn = c.exportFn
				pageRuleForm.isExited = c.isExited
				if len(pageRuleForm.NgingCollectorPage.Charset) < 1 {
					pageRuleForm.NgingCollectorPage.Charset = c.NgingCollectorPage.Charset
				}
				var extraResult []Result
				extraResult, err = pageRuleForm.Collect(
					historyID,
					fetch,
					index,
					extra,
					noticeSender,
					progress,
				)
				if err != nil {
					return result, err
				}
				result = append(result, extraResult...)
				lastRuleForm = pageRuleForm
			}
		}

		if c.NgingCollectorPage.HasChild == `N` { //这是最底层
			if c.exportFn != nil {
				switch collected := collection.(type) {
				case map[string]interface{}:
					c.exportFn(c.NgingCollectorPage.Id, c.result, collected, noticeSender)
				case []interface{}:
					for _, item := range collected {
						collectedMap, ok := item.(map[string]interface{})
						if !ok {
							if sendErr := noticeSender(fmt.Sprintf(`Unsupport export type: %T`, item), 0); sendErr != nil {
								return result, sendErr
							}
							continue
						}
						c.exportFn(c.NgingCollectorPage.Id, c.result, collectedMap, noticeSender)
					}
				default:
					if sendErr := noticeSender(fmt.Sprintf(`Unsupport export type: %T`, collection), 0); sendErr != nil {
						return result, sendErr
					}
				}
			}
		}
		if progress != nil && index == 1 {
			progress.Finish = int64(pageKey + 1)
		}
		if c.debug {
			break
		}
	} //end-for:range urlList
	return result, err
}
