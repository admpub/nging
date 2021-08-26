package flock

import (
	"os"
	"time"

	"github.com/admpub/log"
)

// IsCompleted 等待文件有数据且已写完
// 费时操作 放在子线程中执行
// @param file  文件
// @param start 需要传入 time.Now.Local()，用于兼容遍历的情况
// @return true:已写完 false:外部程序阻塞或者文件不存在
// 翻译自：https://blog.csdn.net/northernice/article/details/115986671
func IsCompleted(file *os.File, start time.Time) bool {
	var (
		fileLength int64
		i          int
		waitTime   = 500 * time.Microsecond
	)
	for {
		fi, err := file.Stat()
		if err != nil {
			log.Error(err)
			return false
		}
		//文件在外部一直在填充数据，每次进入循环体时，文件大小都会改变，一直到不改变时，说明文件数据填充完毕 或者文件大小一直都是0(外部程序阻塞)
		//判断文件大小是否有改变
		if fi.Size() > fileLength { //有改变说明还未写完
			fileLength = fi.Size()
			if i%120 == 0 { //每隔1分钟输出一次日志 (i为120时：120*500/1000=60秒)
				log.Info("文件: " + fi.Name() + " 正在被填充，请稍候...")
			}
			time.Sleep(waitTime) //半秒后再循环一次
		} else { //否则：只能等于 不会小于，等于有两种情况，一种是数据写完了，一种是外部程序阻塞了，导致文件大小一直为0
			if fi.Size() != 0 { //被填充完成则立即输出日志
				return true
			}
			//等待外部程序开始写 只等60秒 120*500/1000=60秒

			//每隔1分钟输出一次日志 (i为120时：120*500/1000=60秒)
			if i%120 == 0 {
				log.Info("文件: " + fi.Name() + " 大小为0，正在等待外部程序填充，已等待：" + time.Since(start).String())
			}

			//如果一直(i为120时：120*500/1000=60秒)等于0，说明外部程序阻塞了
			if i >= 3600 { //120为1分钟 3600为30分钟
				log.Info("文件: " + fi.Name() + " 大小在：" + time.Since(start).String() + " 内始终为0，说明：在[程序监测时间内]文件写入进程依旧在运行，程序监测时间结束") //入库未完成或发生阻塞
				return false
			}

			time.Sleep(waitTime)
		}
		i++
	}
}
