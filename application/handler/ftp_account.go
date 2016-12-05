package handler

import "github.com/webx-top/echo"

func FTPAccountIndex(ctx echo.Context) error {
	return ctx.Render(`ftp/account`, nil)
}

func FTPAccountAdd(ctx echo.Context) error {
	return ctx.Render(`ftp/account_edit`, nil)
}

func FTPAccountEdit(ctx echo.Context) error {
	return ctx.Render(`ftp/account_edit`, nil)
}

func FTPAccountDelete(ctx echo.Context) error {
	return ctx.Redirect(`/ftp/account`)
}
