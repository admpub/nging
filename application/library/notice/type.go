/*
   Nging is a toolbox for webmasters
   Copyright (C) 2018-present  Wenhui Shen <swh@admpub.com>

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU Affero General Public License as published
   by the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Affero General Public License for more details.

   You should have received a copy of the GNU Affero General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package notice

import (
	"fmt"
	"io"

	"github.com/admpub/log"
	"github.com/webx-top/echo/middleware/tplfunc"
)

type (
	Noticer          func(message interface{}, state int, progress ...*Progress) error
	CustomWithWriter func(wOut io.Writer, wErr io.Writer) Noticer
)

func (noticer Noticer) WithProgress(progresses ...*Progress) *NoticeAndProgress {
	return NewWithProgress(noticer, progresses...)
}

var (
	// DefaultNoticer 默认noticer
	// state > 0 为成功；否则为失败
	DefaultNoticer Noticer = func(message interface{}, state int, progs ...*Progress) error {
		if len(progs) > 0 && progs[0] != nil {
			message = `[ ` + tplfunc.NumberFormat(progs[0].CalcPercent().Percent, 2) + `% ] ` + fmt.Sprint(message)
		}
		if state > 0 {
			log.Info(message)
		} else {
			log.Error(message)
		}
		return nil
	}

	CustomOutputNoticer CustomWithWriter = func(wOut io.Writer, wErr io.Writer) Noticer {
		return func(message interface{}, state int, progs ...*Progress) error {
			if len(progs) > 0 && progs[0] != nil {
				message = `[ ` + tplfunc.NumberFormat(progs[0].CalcPercent().Percent, 2) + `% ] ` + fmt.Sprint(message)
			}
			if state > 0 {
				fmt.Fprintln(wOut, message)
			} else {
				fmt.Fprintln(wErr, message)
			}
			return nil
		}
	}
)
