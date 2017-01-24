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
package handler

import (
	"io"
	"os"
	"runtime"

	"os/exec"

	"strings"

	"time"

	"github.com/admpub/log"
	"github.com/admpub/nging/application/library/charset"
	"github.com/admpub/nging/application/library/config"
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

var (
	WebSocketLogger = log.GetLogger(`websocket`)
	IsWindows       bool
)

func init() {
	WebSocketLogger.SetLevel(`Info`)
	IsWindows = runtime.GOOS == `windows`
}

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
			WebSocketLogger.Debug(`Push message: `, message)
			if err := c.Send(message); err != nil {
				WebSocketLogger.Error(`Push error: `, err.Error())
				return
			}
		}
	}()
	timeout := c.Request().URL.Query().Get(`timeout`)
	duration := config.ParseTimeDuration(timeout)
	//echo
	exec := func(session sockjs.Session) error {
		var (
			w   io.WriteCloser
			cmd *exec.Cmd
		)
		for {
			command, err := session.Recv()
			if err != nil {
				return err
			}
			if len(command) == 0 {
				continue
			}
			if w == nil {
				w, cmd, err = CmdRunner(command, send, func() {
					w.Close()
					w = nil
				}, duration)
				if err != nil {
					return err
				}
				continue
			}
			err = CmdContinue(command, w, cmd)
			if err != nil {
				return err
			}
		}
	}
	err := exec(c)
	if err != nil {
		WebSocketLogger.Error(err)
	}
	return nil
}

func CmdRunner(command string, send chan string, onEnd func(), timeout time.Duration) (w io.WriteCloser, cmd *exec.Cmd, err error) {
	cmd = com.CreateCmdStr(command, func(b []byte) (e error) {
		if IsWindows {
			b, e = charset.Convert(`gbk`, `utf-8`, b)
			if e != nil {
				return e
			}
		}
		send <- string(b)
		return nil
	})
	w, err = cmd.StdinPipe()
	if err != nil {
		return
	}
	done := make(chan struct{})
	go func() {
		log.Info(`[command] running: `, command)
		if e := cmd.Run(); e != nil {
			cmd.Stderr.Write([]byte(e.Error()))
		}
		done <- struct{}{}
		onEnd()
	}()
	cmdMaxTimeout := config.DefaultConfig.Sys.CmdTimeoutDuration
	if timeout <= 0 {
		timeout = time.Minute * 5
	}
	if timeout > cmdMaxTimeout {
		timeout = cmdMaxTimeout
	}
	go func() {
		ticker := time.NewTicker(timeout)
		defer ticker.Stop()
		for {
			select {
			case <-done:
				log.Info(`[command] exited: `, command)
				return
			case <-ticker.C:
				if cmd == nil {
					return
				}
				if cmd.Process == nil {
					return
				}
				cmd.Stderr.Write([]byte(`timeout`))
				log.Info(`[command] timeout: `, command)
				err := cmd.Process.Kill()
				if err != nil {
					log.Error(err)
				}
				return
			}
		}
	}()
	return
}

func CmdContinue(command string, w io.WriteCloser, cmd *exec.Cmd) (err error) {
	switch command {
	case `^C`:
		err = cmd.Process.Signal(os.Interrupt)
		if err != nil {
			if !strings.HasPrefix(err.Error(), `not supported by `) {
				WebSocketLogger.Error(err)
			}
			err = cmd.Process.Kill()
			if err != nil {
				WebSocketLogger.Error(err)
			}
		}
	default:
		w.Write([]byte(command + "\n"))
	}
	return nil
}

func ManageWSSendCmd(c *websocket.Conn, ctx echo.Context) error {
	send := make(chan string)
	//push(writer)
	go func() {
		for {
			message := <-send
			WebSocketLogger.Debug(`Push message: `, message)
			if err := c.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
				WebSocketLogger.Error(`Push error: `, err.Error())
				return
			}
		}
	}()

	timeout := ctx.Query(`timeout`)
	duration := config.ParseTimeDuration(timeout)
	//echo
	exec := func(conn *websocket.Conn) error {
		var (
			w   io.WriteCloser
			cmd *exec.Cmd
		)
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				return err
			}
			command := string(message)
			if len(command) == 0 {
				continue
			}
			if w == nil {
				w, cmd, err = CmdRunner(command, send, func() {
					w.Close()
					w = nil
				}, duration)
				if err != nil {
					return err
				}
				continue
			}
			err = CmdContinue(command, w, cmd)
			if err != nil {
				return err
			}
		}
	}
	err := exec(c)
	if err != nil {
		WebSocketLogger.Error(err)
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
