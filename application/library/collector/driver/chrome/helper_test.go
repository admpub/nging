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
	"errors"
	"fmt"
	"io/ioutil"
	"testing"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/stretchr/testify/assert"
)

// 获取服务列表
func getServiceList(res *string, pageURL string) chromedp.Tasks {
	var buf []byte
	return chromedp.Tasks{
		chromedp.ActionFunc(func(ctx context.Context, h cdp.Executor) error {
			expr := cdp.TimeSinceEpoch(time.Now().Add(180 * 24 * time.Hour))
			success, err := network.SetCookie("cookiename", "cookievalue").
				WithExpires(&expr).
				WithDomain("localhost").
				WithHTTPOnly(true).
				Do(ctx, h)
			if err != nil {
				return err
			}
			if !success {
				return errors.New("could not set cookie")
			}
			return nil
		}),
		// 访问服务列表
		chromedp.Navigate(`http://www.admpub.com/`),
		chromedp.WaitReady("content", chromedp.ByID),
		chromedp.Click("#post-246 > h2 > a", chromedp.ByQuery),
		chromedp.WaitReady("comment", chromedp.ByID),
		chromedp.SendKeys("comment", "www.admpub.com", chromedp.ByID),
		chromedp.SendKeys("#commentform input[name=imgcode]", "12345", chromedp.ByQuery),
		chromedp.Screenshot(`#commentform`, &buf, chromedp.ByQuery),
		chromedp.ActionFunc(func(context.Context, cdp.Executor) error {
			return ioutil.WriteFile("testimonials-submit-comment-before.png", buf, 0644)
		}),
		// chromedp.Submit("#commentform", chromedp.ByQuery),
		// chromedp.Sleep(2 * time.Second),
		// chromedp.Screenshot(`content`, &buf, chromedp.ByID),
		// chromedp.ActionFunc(func(context.Context, cdp.Executor) error {
		// 	return ioutil.WriteFile("testimonials-submit-comment-after.png", buf, 0644)
		// }),
		// 访问服务列表
		chromedp.Navigate(pageURL),
		// 等待直到body加载完毕
		chromedp.WaitReady("content", chromedp.ByID),
		// 选择显示可用服务
		chromedp.Click("#content .align_right", chromedp.ByQuery),
		// 等待列表渲染
		chromedp.Sleep(2 * time.Second),
		// 获取文章标题
		chromedp.OuterHTML("#content .post h2", res, chromedp.ByQuery),
		chromedp.Screenshot(`#content`, &buf, chromedp.ByID),
		chromedp.ActionFunc(func(context.Context, cdp.Executor) error {
			return ioutil.WriteFile("testimonials.png", buf, 0644)
		}),
	}
}

func TestServiceList(t *testing.T) {
	var html string
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cdp, err := NewHeadless(ctx, "http://www.admpub.com")
	//cdp, err := NewBrowser(ctx)
	if err != nil {
		panic(err)
	}
	// cdp是chromedp实例
	// ctx是创建cdp时使用的context.Context
	err = cdp.Run(ctx, getServiceList(&html, "http://www.admpub.com/blog/post-255.html"))
	if err != nil {
		panic(err)
	}

	assert.Equal(t, "<h2>安装 Go 第三方包 go-sqlite3</h2>", html)
	// 成功取得HTML内容进行后续处理
	fmt.Println(html)
}
