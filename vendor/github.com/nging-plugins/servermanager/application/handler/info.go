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

package handler

import (
	"context"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
	"github.com/shirou/gopsutil/v3/process"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/param"

	"github.com/admpub/log"
	"github.com/admpub/nging/v4/application/handler"
	"github.com/nging-plugins/servermanager/application/library/system"
)

func Info(ctx echo.Context) error {
	var err error
	cpuInfo, err := cpu.Info()
	if err != nil {
		log.Error(err)
	}
	partitions, err := disk.Partitions(false)
	if err != nil {
		log.Error(err)
	}
	/*
		ioCounter, err := disk.IOCounters()
		if err != nil {
			log.Error(err)
		}
	*/
	hostInfo, err := host.Info()
	if err != nil {
		log.Error(err)
	}
	/*
		avgLoad, err := load.Avg()
		if err != nil {
			log.Error(err)
		}
	*/
	virtualMem, err := mem.VirtualMemory()
	if err != nil {
		log.Error(err)
	}
	swapMem, err := mem.SwapMemory()
	if err != nil {
		log.Error(err)
	}
	if swapMem.UsedPercent == 0 {
		swapMem.UsedPercent = (float64(swapMem.Used) / float64(swapMem.Total)) * 100
	}
	netIOCounter, err := net.IOCounters(false)
	if err != nil {
		log.Error(err)
	}
	cpuPercent, err := cpu.Percent(0, true)
	if err != nil {
		log.Error(err)
	}
	info := &system.SystemInformation{
		CPU:        cpuInfo,
		CPUPercent: cpuPercent,
		Partitions: partitions,
		//DiskIO:         ioCounter,
		Host: hostInfo,
		//Load:       avgLoad,
		Memory: &system.MemoryInformation{Virtual: virtualMem, Swap: swapMem},
		NetIO:  netIOCounter,
		Go:     system.Status(),
	}
	info.DiskUsages = make([]*disk.UsageStat, len(info.Partitions))
	for k, v := range info.Partitions {
		usageStat, err := disk.Usage(v.Mountpoint)
		if err != nil {
			log.Error(err)
			continue
		}
		info.DiskUsages[k] = usageStat
	}

	info.Temp, _ = system.SensorsTemperatures()
	ctx.Data().SetData(info, 1)
	return ctx.Render(`server/sysinfo`, nil)
}

func processInfo(ctx context.Context, pid int32) (echo.H, error) {
	procs, err := process.NewProcessWithContext(ctx, pid)
	if err != nil {
		return nil, err
	}
	cpuPercent, _ := procs.CPUPercentWithContext(ctx)
	memPercent, _ := procs.MemoryPercentWithContext(ctx)
	name, _ := procs.NameWithContext(ctx)
	cmdLine, _ := procs.CmdlineWithContext(ctx)
	exe, _ := procs.ExeWithContext(ctx)
	createTime, _ := procs.CreateTimeWithContext(ctx)
	row := echo.H{
		"name":           name,
		"cmd_line":       cmdLine,
		"exe":            exe,
		"created":        "",
		"cpu_percent":    cpuPercent,
		"memory_percent": memPercent,
	}
	if createTime > 0 {
		row["created"] = com.DateFormat(`Y-m-d H:i:s`, createTime/1000)
	}
	return row, nil
}

func ProcessInfo(ctx echo.Context) error {
	pid := ctx.Paramx(`pid`).Int32()
	row, err := processInfo(ctx, pid)
	data := ctx.Data()
	if err != nil {
		data.SetError(err)
	} else {
		data.SetData(row)
	}
	return ctx.JSON(data)
}

func ProcessKill(ctx echo.Context) error {
	pid := ctx.Paramx(`pid`).Int()
	err := com.CloseProcessFromPid(pid)
	data := ctx.Data()
	if err != nil {
		data.SetError(err)
	} else {
		data.SetData(nil)
	}
	return ctx.JSON(data)
}

func Connections(ctx echo.Context) (err error) {
	var conns []net.ConnectionStat
	var kind string
	switch kind {
	case "tcp", "tcp4", "tcp6", "udp", "udp4", "udp6", "unix", "inet", "inet4", "inet6":
	default:
		kind = "all"
	}
	conns, err = net.Connections(kind)
	if err != nil {
		if err.Error() == "not implemented yet" {
			if runtime.GOOS == "windows" {
				err = nil
				var conn <-chan net.ConnectionStat
				if strings.HasPrefix(kind, `udp`) {
					conn, err = NetStatUDP()
				} else {
					conn, err = NetStatTCP()
				}
				if err != nil {
					return
				}
				done := make(chan bool)
				go func() {
					defer func() {
						done <- true
					}()
					for {
						select {
						case c, r := <-conn:
							if !r {
								return
							}
							conns = append(conns, c)
						}
					}
				}()
				<-done
			}
		}
	}
	ctx.Set(`listData`, conns)
	return ctx.Render(`server/netstat`, nil)
}

var processList []*system.Process
var processLock sync.RWMutex
var processLastQueryTime time.Time
var processQuering bool

func ProcessList(ctx echo.Context) error {
	if ctx.Formx(`status`).Bool() {
		processLock.RLock()
		quering := processQuering
		queryTime := processLastQueryTime
		processLock.RUnlock()
		data := echo.H{`quering`: quering}
		if !quering {
			data.Set(`queryTime`, queryTime.Format(param.DateTimeNormal))
		}
		return ctx.JSON(ctx.Data().SetData(data))
	}
	force := ctx.Formx(`force`).Bool()
	processLock.Lock()
	defer processLock.Unlock()
	var err error
	var list []*system.Process
	var isCached bool
	if !processQuering {
		if force || processLastQueryTime.Before(time.Now().Add(-30*time.Minute)) {
			processQuering = true
			go func() {
				stdCtx := context.Background()
				list, err = system.ProcessList(stdCtx)
				if err != nil {
					log.Error(err)
				}
				processLock.Lock()
				processList = list
				processQuering = false
				processLastQueryTime = time.Now()
				processLock.Unlock()
			}()
		} else {
			list = processList
			isCached = true
		}
	}
	if force {
		return ctx.Redirect(handler.URLFor(`/server/processes`))
	}
	switch ctx.Form(`sort`) {
	case `cpu`:
		sortedList := system.ProcessListSortByCPUPercent(list)
		sort.Sort(sortedList)
		ctx.Set(`listData`, sortedList)
	case `-cpu`:
		sortedList := system.ProcessListSortByCPUPercentReverse(list)
		sort.Sort(sortedList)
		ctx.Set(`listData`, sortedList)
	case `mem`:
		sortedList := system.ProcessListSortByMemPercent(list)
		sort.Sort(sortedList)
		ctx.Set(`listData`, sortedList)
	case `-mem`:
		sortedList := system.ProcessListSortByMemPercentReverse(list)
		sort.Sort(sortedList)
		ctx.Set(`listData`, sortedList)
	case `thread`:
		sortedList := system.ProcessListSortByNumThreads(list)
		sort.Sort(sortedList)
		ctx.Set(`listData`, sortedList)
	case `-thread`:
		sortedList := system.ProcessListSortByNumThreadsReverse(list)
		sort.Sort(sortedList)
		ctx.Set(`listData`, sortedList)
	case `fd`:
		sortedList := system.ProcessListSortByNumFDs(list)
		sort.Sort(sortedList)
		ctx.Set(`listData`, sortedList)
	case `-fd`:
		sortedList := system.ProcessListSortByNumFDsReverse(list)
		sort.Sort(sortedList)
		ctx.Set(`listData`, sortedList)
	case `created`:
		sortedList := system.ProcessListSortByCreated(list)
		sort.Sort(sortedList)
		ctx.Set(`listData`, sortedList)
	case `-created`:
		sortedList := system.ProcessListSortByCreatedReverse(list)
		sort.Sort(sortedList)
		ctx.Set(`listData`, sortedList)
	case `-id`:
		sortedList := system.ProcessListSortByPidReverse(list)
		sort.Sort(sortedList)
		ctx.Set(`listData`, sortedList)
	default:
		ctx.Set(`listData`, list)
	}
	ctx.Set(`lastQueryTime`, processLastQueryTime)
	ctx.Set(`isCached`, isCached)
	ctx.Set(`quering`, processQuering)
	ctx.Set(`activeURL`, `/server/sysinfo`)
	if ctx.Formx(`partial`).Bool() {
		data := ctx.Data()
		if err != nil {
			return ctx.JSON(data.SetInfo(err.Error(), 0))
		}
		var partialBytes []byte
		partialBytes, err = ctx.Fetch(`server/processes_list_partial`, nil)
		data.SetData(echo.H{
			`html`: string(partialBytes),
		})
		return ctx.JSON(data)
	}
	return ctx.Render(`server/processes`, handler.Err(ctx, err))
}
