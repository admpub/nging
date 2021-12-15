/*
   Nging is a toolbox for webmasters
   Copyright (C) 2019-present  Wenhui Shen <swh@admpub.com>

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

package utils

import (
	"context"
	"sync"
	"time"

	"github.com/webx-top/echo"
	"github.com/webx-top/echo/param"
)

// Backgrounds 后台任务集合(包含SQL导入和导出任务)
var Backgrounds = sync.Map{} //后台导入导出任务

// Exec 执行信息
type Exec map[string]*BGExec

// Cancel 取消某个任务
func (e *Exec) Cancel(cacheKey string) {
	if bgExec, ok := (*e)[cacheKey]; ok {
		bgExec.Cancel()()
		delete(*e, cacheKey)
	}
}

// Exists 任务是否存在
func (e *Exec) Exists(cacheKey string) bool {
	_, ok := (*e)[cacheKey]
	return ok
}

// Add 新增任务
func (e *Exec) Add(op OP, cacheKey string, bgExec *BGExec) {
	e.Cancel(cacheKey) // 避免被覆盖后旧任务失去控制，先取消已存在的任务
	(*e)[cacheKey] = bgExec
	Backgrounds.Store(op, *e)
}

// All 所有任务
func All() map[OP]Exec {
	r := map[OP]Exec{}
	Backgrounds.Range(func(key, val interface{}) bool {
		r[OP(param.AsString(key))] = val.(Exec)
		return true
	})
	return r
}

// ListBy 获取某个操作的所有任务
func ListBy(op OP) Exec {
	old, exists := Backgrounds.Load(op)
	if !exists {
		return nil
	}
	exec := old.(Exec)
	return exec
}

// Cancel 取消执行
func Cancel(op OP, cacheKey string) error {
	exec := ListBy(op)
	if exec == nil {
		return nil
	}
	exec.Cancel(cacheKey)
	Backgrounds.Store(op, exec)
	return nil
}

// OP 操作类型
type OP string

func (t OP) String() string {
	return string(t)
}

const (
	// OpExport 导出操作
	OpExport OP = `export`
	// OpImport 导入操作
	OpImport OP = `import`
)

// BGExec 后台执行信息
type BGExec struct {
	ctx     context.Context
	cancel  context.CancelFunc
	Options echo.H
	Started time.Time
	Procs   *FileInfos
}

// Context 暂存上下文信息
func (b *BGExec) Context() context.Context {
	return b.ctx
}

// Cancel 取消执行
func (b *BGExec) Cancel() context.CancelFunc {
	return b.cancel
}

// AddFileInfo 添加文件信息
func (b *BGExec) AddFileInfo(fi *FileInfo) {
	*(b.Procs) = append(*b.Procs, fi)
}

// NewGBExec 新建后台执行信息
func NewGBExec(c context.Context, opt echo.H) *BGExec {
	if c == nil {
		c = context.Background()
	}
	ctx, cancel := context.WithCancel(c)
	return &BGExec{
		ctx:     ctx,
		cancel:  cancel,
		Options: opt,
		Started: time.Now(),
		Procs:   &FileInfos{},
	}
}

// FileInfos 文件信息集合
type FileInfos []*FileInfo

// FileInfo 文件信息
type FileInfo struct {
	Start      time.Time
	End        time.Time
	Elapsed    time.Duration
	Path       string
	Size       int64
	Compressed bool
	Error      string `json:",omitempty"`
}
