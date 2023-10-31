package config

import (
	"context"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/admpub/nging/v5/application/library/config"
	"github.com/nging-plugins/caddymanager/application/library/engine"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
)

var _ engine.EngineConfigFileFixer = New()

const Name = `caddy2`

func New() *Config {
	return &Config{
		CommonConfig: engine.NewConfig(Name, Name),
	}
}

type Config struct {
	*engine.CommonConfig
	vhostConfigDirAbsPath string
}

func (c *Config) Init() *Config {
	return c
}

func DefaultConfigDir() string {
	return filepath.Join(config.FromCLI().ConfDir(), `vhosts-caddy2`)
}

func (c *Config) GetVhostConfigLocalDirAbs() (string, error) {
	if len(c.VhostConfigLocalDir) == 0 {
		c.VhostConfigLocalDir = filepath.Join(DefaultConfigDir(), c.ID)
	}
	return c.VhostConfigLocalDir, nil
}

func (c *Config) GetVhostConfigDirAbs() (string, error) {
	var vhostDir string
	if c.Environ == engine.EnvironContainer {
		vhostDir = c.VhostConfigContainerDir
	} else {
		var err error
		vhostDir, err = c.GetVhostConfigLocalDirAbs()
		if err != nil {
			return vhostDir, err
		}
	}
	return vhostDir, nil
}

func (c *Config) Start(ctx context.Context) error {
	args := []string{`start`}
	if c.CmdWithConfig && len(c.EngineConfigFile()) > 0 {
		args = append(args, `--config`, c.EngineConfigFile())
	}
	err := c.exec(ctx, args...)
	return err
}

func (c *Config) Reload(ctx context.Context) error {
	args := []string{`reload`}
	if c.CmdWithConfig && len(c.EngineConfigFile()) > 0 {
		args = append(args, `--config`, c.EngineConfigFile())
	}
	err := c.exec(ctx, args...)
	return err
}

func (c *Config) TestConfig(ctx context.Context) error {
	args := []string{`validate`}
	if c.CmdWithConfig && len(c.EngineConfigFile()) > 0 {
		args = append(args, `--config`, c.EngineConfigFile())
	}
	err := c.exec(ctx, args...)
	return err
}

func (c *Config) Stop(ctx context.Context) error {
	err := c.exec(ctx, `stop`)
	return err
}

// Caddy 2.4.4 and up supports adding a module directly:
// $ caddy add-package github.com/caddyserver/transform-encoder
func (c *Config) InstallModule(ctx context.Context, module string) error {
	err := c.exec(ctx, `add-package`, module)
	if err != nil {
		return err
	}
	err = c.exec(ctx, `upgrade`)
	if err == nil {
		return nil
	}
	if err = c.Stop(ctx); err == nil {
		err = c.Start(ctx)
	}
	return err
}

func (c *Config) exec(ctx context.Context, args ...string) error {
	if len(c.Command) == 0 {
		c.Command = `caddy`
		if com.IsWindows {
			c.Command += `.exe`
		}
	}
	_, err := c.CommonConfig.Exec(ctx, args...)
	return err
}

func (c *Config) FixEngineConfigFile(deleteMode ...bool) (bool, error) {
	if len(c.EngineConfigLocalFile) == 0 {
		return false, nil
	}
	vhostDir, err := c.GetVhostConfigDirAbs()
	if len(vhostDir) == 0 {
		return false, err
	}
	var delmode bool
	if len(deleteMode) > 0 {
		delmode = deleteMode[0]
	}
	if !delmode && !com.FileExists(c.EngineConfigLocalFile) {
		com.MkdirAll(filepath.Dir(c.EngineConfigLocalFile), os.ModePerm)
		dir := c.FixVhostDirPath(vhostDir)
		err = os.WriteFile(c.EngineConfigLocalFile, []byte("import \""+dir+"*.conf\";\n"), 0644)
		if err != nil {
			return false, err
		}
		return true, nil
	}
	findString := `[\s]*import[\s]+["']?` + regexp.QuoteMeta(vhostDir) + `[\/]?\*(\.conf)?["']?[\s]*`
	re, err := regexp.Compile(findString)
	if err != nil {
		return false, err
	}
	var seekedContent string
	var hasUpdate bool
	err = com.SeekFileLines(c.EngineConfigLocalFile, func(line string) error {
		seekedContent += line + "\n"
		if hasUpdate {
			return nil
		}
		if re.MatchString(line) {
			if delmode {
				seekedContent = strings.TrimSuffix(seekedContent, line+"\n")
				hasUpdate = true
				return nil
			}
			return echo.ErrExit
		}
		return nil
	})
	if err != nil {
		if err != echo.ErrExit {
			return hasUpdate, err
		}
	} else if !hasUpdate {
		dir := c.FixVhostDirPath(vhostDir)
		seekedContent += "import \"" + dir + "*.conf\";\n"
		hasUpdate = true
	}
	if hasUpdate {
		err = com.Copy(c.EngineConfigLocalFile, c.EngineConfigLocalFile+`.`+time.Now().Format(`20060102150405.000`)+`.ngingbak`)
		if err != nil {
			return hasUpdate, err
		}
		seekedContent = strings.TrimRight(seekedContent, "\n ")
		return hasUpdate, os.WriteFile(c.EngineConfigLocalFile, com.Str2bytes(seekedContent), 0644)
	}
	return hasUpdate, nil
}
