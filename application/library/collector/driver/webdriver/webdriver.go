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
	"net/url"
	"strings"
	"time"

	"github.com/admpub/nging/v3/application/library/collector"
	"github.com/tebeka/selenium"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/param"
)

func init() {
	collector.Browsers[`webdriver`] = New()
}

func New() *WebDriver {
	return &WebDriver{
		Base: &collector.Base{},
	}
}

type WebDriver struct {
	Driver  selenium.WebDriver
	service *selenium.Service
	*collector.Base
}

func (s *WebDriver) Start(opt echo.Store) (err error) {
	if err = s.Base.Start(opt); err != nil {
		return
	}
	browserName := opt.String(`browserName`, `chrome`)
	//chromedriver 可从 http://npm.taobao.org/mirrors/chromedriver/ 下载
	driverPath := opt.String(`driverPath`, ChromeDriverDefaultPath())
	servicePort := opt.Int(`servicePort`, defaultPort)

	var opts []selenium.ServiceOption
	switch v := opt.Get(`serviceOption`).(type) {
	case []selenium.ServiceOption:
		opts = v
	case selenium.ServiceOption:
		opts = append(opts, v)
	}

	// set browser as chrome
	caps := selenium.Capabilities{"browserName": browserName}
	switch browserName {
	case `chrome`:
		chromeDriverOption(caps)
		s.service, err = selenium.NewChromeDriverService(driverPath, servicePort, opts...)
	case `firefox`:
		s.service, err = selenium.NewGeckoDriverService(driverPath, servicePort, opts...)
	default:
		s.service, err = selenium.NewSeleniumService(driverPath, servicePort, opts...)
	}
	if err != nil {
		return
	}

	if len(s.Proxy) > 0 {
		proxyConfig := selenium.Proxy{
			Type: selenium.Manual,
		}
		var proxyURL *url.URL
		proxyURL, err = url.Parse(s.Proxy)
		if err != nil {
			return
		}
		switch strings.TrimSuffix(proxyURL.Scheme, `:`) {
		case `http`:
			proxyConfig.HTTP = s.Proxy
			proxyConfig.HTTPPort = param.String(proxyURL.Port()).Int()
		case `https`:
			proxyConfig.SSL = s.Proxy
			proxyConfig.SSLPort = param.String(proxyURL.Port()).Int()
		case `socks5`:
			proxyConfig.SOCKS = proxyURL.Host
			proxyConfig.SocksPort = param.String(proxyURL.Port()).Int()
		}
		if proxyURL.User != nil {
			pwd, ok := proxyURL.User.Password()
			if ok {
				proxyConfig.SOCKSPassword = pwd
			}
			proxyConfig.SOCKSUsername = proxyURL.User.Username()
		}
		caps.AddProxy(proxyConfig)
	}
	// remote to selenium server
	s.Driver, err = selenium.NewRemote(caps, fmt.Sprintf(serivceAPI, servicePort))
	if err == nil {
		if s.Timeout > 0 {
			s.Driver.SetPageLoadTimeout(time.Duration(s.Timeout) * time.Second)
			//s.Driver.SetAsyncScriptTimeout(time.Duration(s.Timeout) * time.Second)
		}
	}

	/*
		//循环执行一段逻辑。
		s.Driver.WaitWithTimeoutAndInterval(func(WebDriver) (bool, error){
			// 如果第一个值返回true，则退出循环，返回false继续循环
			// 第二个参数不返回nil时，会退出循环并返回一个错误
			return true,nil
		},60*time.Second,100*time.Millisecond)
	// */
	return
}

func (s *WebDriver) Close() error {
	err := s.Driver.Quit()
	if err == nil {
		err = s.service.Stop()
	}
	return err
}

func (s *WebDriver) Name() string {
	return `webdriver`
}

func (s *WebDriver) Description() string {
	return ``
}

func (s *WebDriver) Do(pageURL string, data echo.Store) ([]byte, error) {
	p := Page{Driver: s.Driver}
	err := p.Driver.Get(pageURL)
	if err != nil {
		return nil, err
	}
	var html string
	html, err = p.Driver.PageSource()
	if err != nil {
		return nil, err
	}
	s.Sleep()
	return []byte(html), err
}
