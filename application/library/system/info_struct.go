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
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
)

type SystemInformation struct {
	CPU        []cpu.InfoStat                 `json:",omitempty"`
	CPUPercent []float64                      `json:",omitempty"`
	Partitions []disk.PartitionStat           `json:",omitempty"`
	DiskUsages []*disk.UsageStat              `json:",omitempty"`
	DiskIO     map[string]disk.IOCountersStat `json:",omitempty"`
	Host       *host.InfoStat                 `json:",omitempty"`
	Load       *load.AvgStat                  `json:",omitempty"`
	Memory     *MemoryInformation             `json:",omitempty"`
	NetIO      []net.IOCountersStat           `json:",omitempty"`
	Temp       []host.TemperatureStat         `json:",omitempty"`
	Go         *RuntimeStatus                 `json:",omitempty"`
}

type MemoryInformation struct {
	Virtual *mem.VirtualMemoryStat `json:",omitempty"`
	Swap    *mem.SwapMemoryStat    `json:",omitempty"`
}

type DynamicInformation struct {
	CPUPercent []float64
	Load       *load.AvgStat          `json:",omitempty"`
	Memory     *MemoryInformation     `json:",omitempty"`
	NetIO      []net.IOCountersStat   `json:",omitempty"`
	Temp       []host.TemperatureStat `json:",omitempty"`
}

func (d *DynamicInformation) Init() *DynamicInformation {
	d.NetMemoryCPU()
	d.TemperatureStat()
	d.Load, _ = load.Avg()
	return d
}

func (d *DynamicInformation) NetMemoryCPU() *DynamicInformation {
	d.MemoryAndCPU()
	d.NetIO, _ = net.IOCounters(false)
	return d
}

func (d *DynamicInformation) TemperatureStat() *DynamicInformation {
	d.Temp, _ = SensorsTemperatures()
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
