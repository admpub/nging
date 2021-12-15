package system

import (
	"context"

	"github.com/shirou/gopsutil/v3/process"
	"github.com/webx-top/com"
)

type Process struct {
	Name       string  `json:"name"`
	Pid        int32   `json:"pid"`
	Ppid       int32   `json:"ppid"`
	CPUPercent float64 `json:"cpuPercent"`
	MemPercent float32 `json:"memPercent"`
	//Running    bool    `json:"running"`
	CreateTime string `json:"createTime"`
	created    int64
	Exe        string   `json:"exe"`
	Cmdline    string   `json:"cmdline"`
	Cwd        string   `json:"cwd"`
	Status     []string `json:"status"`
	Username   string   `json:"username"`
	NumThreads int32    `json:"numThreads"`
	NumFDs     int32    `json:"numFDs"`
}

func (p *Process) Parse(ctx context.Context, proc *process.Process) *Process {
	p.Pid = proc.Pid
	p.CPUPercent, _ = proc.CPUPercentWithContext(ctx)
	//p.Running, _ = proc.IsRunningWithContext(ctx)
	p.created, _ = proc.CreateTimeWithContext(ctx)
	if p.created > 0 {
		p.CreateTime = com.DateFormat(`Y-m-d H:i:s`, p.created/1000)
	}
	p.MemPercent, _ = proc.MemoryPercentWithContext(ctx)
	p.Ppid, _ = proc.PpidWithContext(ctx)
	p.Name, _ = proc.NameWithContext(ctx)
	p.Exe, _ = proc.ExeWithContext(ctx)
	p.Cmdline, _ = proc.CmdlineWithContext(ctx)
	p.Cwd, _ = proc.CwdWithContext(ctx)
	p.Status, _ = proc.StatusWithContext(ctx)
	p.Username, _ = proc.UsernameWithContext(ctx)
	p.NumThreads, _ = proc.NumThreadsWithContext(ctx)
	p.NumFDs, _ = proc.NumFDsWithContext(ctx)
	return p
}

type processAndIndex struct {
	index int
	proc  *process.Process
}

func ProcessList(ctx context.Context) ([]*Process, error) {
	list, err := process.ProcessesWithContext(ctx)
	if err != nil {
		return nil, err
	}
	processes := make([]*Process, len(list))
	exec := func(idx int, proc *process.Process) {
		p := &Process{}
		processes[idx] = p.Parse(ctx, proc)
	}
	for idx, proc := range list {
		exec(idx, proc)
	}
	return processes, nil
}
