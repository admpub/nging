package oauth2server

import (
	"github.com/admpub/log"
	"github.com/admpub/oauth2/v4/errors"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
)

func InternalErrorHandler(err error) (re *errors.Response) {
	myErr, ok := err.(*echo.Error)
	if !ok {
		return
	}
	re = &errors.Response{}
	re.Error = errors.New(myErr.Code.String())
	re.Description = myErr.Error()
	switch myErr.Code {
	case code.DataNotFound:
		re.StatusCode = errors.StatusCodes[errors.ErrInvalidClient]
	case code.DataUnavailable:
		re.StatusCode = errors.StatusCodes[errors.ErrInvalidClient]
	default:
		log.Debug(`oauth2server.InternalErrorHandler: `, err)
		re.StatusCode = errors.StatusCodes[errors.ErrServerError]
	}
	return
}

func ResponseErrorHandler(re *errors.Response) {
	log.Debugf(`oauth2server.ResponseErrorHandler: %#v`, re)
}
