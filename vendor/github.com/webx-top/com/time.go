// Copyright 2013 com authors
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package com

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

//Date Format unix time int64 to string
func Date(ti int64, format string) string {
	t := time.Unix(int64(ti), 0)
	return DateT(t, format)
}

//DateS Format unix time string to string
func DateS(ts string, format string) string {
	i, _ := strconv.ParseInt(ts, 10, 64)
	return Date(i, format)
}

//DateT Format time.Time struct to string
// MM - month - 01
// M - month - 1, single bit
// DD - day - 02
// D - day 2
// YYYY - year - 2006
// YY - year - 06
// HH - 24 hours - 03
// H - 24 hours - 3
// hh - 12 hours - 03
// h - 12 hours - 3
// mm - minute - 04
// m - minute - 4
// ss - second - 05
// s - second = 5
func DateT(t time.Time, format string) string {
	res := strings.Replace(format, "MM", t.Format("01"), -1)
	res = strings.Replace(res, "M", t.Format("1"), -1)
	res = strings.Replace(res, "DD", t.Format("02"), -1)
	res = strings.Replace(res, "D", t.Format("2"), -1)
	res = strings.Replace(res, "YYYY", t.Format("2006"), -1)
	res = strings.Replace(res, "YY", t.Format("06"), -1)
	res = strings.Replace(res, "HH", fmt.Sprintf("%02d", t.Hour()), -1)
	res = strings.Replace(res, "H", fmt.Sprintf("%d", t.Hour()), -1)
	res = strings.Replace(res, "hh", t.Format("03"), -1)
	res = strings.Replace(res, "h", t.Format("3"), -1)
	res = strings.Replace(res, "mm", t.Format("04"), -1)
	res = strings.Replace(res, "m", t.Format("4"), -1)
	res = strings.Replace(res, "ss", t.Format("05"), -1)
	res = strings.Replace(res, "s", t.Format("5"), -1)
	return res
}

// DateFormat pattern rules.
var datePatterns = []string{
	// year
	"Y", "2006", // A full numeric representation of a year, 4 digits Examples: 1999 or 2003
	"y", "06", //A two digit representation of a year Examples: 99 or 03
	// month
	"m", "01", // Numeric representation of a month, with leading zeros 01 through 12
	"n", "1", // Numeric representation of a month, without leading zeros 1 through 12
	"M", "Jan", // A short textual representation of a month, three letters Jan through Dec
	"F", "January", // A full textual representation of a month, such as January or March January through December
	// day
	"d", "02", // Day of the month, 2 digits with leading zeros 01 to 31
	"j", "2", // Day of the month without leading zeros 1 to 31
	// week
	"D", "Mon", // A textual representation of a day, three letters Mon through Sun
	"l", "Monday", // A full textual representation of the day of the week Sunday through Saturday
	// time
	"g", "3", // 12-hour format of an hour without leading zeros 1 through 12
	"G", "15", // 24-hour format of an hour without leading zeros 0 through 23
	"h", "03", // 12-hour format of an hour with leading zeros 01 through 12
	"H", "15", // 24-hour format of an hour with leading zeros 00 through 23
	"a", "pm", // Lowercase Ante meridiem and Post meridiem am or pm
	"A", "PM", // Uppercase Ante meridiem and Post meridiem AM or PM
	"i", "04", // Minutes with leading zeros 00 to 59
	"s", "05", // Seconds, with leading zeros 00 through 59
	// time zone
	"T", "MST",
	"P", "-07:00",
	"O", "-0700",
	// RFC 2822
	"r", time.RFC1123Z,
}

// DateFormatReplacer .
var DateFormatReplacer = strings.NewReplacer(datePatterns...)

//DateParse Parse Date use PHP time format.
func DateParse(dateString, format string) (time.Time, error) {
	return time.ParseInLocation(ConvDateFormat(format), dateString, time.Local)
}

//ConvDateFormat Convert PHP time format.
func ConvDateFormat(format string) string {
	format = DateFormatReplacer.Replace(format)
	return format
}

//DateFormat 将时间戳格式化为日期字符窜
func DateFormat(format string, timestamp interface{}) (t string) { // timestamp
	switch format {
	case "Y-m-d H:i:s", "":
		format = "2006-01-02 15:04:05"
	case "Y-m-d H:i":
		format = "2006-01-02 15:04"
	case "y-m-d H:i":
		format = "06-01-02 15:04"
	case "m-d H:i":
		format = "01-02 15:04"
	case "Y-m-d":
		format = "2006-01-02"
	case "y-m-d":
		format = "06-01-02"
	case "m-d":
		format = "01-02"
	default:
		format = ConvDateFormat(format)
	}
	sd := Int64(timestamp)
	t = time.Unix(sd, 0).Format(format)
	return
}

//StrToTime 日期字符窜转为时间戳数字
func StrToTime(str string, args ...string) (unixtime int) {
	layout := "2006-01-02 15:04:05"
	if len(args) > 0 {
		layout = args[0]
	}
	t, err := time.Parse(layout, str)
	if err == nil {
		unixtime = int(t.Unix())
	} else {
		fmt.Println(err, str)
	}
	return
}

// RestoreTime 从字符串还原时间
// RestoreTime(`2001-01-01T00:00:03Z`).Format(`2006-01-02 15:04:05`) => 2001-01-01 00:00:03
func RestoreTime(str string, args ...string) time.Time {
	layout := time.RFC3339
	if len(args) > 0 {
		layout = args[0]
	}
	t, _ := time.Parse(layout, str)
	return t
}

// FormatByte 兼容以前的版本，FormatBytes别名
// @param float64 size
// @param int precision
// @param bool trimRightZero
func FormatByte(args ...interface{}) string {
	return FormatBytes(args...)
}

// FormatBytes 格式化字节。 FormatBytes(字节整数，保留小数位数)
// @param float64 size
// @param int precision
// @param bool trimRightZero
func FormatBytes(args ...interface{}) string {
	sizes := [...]string{"YB", "ZB", "EB", "PB", "TB", "GB", "MB", "KB", "B"}
	var (
		total         = len(sizes)
		size          float64
		precision     int
		trimRightZero bool
	)
	switch len(args) {
	case 3:
		trimRightZero = Bool(args[2])
		fallthrough
	case 2:
		precision = Int(args[1])
		fallthrough
	case 1:
		size = Float64(args[0])
	}
	for total--; total > 0 && size > 1024.0; total-- {
		size /= 1024.0
	}
	r := fmt.Sprintf("%.*f", precision, size)
	if trimRightZero {
		r = NumberTrimZero(r)
	}
	return r + sizes[total]
}

//DateFormatShort 格式化耗时
func DateFormatShort(timestamp interface{}) string {
	now := time.Now()
	year := now.Year()
	month := now.Month()
	day := now.Day()
	cTime := StrToTime(fmt.Sprintf(`%d-%.2d-%.2d 00:00:00`, year, month, day)) //月、日始终保持两位
	timestamp2 := Int(timestamp)
	if cTime < timestamp2 {
		return DateFormat("15:04", timestamp)
	}
	cTime = StrToTime(fmt.Sprintf(`%d-01-01 00:00:00`, year))
	if cTime < timestamp2 {
		return DateFormat("01-02", timestamp)
	}
	return DateFormat("06-01-02", timestamp)
}

//FormatPastTime 格式化耗时
func FormatPastTime(timestamp interface{}, args ...string) string {
	duration := time.Now().Sub(time.Unix(Int64(timestamp), 0))
	if u := uint64(duration); u >= uint64(time.Hour)*24 {
		format := "Y-m-d H:i:s"
		if len(args) > 0 {
			format = args[0]
		}
		return DateFormat(format, timestamp)
	}
	return FriendlyTime(duration)
}

// FriendlyTime 对人类友好的经历时间格式
// @param time.Duration d
// @param string suffix
// @param int precision
// @param bool trimRightZero
// @param string language
func FriendlyTime(d time.Duration, args ...interface{}) (r string) {
	var suffix string
	var trimRightZero bool
	var language string
	precision := 2
	switch len(args) {
	case 4:
		language = String(args[3])
		fallthrough
	case 3:
		trimRightZero = Bool(args[2])
		fallthrough
	case 2:
		precision = Int(args[1])
		fallthrough
	case 1:
		suffix = String(args[0])
	}
	var unit string
	var divisor float64
	u := uint64(d)
	if u < uint64(time.Second) {
		switch {
		case u == 0:
			unit = `s`
		case u < uint64(time.Microsecond):
			unit = `ns` //纳秒
			divisor = 1
		case u < uint64(time.Millisecond):
			unit = `us` //微秒
			divisor = 1000
		default:
			unit = `ms` //毫秒
			divisor = 1000 * 1000
		}
	} else {
		divisor = 1000 * 1000 * 1000
		switch {
		case u < uint64(time.Minute):
			unit = `s` //秒
		case u < uint64(time.Hour):
			unit = `m` //分钟
			divisor *= 60
		case u < uint64(time.Hour)*24:
			unit = `h` //小时
			divisor *= 60 * 60
		case u < uint64(time.Hour)*24*7:
			unit = `d` //天
			divisor *= 60 * 60 * 24
		default:
			unit = `w` //周
			divisor *= 60 * 60 * 24 * 7
		}
	}
	if divisor > 0 {
		r = fmt.Sprintf("%.*f", precision, float64(u)/divisor)
		if trimRightZero {
			r = NumberTrimZero(r)
		}
	} else {
		r = `0`
	}
	if len(language) > 0 {
		units, ok := TimeShortUnits[language]
		if ok {
			if ut, ok := units[unit]; ok {
				unit = ut
			}
		}
	}
	r += unit + suffix
	return
}

//StartTime 开始时间
var StartTime = time.Now()

//TotalRunTime 总运行时长
func TotalRunTime() string {
	return FriendlyTime(time.Now().Sub(StartTime))
}

var (
	timeUnits     = []string{"years", "weeks", "days", "hours", "minutes", "seconds", "milliseconds"}
	timeUnitsZhCN = map[string]string{"years": "年", "weeks": "周", "days": "天", "hours": "小时", "minutes": "分", "seconds": "秒", "milliseconds": "毫秒"}
	// TimeUnits 多语言时间单位
	TimeUnits = map[string]map[string]string{`zh-cn`: timeUnitsZhCN}

	// TimeShortUnits 时间单位(简写)
	TimeShortUnits = map[string]map[string]string{
		`zh-cn`: {`s`: `秒`, `ns`: `纳秒`, `us`: `微秒`, `ms`: `毫秒`, `m`: `分钟`, `h`: `小时`, `d`: `天`, `w`: `周`},
	}
)

// Durafmt holds the parsed duration and the original input duration.
type Durafmt struct {
	duration time.Duration
	input    string // Used as reference.
	units    map[string]string
}

func getDurationUnits(args []interface{}) map[string]string {
	var units map[string]string
	if len(args) > 0 {
		switch v := args[0].(type) {
		case map[string]string:
			units = v
		case string:
			var ok bool
			units, ok = TimeUnits[v]
			if ok {
				return units
			}
			switch strings.ToLower(v) {
			case `zh_cn`, `zh-cn`:
				units = timeUnitsZhCN
			}
		}
	}
	return units
}

// ParseDuration creates a new *Durafmt struct, returns error if input is invalid.
func ParseDuration(dinput time.Duration, args ...interface{}) *Durafmt {
	input := dinput.String()
	return &Durafmt{dinput, input, getDurationUnits(args)}
}

// ParseDurationString creates a new *Durafmt struct from a string.
// returns an error if input is invalid.
func ParseDurationString(input string, args ...interface{}) (*Durafmt, error) {
	if input == "0" || input == "-0" {
		return nil, errors.New("durafmt: missing unit in duration " + input)
	}
	duration, err := time.ParseDuration(input)
	if err != nil {
		return nil, err
	}
	return &Durafmt{duration, input, getDurationUnits(args)}, nil
}

func (d *Durafmt) Duration() time.Duration {
	return d.duration
}

// String parses d *Durafmt into a human readable duration.
func (d *Durafmt) String() string {
	var duration string

	// Check for minus durations.
	if string(d.input[0]) == "-" {
		duration += "-"
		d.duration = -d.duration
	}

	// Convert duration.
	seconds := int64(d.duration.Seconds()) % 60
	minutes := int64(d.duration.Minutes()) % 60
	hours := int64(d.duration.Hours()) % 24
	days := int64(d.duration/(24*time.Hour)) % 365 % 7
	weeks := int64(d.duration/(24*time.Hour)) / 7 % 52
	years := int64(d.duration/(24*time.Hour)) / 365
	milliseconds := int64(d.duration/time.Millisecond) -
		(seconds * 1000) - (minutes * 60000) - (hours * 3600000) -
		(days * 86400000) - (weeks * 604800000) - (years * 31536000000)

	// Create a map of the converted duration time.
	durationMap := map[string]int64{
		"milliseconds": milliseconds,
		"seconds":      seconds,
		"minutes":      minutes,
		"hours":        hours,
		"days":         days,
		"weeks":        weeks,
		"years":        years,
	}

	// Construct duration string.
	for _, u := range timeUnits {
		v := durationMap[u]
		if customLable, ok := d.units[u]; ok {
			u = customLable
		}
		strval := strconv.FormatInt(v, 10)
		switch {
		// add to the duration string if v > 1.
		case v > 1:
			duration += strval + " " + u + " "
		// remove the plural 's', if v is 1.
		case v == 1:
			duration += strval + " " + strings.TrimRight(u, "s") + " "
		// omit any value with 0s or 0.
		case d.duration.String() == "0" || d.duration.String() == "0s":
			// note: milliseconds and minutes have the same suffix (m)
			// so we have to check if the units match with the suffix.

			// check for a suffix that is NOT the milliseconds suffix.
			if strings.HasSuffix(d.input, string(u[0])) && !strings.Contains(d.input, "ms") {
				// if it happens that the units are milliseconds, skip.
				if u == "milliseconds" {
					continue
				}
				duration += strval + " " + u
			}
			// process milliseconds here.
			if u == "milliseconds" {
				if strings.Contains(d.input, "ms") {
					duration += strval + " " + u
					break
				}
			}
			break
		// omit any value with 0.
		case v == 0:
			continue
		}
	}
	// trim any remaining spaces.
	duration = strings.TrimSpace(duration)
	return duration
}

// NowStr 当前时间字符串
func NowStr() string {
	return time.Now().Format(`2006-01-02 15:04:05`)
}

func NewTime(t time.Time) *Time {
	return &Time{Time: t}
}

type Time struct {
	time.Time
}

func (t Time) ParseTimestamp(timestamp interface{}) time.Time {
	return time.Unix(Int64(timestamp), 0)
}

func (t Time) SubTimestamp(timestamp interface{}) time.Duration {
	return t.Sub(t.ParseTimestamp(timestamp))
}

func (t Time) IsToday(timestamp interface{}) bool {
	st := t.ParseTimestamp(timestamp)
	return st.Day() == t.Day() && st.Month() == t.Month() && st.Year() == t.Year()
}

func (t Time) IsThisMonth(timestamp interface{}) bool {
	st := t.ParseTimestamp(timestamp)
	return st.Month() == t.Month() && st.Year() == t.Year()
}

func (t Time) IsThisYear(timestamp interface{}) bool {
	st := t.ParseTimestamp(timestamp)
	return st.Year() == t.Year()
}

func (t Time) IsAgo(timestamp interface{}, days int, units ...int) bool {
	unit := 86400
	if len(units) > 0 {
		unit = units[0]
	}
	return t.SubTimestamp(timestamp) > time.Second*time.Duration(days*unit)
}

func (t Time) IsFuture(timestamp interface{}, days int, units ...int) bool {
	unit := 86400
	if len(units) > 0 {
		unit = units[0]
	}
	st := t.ParseTimestamp(timestamp)
	return t.Unix()+int64(days*unit) <= st.Unix()
}

func (t Time) IsAfter(timestamp interface{}, agoDays int, units ...int) bool {
	unit := 86400
	if len(units) > 0 {
		unit = units[0]
	}
	st := t.ParseTimestamp(timestamp)
	return t.Unix()-int64(agoDays*unit) <= st.Unix()
}
