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

package webdriver

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/nging-plugins/collector/application/library/collector"
	"github.com/webx-top/com"

	"github.com/admpub/marmot/miner"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
	"github.com/tebeka/selenium/firefox"
	"github.com/webx-top/echo"
)

var (
	drivers     = sync.Map{}
	serivceAPI  = "http://127.0.0.1:%d/wd/hub"
	defaultPort = 4444
)

func ChromeDriverDefaultPath() string {
	wd := echo.Wd()
	if !com.IsExist(filepath.Join(wd, `support`)) {
		gopath := os.Getenv(`GOPATH`)
		if len(gopath) > 0 {
			wd = filepath.Join(gopath, `src/github.com/admpub/nging`)
		}
	}
	switch runtime.GOOS {
	case `windows`:
		return filepath.Join(wd, `support`, `chromedriver_386.exe`)
	case `linux`, `darwin`:
		return filepath.Join(wd, `support`, `chromedriver_`+runtime.GOOS+`_`+runtime.GOARCH)
	default:
		return ``
	}
}

func StartService(driverPath string, port int, opts ...selenium.ServiceOption) (service *selenium.Service, err error) {
	//chromedriver 可从 http://npm.taobao.org/mirrors/chromedriver/ 下载
	if len(driverPath) == 0 {
		driverPath = `chromedriver`
	}
	// 启动chromedriver，端口号可自定义
	service, err = selenium.NewChromeDriverService(driverPath, port, opts...)
	if err != nil {
		log.Printf("Error starting the ChromeDriver server: %v", err)
	}
	//defer service.Stop()
	return
}

func chromeDriverOption(caps selenium.Capabilities, cfg *collector.Base) {
	// 禁止加载图片，加快渲染速度
	imagCaps := map[string]interface{}{
		"profile.managed_default_content_settings.images": 2,
	}
	var userAgent string
	if len(cfg.UserAgent) > 0 {
		userAgent = cfg.UserAgent
	} else {
		userAgent = miner.RandomUserAgent()
	}
	chromeCaps := chrome.Capabilities{
		Prefs: imagCaps,
		Path:  "",
		Args: []string{
			"--headless", // 设置Chrome无头模式
			"--no-sandbox",
			"--disable-gpu",
			"--user-agent=" + userAgent, // 模拟user-agent，防反爬
			//"window-size=1200x600"
		},
	}
	if cfg.Header != nil {
		for k, v := range cfg.Header {
			chromeCaps.Args = append(chromeCaps.Args, k+`=`+fmt.Sprintf(`%q`, v))
		}
	}
	caps.AddChrome(chromeCaps)
}

func firefoxDriverOption(caps selenium.Capabilities, cfg *collector.Base) {
	imagCaps := map[string]interface{}{}
	var userAgent string
	if len(cfg.UserAgent) > 0 {
		userAgent = cfg.UserAgent
	} else {
		userAgent = miner.RandomUserAgent()
	}
	firefoxCaps := firefox.Capabilities{
		Prefs:  imagCaps,
		Binary: "",
		Args: []string{
			"--headless", // 设置Firefox无头模式
			"--disable-gpu",
			"--user-agent=" + userAgent, // 模拟user-agent，防反爬
			//"window-size=1200x600"
		},
	}
	if cfg.Header != nil {
		for k, v := range cfg.Header {
			firefoxCaps.Args = append(firefoxCaps.Args, k+`=`+fmt.Sprintf(`%q`, v))
		}
	}
	caps.AddFirefox(firefoxCaps)
}

func InitServer(cfg *collector.Base, browserName ...string) (driver selenium.WebDriver, err error) {
	if len(browserName) < 1 {
		browserName = []string{`chrome`}
	}
	if dri, ok := drivers.Load(browserName[0]); ok {
		return dri.(selenium.WebDriver), nil
	}
	// set browser as chrome
	caps := selenium.Capabilities{"browserName": browserName[0]}

	switch browserName[0] {
	case `chrome`:
		chromeDriverOption(caps, cfg)
	case `firefox`:
		firefoxDriverOption(caps, cfg)
	}

	// remote to selenium server
	driver, err = selenium.NewRemote(caps, fmt.Sprintf(serivceAPI, defaultPort))
	if err == nil {
		drivers.Store(browserName[0], driver)
	}
	return
}

func CloseServer(browserName ...string) error {
	if len(browserName) < 1 {
		drivers.Range(func(key, val interface{}) bool {
			val.(selenium.WebDriver).Quit()
			drivers.Delete(key)
			return true
		})
		return nil
	}
	if dri, ok := drivers.Load(browserName[0]); ok {
		drivers.Delete(browserName[0])
		return dri.(selenium.WebDriver).Quit()
	}
	return nil
}

func NewPage(cfg *collector.Base, drivers ...selenium.WebDriver) (page Page) {
	var driver selenium.WebDriver
	if len(drivers) > 0 {
		driver = drivers[0]
	} else {
		var err error
		driver, err = InitServer(cfg)
		if err != nil {
			panic(err)
		}
	}
	return Page{Driver: driver}
}

func Fetch(cfg *collector.Base, pageURL string) ([]byte, error) {
	page := NewPage(cfg)
	err := page.Driver.Get(pageURL)
	if err != nil {
		return nil, nil
	}
	/*
		page.MouseHoverToElement("#nav > ol > li:nth-of-type(1) > a")
		time.Sleep(time.Millisecond * 100)
		page.FindElementByCss("#nav > ol > li:nth-of-type(1) > ul > li:nth-of-type(1) > a").Click()
	*/
	html, err := page.Driver.PageSource()
	if err != nil {
		return nil, err
	}
	return []byte(html), err
}
