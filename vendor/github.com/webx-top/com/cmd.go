//go:build go1.2
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
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

var ErrCmdNotRunning = errors.New(`command is not running`)

// ElapsedMemory 内存占用
func ElapsedMemory() (ret string) {
	memStat := new(runtime.MemStats)
	runtime.ReadMemStats(memStat)
	ret = FormatByte(memStat.Alloc, 3)
	return
}

// ExecCmdDirBytesWithContext executes system command in given directory
// and return stdout, stderr in bytes type, along with possible error.
func ExecCmdDirBytesWithContext(ctx context.Context, dir, cmdName string, args ...string) ([]byte, []byte, error) {
	bufOut := new(bytes.Buffer)
	bufErr := new(bytes.Buffer)

	cmd := exec.CommandContext(ctx, cmdName, args...)
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

// ExecCmdDirBytes executes system command in given directory
// and return stdout, stderr in bytes type, along with possible error.
func ExecCmdDirBytes(dir, cmdName string, args ...string) ([]byte, []byte, error) {
	return ExecCmdDirBytesWithContext(context.Background(), dir, cmdName, args...)
}

// ExecCmdBytes executes system command
// and return stdout, stderr in bytes type, along with possible error.
func ExecCmdBytes(cmdName string, args ...string) ([]byte, []byte, error) {
	return ExecCmdBytesWithContext(context.Background(), cmdName, args...)
}

// ExecCmdBytesWithContext executes system command
// and return stdout, stderr in bytes type, along with possible error.
func ExecCmdBytesWithContext(ctx context.Context, cmdName string, args ...string) ([]byte, []byte, error) {
	return ExecCmdDirBytesWithContext(ctx, "", cmdName, args...)
}

// ExecCmdDir executes system command in given directory
// and return stdout, stderr in string type, along with possible error.
func ExecCmdDir(dir, cmdName string, args ...string) (string, string, error) {
	return ExecCmdDirWithContext(context.Background(), dir, cmdName, args...)
}

// ExecCmdDirWithContext executes system command in given directory
// and return stdout, stderr in string type, along with possible error.
func ExecCmdDirWithContext(ctx context.Context, dir, cmdName string, args ...string) (string, string, error) {
	bufOut, bufErr, err := ExecCmdDirBytesWithContext(ctx, dir, cmdName, args...)
	return string(bufOut), string(bufErr), err
}

// ExecCmd executes system command
// and return stdout, stderr in string type, along with possible error.
func ExecCmd(cmdName string, args ...string) (string, string, error) {
	return ExecCmdWithContext(context.Background(), cmdName, args...)
}

// ExecCmdWithContext executes system command
// and return stdout, stderr in string type, along with possible error.
func ExecCmdWithContext(ctx context.Context, cmdName string, args ...string) (string, string, error) {
	return ExecCmdDirWithContext(ctx, "", cmdName, args...)
}

// WritePidFile writes the process ID to the file at PidFile.
// It does nothing if PidFile is not set.
func WritePidFile(pidFile string, pidNumbers ...int) error {
	if pidFile == "" {
		return nil
	}
	var pidNumber int
	if len(pidNumbers) > 0 {
		pidNumber = pidNumbers[0]
	} else {
		pidNumber = os.Getpid()
	}
	pid := []byte(strconv.Itoa(pidNumber) + "\n")
	return os.WriteFile(pidFile, pid, 0644)
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
		v = ParseWindowsEnvVar(v)
		params[k] = ParseEnvVar(v)
	}
	//fmt.Printf("---> %#v\n", params)
	//params = []string{}
	return
}

func ParseEnvVar(v string) string {
	return envOS.ReplaceAllStringFunc(v, getEnv)
}

func ParseWindowsEnvVar(v string) string {
	return envWin.ReplaceAllStringFunc(v, getWinEnv)
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

func (c CmdResultCapturer) Write(p []byte) (n int, err error) {
	err = c.Do(p)
	n = len(p)
	return
}

func (c CmdResultCapturer) WriteString(p string) (n int, err error) {
	err = c.Do([]byte(p))
	n = len(p)
	return
}

func NewCmdChanReader() *CmdChanReader {
	return &CmdChanReader{ch: make(chan io.Reader)}
}

type CmdChanReader struct {
	ch chan io.Reader
	mu sync.RWMutex
}

func (c *CmdChanReader) getCh() chan io.Reader {
	c.mu.RLock()
	ch := c.ch
	c.mu.RUnlock()
	return ch
}

func (c *CmdChanReader) setCh(ch chan io.Reader) {
	c.mu.Lock()
	c.ch = ch
	c.mu.Unlock()
}

func (c *CmdChanReader) Read(p []byte) (n int, err error) {
	ch := c.getCh()
	if ch == nil {
		ch = make(chan io.Reader)
		c.setCh(ch)
	}
	r := <-ch
	if r == nil {
		return 0, errors.New(`[CmdChanReader] Chan has been closed`)
	}
	return r.Read(p)
}

func (c *CmdChanReader) Close() {
	ch := c.getCh()
	if ch == nil {
		return
	}
	close(ch)
	c.setCh(nil)
}

func (c *CmdChanReader) SendReader(r io.Reader) *CmdChanReader {
	ch := c.getCh()
	if ch == nil {
		return c
	}
	defer recover()
	select {
	case ch <- r:
	default:
	}
	return c
}

func (c *CmdChanReader) SendReaderAndWait(r io.Reader) *CmdChanReader {
	ch := c.getCh()
	if ch == nil {
		return c
	}
	defer recover()
	ch <- r
	return c
}

func (c *CmdChanReader) Send(b []byte) *CmdChanReader {
	return c.SendReader(bytes.NewReader(b))
}

func (c *CmdChanReader) SendString(s string) *CmdChanReader {
	return c.SendReader(strings.NewReader(s))
}

func (c *CmdChanReader) SendAndWait(b []byte) *CmdChanReader {
	return c.SendReaderAndWait(bytes.NewReader(b))
}

func (c *CmdChanReader) SendStringAndWait(s string) *CmdChanReader {
	return c.SendReaderAndWait(strings.NewReader(s))
}

// WatchingStdin Watching os.Stdin
// example: go WatchingStdin(ctx,`name`,fn)
func WatchingStdin(ctx context.Context, name string, fn func(string) error) {
	in := bufio.NewReader(os.Stdin)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			input, err := in.ReadString(LF)
			if err != nil && err != io.EOF {
				log.Printf(`watchingStdin(%s): %s`+StrLF, name, err.Error())
				return
			}
			err = fn(input)
			if err != nil {
				log.Printf(`watchingStdin(%s): %s`+StrLF, name, err.Error())
				return
			}
		}
	}
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

func (c *CmdStartResultCapturer) Write(p []byte) (n int, err error) {
	if time.Since(c.started) < c.duration {
		c.buffer.Write(p)
	}
	return c.writer.Write(p)
}

func (c *CmdStartResultCapturer) Buffer() *bytes.Buffer {
	return c.buffer
}

func (c *CmdStartResultCapturer) Writer() io.Writer {
	return c.writer
}

func CreateCmdStr(command string, recvResult func([]byte) error) *exec.Cmd {
	return CreateCmdStrWithContext(context.Background(), command, recvResult)
}

func CreateCmdStrWithContext(ctx context.Context, command string, recvResult func([]byte) error) *exec.Cmd {
	out := CmdResultCapturer{Do: recvResult}
	return CreateCmdStrWithWriter(command, out)
}

func CreateCmd(params []string, recvResult func([]byte) error) *exec.Cmd {
	return CreateCmdWithContext(context.Background(), params, recvResult)
}

func CreateCmdWithContext(ctx context.Context, params []string, recvResult func([]byte) error) *exec.Cmd {
	out := CmdResultCapturer{Do: recvResult}
	return CreateCmdWriterWithContext(ctx, params, out)
}

func CreateCmdStrWithWriter(command string, writer ...io.Writer) *exec.Cmd {
	return CreateCmdStrWriterWithContext(context.Background(), command, writer...)
}

func CreateCmdStrWriterWithContext(ctx context.Context, command string, writer ...io.Writer) *exec.Cmd {
	params := ParseArgs(command)
	return CreateCmdWriterWithContext(ctx, params, writer...)
}

func CreateCmdWithWriter(params []string, writer ...io.Writer) *exec.Cmd {
	return CreateCmdWriterWithContext(context.Background(), params, writer...)
}

func CreateCmdWriterWithContext(ctx context.Context, params []string, writer ...io.Writer) *exec.Cmd {
	var cmd *exec.Cmd
	length := len(params)
	if length == 0 || len(params[0]) == 0 {
		return cmd
	}
	if length > 1 {
		cmd = exec.CommandContext(ctx, params[0], params[1:]...)
	} else {
		cmd = exec.CommandContext(ctx, params[0])
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
	return RunCmdStrWithContext(context.Background(), command, recvResult)
}

func RunCmdStrWithContext(ctx context.Context, command string, recvResult func([]byte) error) *exec.Cmd {
	out := CmdResultCapturer{Do: recvResult}
	return RunCmdStrWriterWithContext(ctx, command, out)
}

func RunCmd(params []string, recvResult func([]byte) error) *exec.Cmd {
	return RunCmdWithContext(context.Background(), params, recvResult)
}

func RunCmdWithContext(ctx context.Context, params []string, recvResult func([]byte) error) *exec.Cmd {
	out := CmdResultCapturer{Do: recvResult}
	return RunCmdWriterWithContext(ctx, params, out)
}

func RunCmdStrWithWriter(command string, writer ...io.Writer) *exec.Cmd {
	return RunCmdStrWriterWithContext(context.Background(), command, writer...)
}

func RunCmdStrWriterWithContext(ctx context.Context, command string, writer ...io.Writer) *exec.Cmd {
	params := ParseArgs(command)
	return RunCmdWriterWithContext(ctx, params, writer...)
}

var OnCmdExitError = func(params []string, err *exec.ExitError) {
	fmt.Printf("[%v]The process exited abnormally: PID(%d) PARAMS(%v) ERR(%v)\n", time.Now().Format(`2006-01-02 15:04:05`), err.Pid(), params, err)
}

func RunCmdReaderWriterWithContext(ctx context.Context, params []string, reader io.Reader, writer ...io.Writer) *exec.Cmd {
	cmd := CreateCmdWriterWithContext(ctx, params, writer...)
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

func RunCmdWithReaderWriter(params []string, reader io.Reader, writer ...io.Writer) *exec.Cmd {
	return RunCmdReaderWriterWithContext(context.Background(), params, reader, writer...)
}

func RunCmdWithWriter(params []string, writer ...io.Writer) *exec.Cmd {
	return RunCmdWriterWithContext(context.Background(), params, writer...)
}

func RunCmdWriterWithContext(ctx context.Context, params []string, writer ...io.Writer) *exec.Cmd {
	cmd := CreateCmdWriterWithContext(ctx, params, writer...)

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
	return RunCmdWriterxWithContext(context.Background(), params, wait, writer...)
}

func RunCmdWriterxWithContext(ctx context.Context, params []string, wait time.Duration, writer ...io.Writer) (cmd *exec.Cmd, err error, newOut *CmdStartResultCapturer, newErr *CmdStartResultCapturer) {
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
	cmd = CreateCmdWriterWithContext(ctx, params, writer...)
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
	b, err := os.ReadFile(pidFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return
	}
	pid, err := strconv.Atoi(strings.TrimSpace(string(b)))
	if err != nil {
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
