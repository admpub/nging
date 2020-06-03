// +build go1.2

// Copyright 2013 com authors
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

// Package com is an open source project for commonly used functions for the Go programming language.
package com

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/admpub/log"
)

var ErrCmdNotRunning = errors.New(`command is not running`)

// ElapsedMemory 内存占用
func ElapsedMemory() (ret string) {
	memStat := new(runtime.MemStats)
	runtime.ReadMemStats(memStat)
	ret = FormatByte(memStat.Alloc, 3)
	return
}

// ExecCmdDirBytes executes system command in given directory
// and return stdout, stderr in bytes type, along with possible error.
func ExecCmdDirBytes(dir, cmdName string, args ...string) ([]byte, []byte, error) {
	bufOut := new(bytes.Buffer)
	bufErr := new(bytes.Buffer)

	cmd := exec.Command(cmdName, args...)
	cmd.Dir = dir
	cmd.Stdout = bufOut
	cmd.Stderr = bufErr

	err := cmd.Run()
	if err != nil {
		if e, y := err.(*exec.ExitError); y {
			OnCmdExitError(append([]string{cmdName}, args...), e)
		} else {
			cmd.Stderr.Write([]byte(err.Error() + "\n"))
		}
	}
	return bufOut.Bytes(), bufErr.Bytes(), err
}

// ExecCmdBytes executes system command
// and return stdout, stderr in bytes type, along with possible error.
func ExecCmdBytes(cmdName string, args ...string) ([]byte, []byte, error) {
	return ExecCmdDirBytes("", cmdName, args...)
}

// ExecCmdDir executes system command in given directory
// and return stdout, stderr in string type, along with possible error.
func ExecCmdDir(dir, cmdName string, args ...string) (string, string, error) {
	bufOut, bufErr, err := ExecCmdDirBytes(dir, cmdName, args...)
	return string(bufOut), string(bufErr), err
}

// ExecCmd executes system command
// and return stdout, stderr in string type, along with possible error.
func ExecCmd(cmdName string, args ...string) (string, string, error) {
	return ExecCmdDir("", cmdName, args...)
}

// WritePidFile writes the process ID to the file at PidFile.
// It does nothing if PidFile is not set.
func WritePidFile(pidFile string) error {
	if pidFile == "" {
		return nil
	}
	pid := []byte(strconv.Itoa(os.Getpid()) + "\n")
	return ioutil.WriteFile(pidFile, pid, 0644)
}

var (
	equal  = rune('=')
	space  = rune(' ')
	quote  = rune('"')
	slash  = rune('\\')
	envOS  = regexp.MustCompile(`\{\$[a-zA-Z0-9_]+\}`)
	envWin = regexp.MustCompile(`\{%[a-zA-Z0-9_]+%\}`)
)

func ParseArgs(command string) (params []string) {
	item := []rune{}
	hasQuote := false
	hasSlash := false
	maxIndex := len(command) - 1
	//tower.exe -c tower.yaml -p "eee\"ddd" -t aaaa
	for k, v := range command {
		if !hasQuote {
			if v == space || v == equal {
				params = append(params, string(item))
				item = []rune{}
				continue
			}
			if v == quote {
				hasQuote = true
				continue
			}
		} else {
			if !hasSlash && v == quote {
				hasQuote = false
				continue
			}
			if !hasSlash && v == slash && k+1 <= maxIndex && command[k+1] == '"' {
				hasSlash = true
				continue
			}
			hasSlash = false
		}
		item = append(item, v)
	}
	if len(item) > 0 {
		params = append(params, string(item))
	}
	for k, v := range params {
		v = envWin.ReplaceAllStringFunc(v, getWinEnv)
		params[k] = envOS.ReplaceAllStringFunc(v, getEnv)
	}
	//fmt.Printf("---> %#v\n", params)
	//params = []string{}
	return
}

func getWinEnv(s string) string {
	s = strings.TrimPrefix(s, `{%`)
	s = strings.TrimSuffix(s, `%}`)
	return os.Getenv(s)
}

func getEnv(s string) string {
	s = strings.TrimPrefix(s, `{$`)
	s = strings.TrimSuffix(s, `}`)
	return os.Getenv(s)
}

type CmdResultCapturer struct {
	Do func([]byte) error
}

func (this CmdResultCapturer) Write(p []byte) (n int, err error) {
	err = this.Do(p)
	n = len(p)
	return
}

func (this CmdResultCapturer) WriteString(p string) (n int, err error) {
	err = this.Do([]byte(p))
	n = len(p)
	return
}

func NewCmdChanReader(timeouts ...time.Duration) *CmdChanReader {
	timeout := time.Second
	if len(timeouts) > 0 {
		timeout = timeouts[0]
	}
	return &CmdChanReader{ch: make(chan io.Reader), timeout: timeout}
}

type CmdChanReader struct {
	ch      chan io.Reader
	timeout time.Duration
	debug   bool
}

func (c *CmdChanReader) Read(p []byte) (n int, err error) {
	if c.ch == nil {
		c.ch = make(chan io.Reader)
	}
	r := <-c.ch
	if r == nil {
		return 0, errors.New(`CmdChanReader Chan has been closed`)
	}
	return r.Read(p)
}

func (c *CmdChanReader) Debug(on bool) *CmdChanReader {
	c.debug = on
	return c
}

func (c *CmdChanReader) Close() {
	if c.ch == nil {
		return
	}
	close(c.ch)
	c.ch = nil
}

func (c *CmdChanReader) Send(b []byte) *CmdChanReader {
	c.sendWithTimeout(bytes.NewReader(b))
	return c
}

func (c *CmdChanReader) sendWithTimeout(r io.Reader) {
	go func() {
		t := time.NewTicker(c.timeout)
		defer t.Stop()
		for {
			select {
			case c.ch <- r:
				if c.debug {
					println(`CmdChanReader Chan has been sent`)
				}
				return
			case <-t.C:
				if c.debug {
					println(`CmdChanReader Chan has timed out`)
				}
				c.Close()
				return
			}
		}
	}()
}

func (c *CmdChanReader) SendString(s string) *CmdChanReader {
	c.sendWithTimeout(strings.NewReader(s))
	return c
}

func NewCmdStartResultCapturer(writer io.Writer, duration time.Duration) *CmdStartResultCapturer {
	return &CmdStartResultCapturer{
		writer:   writer,
		duration: duration,
		started:  time.Now(),
		buffer:   bytes.NewBuffer(nil),
	}
}

type CmdStartResultCapturer struct {
	writer   io.Writer
	started  time.Time
	duration time.Duration
	buffer   *bytes.Buffer
}

func (this CmdStartResultCapturer) Write(p []byte) (n int, err error) {
	if time.Now().Sub(this.started) < this.duration {
		this.buffer.Write(p)
	}
	return this.writer.Write(p)
}

func (this CmdStartResultCapturer) Buffer() *bytes.Buffer {
	return this.buffer
}

func (this CmdStartResultCapturer) Writer() io.Writer {
	return this.writer
}

func CreateCmdStr(command string, recvResult func([]byte) error) *exec.Cmd {
	out := CmdResultCapturer{Do: recvResult}
	return CreateCmdStrWithWriter(command, out)
}

func CreateCmd(params []string, recvResult func([]byte) error) *exec.Cmd {
	out := CmdResultCapturer{Do: recvResult}
	return CreateCmdWithWriter(params, out)
}

func CreateCmdStrWithWriter(command string, writer ...io.Writer) *exec.Cmd {
	params := ParseArgs(command)
	return CreateCmdWithWriter(params, writer...)
}

func CreateCmdWithWriter(params []string, writer ...io.Writer) *exec.Cmd {
	var cmd *exec.Cmd
	length := len(params)
	if length == 0 || len(params[0]) == 0 {
		return cmd
	}
	if length > 1 {
		cmd = exec.Command(params[0], params[1:]...)
	} else {
		cmd = exec.Command(params[0])
	}
	var wOut, wErr io.Writer = os.Stdout, os.Stderr
	length = len(writer)
	if length > 0 {
		if writer[0] != nil {
			wOut = writer[0]
		}
		if length > 1 && writer[1] != nil {
			wErr = writer[1]
		}
	}
	cmd.Stdout = wOut
	cmd.Stderr = wErr
	return cmd
}

func RunCmdStr(command string, recvResult func([]byte) error) *exec.Cmd {
	out := CmdResultCapturer{Do: recvResult}
	return RunCmdStrWithWriter(command, out)
}

func RunCmd(params []string, recvResult func([]byte) error) *exec.Cmd {
	out := CmdResultCapturer{Do: recvResult}
	return RunCmdWithWriter(params, out)
}

func RunCmdStrWithWriter(command string, writer ...io.Writer) *exec.Cmd {
	params := ParseArgs(command)
	return RunCmdWithWriter(params, writer...)
}

var OnCmdExitError = func(params []string, err *exec.ExitError) {
	fmt.Printf("[%v]The process exited abnormally: PID(%d) PARAMS(%v) ERR(%v)\n", time.Now().Format(`2006-01-02 15:04:05`), err.Pid(), params, err)
}

func RunCmdWithReaderWriter(params []string, reader io.Reader, writer ...io.Writer) *exec.Cmd {
	cmd := CreateCmdWithWriter(params, writer...)
	cmd.Stdin = reader

	go func() {
		err := cmd.Run()
		if err != nil {
			if e, y := err.(*exec.ExitError); y {
				OnCmdExitError(params, e)
			} else {
				cmd.Stderr.Write([]byte(err.Error() + "\n"))
			}
		}
	}()

	return cmd
}

func RunCmdWithWriter(params []string, writer ...io.Writer) *exec.Cmd {
	cmd := CreateCmdWithWriter(params, writer...)

	go func() {
		err := cmd.Run()
		if err != nil {
			if e, y := err.(*exec.ExitError); y {
				OnCmdExitError(params, e)
			} else {
				cmd.Stderr.Write([]byte(err.Error() + "\n"))
			}
		}
	}()

	return cmd
}

func RunCmdWithWriterx(params []string, wait time.Duration, writer ...io.Writer) (cmd *exec.Cmd, err error, newOut *CmdStartResultCapturer, newErr *CmdStartResultCapturer) {
	length := len(writer)
	var wOut, wErr io.Writer = os.Stdout, os.Stderr
	if length > 0 {
		if writer[0] != nil {
			wOut = writer[0]
		}
		if length > 1 {
			if writer[1] != nil {
				wErr = writer[1]
			}
		}
	}
	newOut = NewCmdStartResultCapturer(wOut, wait)
	newErr = NewCmdStartResultCapturer(wErr, wait)
	writer = []io.Writer{newOut, newErr}
	cmd = CreateCmdWithWriter(params, writer...)
	go func() {
		err = cmd.Run()
		if err != nil {
			if e, y := err.(*exec.ExitError); y {
				OnCmdExitError(params, e)
			} else {
				cmd.Stderr.Write([]byte(err.Error() + "\n"))
			}
		}
	}()
	time.Sleep(wait)
	return
}

func CloseProcessFromPidFile(pidFile string) (err error) {
	if pidFile == `` {
		return
	}
	b, err := ioutil.ReadFile(pidFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return
	}
	pid, err := strconv.Atoi(strings.TrimSpace(string(b)))
	if err != nil {
		log.Error(err.Error())
		return nil
	}
	return CloseProcessFromPid(pid)
}

func CloseProcessFromPid(pid int) (err error) {
	if pid <= 0 {
		return nil
	}
	procs, err := os.FindProcess(pid)
	if err == nil {
		return procs.Kill()
	}
	return
}

func CloseProcessFromCmd(cmd *exec.Cmd) error {
	if cmd == nil {
		return nil
	}
	if cmd.Process == nil {
		return nil
	}
	err := cmd.Process.Kill()
	if cmd.ProcessState == nil || cmd.ProcessState.Exited() {
		return nil
	}
	return err
}

func CmdIsRunning(cmd *exec.Cmd) bool {
	return cmd != nil && cmd.ProcessState == nil
}
