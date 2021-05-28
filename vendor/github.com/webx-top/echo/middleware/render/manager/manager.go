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
	caches   map[string][]byte
	firstDir string
	mutex    *sync.RWMutex
	ignores  map[string]bool
	allows   map[string]bool
	Logger   logger.Logger
	callback map[string]func(string, string, string) //参数为：目标名称，类型(file/dir)，事件名(create/delete/modify/rename)
	done     chan bool
	watcher  *fsnotify.Watcher
}

func (m *Manager) closeMoniter() {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.firstDir = ``
	if m.done == nil {
		return
	}
	close(m.done)
	m.done = nil
	if m.watcher != nil {
		m.watcher.Close()
	}
}

func (m *Manager) getWatcher() *fsnotify.Watcher {
	if m.watcher == nil {
		var err error
		m.watcher, err = fsnotify.NewWatcher()
		if err != nil {
			m.Logger.Error(err)
		}
	}
	return m.watcher
}

func (m *Manager) AddCallback(rootDir string, callback func(name, typ, event string)) {
	m.mutex.Lock()
	m.callback[rootDir] = callback
	m.mutex.Unlock()
}

func (m *Manager) ClearCallback() {
	m.callback = map[string]func(string, string, string){}
}

func (m *Manager) DelCallback(rootDir string) {
	m.mutex.Lock()
	if _, ok := m.callback[rootDir]; ok {
		delete(m.callback, rootDir)
	}
	m.mutex.Unlock()
}

func (m *Manager) ClearAllows() {
	m.allows = map[string]bool{}
}

func (m *Manager) AddAllow(allows ...string) {
	for _, allow := range allows {
		m.allows[allow] = true
	}
}

func (m *Manager) DelAllow(allow string) {
	if _, ok := m.allows[allow]; ok {
		delete(m.allows, allow)
	}
}

func (m *Manager) ClearIgnores() {
	m.ignores = map[string]bool{}
}

func (m *Manager) AddIgnore(ignores ...string) {
	for _, ignore := range ignores {
		m.allows[ignore] = false
	}
}

func (m *Manager) DelIgnore(ignore string) {
	if _, ok := m.ignores[ignore]; ok {
		delete(m.ignores, ignore)
	}
}

func (m *Manager) SetLogger(logger logger.Logger) {
	m.Logger = logger
}

func (m *Manager) allowCached(name string) bool {
	ok := len(m.allows) == 0
	if !ok {
		_, ok = m.allows[`*`+filepath.Ext(name)]
		if !ok {
			ok = m.allows[filepath.Base(name)]
		}
	}
	return ok
}

func (m *Manager) AddWatchDir(ppath string) (err error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	ppath, err = filepath.Abs(ppath)
	if err != nil {
		return
	}
	if len(m.firstDir) == 0 {
		m.firstDir = ppath
	}
	err = m.getWatcher().Add(ppath)
	if err != nil {
		m.Logger.Error(err.Error())
		return
	}

	err = filepath.Walk(ppath, func(f string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return m.getWatcher().Add(f)
		}
		return nil
	})

	//err = m.cacheAll(ppath)
	return
}

func (m *Manager) CancelWatchDir(oldDir string) (err error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	oldDir, err = filepath.Abs(oldDir)
	if err != nil {
		return
	}
	for tmpl := range m.caches {
		if strings.HasPrefix(tmpl, oldDir) {
			if err != nil {
				return
			}
			delete(m.caches, tmpl)
		}
	}
	filepath.Walk(oldDir, func(f string, info os.FileInfo, err error) error {
		if info.IsDir() {
			m.getWatcher().Remove(f)
			return nil
		}
		return nil
	})
	m.getWatcher().Remove(oldDir)
	return
}

func (m *Manager) ChangeWatchDir(oldDir string, newDir string) (err error) {
	err = m.CancelWatchDir(oldDir)
	if err != nil {
		return err
	}
	err = m.AddWatchDir(newDir)
	return
}

func (m *Manager) Start() error {
	go m.watch()
	return nil
}

func (m *Manager) watch() error {
	watcher := m.getWatcher()
	var logSuffix string
	if len(m.firstDir) > 0 {
		logSuffix = ": " + m.firstDir + " etc"
	}
	m.Logger.Debug("TemplateMgr watcher is start" + logSuffix + ".")
	defer func() {
		watcher.Close()
		m.watcher = nil
	}()
	go func() {
		for {
			select {
			case ev := <-watcher.Events:
				if _, ok := m.ignores[filepath.Base(ev.Name)]; ok {
					return
				}
				if _, ok := m.ignores[`*`+filepath.Ext(ev.Name)]; ok {
					return
				}
				d, err := os.Stat(ev.Name)
				if err != nil {
					return
				}
				if ev.Op&fsnotify.Create == fsnotify.Create {
					if d.IsDir() {
						watcher.Add(ev.Name)
						m.onChange(ev.Name, "dir", "create")
						return
					}
					m.onChange(ev.Name, "file", "create")
					if m.allowCached(ev.Name) {
						content, err := ioutil.ReadFile(ev.Name)
						if err != nil {
							m.Logger.Infof("loaded template %v failed: %v", ev.Name, err)
							return
						}
						m.Logger.Infof("loaded template file %v success", ev.Name)
						m.CacheTemplate(ev.Name, content)
					}
				} else if ev.Op&fsnotify.Remove == fsnotify.Remove {
					if d.IsDir() {
						watcher.Remove(ev.Name)
						m.onChange(ev.Name, "dir", "delete")
						return
					}
					m.onChange(ev.Name, "file", "delete")
					if m.allowCached(ev.Name) {
						m.CacheDelete(ev.Name)
					}
				} else if ev.Op&fsnotify.Write == fsnotify.Write {
					if d.IsDir() {
						m.onChange(ev.Name, "dir", "modify")
						return
					}
					m.onChange(ev.Name, "file", "modify")
					if m.allowCached(ev.Name) {
						content, err := ioutil.ReadFile(ev.Name)
						if err != nil {
							m.Logger.Errorf("reloaded template %v failed: %v", ev.Name, err)
							return
						}
						m.CacheTemplate(ev.Name, content)
						m.Logger.Infof("reloaded template %v success", ev.Name)
					}
				} else if ev.Op&fsnotify.Rename == fsnotify.Rename {
					if d.IsDir() {
						watcher.Remove(ev.Name)
						m.onChange(ev.Name, "dir", "rename")
						return
					}
					m.onChange(ev.Name, "file", "rename")
					if m.allowCached(ev.Name) {
						m.CacheDelete(ev.Name)
					}
				}
			case err := <-watcher.Errors:
				if err != nil {
					m.Logger.Error("error:", err)
				}
			}
		}
	}()

	<-m.done
	m.Logger.Debug("TemplateMgr watcher is closed" + logSuffix + ".")
	return nil
}

func (m *Manager) onChange(name, typ, event string) {
	for _, callback := range m.callback {
		callback(name, typ, event)
	}
}

func (m *Manager) cacheAll(rootDir string) error {
	fmt.Print(rootDir + ": Reading the contents of the template files, please wait... ")
	err := filepath.Walk(rootDir, func(f string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if _, ok := m.ignores[filepath.Base(f)]; !ok {
			content, err := ioutil.ReadFile(f)
			if err != nil {
				m.Logger.Debugf("load template %s error: %v", f, err)
				return err
			}
			m.Logger.Debugf("loaded template", f)
			m.caches[f] = content
		}
		return nil
	})
	fmt.Println(rootDir + ": Complete.")
	return err
}

func (m *Manager) Close() {
	m.closeMoniter()
}

func (m *Manager) GetTemplate(tmpl string) ([]byte, error) {
	tmplPath, err := filepath.Abs(tmpl)
	if err != nil {
		return nil, err
	}
	if !m.allowCached(tmplPath) {
		return ioutil.ReadFile(tmplPath)
	}
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if content, ok := m.caches[tmplPath]; ok {
		m.Logger.Debugf("load template %v from cache", tmplPath)
		return content, nil
	}
	content, err := ioutil.ReadFile(tmplPath)
	if err != nil {
		return nil, err
	}
	m.Logger.Debugf("load template %v from the file", tmplPath)
	m.caches[tmplPath] = content
	return content, err
}

func (m *Manager) SetTemplate(tmpl string, content []byte) error {
	tmplPath, err := filepath.Abs(tmpl)
	if err != nil {
		return err
	}
	m.mutex.Lock()
	defer m.mutex.Unlock()
	err = ioutil.WriteFile(tmplPath, content, 0666)
	if err != nil {
		return err
	}
	if m.allowCached(tmplPath) {
		m.Logger.Debugf("load template %v from the file", tmplPath)
		m.caches[tmplPath] = content
	}
	return err
}

func (m *Manager) CacheTemplate(tmpl string, content []byte) {
	m.Logger.Debugf("update template %v on cache", tmpl)
	m.mutex.Lock()
	m.caches[tmpl] = content
	m.mutex.Unlock()
}

func (m *Manager) CacheDelete(tmpl string) {
	m.mutex.Lock()
	if _, ok := m.caches[tmpl]; ok {
		m.Logger.Infof("delete template %v from cache", tmpl)
		delete(m.caches, tmpl)
	}
	m.mutex.Unlock()
}

func (m *Manager) ClearCache() {
	m.caches = make(map[string][]byte)
}
