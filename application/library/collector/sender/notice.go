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

package sender

import (
	"fmt"
	"io"

	"github.com/webx-top/echo"
	"github.com/webx-top/echo/middleware/tplfunc"

	"github.com/admpub/log"
	"github.com/admpub/nging/application/library/notice"
)

var Default Notice = func(message interface{}, statusCode int, progs ...*notice.Progress) error {
	if len(progs) > 0 && progs[0] != nil {
		message = `[ ` + tplfunc.NumberFormat(progs[0].Percent, 2) + `% ] ` + echo.Dump(message, false)
	}
	if statusCode > 0 {
		log.Info(message)
	} else {
		log.Error(message)
	}
	return nil
}

var CustomOutput func(wOut io.Writer, wErr io.Writer) Notice = func(wOut io.Writer, wErr io.Writer) Notice {
	return func(message interface{}, statusCode int, progs ...*notice.Progress) error {
		if len(progs) > 0 && progs[0] != nil {
			message = `[ ` + tplfunc.NumberFormat(progs[0].Percent, 2) + `% ] ` + echo.Dump(message, false)
		}
		if statusCode > 0 {
			fmt.Fprintln(wOut, message)
		} else {
			fmt.Fprintln(wErr, message)
		}
		return nil
	}
}

type Notice func(message interface{}, statusCode int, progress ...*notice.Progress) error
