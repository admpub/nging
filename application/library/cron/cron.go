package cron

import (
	"sync"

	"github.com/admpub/cron"
	"github.com/admpub/log"
	"github.com/admpub/nging/application/library/config"
)

var (
	mainCron *cron.Cron
	workPool chan bool
	lock     sync.Mutex
)

func Initial(sizes ...int) {
	var size int
	if len(sizes) > 0 {
		size = sizes[0]
	} else {
		size = config.DefaultConfig.Cron.PoolSize
	}
	if size <= 0 {
		size = 1
	}
	if mainCron != nil {
		mainCron.Stop()
		mainCron = nil
		close(workPool)
	}
	workPool = make(chan bool, size)
	mainCron = cron.New()
	mainCron.Start()
}

func MainCron() *cron.Cron {
	if mainCron == nil {
		Initial()
	}
	return mainCron
}

func Parse(spec string) error {
	_, err := cron.Parse(spec)
	return err
}

func AddJob(spec string, job *Job) bool {
	lock.Lock()
	defer lock.Unlock()

	if GetEntryById(job.id) != nil {
		return false
	}
	err := MainCron().AddJob(spec, job)
	if err != nil {
		log.Error("AddJob: ", err.Error())
		return false
	}
	return true
}

func RemoveJob(id uint) {
	MainCron().RemoveJob(func(e *cron.Entry) bool {
		if v, ok := e.Job.(*Job); ok {
			if v.id == id {
				return true
			}
		}
		return false
	})
}

func GetEntryById(id uint) *cron.Entry {
	entries := MainCron().Entries()
	for _, e := range entries {
		if v, ok := e.Job.(*Job); ok {
			if v.id == id {
				return e
			}
		}
	}
	return nil
}

func GetEntries(size int) []*cron.Entry {
	ret := MainCron().Entries()
	if len(ret) > size {
		return ret[:size]
	}
	return ret
}
