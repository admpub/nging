package cloud

import (
	"time"

	"github.com/admpub/nging/v5/application/handler"
	"github.com/admpub/nging/v5/application/model"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
)

func BackupStart(ctx echo.Context) error {
	id := ctx.Formx(`id`).Uint()
	m := model.NewCloudBackup(ctx)
	err := m.Get(nil, `id`, id)
	if err != nil {
		if err == db.ErrNoMoreRows {
			err = ctx.NewError(code.DataNotFound, `数据不存在`)
		}
		return err
	}
	switch ctx.Form(`op`) {
	case "full":
		err = fullBackupStart(m.NgingCloudBackup)
		if err != nil {
			if err == ErrRunningPleaseWait {
				err = ctx.NewError(code.OperationProcessing, `运行中，请稍候，如果文件很多可能会需要多等一会儿`)
			}
		}
	default:
		err = monitorBackupStart(m.NgingCloudBackup)
	}
	if err != nil {
		return err
	}
	err = m.UpdateField(nil, `last_executed`, time.Now().Unix(), `id`, m.Id)
	if err != nil {
		return err
	}
	handler.SendOk(ctx, ctx.T(`操作成功`))
	return ctx.Redirect(handler.URLFor(`/cloud/backup`))
}

func BackupStop(ctx echo.Context) error {
	id := ctx.Formx(`id`).Uint()
	m := model.NewCloudBackup(ctx)
	err := m.Get(nil, db.Cond{`id`: id})
	if err != nil {
		if err == db.ErrNoMoreRows {
			err = ctx.NewError(code.DataNotFound, `数据不存在`)
		}
		return err
	}
	switch ctx.Form(`op`) {
	case "full":
		if fullBackupIsRunning(m.Id) {
			fullBackupExit.Store(true)
		}
	default:
		err = monitorBackupStop(m.Id)
	}
	if err != nil {
		return err
	}
	handler.SendOk(ctx, ctx.T(`操作成功`))
	return ctx.Redirect(handler.URLFor(`/cloud/backup`))
}

func allBackupStop(id uint) error {
	if fullBackupIsRunning(id) {
		fullBackupExit.Store(true)
	}
	return monitorBackupStop(id)
}
