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
	lock    *sync.Once
	watcher *fsnotify.Watcher
}

func (m *MonitorEvent) Watch(rootDir string, args ...func(string) bool) {
	go func() {
		err := Monitor(rootDir, m, args...)
		if err != nil {
			log.Println(err.Error())
		}
	}()
}

func (m *MonitorEvent) Watcher() *fsnotify.Watcher {
	if m.watcher == nil {
		var err error
		m.watcher, err = fsnotify.NewWatcher()
		if err != nil {
			log.Panic(err)
		}
	}
	return m.watcher
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
	var filter func(string) bool
	if len(args) > 0 {
		filter = args[0]
	}
	f, err := os.Stat(rootDir)
	if err != nil {
		return err
	}
	if !f.IsDir() {
		return errors.New(rootDir + ` is not dir.`)
	}
	watcher := callback.Watcher()
	defer watcher.Close()
	callback.lock = &sync.Once{}
	callback.Channel = make(chan bool)
	go func() {
		for {
			select {
			case ev := <-watcher.Events:
				if callback.Debug {
					log.Println(`[Monitor]`, `Trigger Event:`, ev)
				}
				if filter != nil {
					if !filter(ev.Name) {
						break
					}
				}
				callback.lock.Do(func() {
					switch ev.Op {
					case fsnotify.Create:
						if callback.IsDir(ev.Name) {
							watcher.Add(ev.Name)
						}
						if callback.Create != nil {
							callback.Create(ev.Name)
						}
					case fsnotify.Remove:
						if callback.Delete != nil {
							callback.Delete(ev.Name)
						}
					case fsnotify.Write:
						if callback.Modify != nil {
							callback.Modify(ev.Name)
						}
					case fsnotify.Rename:
						if callback.Rename != nil {
							callback.Rename(ev.Name)
						}
					case fsnotify.Chmod:
						if callback.Chmod != nil {
							callback.Chmod(ev.Name)
						}
					}
					callback.lock = &sync.Once{}
				})
			case err := <-watcher.Errors:
				if err != nil {
					log.Println("Watcher error:", err)
				}
			}
		}
	}()

	err = filepath.Walk(rootDir, func(f string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if callback.Debug {
				log.Println(`[Monitor]`, `Add Watch:`, f)
			}
			return watcher.Add(f)
		}
		return nil
	})

	if err != nil {
		close(callback.Channel)
		return err
	}

	<-callback.Channel
	return nil
}
