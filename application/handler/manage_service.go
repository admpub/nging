package handler

import (
	"io"

	"github.com/admpub/caddyui/application/library/config"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
)

func ManageWebRestart(ctx echo.Context) error {
	wOut, wErr, err := NoticeWriter(ctx, ctx.T(`重启Web服务`))
	if err != nil {
		return err
	}
	if err := config.DefaultCLIConfig.CaddyRestart(wOut, wErr); err != nil {
		return err
	}
	return ctx.String(ctx.T(`已经执行重启命令`))
}

func ManageWebStop(ctx echo.Context) error {
	if err := config.DefaultCLIConfig.CaddyStop(); err != nil {
		return err
	}
	return ctx.String(ctx.T(`已经执行关闭命令`))
}

func NoticeWriter(ctx echo.Context, noticeType string) (wOut io.Writer, wErr io.Writer, err error) {
	user, ok := ctx.Get(`user`).(string)
	if !ok {
		return nil, nil, ctx.Redirect(`/login`)
	}
	wOut = &com.CmdResultCapturer{Do: func(b []byte) error {
		SendNotice(user, `<span class="badge badge-success">`+noticeType+`</span>`+string(b))
		return nil
	}}
	wErr = &com.CmdResultCapturer{Do: func(b []byte) error {
		SendNotice(user, `<span class="badge badge-danger">`+noticeType+`</span>`+string(b))
		return nil
	}}
	return
}

func ManageFTPRestart(ctx echo.Context) error {
	wOut, wErr, err := NoticeWriter(ctx, ctx.T(`重启FTP服务`))
	if err != nil {
		return err
	}
	if err := config.DefaultCLIConfig.FTPRestart(wOut, wErr); err != nil {
		return err
	}
	return ctx.String(ctx.T(`已经执行重启命令`))
}

func ManageFTPStop(ctx echo.Context) error {
	if err := config.DefaultCLIConfig.FTPStop(); err != nil {
		return err
	}
	return ctx.String(ctx.T(`已经执行关闭命令`))
}

func ManageService(ctx echo.Context) error {
	return ctx.Render(`manage/service`, nil)
}
