package handler

import (
	"io"

	"os"

	"github.com/admpub/nging/application/library/config"
	"github.com/admpub/nging/application/library/notice"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
)

func ManageWebRestart(ctx echo.Context) error {
	wOut, wErr, err := NoticeWriter(ctx, ctx.T(`Web服务`))
	if err != nil {
		return ctx.String(err.Error())
	}
	if err := config.DefaultCLIConfig.CaddyRestart(wOut, wErr); err != nil {
		return ctx.String(err.Error())
	}
	return ctx.String(ctx.T(`已经重启Web服务`))
}

func ManageWebLog(ctx echo.Context) error {
	on := ctx.Formx(`on`).Bool()
	if on {
		wOut, wErr, err := NoticeWriter(ctx, ctx.T(`Web服务`))
		if err != nil {
			return ctx.String(err.Error())
		}
		err = config.DefaultCLIConfig.SetLogWriter(`caddy`, wOut, wErr)
		if err != nil {
			return ctx.String(err.Error())
		}
		return ctx.String(ctx.T(`已经开始直播Web服务状态`))
	}
	err := config.DefaultCLIConfig.SetLogWriter(`caddy`, os.Stdout, os.Stderr)
	if err != nil {
		return ctx.String(err.Error())
	}
	user, ok := ctx.Get(`user`).(string)
	if !ok {
		return ctx.String(ctx.T(`请先登录`))
	}
	typ := `service:` + ctx.T(`Web服务`)
	notice.CloseMessage(user, typ)
	return ctx.String(ctx.T(`已经停止直播Web服务状态`))
}

func ManageWebStop(ctx echo.Context) error {
	if err := config.DefaultCLIConfig.CaddyStop(); err != nil {
		return ctx.String(err.Error())
	}
	return ctx.String(ctx.T(`已经关闭Web服务`))
}

func NoticeWriter(ctx echo.Context, noticeType string) (wOut io.Writer, wErr io.Writer, err error) {
	user, ok := ctx.Get(`user`).(string)
	if !ok {
		return nil, nil, ctx.Redirect(`/login`)
	}
	typ := `service:` + noticeType
	notice.OpenMessage(user, typ)

	wOut = &com.CmdResultCapturer{Do: func(b []byte) error {
		os.Stdout.Write(b)
		notice.Send(user, notice.NewMessageWithValue(typ, noticeType, string(b), notice.Succeed))
		return nil
	}}
	wErr = &com.CmdResultCapturer{Do: func(b []byte) error {
		os.Stderr.Write(b)
		notice.Send(user, notice.NewMessageWithValue(typ, noticeType, string(b), notice.Failed))
		return nil
	}}
	return
}

func ManageFTPRestart(ctx echo.Context) error {
	wOut, wErr, err := NoticeWriter(ctx, ctx.T(`FTP服务`))
	if err != nil {
		return ctx.String(err.Error())
	}
	if err := config.DefaultCLIConfig.FTPRestart(wOut, wErr); err != nil {
		return ctx.String(err.Error())
	}
	return ctx.String(ctx.T(`已经重启FTP服务`))
}

func ManageFTPStop(ctx echo.Context) error {
	if err := config.DefaultCLIConfig.FTPStop(); err != nil {
		return ctx.String(err.Error())
	}
	return ctx.String(ctx.T(`已经关闭FTP服务`))
}

func ManageFTPLog(ctx echo.Context) error {
	on := ctx.Formx(`on`).Bool()
	if on {
		wOut, wErr, err := NoticeWriter(ctx, ctx.T(`FTP服务`))
		if err != nil {
			return ctx.String(err.Error())
		}
		err = config.DefaultCLIConfig.SetLogWriter(`ftp`, wOut, wErr)
		if err != nil {
			return ctx.String(err.Error())
		}
		return ctx.String(ctx.T(`已经开始直播FTP服务状态`))
	}
	err := config.DefaultCLIConfig.SetLogWriter(`ftp`, os.Stdout, os.Stderr)
	if err != nil {
		return ctx.String(err.Error())
	}
	user, ok := ctx.Get(`user`).(string)
	if !ok {
		return ctx.String(ctx.T(`请先登录`))
	}
	typ := `service:` + ctx.T(`FTP服务`)
	notice.CloseMessage(user, typ)
	return ctx.String(ctx.T(`已经停止直播FTP服务状态`))
}

func ManageService(ctx echo.Context) error {
	return ctx.Render(`manage/service`, nil)
}
