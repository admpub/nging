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
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/nging-plugins/caddymanager/application/library/logcool"
)

var (
	regexPlaceholder = regexp.MustCompile(`\{[<>]?([^\{\}]+)\}`)
	regexTimeZone    = regexp.MustCompile(` [+-][0-9]{4}$`)
)

// 私有函数, 使用正则式解析日志, 扩展性好, 速度很慢
func (l *AccessLog) parseWithPattern(line string, regex *regexp.Regexp) error {
	match := regex.FindStringSubmatch(line)
	if match == nil {
		return fmt.Errorf("NOT MATCH: %s", line)
	}

	for i, name := range regex.SubexpNames() {
		if i == 0 || len(name) == 0 {
			continue
		}
		//{remote} - {user} [{when}] "{method} {uri} {proto}" {status} {size} {latency} "{>Referer}" "{>User-Agent}"
		//(?P<XRealIP>\S+) \| \[(?P<TimeLocal>\S+) \+0800] \| (?P<Host>\S+) \| "(?P<Method>\S+) (?P<Uri>\S+) (?P<Version>\S+?)" \| (?P<StatusCode>\d+) \| (?P<BodyBytes>\S+) \| "(?P<Referer>\S+)" \| "(?P<UserAgent>.*?)" \| "(?P<XFowardFor>.*?)" \| (?P<BackendHost>.+?) \| (?P<BackendStatus>.+?) \| (?P<BackendTimeSeconds>.+?) \| (?P<LocalAddr>\S+) \| (?P<Ignore>\S+)
		switch name {
		case "RemoteAddr", "remote":
			l.RemoteAddr = match[i]
		case "XRealIP":
			l.XRealIp = match[i]
		case "user", "request_id":
			l.User = match[i]
		case "when": //Timestamp in the format 02/Jan/2006:15:04:05 -0700 in local time
			if t, err := time.Parse(`02/Jan/2006:15:04:05 -0700`, match[i]); err == nil {
				l.TimeLocal = t.Format(`2006-01-02 15:04:05`)
				l.Minute = t.Format(`15:04`)
			}
		case "when_iso": //Timestamp in the format 2006-01-02T15:04:05Z in UTC
			if t, err := time.Parse(`2006-01-02T15:04:05Z`, match[i]); err == nil {
				l.TimeLocal = t.Format(`2006-01-02 15:04:05`)
				l.Minute = t.Format(`15:04`)
			}
		case "when_unix":
			if sec, err := strconv.ParseInt(match[i], 10, 64); err == nil {
				t := time.Unix(sec, 0)
				l.TimeLocal = t.Format(`2006-01-02 15:04:05`)
				l.Minute = t.Format(`15:04`)
			}
		case "TimeLocal":
			l.TimeLocal = match[i]
			l.Minute = match[i][12:17]
		case "Uri", "URI", "uri":
			l.Uri = match[i]
		case "proto":
			l.Version = match[i]
		case "StatusCode", "status":
			i, _ := strconv.ParseUint(match[i], 10, 64)
			l.StatusCode = uint(i)
		case "BodyBytes", "size":
			l.BodyBytes, _ = strconv.ParseUint(match[i], 10, 64)
		case "latency":
			dur, err := time.ParseDuration(match[i])
			if err == nil {
				l.Elapsed = dur.Seconds()
			}
		case "TimeSeconds", "Elapsed":
			l.Elapsed, _ = strconv.ParseFloat(match[1], 64)
		case "TimeMicroSeconds":
			l.Elapsed, _ = strconv.ParseFloat(match[i], 64)
			l.Elapsed = l.Elapsed / 1000000
		case "TimeMilliSeconds", "latency_ms":
			l.Elapsed, _ = strconv.ParseFloat(match[i], 64)
			l.Elapsed = l.Elapsed / 1000
		case "Referer":
			l.Referer = match[i]
		case "AtsUri":
			requestLine := strings.SplitN(match[i], "/", 4)
			if len(requestLine) < 4 {
				l.Uri = match[i]
				l.Host = match[i]
			} else {
				l.Uri = "/" + requestLine[3]
				l.Host = requestLine[2]
			}
		case "UserAgent":
			l.UserAgent = match[i]
		case "Host", "host":
			l.Host = match[i]
		case "Method", "method":
			l.Method = match[i]
		case "XForwardFor":
			l.XForwardFor = match[i]
		case "HitStatus":
			i, _ := strconv.ParseUint(match[i], 10, 64)
			l.HitStatus = uint(i)
		case "LocalAddr":
			l.LocalAddr = match[i]
		case "scheme":
			l.Scheme = match[i]
		default:
		}
	}
	if len(l.UserAgent) > 0 {
		l.BrowerType, l.BrowerName = logcool.BrowserList.Get(l.UserAgent)
	}
	return nil
}
