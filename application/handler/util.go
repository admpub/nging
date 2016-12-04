package handler

import (
	"github.com/webx-top/echo"

	"regexp"

	"os"

	"strings"

	"github.com/admpub/caddyui/application/library/errors"
)

var PageMaxSize = 1000

func Paging(ctx echo.Context) (page int, size int) {
	page = ctx.Formx(`page`).Int()
	size = ctx.Formx(`size`).Int()
	if page < 1 {
		page = 1
	}
	if size < 1 || size > PageMaxSize {
		size = 50
	}
	return
}

func Ok(v string) errors.Successor {
	return errors.NewOk(v)
}

var (
	space  = rune(' ')
	quote  = rune('"')
	slash  = rune('\\')
	envOS  = regexp.MustCompile(`\{\$[a-zA-Z0-9_]+\}`)
	envWin = regexp.MustCompile(`\{%[a-zA-Z0-9_]+%\}`)
)

func ParseArgs(command string) (params []string) {
	item := []rune{}
	hasQuote := false
	hasSlash := false
	maxIndex := len(command) - 1
	//tower.exe -c tower.yaml -p "eee\"ddd" -t aaaa
	for k, v := range command {
		if !hasQuote {
			if v == space {
				params = append(params, string(item))
				item = []rune{}
				continue
			}
			if v == quote {
				hasQuote = true
				continue
			}
		} else {
			if !hasSlash && v == quote {
				hasQuote = false
				continue
			}
			if !hasSlash && v == slash && k+1 <= maxIndex && command[k+1] == '"' {
				hasSlash = true
				continue
			}
			hasSlash = false
		}
		item = append(item, v)
	}
	if len(item) > 0 {
		params = append(params, string(item))
	}
	for k, v := range params {
		v = envWin.ReplaceAllStringFunc(v, getWinEnv)
		params[k] = envOS.ReplaceAllStringFunc(v, getEnv)
	}
	//fmt.Printf("---> %#v\n", params)
	//params = []string{}
	return
}

func getWinEnv(s string) string {
	s = strings.TrimPrefix(s, `{%`)
	s = strings.TrimSuffix(s, `%}`)
	return os.Getenv(s)
}

func getEnv(s string) string {
	s = strings.TrimPrefix(s, `{$`)
	s = strings.TrimSuffix(s, `}`)
	return os.Getenv(s)
}
