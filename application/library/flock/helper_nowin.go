// +build !windows

/*
   Nging is a toolbox for webmasters
   Copyright (C) 2021-present Wenhui Shen <swh@admpub.com>

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

package flock

import (
	"os"
	"syscall"
)

// LockEx 独占锁(独占读和写)，非阻塞模式，获取失败直接返回错误
func LockEx(f *os.File) error {
	return syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
}

// LockExBlock 独占锁，阻塞模式，获取不到时阻塞等待直至成功
func LockExBlock(f *os.File) error {
	return syscall.Flock(int(f.Fd()), syscall.LOCK_EX)
}

// Lock 共享锁(支持多读)，非阻塞模式，获取失败直接返回错误
func Lock(f *os.File) error {
	return syscall.Flock(int(f.Fd()), syscall.LOCK_SH|syscall.LOCK_NB)
}

// LockBlock 共享锁，阻塞模式，获取不到时阻塞等待直至成功
func LockBlock(f *os.File) error {
	return syscall.Flock(int(f.Fd()), syscall.LOCK_SH)
}

// Unlock 解锁
func Unlock(f *os.File) error {
	return syscall.Flock(int(f.Fd()), syscall.LOCK_UN)
}

// UnlockAndClose 解锁并关闭文件
func UnlockAndClose(f *os.File) error {
	defer f.Close()
	return syscall.Flock(int(f.Fd()), syscall.LOCK_UN)
}
