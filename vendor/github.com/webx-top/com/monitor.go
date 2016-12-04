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
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/howeyc/fsnotify"
)

//监控事件函数
type MonitorEvent struct {
	Create  func(string) //创建
	Delete  func(string) //删除
	Modify  func(string) //修改
	Rename  func(string) //重命名
	Channel chan bool    //管道
	lock    *sync.Once
}

func (m *MonitorEvent) Watch(rootDir string, args ...func(string) bool) {
	go Monitor(rootDir, m, args...)
}

//文件监测
func Monitor(rootDir string, callback *MonitorEvent, args ...func(string) bool) error {
	var filter func(string) bool
	if len(args) > 0 {
		filter = args[0]
	}
	f, err := os.Stat(rootDir)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if !f.IsDir() {
		return errors.New(rootDir + ` is not dir.`)
	}
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()
	callback.lock = &sync.Once{}
	callback.Channel = make(chan bool)
	go func() {
		for {
			select {
			case ev := <-watcher.Event:
				if ev == nil {
					break
				}
				if filter != nil {
					if !filter(ev.Name) {
						break
					}
				}
				d, err := os.Stat(ev.Name)
				if err != nil {
					break
				}
				callback.lock.Do(func() {
					if callback.Create != nil && ev.IsCreate() {
						if d.IsDir() {
							watcher.Watch(ev.Name)
						} else {
							callback.Create(ev.Name)
						}
					} else if callback.Delete != nil && ev.IsDelete() {
						if d.IsDir() {
							watcher.RemoveWatch(ev.Name)
						} else {
							callback.Delete(ev.Name)
						}
					} else if callback.Modify != nil && ev.IsModify() {
						if d.IsDir() {
						} else {
							callback.Modify(ev.Name)
						}
					} else if callback.Rename != nil && ev.IsRename() {
						if d.IsDir() {
							watcher.RemoveWatch(ev.Name)
						} else {
							callback.Rename(ev.Name)
						}
					}
					callback.lock = &sync.Once{}
				})
			case err := <-watcher.Error:
				fmt.Println("Watcher error:", err)
			}
		}
	}()

	err = filepath.Walk(rootDir, func(f string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return watcher.Watch(f)
		}
		return nil
	})

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	<-callback.Channel
	return nil
}
