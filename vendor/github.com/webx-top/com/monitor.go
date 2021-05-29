/*

   Copyright 2016 Wenhui Shen <www.webx.top>

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.

*/

package com

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/admpub/fsnotify"
)

var (
	DefaultMonitor      = NewMonitor()
	MonitorEventEmptyFn = func(string) {}
)

func NewMonitor() *MonitorEvent {
	return &MonitorEvent{
		Create:  MonitorEventEmptyFn,
		Delete:  MonitorEventEmptyFn,
		Modify:  MonitorEventEmptyFn,
		Chmod:   MonitorEventEmptyFn,
		Rename:  MonitorEventEmptyFn,
		filters: []func(string) bool{},
	}
}

//MonitorEvent 监控事件函数
type MonitorEvent struct {
	//文件事件
	Create func(string) //创建
	Delete func(string) //删除（包含文件夹和文件。因为已经删除，无法确定是文件夹还是文件）
	Modify func(string) //修改（包含修改权限。如果是文件夹，则内部的文件被更改也会触发此事件）
	Chmod  func(string) //修改权限（windows不支持）
	Rename func(string) //重命名

	//其它
	Channel chan bool //管道
	Debug   bool
	watcher *fsnotify.Watcher
	filters []func(string) bool
	lock    sync.RWMutex
}

func (m *MonitorEvent) AddFilter(args ...func(string) bool) *MonitorEvent {
	if m.filters == nil {
		m.filters = []func(string) bool{}
	}
	m.filters = append(m.filters, args...)
	return m
}

func (m *MonitorEvent) SetFilters(args ...func(string) bool) *MonitorEvent {
	m.filters = args
	return m
}

func (m *MonitorEvent) Watch(args ...func(string) bool) *MonitorEvent {
	m.SetFilters(args...)
	go func() {
		m.backendListen()
		<-m.Channel
	}()
	return m
}

func (m *MonitorEvent) Close() error {
	if m.Channel != nil {
		close(m.Channel)
	}
	if m.watcher != nil {
		return m.watcher.Close()
	}
	return nil
}

func (m *MonitorEvent) Watcher() *fsnotify.Watcher {
	m.lock.Lock()
	defer m.lock.Unlock()
	if m.watcher == nil {
		var err error
		m.watcher, err = fsnotify.NewWatcher()
		m.Channel = make(chan bool)
		m.filters = []func(string) bool{}
		if err != nil {
			log.Panic(err)
		}
	}
	return m.watcher
}

func (m *MonitorEvent) backendListen() *MonitorEvent {
	go m.listen()
	return m
}

func (m *MonitorEvent) AddDir(dir string) error {
	f, err := os.Stat(dir)
	if err != nil {
		return err
	}
	if !f.IsDir() {
		return errors.New(dir + ` is not dir.`)
	}
	err = filepath.Walk(dir, func(f string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if m.Debug {
				log.Println(`[Monitor]`, `Add Watch:`, f)
			}
			return m.Watcher().Add(f)
		}
		return nil
	})
	return err
}

func (m *MonitorEvent) AddFile(file string) error {
	return m.Watcher().Add(file)
}

func (m *MonitorEvent) Remove(fileOrDir string) error {
	if m.watcher != nil {
		return m.watcher.Remove(fileOrDir)
	}
	return nil
}

func (m *MonitorEvent) listen() {
	for {
		watcher := m.Watcher()
		select {
		case ev, ok := <-watcher.Events:
			if !ok {
				return
			}
			if m.Debug {
				log.Println(`[Monitor]`, `Trigger Event:`, ev)
			}
			if m.filters != nil {
				var skip bool
				for _, filter := range m.filters {
					if !filter(ev.Name) {
						skip = true
						break
					}
				}
				if skip {
					break
				}
			}
			switch {
			case ev.Op&fsnotify.Create == fsnotify.Create:
				if m.IsDir(ev.Name) {
					watcher.Add(ev.Name)
				}
				if m.Create != nil {
					m.Create(ev.Name)
				}
			case ev.Op&fsnotify.Remove == fsnotify.Remove:
				if m.IsDir(ev.Name) {
					watcher.Remove(ev.Name)
				}
				if m.Delete != nil {
					m.Delete(ev.Name)
				}
			case ev.Op&fsnotify.Write == fsnotify.Write:
				if m.Modify != nil {
					m.Modify(ev.Name)
				}
			case ev.Op&fsnotify.Rename == fsnotify.Rename:
				watcher.Remove(ev.Name)
				if m.Rename != nil {
					m.Rename(ev.Name)
				}
			case ev.Op&fsnotify.Chmod == fsnotify.Chmod:
				if m.Chmod != nil {
					m.Chmod(ev.Name)
				}
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			if err != nil {
				log.Println("Watcher error:", err)
			}
		}
	}
}

func (m *MonitorEvent) IsDir(path string) bool {
	d, e := os.Stat(path)
	if e != nil {
		return false
	}
	return d.IsDir()
}

//Monitor 文件监测
func Monitor(rootDir string, callback *MonitorEvent, args ...func(string) bool) error {
	watcher := callback.Watcher()
	defer watcher.Close()
	callback.Watch(args...)
	err := callback.AddDir(rootDir)
	if err != nil {
		callback.Close()
		return err
	}
	return nil
}
