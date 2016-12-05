package handler

import (
	"github.com/admpub/log"
	"github.com/admpub/sockjs-go/sockjs"
	"github.com/admpub/websocket"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"github.com/shirou/gopsutil/process"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
)

func ManageSysCmd(ctx echo.Context) error {
	var err error
	return ctx.Render(`manage/execmd`, err)
}

func ManageSockJSSendCmd(c sockjs.Session) error {
	send := make(chan string)
	//push(writer)
	go func() {
		for {
			message := <-send
			log.Info(`Push message: `, message)
			if err := c.Send(message); err != nil {
				log.Error(`Push error: `, err.Error())
				return
			}
		}
	}()

	//echo
	var execute = func(session sockjs.Session) error {
		for {
			command, err := session.Recv()
			if err != nil {
				return err
			}
			if len(command) > 0 {
				com.RunCmdStr(command, func(b []byte) error {
					send <- string(b)
					return nil
				})
			}
			err = session.Send(command)
			if err != nil {
				return err
			}
		}
	}
	err := execute(c)
	if err != nil {
		log.Error(err)
	}
	return nil
}

func ManageWSSendCmd(c *websocket.Conn, ctx echo.Context) error {
	send := make(chan string)
	//push(writer)
	go func() {
		for {
			message := <-send
			log.Info(`Push message: `, message)
			if err := c.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
				log.Error(`Push error: `, err.Error())
				return
			}
		}
	}()

	//echo
	var execute = func(conn *websocket.Conn) error {
		for {
			mt, message, err := conn.ReadMessage()
			if err != nil {
				return err
			}
			command := string(message)
			if len(command) > 0 {
				com.RunCmdStr(command, func(b []byte) error {
					send <- string(b)
					return nil
				})
			}

			if err = conn.WriteMessage(mt, message); err != nil {
				return err
			}
		}
	}
	err := execute(c)
	if err != nil {
		log.Error(err)
	}
	return nil
}

var data *echo.Data

func ManageSysInfo(ctx echo.Context) error {
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
