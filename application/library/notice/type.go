package notice

import (
	"fmt"
	"io"

	"github.com/admpub/log"
	"github.com/webx-top/echo/middleware/tplfunc"
)

type (
	Noticer          func(message interface{}, statusCode int, progress ...*Progress) error
	CustomWithWriter func(wOut io.Writer, wErr io.Writer) Noticer
)

func (noticer Noticer) WithProgress(progresses ...*Progress) *NoticeAndProgress {
	return NewWithProgress(noticer, progresses...)
}

var (
	// DefaultNoticer 默认noticer
	// statusCode > 0 为成功；否则为失败
	DefaultNoticer Noticer = func(message interface{}, statusCode int, progs ...*Progress) error {
		if len(progs) > 0 && progs[0] != nil {
			message = `[ ` + tplfunc.NumberFormat(progs[0].CalcPercent().Percent, 2) + `% ] ` + fmt.Sprint(message)
		}
		if statusCode > 0 {
			log.Info(message)
		} else {
			log.Error(message)
		}
		return nil
	}

	CustomOutputNoticer CustomWithWriter = func(wOut io.Writer, wErr io.Writer) Noticer {
		return func(message interface{}, statusCode int, progs ...*Progress) error {
			if len(progs) > 0 && progs[0] != nil {
				message = `[ ` + tplfunc.NumberFormat(progs[0].CalcPercent().Percent, 2) + `% ] ` + fmt.Sprint(message)
			}
			if statusCode > 0 {
				fmt.Fprintln(wOut, message)
			} else {
				fmt.Fprintln(wErr, message)
			}
			return nil
		}
	}
)
