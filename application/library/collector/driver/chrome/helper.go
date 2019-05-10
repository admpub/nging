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

	"github.com/admpub/log"
	"github.com/chromedp/chromedp"

	// runner用于配置headless chrome
	"github.com/chromedp/chromedp/runner"
)

// NewHeadless 创建headless chrome实例
// chromedp内部有自己的超时设置，你也可以通过ctx来设置更短的超时
func NewHeadless(ctx context.Context, starturl string, extra ...runner.CommandLineOption) (*chromedp.CDP, error) {
	// runner.Flag设置启动headless chrome时的命令行参数
	// runner.URL设置启动时打开的URL
	// Windows用户需要设置runner.Flag("disable-gpu", true)，具体信息参见文档的FAQ
	options := []runner.CommandLineOption{
		runner.Flag("headless", true),
		runner.Flag("disable-gpu", true),
		//runner.WindowSize(800, 600),
		//runner.ProxyServer("代理服务器地址"),
		runner.URL(starturl),
		//runner.ExecPath(runner.LookChromeNames("设置新的Chrome浏览器启动程序路径")),
	}
	for _, option := range extra {
		options = append(options, option)
	}
	run, err := runner.New(options...)

	if err != nil {
		return nil, err
	}

	// run.Start启动实例
	err = run.Start(ctx)
	if err != nil {
		return nil, err
	}

	// 默认情况chromedp会输出大量log，因为是示例所以选择屏蔽
	// 使用runner初始化chromedp实例
	// 实例在使用完毕后需要调用c.Shutdown()来释放资源
	c, err := chromedp.New(ctx, chromedp.WithRunner(run), chromedp.WithErrorf(log.Errorf))
	if err != nil {
		return nil, err
	}

	return c, nil
}

// NewBrowser 创建chrome实例
// chromedp内部有自己的超时设置，你也可以通过ctx来设置更短的超时
func NewBrowser(ctx context.Context, extra ...chromedp.Option) (*chromedp.CDP, error) {
	options := []chromedp.Option{
		chromedp.WithErrorf(log.Errorf),
	}
	// 默认情况chromedp会输出大量log，因为是示例所以选择屏蔽
	// 实例在使用完毕后需要调用c.Shutdown()来释放资源
	c, err := chromedp.New(ctx, append(options, extra...)...)
	if err != nil {
		return nil, err
	}
	return c, nil
}
