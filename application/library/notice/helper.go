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

import "github.com/webx-top/echo"

func NewP(ctx echo.Context, noticeType string, user string) *NoticeAndProgress {
	return New(ctx, noticeType, user).WithProgress(NewProgress())
}

func New(ctx echo.Context, noticeType string, user string) Noticer {
	clientID := ctx.Form(`clientID`)
	var noticer Noticer
	if len(user) > 0 && len(clientID) > 0 {
		noticeID := ctx.Form(`noticeID`)
		if len(noticeID) == 0 {
			noticeID = ctx.Form(`id`)
		}
		noticeMode := ctx.Form(`noticeMode`)
		if len(noticeID) == 0 {
			noticeID = ctx.Form(`mode`, `element`)
		}
		noticerConfig := &HTTPNoticerConfig{
			User:     user,
			Type:     noticeType,
			ClientID: clientID,
			ID:       noticeID,
			Mode:     noticeMode,
		}
		noticer = noticerConfig.Noticer(ctx)
	} else {
		noticer = DefaultNoticer
	}
	return noticer
}
