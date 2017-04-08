package cron

import (
	"os/exec"
	"time"

	"github.com/admpub/log"
	"github.com/admpub/nging/application/dbschema"
)

var historyJobsRunning bool

func InitJobs() error {
	m := new(dbschema.Task)
	limit := 1000
	cnt, err := m.ListByOffset(nil, nil, 0, limit, "disabled", `N`)
	if err != nil {
		return err
	}
	total := int(cnt())
	for offset := 0; offset < total; offset += limit {
		if offset > 0 {
			_, err := m.ListByOffset(nil, nil, offset, limit, "disabled", `N`)
			if err != nil {
				return err
			}
		}
		for _, task := range m.Objects() {
			job, err := NewJobFromTask(task)
			if err != nil {
				log.Error("InitJobs: ", err.Error())
				continue
			}
			if AddJob(task.CronSpec, job) {
				log.Infof("InitJobs: 添加任务[%d]", task.Id)
				continue
			}
		}
	}
	historyJobsRunning = true
	return nil
}

func HistoryJobsRunning() bool {
	return historyJobsRunning
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
