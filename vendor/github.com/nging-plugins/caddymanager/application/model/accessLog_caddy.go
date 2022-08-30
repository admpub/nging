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

package model

import (
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/nging-plugins/caddymanager/application/library/logcool"
)

func (l *AccessLog) parseCaddyHelper2(part string) string {
	var regexpRule string
	var hasQuote bool
	p := strings.Index(part, `"`)
	for i := 0; p > -1; i++ {
		hasQuote = true
		if i%2 == 0 {
			regexpRule += regexPlaceholder.ReplaceAllStringFunc(part[0:p], func(r string) string {
				r = strings.Replace(r, `-`, ``, -1)
				return regexPlaceholder.ReplaceAllString(r, `(?P<$1>\S+)`)
			})
		} else {
			if strings.Contains(part[0:p], ` `) {
				regexpRule += `"` + regexPlaceholder.ReplaceAllStringFunc(part[0:p], func(r string) string {
					r = strings.Replace(r, `-`, ``, -1)
					return regexPlaceholder.ReplaceAllString(r, `(?P<$1>\S+)`)
				}) + `"`
			} else {
				regexpRule += `"` + regexPlaceholder.ReplaceAllStringFunc(part[0:p], func(r string) string {
					r = strings.Replace(r, `-`, ``, -1)
					return regexPlaceholder.ReplaceAllString(r, `(?P<$1>[^"]+)`)
				}) + `"`
			}
		}
		if len(part) > p+1 {
			part = part[p+1:]
			p = strings.Index(part, `"`)
		} else {
			part = ``
			break
		}
	}
	if !hasQuote || len(part) > 0 {
		regexpRule += regexPlaceholder.ReplaceAllString(part, `(?P<$1>\S+)`)
	}
	return regexpRule
}

func (l *AccessLog) parseCaddyHelper(layout string) string {
	var regexpRule string
	pos := strings.Index(layout, `[`)
	if pos > -1 {
		regexpRule += l.parseCaddyHelper2(layout[0:pos])
		layout = layout[pos+1:]
		pos = strings.Index(layout, `]`)
		if pos > -1 {
			regexpRule += regexPlaceholder.ReplaceAllString(layout[0:pos], `\[(?P<$1>[^\]]+)\]`)
			layout = layout[pos+1:]
			regexpRule += l.parseCaddyHelper(layout)
		}
	} else {
		regexpRule += l.parseCaddyHelper2(layout)
	}
	return regexpRule
}

func (l *AccessLog) parseCaddyHardCode(line string) error {

	pos := strings.Index(line, ` `)
	if pos > -1 {
		l.RemoteAddr = line[0:pos] //{remote}
		line = line[pos+1:]
	}
	pos = strings.Index(line, ` `)
	if pos > -1 { // -
		line = line[pos+1:]
	}
	pos = strings.Index(line, ` `)
	if pos > -1 {
		l.User = line[0:pos] //{user}
		line = line[pos+1:]
	}
	pos = strings.Index(line, `[`)
	if pos > -1 {
		line = line[pos+1:]
	}
	pos = strings.Index(line, `]`)
	if pos > -1 { // {when}
		if t, err := time.Parse(`02/Jan/2006:15:04:05 -0700`, line[0:pos]); err == nil {
			l.TimeLocal = t.Format(`2006-01-02 15:04:05`)
			l.Minute = t.Format(`15:04`)
		}
		line = line[pos+1:]
	}
	pos = strings.Index(line, `"`)
	if pos > -1 {
		line = line[pos+1:]
	}
	pos = strings.Index(line, `"`)
	if pos > -1 { // {method} {uri} {proto}
		urlInfo := strings.Split(line[0:pos], " ")
		switch len(urlInfo) {
		case 3: // {method} {uri} {proto}
			l.Method = urlInfo[0]
			l.Uri = urlInfo[1]
			l.Version = urlInfo[2]
		case 5: // {method} {scheme} {host} {uri} {proto}
			l.Method = urlInfo[0]
			l.Scheme = urlInfo[1]
			l.Host = urlInfo[2]
			l.Uri = urlInfo[3]
			l.Version = urlInfo[4]
		}
		line = line[pos+1:]
	}
	pos = strings.Index(line, `"`)
	if pos > -1 { // {status} {size} {latency}
		urlInfo := strings.SplitN(strings.TrimSpace(line[0:pos]), " ", 2)
		i, _ := strconv.ParseUint(urlInfo[0], 10, 64)
		l.StatusCode = uint(i)
		l.BodyBytes, _ = strconv.ParseUint(urlInfo[1], 10, 64)
		line = line[pos+1:]
	}
	pos = strings.Index(line, `"`)
	if pos > -1 {
		l.Referer = line[0:pos] //{>Referer}
		line = line[pos+1:]
	}
	pos = strings.Index(line, `"`)
	if pos > -1 {
		line = line[pos+1:]
	}
	pos = strings.Index(line, `"`)
	if pos > -1 {
		l.UserAgent = line[0:pos] //{>User-Agent}
		line = line[pos+1:]
	}
	pos = strings.Index(line, ` `)
	if pos > -1 {
		line = line[pos+1:]
	}
	if dur, err := time.ParseDuration(line); err == nil {
		l.Elapsed = dur.Seconds()
	}
	if len(l.UserAgent) > 0 {
		l.BrowerType, l.BrowerName = logcool.BrowserList.Get(l.UserAgent)
	}
	return nil
}

var CaddyLogRegexpList = sync.Map{}

func (l *AccessLog) parseCaddy(line string, layout string) error {
	var err error
	if len(layout) == 0 || layout == `{common}` || layout == `{combined}` || layout == `{combined} {latency}` || layout == `{remote} - {user} [{when}] "{method} {uri} {proto}" {status} {size} "{>Referer}" "{>User-Agent}" {latency}` {
		return l.parseCaddyHardCode(line)
	}
	if layout == `{remote} - {user} [{when}] "{method} {scheme} {host} {uri} {proto}" {status} {size} "{>Referer}" "{>User-Agent}" {latency}` {
		return l.parseCaddyHardCode(line)
	}

	var re *regexp.Regexp
	cc, ok := CaddyLogRegexpList.Load(layout)
	if !ok {
		regexpRule := l.parseCaddyHelper(layout)
		//panic(regexpRule)
		re, err = regexp.Compile(regexpRule)
		if err != nil {
			return err
		}
		CaddyLogRegexpList.Store(layout, re)
	} else {
		re = cc.(*regexp.Regexp)
	}
	err = l.parseWithPattern(line, re)
	return err
}
