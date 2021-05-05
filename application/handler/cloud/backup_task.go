package cloud

import (
	"time"

	"github.com/admpub/nging/application/handler"
	"github.com/admpub/nging/application/model"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
)

func BackupStart(ctx echo.Context) error {
	id := ctx.Formx(`id`).Uint()
	m := model.NewCloudBackup(ctx)
	recv := &model.CloudBackupExt{}
	err := m.NewParam().SetArgs(db.Cond{`id`: id}).SetRecv(recv).One()
	if err != nil {
		if err == db.ErrNoMoreRows {
			err = ctx.NewError(code.DataNotFound, ctx.T(`数据不存在`))
		}
		return err
	}
	if len(recv.Storage.Endpoint) == 0 {
		return ctx.NewError(code.InvalidParameter, ctx.T(`Endpoint无效`))
	}
	switch ctx.Form(`op`) {
	case "full":
		err = fullBackupStart(recv)
		if err != nil {
			if err == ErrRunningPleaseWait {
				err = ctx.NewError(code.OperationProcessing, ctx.T(`运行中，请稍候，如果文件很多可能需要会多等一会儿`))
			}
		}
	default:
		err = monitorBackupStart(recv)
	}
	if err != nil {
		return err
	}
	err = m.SetField(nil, `last_executed`, time.Now().Local().Unix(), `id`, m.Id)
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
			err = ctx.NewError(code.DataNotFound, ctx.T(`数据不存在`))
		}
		return err
	}
	switch ctx.Form(`op`) {
	case "full":
		if fullBackupIsRunning(m.Id) {
			fullBackupExit = true
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
		fullBackupExit = true
	}
	return monitorBackupStop(id)
}
