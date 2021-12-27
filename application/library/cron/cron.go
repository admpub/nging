/*
   Nging is a toolbox for webmasters
   Copyright (C) 2018-present Wenhui Shen <swh@admpub.com>

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU Affero General Public License as published
   by the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Affero General Public License for more details.

   You should have received a copy of the GNU Affero General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

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
	PoolSize = 50 //连接池容量
)

func Initial(sizes ...int) {
	var size int
	if len(sizes) > 0 {
		size = sizes[0]
	}
	if size <= 0 {
		size = PoolSize
	}
	Close()
	workPool = make(chan bool, size)
	mainCron = cron.New(cron.WithSeconds())
}

func Running() bool {
	return mainCron != nil && mainCron.Running()
}

func Close() {
	if Running() {
		mainCron.Stop()
		mainCron = nil
		close(workPool)
		workPool = nil
		historyJobsRunning = false
		log.Info(`退出任务处理`)
	}
}

func MainCron(mustStart bool) *cron.Cron {
	if mainCron == nil {
		Initial()
	}
	if mustStart {
		if !mainCron.Running() {
			mainCron.Start()
		}
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
	_, err := MainCron(true).AddJob(spec, job)
	if err != nil {
		log.Error("AddJob: ", err.Error())
		return false
	}
	return true
}

func RemoveJob(id uint) {
	MainCron(false).RemoveJob(func(e *cron.Entry) bool {
		if v, ok := e.Job.(*Job); ok {
			if v.id == id {
				return true
			}
		}
		return false
	})
}

func GetEntryById(id uint) *cron.Entry {
	entries := MainCron(false).Entries()
	for _, e := range entries {
		if v, ok := e.Job.(*Job); ok {
			if v.id == id {
				return &e
			}
		}
	}
	return nil
}

func GetEntries(size int) []cron.Entry {
	ret := MainCron(false).Entries()
	if len(ret) > size {
		return ret[:size]
	}
	return ret
}
