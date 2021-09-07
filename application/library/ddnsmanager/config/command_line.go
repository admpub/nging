package config

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/webx-top/com"
)

type CommandLine struct {
	Command string // 命令
	EnvVars string // 环境变量
	WorkDir string // 工作目录
	Timeout time.Duration
}

var CommandDefaultTimeout = 10 * time.Second

var ErrCommandRequired = errors.New(`命令行命令不能为空`)

func (c *CommandLine) Exec() ([]byte, error) {
	if len(c.Command) == 0 {
		return nil, ErrCommandRequired
	}
	parts := com.ParseArgs(c.Command)
	cmdName := parts[0]
	var args []string
	if len(parts) > 1 {
		args = append(args, parts[1:]...)
	}
	timeout := c.Timeout
	if timeout < time.Second {
		timeout = CommandDefaultTimeout
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, cmdName, args...)
	cmd.Dir = c.WorkDir
	cmd.Env = os.Environ()
	envVars := strings.TrimSpace(c.EnvVars)
	for _, envs := range strings.Split(envVars, "\n") {
		envs = strings.TrimSpace(envs)
		if len(envs) == 0 {
			continue
		}
		cmd.Env = append(cmd.Env, envs)
	}
	return cmd.CombinedOutput()
}
