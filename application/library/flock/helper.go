package flock

import (
	"os"
	"time"

	"github.com/admpub/log"
	"github.com/webx-top/com"
)

// IsCompleted 等待文件有数据且已写完
// 费时操作 放在子线程中执行
// @param file  文件
// @param start 需要传入 time.Now.Local()，用于兼容遍历的情况
// @return true:已写完 false:外部程序阻塞或者文件不存在
func IsCompleted(file *os.File, start time.Time) bool {
	ok, err := com.FileIsCompleted(file, start)
	if err != nil {
		log.Error(err)
	}
	return ok
}
