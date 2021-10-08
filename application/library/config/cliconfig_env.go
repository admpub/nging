package config

import (
	"os"
	"path/filepath"

	"github.com/admpub/godotenv"
	"github.com/admpub/log"
	"github.com/admpub/nging/v3/application/library/common"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
)

func findEnvFile() []string {
	var envFiles []string
	envFile := filepath.Join(echo.Wd(), `.env`)
	if fi, err := os.Stat(envFile); err == nil && !fi.IsDir() {
		envFiles = append(envFiles, envFile)
	}
	return envFiles
}

func (c *CLIConfig) InitEnviron(needFindEnvFile ...bool) (err error) {
	if len(needFindEnvFile) > 0 && needFindEnvFile[0] {
		c.envFiles = findEnvFile()
	}
	if c.envVars != nil {
		for k := range c.envVars {
			os.Unsetenv(k)
		}
	}
	if len(c.envFiles) > 0 {
		log.Infof(`Loading env file: %#v`, c.envFiles)
		c.envVars, err = godotenv.Read(c.envFiles...)
		if err != nil {
			return
		}
		currentEnv := godotenv.CurrentEnvKeys()
		for k, v := range c.envVars {
			if !currentEnv[k] {
				log.Infof(`Set env var: %s`, k)
				os.Setenv(k, v)
			} else {
				log.Infof(`Skip env var: %s`, k)
				delete(c.envVars, k)
			}
		}
	}
	return
}

func (c *CLIConfig) WatchEnvConfig() {
	if c.envMonitor != nil {
		c.envMonitor.Close()
		c.envMonitor = nil
	}
	if len(c.envFiles) == 0 {
		return
	}
	c.envMonitor = &com.MonitorEvent{
		Modify: func(file string) {
			log.Info(`Start reloading env file: ` + file)
			err := c.InitEnviron()
			if err == nil {
				log.Info(`Succcessfully reload the env file: ` + file)
				return
			}
			if err == common.ErrIgnoreConfigChange {
				log.Info(`No need to reload the env file: ` + file)
				return
			}
			log.Error(err)
		},
	}
	for _, envFile := range c.envFiles {
		err := c.envMonitor.AddFile(envFile)
		if err != nil {
			log.Error(err)
		}
	}
	c.envMonitor.Watch()
}
