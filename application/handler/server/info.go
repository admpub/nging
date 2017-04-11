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
package server

import (
	"runtime"
	"strings"

	"time"

	"github.com/admpub/log"
	"github.com/admpub/nging/application/handler"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"github.com/shirou/gopsutil/process"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/middleware/render"
)

func init() {
	handler.RegisterToGroup(`/manage`, func(g *echo.Group) {
		g.Route("GET", `/sysinfo`, Info, render.AutoOutput(nil))
		g.Route("GET", `/netstat`, Connections, render.AutoOutput(nil))
		g.Route("GET", `/process/:pid`, ProcessInfo)
		g.Route("GET", `/procskill/:pid`, ProcessKill)
	})
}

var data *echo.Data

func Info(ctx echo.Context) error {
	ctx.Set(`tmpl`, `manage/sysinfo`)
	if data != nil {
		ctx.Set(`data`, data)
		return nil
	}
	var err error
	cpuInfo, err := cpu.Info()
	if err != nil {
		log.Error(err)
	}
	partitions, err := disk.Partitions(true)
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
	netIOCounter, err := net.IOCounters(false)
	if err != nil {
		log.Error(err)
	}
	/*
		pids, err := process.Pids()
		if err != nil {
			log.Error(err)
		}
		procses := []*process.Process{}
		for _, pid := range pids {
			procs, err := process.NewProcess(pid)
			if err != nil {
				log.Error(err)
			}
			procses = append(procses, procs)
		}
		//*/
	info := &SystemInformation{
		CPU:        cpuInfo,
		Partitions: partitions,
		//DiskIO:         ioCounter,
		Host: hostInfo,
		//Load:       avgLoad,
		Memory: &MemoryInformation{Virtual: virtualMem, Swap: swapMem},
		NetIO:  netIOCounter,
		/*
			PIDs:       pids,
			Process:    procses,
		//*/
	}
	info.DiskUsages = make([]*disk.UsageStat, len(info.Partitions))
	for k, v := range info.Partitions {
		usageStat, err := disk.Usage(v.Mountpoint)
		if err != nil {
			log.Error(err)
		}
		info.DiskUsages[k] = usageStat
	}
	data := ctx.NewData().SetData(info).SetCode(1)

	ctx.Set(`data`, data)
	return nil
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
	data := ctx.NewData()
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
	data := ctx.NewData()
	if err != nil {
		data.SetError(err)
	} else {
		data.SetData(nil)
	}
	return ctx.JSON(data)
}

func Connections(ctx echo.Context) (err error) {
	ctx.Set(`tmpl`, `manage/netstat`)
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
	return
}

type SystemInformation struct {
	CPU        []cpu.InfoStat
	Partitions []disk.PartitionStat
	DiskUsages []*disk.UsageStat
	DiskIO     map[string]disk.IOCountersStat
	Host       *host.InfoStat
	Load       *load.AvgStat
	Memory     *MemoryInformation
	NetIO      []net.IOCountersStat
	PIDs       []int32
	Process    []*process.Process
}

type MemoryInformation struct {
	Virtual *mem.VirtualMemoryStat
	Swap    *mem.SwapMemoryStat
}
