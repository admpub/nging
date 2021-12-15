package formfilter

import (
	"fmt"
	"strings"
	"time"

	"github.com/admpub/log"
	"github.com/webx-top/echo/param"
)

// DateToTimestamp 日期转时间戳
func DateToTimestamp(field string, layouts ...string) Options {
	layout := `2006-01-02`
	if len(layouts) > 0 && len(layouts[0]) > 0 {
		layout = layouts[0]
	}
	return func() (string, Filter) {
		return field, func(data *Data) {
			if len(data.Value) > 0 && len(data.Value[0]) > 0 {
				t, e := time.ParseInLocation(layout, data.Value[0], time.Local)
				if e == nil {
					data.Value = []string{fmt.Sprint(t.Unix())}
					return
				}
				log.Debug(`Form field: `, data.Key, `: `, e)
			}

			data.Value = []string{`0`}
		}
	}
}

// TimestampToDate 时间戳转日期
func TimestampToDate(field string, layouts ...string) Options {
	layout := `2006-01-02 15:04:05`
	if len(layouts) > 0 && len(layouts[0]) > 0 {
		layout = layouts[0]
	}
	return func() (string, Filter) {
		return field, func(data *Data) {
			if len(data.Value) > 0 && len(data.Value[0]) > 0 {
				tsi := strings.SplitN(data.Value[0], `.`, 2)
				var sec, nsec int64
				switch len(tsi) {
				case 2:
					nsec = param.AsInt64(tsi[1])
					fallthrough
				case 1:
					sec = param.AsInt64(tsi[0])
				}
				t := time.Unix(sec, nsec)
				if t.IsZero() {
					data.Value = []string{``}
					return
				}
				data.Value = []string{t.Local().Format(layout)}
			}
		}
	}
}

// StartDateToTimestamp 起始日期(当天的零点)转时间戳
func StartDateToTimestamp(field string, layouts ...string) Options {
	layout := `2006-01-02 15:04:05`
	if len(layouts) > 0 && len(layouts[0]) > 0 {
		layout = layouts[0]
	}
	return func() (string, Filter) {
		return field, func(data *Data) {
			if len(data.Value) > 0 && len(data.Value[0]) > 0 {
				if !strings.Contains(data.Value[0], `:`) {
					data.Value[0] += ` 00:00:00`
				}
				t, e := time.ParseInLocation(layout, data.Value[0], time.Local)
				if e == nil {
					data.Value = []string{fmt.Sprint(t.Unix())}
					return
				}
				log.Debug(`Form field: `, data.Key, `: `, e)
			}

			data.Value = []string{`0`}
		}
	}
}

// EndDateToTimestamp 结束日期(当天的最后一秒)转时间戳
func EndDateToTimestamp(field string, layouts ...string) Options {
	layout := `2006-01-02 15:04:05`
	if len(layouts) > 0 && len(layouts[0]) > 0 {
		layout = layouts[0]
	}
	return func() (string, Filter) {
		return field, func(data *Data) {
			if len(data.Value) > 0 && len(data.Value[0]) > 0 {
				if !strings.Contains(data.Value[0], `:`) {
					data.Value[0] += ` 23:59:59`
				}
				t, e := time.ParseInLocation(layout, data.Value[0], time.Local)
				if e == nil {
					data.Value = []string{fmt.Sprint(t.Unix())}
					return
				}
				log.Debug(`Form field: `, data.Key, `: `, e)
			}

			data.Value = []string{`0`}
		}
	}
}
