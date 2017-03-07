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
	"github.com/admpub/log"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"github.com/shirou/gopsutil/process"
	"github.com/webx-top/echo"
)

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
