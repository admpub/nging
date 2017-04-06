package cron

import (
	"sync"

	"github.com/admpub/cron"
	"github.com/admpub/log"
)

var (
	mainCron *cron.Cron
	workPool chan bool
	lock     sync.Mutex
)

func Initial(size int) {
	workPool = make(chan bool, size)
	mainCron = cron.New()
	mainCron.Start()
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
	err := mainCron.AddJob(spec, job)
	if err != nil {
		log.Error("AddJob: ", err.Error())
		return false
	}
	return true
}

func RemoveJob(id uint) {
	mainCron.RemoveJob(func(e *cron.Entry) bool {
		if v, ok := e.Job.(*Job); ok {
			if v.id == id {
				return true
			}
		}
		return false
	})
}

func GetEntryById(id uint) *cron.Entry {
	entries := mainCron.Entries()
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
	ret := mainCron.Entries()
	if len(ret) > size {
		return ret[:size]
	}
	return ret
}
