//go:build windows

package goforever

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/shirou/gopsutil/process"
	"github.com/webx-top/com"
)

var debug bool

func buildOption(options map[string]interface{}) map[string]interface{} {
	if options == nil {
		options = map[string]interface{}{}
	}
	options[`HideWindow`] = false
	options[`Password`] = ``
	return options
}

func SetOption(options map[string]interface{}, name string, value interface{}) map[string]interface{} {
	name := com.PascalCase(name)
	switch name {
	case `HideWindow`:
		options[name] = com.Bool(value)
	case `Password`:
		options[name] = com.Str(value)
	}
	return options
}

func SetSysProcAttr(attr *syscall.SysProcAttr, userName string, options map[string]interface{}) (func(), error) {
	parts := strings.SplitN(userName, `\`, 2)
	var system string
	var err error
	if len(parts) != 2 {
		userName = parts[0]
	} else {
		system = parts[0]
		userName = parts[1]
	}
	var token syscall.Token
	if v, y := options[`Password`]; y {
		password := com.String(v)
		token, err = LogonUser(userName, password, Logon32LogonInteractive)
	} else {
		token, err = getToken(system, userName)
	}
	if err != nil {
		return nil, err
	}
	if v, y := options[`HideWindow`]; y {
		attr.HideWindow = com.Bool(v)
	}
	//attr.CreationFlags = syscall.CREATE_NEW_PROCESS_GROUP
	attr.Token = token
	var closed bool
	return func() {
		if !closed {
			closed = true
			token.Close()
		}
	}, nil
}

func getToken(system, user string) (token syscall.Token, err error) {
	if len(system) == 0 {
		system, err = os.Hostname()
		if err != nil {
			err = fmt.Errorf(`failed to query os.Hostname(): %w`, err)
			return
		}
	}

	// 仅用于此用户已经登录系统的情况
	pid, err := getPidByUsername(system + `\` + user)
	if err != nil {
		return 0, err
	}
	return getTokenByPid(uint32(pid))
}

var ErrUsersProcessNotFound = errors.New(`the user's process not found`)

func getPidByUsername(username string, exename ...string) (int32, error) {
	var name string
	if len(exename) > 0 {
		name = exename[0]
	}
	pids, err := process.Pids()
	if err != nil {
		return 0, err
	}
	var pname, pusername string
	for _, pid := range pids {
		var proc *process.Process
		proc, err = process.NewProcess(pid)
		if err != nil {
			return 0, err
		}
		if len(name) > 0 || debug {
			pname, err = proc.Name()
			if err != nil {
				return 0, fmt.Errorf(`failed to query proc.Name(): %w`, err)
			}
			if len(name) > 0 && !strings.EqualFold(pname, name) {
				continue
			}
		}
		pusername, err = proc.Username()
		if err != nil {
			err = fmt.Errorf(`failed to query proc.Username(): %w`, err)
			continue
		}
		if debug {
			fmt.Println(`pname:`, pname, `pusername:`, pusername)
		}
		if strings.EqualFold(pusername, username) {
			return pid, nil
		}
	}
	if err != nil {
		return 0, err
	}
	if len(name) == 0 {
		name = `<any>`
	}
	err = fmt.Errorf(`%w: process(name: %v, username: %v) not found`, ErrUsersProcessNotFound, name, username)
	return 0, err
}

func getTokenByPid(pid uint32) (syscall.Token, error) {
	var err error
	var token syscall.Token

	handle, err := syscall.OpenProcess(syscall.PROCESS_QUERY_INFORMATION, false, pid)
	if err != nil {
		return 0, fmt.Errorf("failed to OpenProcess(%d): %w", pid, err)
	}
	defer syscall.CloseHandle(handle)
	// Find process token via win32
	// 仅用于在以服务的方式启动的程序内调用，否则会报错
	err = syscall.OpenProcessToken(handle, syscall.TOKEN_ALL_ACCESS, &token)
	if err != nil {
		return 0, fmt.Errorf("failed to OpenProcessToken(%d): %w", handle, err)
	}
	return token, err
}
