package config

import (
	"net/http"

	"github.com/webx-top/echo/code"
)

func init() {
	dict := code.CodeDict[code.NonPrivileged]
	dict.HTTPCode = http.StatusForbidden
	code.CodeDict[code.NonPrivileged] = dict
}
