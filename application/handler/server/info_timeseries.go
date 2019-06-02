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
	"context"
	"time"

	"github.com/admpub/log"
	"github.com/shirou/gopsutil/net"
)

var (
	realTimeStatus                 *RealTimeStatus
	CancelRealTimeStatusCollection context.CancelFunc
)

func ListenRealTimeStatus() {
	if realTimeStatus == nil {
		realTimeStatus = NewRealTimeStatus(time.Second*2, 80)
	}
	ctx, cancel := context.WithCancel(context.Background())
	go realTimeStatus.Listen(ctx)
	CancelRealTimeStatusCollection = cancel
}

func NewRealTimeStatus(interval time.Duration, maxSize int) *RealTimeStatus {
	return &RealTimeStatus{
		max:      maxSize,
		interval: interval,
		CPU:      TimeSeries{},
		Mem:      TimeSeries{},
		Net:      NewNetIOTimeSeries(),
	}
}

func NewNetIOTimeSeries() NetIOTimeSeries {
	return NetIOTimeSeries{
		BytesSent:   TimeSeries{},
		BytesRecv:   TimeSeries{},
		PacketsSent: TimeSeries{},
		PacketsRecv: TimeSeries{},
	}
}

type NetIOTimeSeries struct {
	BytesSent   TimeSeries
	BytesRecv   TimeSeries
	PacketsSent TimeSeries
	PacketsRecv TimeSeries
}

type RealTimeStatus struct {
	max      int
	interval time.Duration
	CPU      TimeSeries
	Mem      TimeSeries
	Net      NetIOTimeSeries
}

func (r *RealTimeStatus) Listen(ctx context.Context) *RealTimeStatus {
	info := &DynamicInformation{}
	t := time.NewTicker(r.interval)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			log.Info(`Exit server real-time status collection`)
			return r
		case <-t.C:
			info.NetMemoryCPU()
			if len(info.CPUPercent) > 0 {
				r.CPUAdd(info.CPUPercent[0])
			} else {
				r.CPUAdd(0)
			}
			r.MemAdd(info.Memory.Virtual.UsedPercent)
			//r.NetAdd(info.NetIO[0])
			//log.Info(`Collect server status`)
		}
	}
	return r
}

func (r *RealTimeStatus) CPUAdd(y float64) *RealTimeStatus {
	if r.max <= 0 {
		return r
	}
	l := len(r.CPU)
	if l >= r.max {
		r.CPU = r.CPU[1+l-r.max:]
	}
	r.CPU = append(r.CPU, NewXY(y))
	return r
}

func (r *RealTimeStatus) MemAdd(y float64) *RealTimeStatus {
	if r.max <= 0 {
		return r
	}
	l := len(r.Mem)
	if l >= r.max {
		r.Mem = r.Mem[1+l-r.max:]
	}
	r.Mem = append(r.Mem, NewXY(y))
	return r
}

func (r *RealTimeStatus) NetAdd(stat net.IOCountersStat) *RealTimeStatus {
	if r.max <= 0 {
		return r
	}
	l := len(r.Net.BytesRecv)
	if l >= r.max {
		r.Net.BytesRecv = r.Net.BytesRecv[1+l-r.max:]
	}
	r.Net.BytesRecv = append(r.Net.BytesRecv, NewXY(float64(stat.BytesRecv)))

	l = len(r.Net.BytesSent)
	if l >= r.max {
		r.Net.BytesSent = r.Net.BytesSent[1+l-r.max:]
	}
	r.Net.BytesSent = append(r.Net.BytesSent, NewXY(float64(stat.BytesSent)))

	l = len(r.Net.PacketsRecv)
	if l >= r.max {
		r.Net.PacketsRecv = r.Net.PacketsRecv[1+l-r.max:]
	}
	r.Net.PacketsRecv = append(r.Net.PacketsRecv, NewXY(float64(stat.PacketsRecv)))

	l = len(r.Net.PacketsSent)
	if l >= r.max {
		r.Net.PacketsSent = r.Net.PacketsSent[1+l-r.max:]
	}
	r.Net.PacketsSent = append(r.Net.PacketsSent, NewXY(float64(stat.PacketsSent)))
	return r
}

type (
	TimeSeries []XY
	XY         [2]interface{}
)

func NewXY(y float64) XY {
	x := time.Now().UnixNano() / 1e6 //毫秒
	return XY{x, y}
}
