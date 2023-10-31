package form

import (
	"regexp"
	"strings"
	"time"

	"github.com/webx-top/echo/param"
)

var spaceSeperate = regexp.MustCompile(`[\s]+`)
var durationRegex = regexp.MustCompile(`(?P<years>\d+y)?(?P<months>\d+m)?(?P<days>\d+d)?T?(?P<hours>\d+h)?(?P<minutes>\d+i)?(?P<seconds>\d+s)?`)

func SplitBySpace(value string, formatter ...func(string) string) []string {
	value = strings.TrimSpace(value)
	if len(value) == 0 {
		return nil
	}
	values := spaceSeperate.Split(value, -1)
	if len(formatter) == 0 {
		return values
	}
	for index, value := range values {
		for _, format := range formatter {
			value = format(value)
		}
		values[index] = value
	}
	return values
}

func ExplodeCombinedLogFormat(value string) string {
	return strings.Replace(value, `{combined}`, `{remote} - {user} [{when}] "{method} {uri} {proto}" {status} {size}`, 1)
}

func ParseDuration(str string) time.Duration {
	matches := durationRegex.FindStringSubmatch(str)
	if len(matches) < 7 {
		return 0
	}

	years := parseInt64ForDuration(matches[1])
	months := parseInt64ForDuration(matches[2])
	days := parseInt64ForDuration(matches[3])
	hours := parseInt64ForDuration(matches[4])
	minutes := parseInt64ForDuration(matches[5])
	seconds := parseInt64ForDuration(matches[6])

	hour := int64(time.Hour)
	minute := int64(time.Minute)
	second := int64(time.Second)
	return time.Duration(years*24*365*hour + months*30*24*hour + days*24*hour + hours*hour + minutes*minute + seconds*second)
}

func parseInt64ForDuration(str string) int64 {
	if len(str) == 0 {
		return 0
	}
	return param.AsInt64(str[:len(str)-1])
}

func AddCSlashesIngoreSlash(s string, b ...rune) string {
	var builder strings.Builder
	var cnt int
	for _, v := range s {
		if v == '\\' {
			cnt++
		} else {
			cnt = 0
		}
		for _, f := range b {
			if v == f {
				builder.WriteRune('\\')
				break
			}
		}
		builder.WriteRune(v)
	}
	if cnt%2 != 0 {
		builder.WriteRune('\\')
	}
	return builder.String()
}
