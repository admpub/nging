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
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
)

type SystemInformation struct {
	CPU        []cpu.InfoStat
	CPUPercent []float64
	Partitions []disk.PartitionStat
	DiskUsages []*disk.UsageStat
	DiskIO     map[string]disk.IOCountersStat
	Host       *host.InfoStat
	Load       *load.AvgStat
	Memory     *MemoryInformation
	NetIO      []net.IOCountersStat
}

type MemoryInformation struct {
	Virtual *mem.VirtualMemoryStat
	Swap    *mem.SwapMemoryStat
}

type DynamicInformation struct {
	CPUPercent []float64
	Load       *load.AvgStat `json:",omitempty"`
	Memory     *MemoryInformation
	NetIO      []net.IOCountersStat `json:",omitempty"`
}

func (d *DynamicInformation) Init() *DynamicInformation {
	d.NetMemoryCPU()
	d.Load, _ = load.Avg()
	return d
}

func (d *DynamicInformation) NetMemoryCPU() *DynamicInformation {
	d.MemoryAndCPU()
	d.NetIO, _ = net.IOCounters(false)
	return d
}

func (d *DynamicInformation) MemoryAndCPU() *DynamicInformation {
	d.Memory = &MemoryInformation{}
	d.Memory.Virtual, _ = mem.VirtualMemory()
	d.Memory.Swap, _ = mem.SwapMemory()
	if d.Memory.Swap != nil && d.Memory.Swap.UsedPercent == 0 {
		d.Memory.Swap.UsedPercent = (float64(d.Memory.Swap.Used) / float64(d.Memory.Swap.Total)) * 100
	}
	d.CPUPercent, _ = cpu.Percent(0, false)
	return d
}
