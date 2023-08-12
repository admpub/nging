//go:build !windows

package goforever

import (
	"errors"
	"fmt"
	"os"
	"os/user"
	"strconv"
	"syscall"
)

func (p *Process) StartProcess(name string, argv []string, attr *os.ProcAttr) (Processer, error) {
	if len(p.User) > 0 {
		attr.Sys = &syscall.SysProcAttr{}
		userInfo, err := user.Lookup(p.User)
		if err != nil {
			return nil, errors.New("failed to get user: " + err.Error())
		}
		uid, err := strconv.ParseUint(userInfo.Uid, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("failed to ParseUint(userInfo.Uid=%q): %w", userInfo.Uid, err)
		}
		gid, err := strconv.ParseUint(userInfo.Gid, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("failed to ParseUint(userInfo.Gid=%q): %w", userInfo.Gid, err)
		}
		attr.Sys.Credential = &syscall.Credential{
			Uid:         uint32(uid),
			Gid:         uint32(gid),
			NoSetGroups: true,
		}
	}
	process, err := os.StartProcess(name, argv, attr)
	if err != nil {
		return nil, err
	}
	return &osProcess{Process: process}, nil
}

type osProcess struct {
	*os.Process
}

func (p *osProcess) Pid() int {
	return p.Process.Pid
}

func buildOption(options map[string]interface{}) map[string]interface{} {
	return options
}

func SetOption(options map[string]interface{}, name string, value interface{}) map[string]interface{} {
	return options
}
