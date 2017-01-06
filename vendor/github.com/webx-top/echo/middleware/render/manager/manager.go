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
package manager

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/howeyc/fsnotify"
	"github.com/webx-top/echo/logger"
)

func New(logger logger.Logger, tmplDir string, allows []string, callback func(name, typ, event string), cached ...bool) *Manager {
	m := new(Manager)
	ln := len(cached)
	if ln < 1 || !cached[0] {
		return m
	}
	reloadTemplates := true
	if ln > 1 {
		reloadTemplates = cached[1]
	}
	m.SetOnChangeCallback(callback)
	m.Init(logger, tmplDir, reloadTemplates, allows...)
	return m
}

type Manager struct {
	Caches           map[string][]byte
	lock             *sync.Once
	RootDir          string
	NewRoorDir       string
	Ignores          map[string]bool
	CachedAllows     map[string]bool
	Logger           logger.Logger
	Preprocessor     func([]byte) []byte
	timerCallback    func() bool
	TimerCallback    func() bool
	initialized      bool
	OnChangeCallback func(string, string, string) //参数为：目标名称，类型(file/dir)，事件名(create/delete/modify/rename)
	done             chan bool
}

func (self *Manager) closeMoniter() {
	close(self.done)
}

func (self *Manager) SetOnChangeCallback(callback func(name, typ, event string)) {
	self.OnChangeCallback = callback
}
func (self *Manager) SetLogger(logger logger.Logger) {
	self.Logger = logger
}

func (self *Manager) allowCached(name string) bool {
	_, ok := self.CachedAllows["*.*"]
	if !ok {
		_, ok = self.CachedAllows[`*`+filepath.Ext(name)]
		if !ok {
			ok = self.CachedAllows[filepath.Base(name)]
		}
	}
	return ok
}

func (self *Manager) Moniter(rootDir string) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	//fmt.Println("[webx] TemplateMgr watcher is start.")
	defer watcher.Close()
	self.done = make(chan bool)
	self.lock = &sync.Once{}
	go func() {
		for {
			select {
			case ev := <-watcher.Event:
				self.lock.Do(func() {
					defer func() {
						self.lock = &sync.Once{}
					}()
					if ev == nil {
						return
					}
					if _, ok := self.Ignores[filepath.Base(ev.Name)]; ok {
						return
					}
					if _, ok := self.Ignores[`*`+filepath.Ext(ev.Name)]; ok {
						return
					}
					d, err := os.Stat(ev.Name)
					if err != nil {
						return
					}
					if ev.IsCreate() {
						if d.IsDir() {
							watcher.Watch(ev.Name)
							self.onChange(ev.Name, "dir", "create")
							return
						}
						self.onChange(ev.Name, "file", "create")
						if self.allowCached(ev.Name) {
							tmpl := ev.Name[len(self.RootDir)+1:]
							content, err := ioutil.ReadFile(ev.Name)
							if err != nil {
								self.Logger.Infof("loaded template %v failed: %v", tmpl, err)
								return
							}
							self.Logger.Infof("loaded template file %v success", tmpl)
							self.CacheTemplate(tmpl, content)
						}
					} else if ev.IsDelete() {
						if d.IsDir() {
							watcher.RemoveWatch(ev.Name)
							self.onChange(ev.Name, "dir", "delete")
							return
						}
						self.onChange(ev.Name, "file", "delete")
						if self.allowCached(ev.Name) {
							tmpl := ev.Name[len(self.RootDir)+1:]
							self.CacheDelete(tmpl)
						}
					} else if ev.IsModify() {
						if d.IsDir() {
							self.onChange(ev.Name, "dir", "modify")
							return
						}
						self.onChange(ev.Name, "file", "modify")
						if self.allowCached(ev.Name) {
							tmpl := ev.Name[len(self.RootDir)+1:]
							content, err := ioutil.ReadFile(ev.Name)
							if err != nil {
								self.Logger.Errorf("reloaded template %v failed: %v", tmpl, err)
								return
							}
							self.CacheTemplate(tmpl, content)
							self.Logger.Infof("reloaded template %v success", tmpl)
						}
					} else if ev.IsRename() {
						if d.IsDir() {
							watcher.RemoveWatch(ev.Name)
							self.onChange(ev.Name, "dir", "rename")
							return
						}
						self.onChange(ev.Name, "file", "rename")
						if self.allowCached(ev.Name) {
							tmpl := ev.Name[len(self.RootDir)+1:]
							self.CacheDelete(tmpl)
						}
					}

				})
			case err := <-watcher.Error:
				self.Logger.Error("error:", err)
			case <-time.After(time.Second * 2):
				if self.timerCallback != nil {
					if self.timerCallback() == false {
						close(self.done)
						return
					}
				}
				//fmt.Printf("TemplateMgr timer operation: %v.\n", time.Now())
			}
		}
	}()

	err = filepath.Walk(self.RootDir, func(f string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return watcher.Watch(f)
		}
		return nil
	})

	if err != nil {
		self.Logger.Error(err.Error())
		return err
	}

	<-self.done
	//fmt.Println("[webx] TemplateMgr watcher is closed.")
	return nil
}

func (self *Manager) onChange(name, typ, event string) {
	if self.OnChangeCallback != nil {
		name = FixDirSeparator(name)
		self.OnChangeCallback(name[len(self.RootDir)+1:], typ, event)
	}
}

func (self *Manager) cacheAll(rootDir string) error {
	fmt.Print("Reading the contents of the template files, please wait... ")
	err := filepath.Walk(rootDir, func(f string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		tmpl := f[len(rootDir)+1:]
		tmpl = FixDirSeparator(tmpl)
		if _, ok := self.Ignores[filepath.Base(tmpl)]; !ok {
			fpath := filepath.Join(self.RootDir, tmpl)
			content, err := ioutil.ReadFile(fpath)
			if err != nil {
				self.Logger.Debugf("load template %s error: %v", fpath, err)
				return err
			}
			self.Logger.Debugf("loaded template", fpath)
			self.Caches[tmpl] = content
		}
		return nil
	})
	fmt.Println("Complete.")
	return err
}

func (self *Manager) defaultTimerCallback() func() bool {
	return func() bool {
		if self.TimerCallback != nil {
			return self.TimerCallback()
		}
		//更改模板主题后，关闭当前监控，重新监控新目录
		if self.NewRoorDir == "" || self.NewRoorDir == self.RootDir {
			return true
		}
		self.ClearCache()
		self.Ignores = make(map[string]bool)
		self.RootDir = self.NewRoorDir
		go self.Moniter(self.RootDir)
		return false
	}
}

func (self *Manager) Close() {
	self.TimerCallback = func() bool {
		self.ClearCache()
		self.Ignores = make(map[string]bool)
		self.TimerCallback = nil
		return false
	}
	self.initialized = false
}

func (self *Manager) Init(logger logger.Logger, rootDir string, reload bool, allows ...string) {
	if self.initialized {
		if rootDir == self.RootDir {
			return
		}
		self.TimerCallback = func() bool {
			self.ClearCache()
			self.Ignores = make(map[string]bool)
			self.CachedAllows = make(map[string]bool)
			self.TimerCallback = nil
			return false
		}
	} else if !reload {
		self.TimerCallback = func() bool {
			self.TimerCallback = nil
			return false
		}
	}
	self.RootDir = rootDir
	self.Caches = make(map[string][]byte)
	self.Ignores = make(map[string]bool)
	self.CachedAllows = make(map[string]bool)
	for _, allow := range allows {
		self.CachedAllows[allow] = true
	}
	self.Logger = logger
	if dirExists(rootDir) {
		//self.cacheAll(rootDir)
		if reload {
			self.timerCallback = self.defaultTimerCallback()
			go self.Moniter(rootDir)
		}
	}

	if len(self.Ignores) == 0 {
		self.Ignores["*.tmp"] = false
		self.Ignores["*.TMP"] = false
	}
	if len(self.CachedAllows) == 0 {
		self.CachedAllows["*.*"] = true
	}
	self.initialized = true
}

func (self *Manager) GetTemplate(tmpl string) ([]byte, error) {
	tmpl = FixDirSeparator(tmpl)
	if tmpl[0] == '/' {
		tmpl = tmpl[1:]
	}

	if content, ok := self.Caches[tmpl]; ok {
		self.Logger.Debugf("load template %v from cache", tmpl)
		return content, nil
	}

	content, err := ioutil.ReadFile(filepath.Join(self.RootDir, tmpl))
	if err == nil {
		self.Logger.Debugf("load template %v from the file", tmpl)
		self.Caches[tmpl] = content
	}
	return content, err
}

func (self *Manager) CacheTemplate(tmpl string, content []byte) {
	if self.Preprocessor != nil {
		content = self.Preprocessor(content)
	}

	tmpl = FixDirSeparator(tmpl)
	self.Logger.Debugf("update template %v on cache", tmpl)
	self.Caches[tmpl] = content
	return
}

func (self *Manager) CacheDelete(tmpl string) {
	tmpl = FixDirSeparator(tmpl)
	if _, ok := self.Caches[tmpl]; ok {
		self.Logger.Infof("delete template %v from cache", tmpl)
		delete(self.Caches, tmpl)
	}
	return
}

func (self *Manager) ClearCache() {
	self.Caches = make(map[string][]byte)
}
