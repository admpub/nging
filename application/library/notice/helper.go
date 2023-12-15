/*
   Nging is a toolbox for webmasters
   Copyright (C) 2018-present Wenhui Shen <swh@admpub.com>

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
	"context"

	"github.com/webx-top/echo"
)

func NewP(eCtx echo.Context, noticeType string, user string, ctx context.Context, opts ...func(*HTTPNoticerConfig)) *NoticeAndProgress {
	return New(eCtx, noticeType, user, ctx, opts...).WithProgress(NewProgress())
}

func GetClientID(eCtx echo.Context) string {
	clientID := eCtx.Form(`notifyClientID`)
	if len(clientID) == 0 {
		clientID = eCtx.Form(`clientID`)
	}
	return clientID
}

func GetNoticeID(eCtx echo.Context) string {
	noticeID := eCtx.Form(`noticeID`)
	if len(noticeID) == 0 {
		noticeID = eCtx.Form(`id`)
	}
	return noticeID
}

func GetNoticeMode(eCtx echo.Context) string {
	noticeMode := eCtx.Form(`noticeMode`)
	if len(noticeMode) == 0 {
		noticeMode = eCtx.Form(`mode`, `element`)
	}
	return noticeMode
}

func New(eCtx echo.Context, noticeType string, user string, ctx context.Context, opts ...func(*HTTPNoticerConfig)) Noticer {
	clientID := GetClientID(eCtx)
	var noticer Noticer
	if len(user) > 0 && len(clientID) > 0 {
		noticeID := GetNoticeID(eCtx)
		noticeMode := GetNoticeMode(eCtx)
		noticerConfig := &HTTPNoticerConfig{
			User:     user,
			Type:     noticeType,
			ClientID: clientID,
			ID:       noticeID,
			Mode:     noticeMode,
		}
		for _, opt := range opts {
			opt(noticerConfig)
		}
		noticer = noticerConfig.Noticer(ctx)
	} else {
		noticer = DefaultNoticer
	}
	return noticer
}
