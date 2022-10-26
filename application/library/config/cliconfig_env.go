package config

import (
	"os"
	"path/filepath"

	"github.com/admpub/godotenv"
	"github.com/admpub/log"
	"github.com/admpub/nging/v5/application/library/common"
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
	c.envLock.Lock()
	defer c.envLock.Unlock()
	if len(needFindEnvFile) > 0 && needFindEnvFile[0] {
		c.envFiles = findEnvFile()
	}
	var newEnvVars map[string]string
	if len(c.envFiles) > 0 {
		log.Infof(`Loading env file: %#v`, c.envFiles)
		newEnvVars, err = godotenv.Read(c.envFiles...)
		if err != nil {
			return
		}
	}
	if newEnvVars != nil {
		if c.envVars != nil {
			for k, v := range c.envVars {
				newV, ok := newEnvVars[k]
				if !ok {
					log.Infof(`Unset env var: %s`, k)
					os.Unsetenv(k)
					delete(c.envVars, k)
					continue
				}
				if v != newV {
					log.Infof(`Set env var: %s`, k)
					os.Setenv(k, v)
					c.envVars[k] = newV
				}
				delete(newEnvVars, k)
			}
		} else {
			c.envVars = map[string]string{}
		}
		for k, v := range newEnvVars {
			log.Infof(`Set env var: %s`, k)
			os.Setenv(k, v)
			c.envVars[k] = v
		}
	} else {
		if c.envVars != nil {
			for k := range c.envVars {
				log.Infof(`Unset env var: %s`, k)
				os.Unsetenv(k)
			}
			c.envVars = nil
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
