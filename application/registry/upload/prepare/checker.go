package prepare

import (
	"io"

	uploadClient "github.com/webx-top/client/upload"
)

var NopChecker uploadClient.Checker = func(rs *uploadClient.Result, rd io.Reader) error {
	return nil
}
