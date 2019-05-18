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
	"strings"
	"sync"

	"github.com/admpub/fsnotify"
	"github.com/admpub/log"
	"github.com/webx-top/echo/logger"
	"github.com/webx-top/echo/middleware/render/driver"
)

var Default driver.Manager = New()

func New() *Manager {
	m := &Manager{
		caches: make(map[string][]byte),
		mutex:  &sync.RWMutex{},
		ignores: map[string]bool{
			"*.tmp": false,
			"*.TMP": false,
		},
		allows:   map[string]bool{},
		callback: map[string]func(string, string, string){},
		Logger:   log.GetLogger(`watcher`),
		done:     make(chan bool),
	}
	return m
}

// Manager Tempate manager
type Manager struct {
	caches       map[string][]byte
	mapping      map[string]string
	firstDir     string
	mutex        *sync.RWMutex
	ignores      map[string]bool
	allows       map[string]bool
	Logger       logger.Logger
	preprocessor func([]byte) []byte
	callback     map[string]func(string, string, string) //参数为：目标名称，类型(file/dir)，事件名(create/delete/modify/rename)
	done         chan bool
	watcher      *fsnotify.Watcher
	fallback     []func(string) (string, bool)
}

func (self *Manager) closeMoniter() {
	self.mutex.Lock()
	defer self.mutex.Unlock()
	self.firstDir = ``
	if self.done == nil {
		return
	}
	close(self.done)
	self.done = nil
	if self.watcher != nil {
		self.watcher.Close()
	}
}

func (self *Manager) getWatcher() *fsnotify.Watcher {
	if self.watcher == nil {
		var err error
		self.watcher, err = fsnotify.NewWatcher()
		if err != nil {
			self.Logger.Error(err)
		}
	}
	return self.watcher
}

func (self *Manager) AddCallback(rootDir string, callback func(name, typ, event string)) {
	self.mutex.Lock()
	self.callback[rootDir] = callback
	self.mutex.Unlock()
}

func (self *Manager) ClearCallback() {
	self.callback = map[string]func(string, string, string){}
}

func (self *Manager) DelCallback(rootDir string) {
	self.mutex.Lock()
	if _, ok := self.callback[rootDir]; ok {
		delete(self.callback, rootDir)
	}
	self.mutex.Unlock()
}

func (self *Manager) ClearAllows() {
	self.allows = map[string]bool{}
}

func (self *Manager) AddAllow(allows ...string) {
	for _, allow := range allows {
		self.allows[allow] = true
	}
}

func (self *Manager) DelAllow(allow string) {
	if _, ok := self.allows[allow]; ok {
		delete(self.allows, allow)
	}
}

func (self *Manager) ClearIgnores() {
	self.ignores = map[string]bool{}
}

func (self *Manager) AddIgnore(ignores ...string) {
	for _, ignore := range ignores {
		self.allows[ignore] = false
	}
}

func (self *Manager) DelIgnore(ignore string) {
	if _, ok := self.ignores[ignore]; ok {
		delete(self.ignores, ignore)
	}
}

func (self *Manager) SetLogger(logger logger.Logger) {
	self.Logger = logger
}

func (self *Manager) allowCached(name string) bool {
	ok := len(self.allows) == 0
	if !ok {
		_, ok = self.allows[`*`+filepath.Ext(name)]
		if !ok {
			ok = self.allows[filepath.Base(name)]
		}
	}
	return ok
}

func (self *Manager) AddWatchDir(ppath string) (err error) {
	self.mutex.Lock()
	defer self.mutex.Unlock()
	ppath, err = filepath.Abs(ppath)
	if err != nil {
		return
	}
	if len(self.firstDir) == 0 {
		self.firstDir = ppath
	}
	err = self.getWatcher().Add(ppath)
	if err != nil {
		self.Logger.Error(err.Error())
		return
	}

	err = filepath.Walk(ppath, func(f string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return self.getWatcher().Add(f)
		}
		return nil
	})

	//err = self.cacheAll(ppath)
	return
}

func (self *Manager) CancelWatchDir(oldDir string) (err error) {
	self.mutex.Lock()
	defer self.mutex.Unlock()
	oldDir, err = filepath.Abs(oldDir)
	if err != nil {
		return
	}
	for tmpl := range self.caches {
		if strings.HasPrefix(tmpl, oldDir) {
			if err != nil {
				return
			}
			delete(self.caches, tmpl)
		}
	}
	filepath.Walk(oldDir, func(f string, info os.FileInfo, err error) error {
		if info.IsDir() {
			self.getWatcher().Remove(f)
			return nil
		}
		return nil
	})
	self.getWatcher().Remove(oldDir)
	return
}

func (self *Manager) ChangeWatchDir(oldDir string, newDir string) (err error) {
	err = self.CancelWatchDir(oldDir)
	if err != nil {
		return err
	}
	err = self.AddWatchDir(newDir)
	return
}

func (self *Manager) Start() error {
	go self.watch()
	return nil
}

func (self *Manager) watch() error {
	watcher := self.getWatcher()
	var logSuffix string
	if len(self.firstDir) > 0 {
		logSuffix = ": " + self.firstDir + " etc"
	}
	self.Logger.Debug("TemplateMgr watcher is start" + logSuffix + ".")
	defer func() {
		watcher.Close()
		self.watcher = nil
	}()
	go func() {
		for {
			select {
			case ev := <-watcher.Events:
				if _, ok := self.ignores[filepath.Base(ev.Name)]; ok {
					return
				}
				if _, ok := self.ignores[`*`+filepath.Ext(ev.Name)]; ok {
					return
				}
				d, err := os.Stat(ev.Name)
				if err != nil {
					return
				}
				if ev.Op&fsnotify.Create == fsnotify.Create {
					if d.IsDir() {
						watcher.Add(ev.Name)
						self.onChange(ev.Name, "dir", "create")
						return
					}
					self.onChange(ev.Name, "file", "create")
					if self.allowCached(ev.Name) {
						content, err := ioutil.ReadFile(ev.Name)
						if err != nil {
							self.Logger.Infof("loaded template %v failed: %v", ev.Name, err)
							return
						}
						self.Logger.Infof("loaded template file %v success", ev.Name)
						self.CacheTemplate(ev.Name, content)
					}
				} else if ev.Op&fsnotify.Remove == fsnotify.Remove {
					if d.IsDir() {
						watcher.Remove(ev.Name)
						self.onChange(ev.Name, "dir", "delete")
						return
					}
					self.onChange(ev.Name, "file", "delete")
					if self.allowCached(ev.Name) {
						self.CacheDelete(ev.Name)
					}
				} else if ev.Op&fsnotify.Write == fsnotify.Write {
					if d.IsDir() {
						self.onChange(ev.Name, "dir", "modify")
						return
					}
					self.onChange(ev.Name, "file", "modify")
					if self.allowCached(ev.Name) {
						content, err := ioutil.ReadFile(ev.Name)
						if err != nil {
							self.Logger.Errorf("reloaded template %v failed: %v", ev.Name, err)
							return
						}
						self.CacheTemplate(ev.Name, content)
						self.Logger.Infof("reloaded template %v success", ev.Name)
					}
				} else if ev.Op&fsnotify.Rename == fsnotify.Rename {
					if d.IsDir() {
						watcher.Remove(ev.Name)
						self.onChange(ev.Name, "dir", "rename")
						return
					}
					self.onChange(ev.Name, "file", "rename")
					if self.allowCached(ev.Name) {
						self.CacheDelete(ev.Name)
					}
				}
			case err := <-watcher.Errors:
				if err != nil {
					self.Logger.Error("error:", err)
				}
			}
		}
	}()

	<-self.done
	self.Logger.Debug("TemplateMgr watcher is closed" + logSuffix + ".")
	return nil
}

func (self *Manager) onChange(name, typ, event string) {
	for _, callback := range self.callback {
		callback(name, typ, event)
	}
}

func (self *Manager) cacheAll(rootDir string) error {
	fmt.Print(rootDir + ": Reading the contents of the template files, please wait... ")
	err := filepath.Walk(rootDir, func(f string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if _, ok := self.ignores[filepath.Base(f)]; !ok {
			content, err := ioutil.ReadFile(f)
			if err != nil {
				self.Logger.Debugf("load template %s error: %v", f, err)
				return err
			}
			self.Logger.Debugf("loaded template", f)
			self.caches[f] = content
		}
		return nil
	})
	fmt.Println(rootDir + ": Complete.")
	return err
}

func (self *Manager) Close() {
	self.closeMoniter()
}

func (self *Manager) GetTemplate(tmpl string) ([]byte, error) {
	tmplPath, err := filepath.Abs(tmpl)
	if err != nil {
		return nil, err
	}
	self.mutex.Lock()
	defer self.mutex.Unlock()
	if content, ok := self.caches[tmplPath]; ok {
		self.Logger.Debugf("load template %v from cache", tmplPath)
		return content, nil
	}
	content, err := ioutil.ReadFile(tmplPath)
	if err != nil {
		return nil, err
	}
	self.Logger.Debugf("load template %v from the file", tmplPath)
	self.caches[tmplPath] = content
	return content, err
}

func (self *Manager) CacheTemplate(tmpl string, content []byte) {
	if self.preprocessor != nil {
		content = self.preprocessor(content)
	}
	self.Logger.Debugf("update template %v on cache", tmpl)
	self.mutex.Lock()
	self.caches[tmpl] = content
	self.mutex.Unlock()
	return
}

func (self *Manager) CacheDelete(tmpl string) {
	self.mutex.Lock()
	if _, ok := self.caches[tmpl]; ok {
		self.Logger.Infof("delete template %v from cache", tmpl)
		delete(self.caches, tmpl)
	}
	self.mutex.Unlock()
	return
}

func (self *Manager) ClearCache() {
	self.caches = make(map[string][]byte)
}
