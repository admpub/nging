//go:build !windows

package goforever

import (
	"errors"
	"fmt"
	"os/user"
	"strconv"
	"syscall"
)

func buildOption(options map[string]interface{}) map[string]interface{} {
	return options
}

func SetOption(options map[string]interface{}, name string, value interface{}) map[string]interface{} {
	return options
}

func SetSysProcAttr(attr *syscall.SysProcAttr, userName string, options map[string]interface{}) (func(), error) {
	userInfo, err := user.Lookup(userName)
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
	attr.Credential = &syscall.Credential{
		Uid:         uint32(uid),
		Gid:         uint32(gid),
		NoSetGroups: true,
	}
	return nil, err
}
