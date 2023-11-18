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

package dockerclient

import (
	"context"
	"io"
	"log"

	"github.com/admpub/nging/v5/application/library/background"
	"github.com/admpub/nging/v5/application/library/notice"
	"github.com/nging-plugins/dockermanager/application/library/utils"
	"github.com/webx-top/echo"
)

func StartBackgroundRun(ctx echo.Context, username string, actionIdent string, bgKey string, runner func(context.Context) (io.ReadCloser, error)) error {
	bg := background.New(context.Background(), nil)
	group, err := background.Register(ctx, actionIdent, bgKey, bg)
	if err != nil {
		return err
	}
	go func() {
		defer group.Cancel(bgKey)
		noticer := notice.NewP(ctx, actionIdent, username, bg.Context()).Add(100).AutoComplete(true)
		reader, err := runner(bg.Context())
		if err != nil {
			noticer.Send(err.Error(), notice.StateFailure)
			return
		}
		SyncReaderToNotice(noticer, reader)
	}()
	return nil
}

func SyncReaderToNotice(noticer notice.NProgressor, reader io.ReadCloser) {
	var lastMessage string
	defer reader.Close()
	buf := utils.NewResultWriter(func(r *utils.Result) error {
		if r.Completed() {
			noticer.Complete()
		}
		lastMessage = r.Status
		return noticer.Send(r.Status, notice.StateSuccess)
	})
	w := io.MultiWriter(buf, log.Writer())
	_, err := io.Copy(w, reader)
	if err != nil {
		noticer.Complete().Send(err.Error(), notice.StateFailure)
		return
	}
	buf.Flush()
	noticer.Done(100).Complete().Send(lastMessage, notice.StateSuccess)
}
