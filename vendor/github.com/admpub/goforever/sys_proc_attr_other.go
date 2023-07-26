//go:build !windows

package goforever

import (
	"errors"
	"fmt"
	"os/user"
	"strconv"
	"syscall"
)

func (p *Process) setSysProcAttr(attr *syscall.SysProcAttr) error {
	userInfo, err := user.Lookup(p.User)
	if err != nil {
		return errors.New("failed to get user: " + err.Error())
	}
	uid, err := strconv.ParseUint(userInfo.Uid, 10, 32)
	if err != nil {
		return fmt.Errorf("failed to ParseUint(userInfo.Uid=%q): %w", userInfo.Uid, err)
	}
	gid, err := strconv.ParseUint(userInfo.Gid, 10, 32)
	if err != nil {
		return fmt.Errorf("failed to ParseUint(userInfo.Gid=%q): %w", userInfo.Gid, err)
	}
	attr.Credential = &syscall.Credential{
		Uid:         uint32(uid),
		Gid:         uint32(gid),
		NoSetGroups: true,
	}
	return err
}
