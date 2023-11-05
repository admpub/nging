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
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/tebeka/selenium"
	"github.com/webx-top/echo/param"

	"github.com/nging-plugins/collector/application/library/collector"
)

func InitClient(cfg *collector.Base) (client selenium.WebDriver, err error) {
	opt := cfg.Extra
	browserName := opt.String(`browserName`, `chrome`)
	if len(browserName) < 1 {
		browserName = `chrome`
	}
	servicePort := opt.Int(`servicePort`, defaultPort)
	// set browser as chrome
	caps := selenium.Capabilities{"browserName": browserName}

	switch browserName {
	case `chrome`:
		chromeDriverOption(caps, cfg)
	case `firefox`:
		firefoxDriverOption(caps, cfg)
	}

	if len(cfg.Proxy) > 0 {
		proxyConfig := selenium.Proxy{
			Type: selenium.Manual,
		}
		var proxyURL *url.URL
		proxyURL, err = url.Parse(cfg.Proxy)
		if err != nil {
			return
		}
		switch strings.TrimSuffix(proxyURL.Scheme, `:`) {
		case `http`:
			proxyConfig.HTTP = cfg.Proxy
			proxyConfig.HTTPPort = param.String(proxyURL.Port()).Int()
		case `https`:
			proxyConfig.SSL = cfg.Proxy
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
	client, err = selenium.NewRemote(caps, fmt.Sprintf(serivceAPI, servicePort))
	if err != nil {
		return
	}
	if cfg.Timeout > 0 {
		client.SetPageLoadTimeout(time.Duration(cfg.Timeout) * time.Second)
		//client.SetAsyncScriptTimeout(time.Duration(cfg.Timeout) * time.Second)
	}
	// if cfg.Header != nil {
	// 	for k, v := range cfg.Header {
	// 		//TODO
	// 	}
	// }
	if len(cfg.Cookies) > 0 {
		for _, c := range cfg.Cookies {
			cookie := &selenium.Cookie{
				Name:   c.Name,
				Value:  c.Value,
				Path:   c.Path,
				Domain: c.Domain,
				Secure: c.Secure,
			}
			if c.MaxAge > 0 {
				cookie.Expiry = uint(c.MaxAge)
			}
			client.AddCookie(cookie)
		}
	} else if len(cfg.CookieString) > 0 {
		header := http.Header{}
		header.Add("Cookie", cfg.CookieString)
		request := http.Request{Header: header}
		for _, c := range request.Cookies() {
			cookie := &selenium.Cookie{
				Name:   c.Name,
				Value:  c.Value,
				Path:   c.Path,
				Domain: c.Domain,
				Secure: c.Secure,
			}
			if c.MaxAge > 0 {
				cookie.Expiry = uint(c.MaxAge)
			}
			client.AddCookie(cookie)
		}
	}
	/*
		//循环执行一段逻辑。
		client.WaitWithTimeoutAndInterval(func(WebDriver) (bool, error){
			// 如果第一个值返回true，则退出循环，返回false继续循环
			// 第二个参数不返回nil时，会退出循环并返回一个错误
			return true,nil
		},60*time.Second,100*time.Millisecond)
	// */
	return
}

func NewPage(cfg *collector.Base, clients ...selenium.WebDriver) (page Page) {
	var client selenium.WebDriver
	if len(clients) > 0 {
		client = clients[0]
	} else {
		var err error
		client, err = InitClient(cfg)
		if err != nil {
			panic(err)
		}
	}
	return Page{client: client}
}

func Fetch(cfg *collector.Base, pageURL string) ([]byte, error) {
	page := NewPage(cfg)
	err := page.client.Get(pageURL)
	if err != nil {
		return nil, nil
	}
	/*
		page.MouseHoverToElement("#nav > ol > li:nth-of-type(1) > a")
		time.Sleep(time.Millisecond * 100)
		page.FindElementByCss("#nav > ol > li:nth-of-type(1) > ul > li:nth-of-type(1) > a").Click()
	*/
	html, err := page.client.PageSource()
	if err != nil {
		return nil, err
	}
	return []byte(html), err
}
