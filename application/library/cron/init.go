package cron

import (
	"os/exec"
	"time"

	"github.com/admpub/log"
	"github.com/admpub/nging/application/dbschema"
)

func InitJobs() {
	m := new(dbschema.Task)
	cnt, err := m.ListByOffset(nil, nil, 0, -1, "disabled", 0)
	if err != nil {
		log.Error(err)
		return
	}
	_ = cnt
	for _, task := range m.Objects() {
		job, err := NewJobFromTask(task)
		if err != nil {
			log.Error("InitJobs: ", err.Error())
			continue
		}
		AddJob(task.CronSpec, job)
	}
}

func runCmdWithTimeout(cmd *exec.Cmd, timeout time.Duration) (error, bool) {
	done := make(chan error)
	go func() {
		done <- cmd.Wait()
	}()

	var err error
	select {
	case <-time.After(timeout):
		log.Warnf("任务执行时间超过%d秒，强制关闭进程: %d", int(timeout/time.Second), cmd.Process.Pid)
		go func() {
			<-done // 读出上面的goroutine数据，避免阻塞导致无法退出
		}()
		if err = cmd.Process.Kill(); err != nil {
			log.Errorf("进程无法关闭: %d, 错误信息: %s", cmd.Process.Pid, err)
		}
		return err, true
	case err = <-done:
		return err, false
	}
}
