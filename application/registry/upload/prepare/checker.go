package prepare

import uploadClient "github.com/webx-top/client/upload"

var NopChecker uploadClient.Checker = func(r *uploadClient.Result) error {
	return nil
}
