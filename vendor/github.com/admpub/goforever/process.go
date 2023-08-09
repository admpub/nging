// goforever - processes management
// Copyright (c) 2013 Garrett Woodworth (https://github.com/gwoo).

package goforever

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"

	"github.com/webx-top/com"

	ps "github.com/admpub/go-ps"
)

var ping = "1m"

// RunProcess Run the process
func RunProcess(name string, p *Process) chan *Process {
	ch := make(chan *Process)
	go func() {
		proc, _, _ := p.Find()
		// proc, err := ps.FindProcess(p.Pid)
		if proc == nil {
			p.Start(name)
		}
		p.ping(ping, func(time time.Duration, p *Process) {
			if p.pid > 0 {
				p.respawns = 0
				log.Println(p.logPrefix()+"refreshed after", time)
				p.status = StatusRunning
				p.RunHook(p.status)
			}
		})
		go p.watch()
		ch <- p
	}()
	return ch
}

const (
	StatusStarted   = `started`
	StatusRunning   = `running`
	StatusStopped   = `stopped`
	StatusRestarted = `restarted`
	StatusExited    = `exited`
	StatusKilled    = `killed`
)

type Process struct {
	Name    string
	Command string
	Env     []string
	Dir     string
	Args    []string
	User    string
	Pidfile Pidfile
	Logfile string
	Errfile string
	Respawn int
	Delay   string
	Ping    string
	Debug   bool

	// Options
	// HideWindow bool // for windows
	Options map[string]interface{} `json:",omitempty" xml:",omitempty"`

	pid      int
	status   string
	x        *os.Process
	respawns int
	children Children
	hooks    map[string][]func(procs *Process)
	cleanup  func()
	err      error
}

func (p *Process) Init() {
	p.SetChildren(map[string]*Process{})
	p.Options = buildOption(p.Options)
}

func (p *Process) SetChildren(children map[string]*Process) {
	p.children = children
}

func (p *Process) SetChildKV(name string, proc *Process) {
	p.children[name] = proc
}

func (p *Process) Status() string {
	return p.status
}

func (p *Process) Pid() int {
	return p.pid
}

func (p *Process) Reset() error {
	log.Println(p.Stop())
	p.hooks = make(map[string][]func(procs *Process), 0)
	p.x = nil
	p.respawns = 0
	return p.err
}

func (p *Process) Error() error {
	return p.err
}

func (p *Process) String() string {
	js, err := json.Marshal(p)
	if err != nil {
		log.Println(p.logPrefix(), err)
		return ""
	}
	return string(js)
}

func (p *Process) RunHook(status string) {
	if p.hooks == nil {
		return
	}
	if fnList, ok := p.hooks[status]; ok {
		for _, f := range fnList {
			f(p)
		}
	}
}

func (p *Process) SetHook(status string, hooks ...func(procs *Process)) *Process {
	if p.hooks == nil {
		p.hooks = map[string][]func(procs *Process){}
	}
	p.hooks[status] = hooks
	return p
}

func (p *Process) AddHook(status string, hooks ...func(procs *Process)) *Process {
	if p.hooks == nil {
		p.hooks = map[string][]func(procs *Process){}
	}
	if _, ok := p.hooks[status]; !ok {
		p.hooks[status] = []func(procs *Process){}
	}
	p.hooks[status] = append(p.hooks[status], hooks...)
	return p
}

var ErrPidfileEmpty = errors.New("Pidfile is empty")
var ErrCouldNotFindProcess = errors.New("could not find process")

// Find a process by name
func (p *Process) Find() (*os.Process, string, error) {
	if len(p.Pidfile) == 0 {
		return nil, "", fmt.Errorf(p.logPrefix()+"%w", ErrPidfileEmpty)
	}
	if pid := p.Pidfile.Read(); pid > 0 {
		proc, err := ps.FindProcess(pid)
		if err != nil || proc == nil {
			return nil, "", err
		}
		process, err := os.FindProcess(pid)
		if err != nil {
			return nil, "", err
		}
		p.x = process
		p.pid = process.Pid
		p.status = StatusRunning
		p.RunHook(p.status)
		message := fmt.Sprintf(p.logPrefix()+"%s is %#v", p.Name, process.Pid)
		return process, message, nil
	}
	message := fmt.Sprintf(p.logPrefix()+"%s not running.", p.Name)
	return nil, message, fmt.Errorf("%w %s", ErrCouldNotFindProcess, p.Name)
}

// Start the process
func (p *Process) Start(name string) string {
	p.Name = name
	p.err = nil
	logPrefix := p.logPrefix()
	if p.Debug {
		log.Println(logPrefix+`Dir:`, p.Dir)
	}
	files := []*os.File{
		os.Stdin,
		os.Stdout,
		os.Stderr,
	}
	if len(p.Logfile) > 0 {
		logDir := filepath.Dir(p.Logfile)
		os.MkdirAll(logDir, os.ModePerm)
		files[1] = NewLog(p.Logfile)
	}
	if len(p.Errfile) > 0 {
		logDir := filepath.Dir(p.Errfile)
		os.MkdirAll(logDir, os.ModePerm)
		files[2] = NewLog(p.Errfile)
	}
	proc := &os.ProcAttr{
		Dir:   p.Dir,
		Env:   append(os.Environ()[:], p.Env...),
		Files: files,
	}
	if len(p.User) > 0 {
		proc.Sys = &syscall.SysProcAttr{}
		cleanup, err := SetSysProcAttr(proc.Sys, p.User, p.Options)
		if err != nil {
			p.err = errors.New(logPrefix + err.Error())
			log.Println(p.err.Error())
			return ""
		}
		p.cleanup = cleanup
	}
	args := com.ParseArgs(p.Command)
	args = append(args, p.Args...)
	if filepath.Base(args[0]) == args[0] {
		if lp, err := exec.LookPath(args[0]); err != nil {
			p.err = err
			log.Println(logPrefix+"LookPath:", err.Error())
		} else {
			args[0] = lp
		}
	}
	if p.Debug {
		b, _ := json.MarshalIndent(args, ``, `  `)
		log.Println(logPrefix+"Args:", string(b))
		b, _ = json.MarshalIndent(proc, ``, `  `)
		log.Println(logPrefix+"Attr:", string(b))
	}
	process, err := os.StartProcess(args[0], args, proc)
	if err != nil {
		p.err = errors.New(logPrefix + "failed. " + err.Error())
		log.Println(p.err.Error())
		return ""
	}
	err = p.Pidfile.Write(process.Pid)
	if err != nil {
		log.Printf(logPrefix+"pidfile error: %v", err)
		return ""
	}
	p.x = process
	p.pid = process.Pid
	p.status = StatusStarted
	p.RunHook(p.status)
	return fmt.Sprintf(logPrefix+"%s is %#v", p.Name, process.Pid)
}

func (p *Process) logPrefix() string {
	return `[Process:` + p.Name + `][` + time.Now().Format(time.RFC3339) + `]`
}

// Stop the process
func (p *Process) Stop() string {
	p.err = nil
	logPrefix := p.logPrefix()
	if p.x != nil {
		// Initial code has the following comment: "p.x.Kill() this seems to cause trouble"
		// I want this to work on windows where AFAIK the existing code was not portable
		if err := p.x.Kill(); err != nil { //err := syscall.Kill(p.x.Pid, syscall.SIGTERM)
			p.err = errors.New(logPrefix + err.Error())
			log.Println(p.err.Error())
		} else {
			log.Println(logPrefix + "Stop command seemed to work")
		}
		p.children.Stop()
	}
	p.release(StatusStopped)
	message := fmt.Sprintf(logPrefix + "stopped.")
	return message
}

// Release process and remove pidfile
func (p *Process) release(status string) {
	// debug.PrintStack()
	if p.x != nil {
		p.x.Release()
	}
	p.pid = 0
	// 去掉删除pid文件的动作，用于goforever进程重启后继续监控，防止启动重复进程
	//p.Pidfile.Delete()
	p.status = status
	p.RunHook(p.status)
	if p.cleanup != nil {
		p.cleanup()
	}
}

// Restart the process
func (p *Process) Restart() (chan *Process, string) {
	p.Stop()
	message := p.logPrefix() + "restarted."
	ch := RunProcess(p.Name, p)
	return ch, message
}

// Run callback on the process after given duration.
func (p *Process) ping(duration string, f func(t time.Duration, p *Process)) {
	if len(p.Ping) > 0 {
		duration = p.Ping
	}
	t, err := time.ParseDuration(duration)
	if err != nil {
		t, _ = time.ParseDuration(ping)
	}
	go func() {
		ticker := time.NewTicker(t)
		defer ticker.Stop()
		select {
		case <-ticker.C:
			f(t, p)
		}
	}()
}

// Watch the process
func (p *Process) watch() {
	if p.x == nil {
		p.release(StatusStopped)
		return
	}
	status := make(chan *os.ProcessState)
	died := make(chan error)
	go func() {
		var err error
		// state, err := p.x.Wait()
		proc, _ := ps.FindProcess(p.pid)
		var ppid int
		var state = &os.ProcessState{}
		if proc != nil {
			ppid = proc.PPid()
		}
		// 如果是当前进程fork的子进程，则阻塞等待获取子进程状态，否则循环检测进程状态（1s一次，直到状态变更）
		if ppid == os.Getpid() {
			state, err = p.x.Wait()
		} else {
			for {
				time.Sleep(1 * time.Second)
				proc, err = ps.FindProcess(p.pid)
				if err != nil || proc == nil {
					break
				}
			}
		}
		if err != nil {
			died <- err
			return
		}
		status <- state
	}()
	select {
	case s := <-status:
		if p.status == StatusStopped {
			p.RunHook(p.status)
			return
		}
		logPrefix := p.logPrefix() + ` `
		fmt.Fprintf(os.Stderr, logPrefix+"%s\n", s)
		fmt.Fprintf(os.Stderr, logPrefix+"success = %#v\n", s.Success())
		fmt.Fprintf(os.Stderr, logPrefix+"exited = %#v\n", s.Exited())
		p.respawns++
		if p.respawns > p.Respawn {
			p.release(StatusExited)
			log.Println(logPrefix + "respawn limit reached.")
			return
		}
		fmt.Fprintf(os.Stderr, logPrefix+"respawns = %#v\n", p.respawns)
		if len(p.Delay) > 0 {
			t, _ := time.ParseDuration(p.Delay)
			time.Sleep(t)
		}
		p.Restart()
		p.status = StatusRestarted
		p.RunHook(p.status)
	case err := <-died:
		p.release(StatusKilled)
		log.Printf(p.logPrefix()+"%d %s killed = %#v\n", p.x.Pid, p.Name, err)
	}
}

// Run child processes
func (p *Process) Run() {
	for name, p := range p.children {
		RunProcess(name, p)
	}
}

func (p *Process) StartChild(name string) (*Process, error) {
	cp := p.Child(name)
	if cp == nil {
		return nil, fmt.Errorf("%s does not exist", name)
	}
	cpp, _, err := cp.Find()
	if err != nil && !errors.Is(err, ErrCouldNotFindProcess) {
		return nil, err
	}
	if cpp != nil {
		return nil, fmt.Errorf("%s already running", name)
	}
	procs := <-RunProcess(name, cp)
	return procs, nil
}

func (p *Process) RestartChild(name string) (*Process, error) {
	cp := p.Child(name)
	if p == nil {
		return nil, fmt.Errorf("%s does not exist", name)
	}
	cp.Find()
	ch, _ := cp.Restart()
	procs := <-ch
	return procs, nil
}

func (p *Process) StopChild(name string) error {
	cp := p.Child(name)
	if cp == nil {
		return fmt.Errorf("%s does not exist", name)
	}
	cp.Find()
	cp.Stop()
	return nil
}

func (p *Process) Child(name string) *Process {
	return p.children.Get(name)
}

func (p *Process) Add(name string, procs *Process, run ...bool) *Process {
	p.StopChild(name)
	p.children[name] = procs
	if len(run) > 0 && run[0] {
		RunProcess(name, procs)
	}
	return p
}

func (p *Process) ChildKeys() []string {
	return p.children.Keys()
}
