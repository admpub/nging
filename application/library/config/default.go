package config

import (
	"errors"
	"os/exec"
)

var (
	DefaultConfig         = &Config{}
	DefaultCLIConfig      = &CLIConfig{cmds: map[string]*exec.Cmd{}}
	ErrUnknowDatabaseType = errors.New(`unkown database type`)
)
