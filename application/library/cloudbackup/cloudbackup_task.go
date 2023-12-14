package cloudbackup

import (
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/param"

	"github.com/admpub/nging/v5/application/dbschema"
	"github.com/admpub/nging/v5/application/library/msgbox"
	"github.com/admpub/nging/v5/application/library/notice"
)

var (
	BackupTasks = param.NewMap()
)

func NewTask(monitor *com.MonitorEvent, storager Storager) *Task {
	return &Task{
		Monitor:  monitor,
		Storager: storager,
	}
}

type Task struct {
	Monitor  *com.MonitorEvent
	Storager Storager
}

func (t *Task) Close() {
	t.Monitor.Close()
	t.Storager.Close()
}

func MonitorBackupStop(id uint) error {
	if task, ok := BackupTasks.Get(id).(*Task); ok {
		task.Close()
		BackupTasks.Delete(id)
		LevelDB().CloseDB(id)
		msgbox.Success(`Cloud-Backup`, `Close: `+com.String(id))
	}
	return nil
}

func Restore(ctx echo.Context, cfg dbschema.NgingCloudBackup, callback func(from, to string), prog notice.Progressor) error {
	mgr, err := NewStorage(ctx, cfg)
	if err != nil {
		return err
	}
	if err := mgr.Connect(); err != nil {
		return err
	}
	defer mgr.Close()
	if prog != nil {
		if st, ok := mgr.(ProgressorSetter); ok {
			st.SetProgressor(prog)
		}
	}
	return mgr.Restore(ctx, cfg.DestPath, cfg.SourcePath, callback) // 从云存储服务器路径还原文件到本机源路径
}
