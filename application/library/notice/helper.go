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
