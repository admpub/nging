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
package server

import (
	"runtime"
	"strings"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"github.com/shirou/gopsutil/process"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"

	"github.com/admpub/log"
	collectd "github.com/admpub/logcool/input/collectd"
	"github.com/admpub/nging/application/library/system"
)

var _ = collectd.SystemInfo

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
	cpuPercent, err := cpu.Percent(0, false)
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
	}
	info.DiskUsages = make([]*disk.UsageStat, len(info.Partitions))
	for k, v := range info.Partitions {
		usageStat, err := disk.Usage(v.Mountpoint)
		if err != nil {
			log.Error(err)
		}
		info.DiskUsages[k] = usageStat
	}
	ctx.Data().SetData(info, 1)
	return ctx.Render(`server/sysinfo`, nil)
}

func processList() ([]int32, []*process.Process) {
	pids, err := process.Pids()
	if err != nil {
		log.Error(err)
	}
	processes := []*process.Process{}
	for _, pid := range pids {
		procs, err := process.NewProcess(pid)
		if err != nil {
			log.Error(err)
		}
		processes = append(processes, procs)
	}
	return pids, processes
}

func processInfo(pid int32) (echo.H, error) {
	procs, err := process.NewProcess(pid)
	if err != nil {
		return nil, err
	}
	cpuPercent, _ := procs.Percent(time.Second * 5)
	memPercent, _ := procs.MemoryPercent()
	name, _ := procs.Name()
	cmdLine, _ := procs.Cmdline()
	exe, _ := procs.Exe()
	createTime, _ := procs.CreateTime()
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
	row, err := processInfo(pid)
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
