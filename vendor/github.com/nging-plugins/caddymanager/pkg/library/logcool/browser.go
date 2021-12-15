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

package logcool

import (
	"strings"
)

var BrowserList = NewBrowsers()

func init() {
	BrowserList.AddSpider(`Baidu`, `Baiduspider`)
	BrowserList.AddSpider(`360`, `360Spider`)
	BrowserList.AddSpider(`Sogou`, `Sogou web spider`)
	BrowserList.AddSpider(`Google`, `Googlebot`)
	BrowserList.AddSpider(`Soso`, `Sosospider`)
	BrowserList.AddSpider(`Yisou`, `YisouSpider`)
	BrowserList.AddMobile(`UCBrowser`, `UCBrowser`)
	BrowserList.AddMobile(`MicroMessenger`, `MicroMessenger`)
	BrowserList.AddMobile(`MQQBrowser`, `MQQBrowser`)
	BrowserList.AddPC(`Maxthon`, `Maxthon`)
	BrowserList.AddPC(`QQBrowser`, `QQBrowser`)
	BrowserList.AddPC(`LBBrowser`, `LBBROWSER`)
	BrowserList.AddPC(`360Browser`, `360SE`)
	BrowserList.AddPC(`360Browser`, `360EE`)
	BrowserList.AddPC(`IE`, `MSIE`)
	BrowserList.AddPC(`Chrome`, `Chrome`)
	BrowserList.AddPC(`Firefox`, `Firefox`)
	BrowserList.AddPC(`Safari`, `Safari`)
}

func NewBrowsers() *Browers {
	return &Browers{
		Spider: map[string][]string{},
		Mobile: map[string][]string{},
		PC:     map[string][]string{},
		Other:  map[string][]string{},
	}
}

type Browers struct {
	Spider map[string][]string
	Mobile map[string][]string
	PC     map[string][]string
	Other  map[string][]string
}

func (b *Browers) AddSpider(name string, value string) {
	if _, ok := b.Spider[name]; !ok {
		b.Spider[name] = []string{}
	}
	b.Spider[name] = append(b.Spider[name], value)
}

func (b *Browers) AddMobile(name string, value string) {
	if _, ok := b.Mobile[name]; !ok {
		b.Mobile[name] = []string{}
	}
	b.Mobile[name] = append(b.Mobile[name], value)
}

func (b *Browers) AddPC(name string, value string) {
	if _, ok := b.PC[name]; !ok {
		b.PC[name] = []string{}
	}
	b.PC[name] = append(b.PC[name], value)
}

func (b *Browers) AddOther(name string, value string) {
	if _, ok := b.Other[name]; !ok {
		b.Other[name] = []string{}
	}
	b.Other[name] = append(b.Other[name], value)
}

func (b *Browers) Get(userAgent string) (string, string) {
	for name, finds := range b.Spider {
		for _, find := range finds {
			if strings.Contains(userAgent, find) {
				return `spider`, name
			}
		}
	}
	for name, finds := range b.Mobile {
		for _, find := range finds {
			if strings.Contains(userAgent, find) {
				return `mobile`, name
			}
		}
	}
	for name, finds := range b.PC {
		for _, find := range finds {
			if strings.Contains(userAgent, find) {
				return `pc`, name
			}
		}
	}
	return ``, ``
}
