/*
   Nging is a toolbox for webmasters
   Copyright (C) 2018-present  Wenhui Shen <swh@admpub.com>

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

package system

import (
	"fmt"
	"runtime"
	"time"
)

// RuntimeStatus 运行时信息
type RuntimeStatus struct {
	NumGoroutine int
	// General statistics.
	MemAllocated uint64 // bytes allocated and still in use
	MemTotal     uint64 // bytes allocated (even if freed)
	MemSys       uint64 // bytes obtained from system (sum of XxxSys below)
	Lookups      uint64 // number of pointer lookups
	MemMallocs   uint64 // number of mallocs
	MemFrees     uint64 // number of frees
	// Main allocation heap statistics.
	HeapAlloc    uint64 // bytes allocated and still in use
	HeapSys      uint64 // bytes obtained from system
	HeapIdle     uint64 // bytes in idle spans
	HeapInuse    uint64 // bytes in non-idle span
	HeapReleased uint64 // bytes released to the OS
	HeapObjects  uint64 // total number of allocated objects
	// Low-level fixed-size structure allocator statistics.
	// Inuse is bytes used now.
	// Sys is bytes obtained from system.
	StackInuse  uint64 // bootstrap stacks
	StackSys    uint64
	MSpanInuse  uint64 // mspan structures
	MSpanSys    uint64
	MCacheInuse uint64 // mcache structures
	MCacheSys   uint64
	BuckHashSys uint64 // profiling bucket hash table
	GCSys       uint64 // GC metadata
	OtherSys    uint64 // other system allocations
	// Garbage collector statistics.
	NextGC       uint64 // next run in HeapAlloc time (bytes)
	LastGC       uint64 // last run in absolute time (ns)
	PauseTotalNs uint64
	PauseNs      string // circular buffer of recent GC pause times, most recent at [(NumGC+255)%256]
	NumGC        uint32
}

// LastGCString LastGC
func (r *RuntimeStatus) LastGCString() string {
	return fmt.Sprintf("%.1fs ago", float64(time.Now().UnixNano()-int64(r.LastGC))/1000/1000/1000)
}

// PauseTotalNsString PauseTotalNs
func (r *RuntimeStatus) PauseTotalNsString() string {
	return fmt.Sprintf("%.1fs", float64(r.PauseTotalNs)/1000/1000/1000)
}

// Status go运行时信息
func Status() *RuntimeStatus {
	m := new(runtime.MemStats)
	runtime.ReadMemStats(m)
	ms := &RuntimeStatus{}
	ms.NumGoroutine = runtime.NumGoroutine()
	ms.MemAllocated = m.Alloc
	ms.MemTotal = m.TotalAlloc
	ms.MemSys = m.Sys
	ms.Lookups = m.Lookups
	ms.MemMallocs = m.Mallocs
	ms.MemFrees = m.Frees
	ms.HeapAlloc = m.HeapAlloc
	ms.HeapSys = m.HeapSys
	ms.HeapIdle = m.HeapIdle
	ms.HeapInuse = m.HeapInuse
	ms.HeapReleased = m.HeapReleased
	ms.HeapObjects = m.HeapObjects
	ms.StackInuse = m.StackInuse
	ms.StackSys = m.StackSys
	ms.MSpanInuse = m.MSpanInuse
	ms.MSpanSys = m.MSpanSys
	ms.MCacheInuse = m.MCacheInuse
	ms.MCacheSys = m.MCacheSys
	ms.BuckHashSys = m.BuckHashSys
	ms.GCSys = m.GCSys
	ms.OtherSys = m.OtherSys
	ms.NextGC = m.NextGC
	ms.LastGC = m.LastGC
	ms.PauseTotalNs = m.PauseTotalNs
	ms.PauseNs = fmt.Sprintf("%.3fs", float64(m.PauseNs[(m.NumGC+255)%256])/1000/1000/1000)
	ms.NumGC = m.NumGC
	return ms
}
