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
)

// NewHeadless 创建headless chrome实例
// chromedp内部有自己的超时设置，你也可以通过ctx来设置更短的超时
func NewHeadless(ctx context.Context, starturl string, extra ...chromedp.ExecAllocatorOption) (context.Context, error) {
	// runner.Flag设置启动headless chrome时的命令行参数
	// runner.URL设置启动时打开的URL
	// Windows用户需要设置runner.Flag("disable-gpu", true)，具体信息参见文档的FAQ
	options := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.DisableGPU,
		chromedp.Headless,
		//chromedp.WindowSize(800, 600),
		//chromedp.ProxyServer("代理服务器地址"),
		//chromedp.ExecPath(runner.LookChromeNames("设置新的Chrome浏览器启动程序路径")),
	)
	for _, option := range extra {
		options = append(options, option)
	}
	allocCtx, _ := chromedp.NewExecAllocator(ctx, options...)

	// also set up a custom logger
	taskCtx, _ := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Errorf))

	// ensure that the browser process is started
	if err := chromedp.Run(taskCtx, chromedp.Navigate(starturl)); err != nil {
		return taskCtx, err
	}
	return taskCtx, nil
}
