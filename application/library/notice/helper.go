package notice

import "github.com/webx-top/echo"

func New(ctx echo.Context, noticeType string, user string) *NoticeAndProgress {
	noticeID := ctx.Form(`noticeID`)
	if len(noticeID) == 0 {
		noticeID = ctx.Form(`id`)
	}
	clientID := ctx.Form(`clientID`)
	noticeMode := ctx.Form(`noticeMode`)
	if len(noticeID) == 0 {
		noticeID = ctx.Form(`mode`, `element`)
	}
	var noticer *NoticeAndProgress
	prog := NewProgress()
	if len(user) > 0 && len(clientID) > 0 {
		noticerConfig := &HTTPNoticerConfig{
			User:     user,
			Type:     noticeType,
			ClientID: clientID,
			ID:       noticeID,
			Mode:     noticeMode,
		}
		noticer = noticerConfig.Noticer(ctx).WithProgress(prog)
	} else {
		noticer = DefaultNoticer.WithProgress(prog)
	}
	return noticer
}

func NewSimple(ctx echo.Context, user string) Noticer {
	clientID := ctx.Form(`clientID`)
	var noticer Noticer
	if len(user) > 0 && len(clientID) > 0 {
		noticerConfig := &HTTPNoticerConfig{
			User:     user,
			Type:     `databaseExport`,
			ClientID: clientID,
		}
		noticer = noticerConfig.Noticer(ctx)
	} else {
		noticer = DefaultNoticer
	}
	return noticer
}
