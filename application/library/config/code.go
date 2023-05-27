package config

import (
	"net/http"

	"github.com/webx-top/echo/code"
)

func init() {
	code.CodeDict.SetHTTPCodeToExists(code.NonPrivileged, http.StatusForbidden)
}
