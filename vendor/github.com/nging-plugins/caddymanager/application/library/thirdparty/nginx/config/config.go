package config

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"

	"github.com/admpub/nging/v5/application/library/config"
	"github.com/nging-plugins/caddymanager/application/dbschema"
	"github.com/nging-plugins/caddymanager/application/library/engine"
	"github.com/nging-plugins/caddymanager/application/library/engine/enginevent"
	"github.com/nging-plugins/caddymanager/application/library/form"
	"github.com/nging-plugins/caddymanager/application/library/htpasswd"
)

const Name = `nginx`

var (
	regexConfigFile     = regexp.MustCompile(`[\s]+configuration file (.+\.conf)[\s]+`)
	regexConfigInclude  = regexp.MustCompile(`^include[\s]+([\S]+(?:\*|\*\.conf))[\s]*;(?:[\s]*#.*)?[\s]*$`)
	regexVersion        = regexp.MustCompile(`[\d]+\.[\d]+\.[\d]+`)
	regexHttpBlockStart = regexp.MustCompile(`^http[\s]*\{$`)

	_ engine.EngineConfigFileFixer   = (*Config)(nil)
	_ engine.VhostConfigRemover      = (*Config)(nil)
	_ enginevent.OnVhostConfigSaving = (*Config)(nil)
	_ engine.CertRenewaler           = (*Config)(nil)
	_ engine.CertFileRemover         = (*Config)(nil)
	_ engine.CertPathFormatGetter    = (*Config)(nil)
	_ engine.CertUpdaterGetter       = (*Config)(nil)
)

func New() *Config {
	return &Config{
		CommonConfig: engine.NewConfig(Name, Name),
	}
}

type Config struct {
	Version        string
	CertPathFormat engine.CertPathFormatWithUpdater
	*engine.CommonConfig
}

func (c *Config) Init() error {
	var err error
	/*
		ctx := context.Background()
		if len(c.Version) == 0 {
			c.Version, err = c.getVersion(ctx)
			if err != nil {
				return err
			}
		}
		if len(c.VhostConfigLocalDir) == 0 {
			if len(c.EngineConfigLocalFile) == 0 {
				c.EngineConfigLocalFile, err = c.getEngineConfigLocalFile(ctx)
				if err != nil {
					return err
				}
			}
			c.VhostConfigLocalDir, err = c.getVhostConfigLocalDir(c.EngineConfigLocalFile)
		}
	*/
	return err
}

func DefaultConfigDir() string {
	return filepath.Join(config.FromCLI().ConfDir(), `vhosts-nginx`)
}

func (c *Config) CopyFrom(m *dbschema.NgingVhostServer) {
	c.CertPathFormat.CopyFrom(m)
	c.CommonConfig.CopyFrom(m)
}

func (c *Config) GetCertPathFormat(ctx echo.Context) engine.CertPathFormat {
	c.CertPathFormat.AutoDetect(ctx)
	return c.CertPathFormat.CertPathFormat
}

func (c *Config) CertUpdater() string {
	if len(c.CertLocalDir) > 0 && len(c.CertPathFormat.CertLocalUpdater()) > 0 {
		return c.CertPathFormat.CertLocalUpdater()
	}
	return c.CertPathFormat.CertContainerUpdater()
}

func (c *Config) GetVhostConfigLocalDirAbs() (string, error) {
	if len(c.VhostConfigLocalDir) == 0 {
		c.VhostConfigLocalDir = filepath.Join(DefaultConfigDir(), c.ID)
	}
	return c.VhostConfigLocalDir, nil
}

func (c *Config) Start(ctx context.Context) error {
	args := []string{}
	if c.CmdWithConfig && len(c.EngineConfigFile()) > 0 && strings.HasSuffix(c.EngineConfigFile(), `.conf`) {
		args = append(args, `-c`, c.EngineConfigFile())
	}
	_, err := c.exec(ctx)
	return err
}

func (c *Config) TestConfig(ctx context.Context) error {
	args := []string{`-t`}
	if c.CmdWithConfig && len(c.EngineConfigFile()) > 0 {
		args = append(args, `-c`, c.EngineConfigFile())
	}
	//nginx: the configuration file /etc/nginx/nginx.conf syntax is ok
	//nginx: configuration file /etc/nginx/nginx.conf test is successful
	_, err := c.exec(ctx, args...)
	return err
}

func (c *Config) getVersion(ctx context.Context) (string, error) {
	result, err := c.exec(ctx, `-v`)
	if err != nil {
		return ``, err
	}
	result = bytes.TrimSpace(result)
	parts := bytes.SplitN(result, []byte(`:`), 2)
	if len(parts) != 2 {
		matches := regexVersion.FindStringSubmatch(string(result))
		if len(matches) > 0 {
			return matches[0], err
		}
		return ``, err
	}
	result = parts[1]
	parts = bytes.SplitN(result, []byte(`/`), 2)
	if len(parts) != 2 {
		return ``, err
	}
	parts[1] = bytes.TrimSpace(bytes.SplitN(parts[1], []byte(` `), 2)[0])
	return string(parts[1]), err
}

func (c *Config) getEngineConfigLocalFile(ctx context.Context) (string, error) {
	result, err := c.exec(ctx, `-t`)
	if err != nil {
		return ``, err
	}
	result = bytes.TrimSpace(result)
	lines := bytes.Split(result, []byte{com.LF})
	var configFilePath string
	for _, line := range lines {
		matches := regexConfigFile.FindAllSubmatch(line, 1)
		if len(matches) > 0 && len(matches[0]) > 1 {
			configFilePath = string(matches[0][1])
			break
		}
	}
	return configFilePath, err
}

func (c *Config) getVhostConfigLocalDir(confPath string) (string, error) {
	if len(confPath) == 0 {
		return ``, nil
	}
	var includeConfD string
	var includeSitesEnabled string
	var includeDir string
	err := com.SeekFileLines(confPath, func(line string) error {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			return nil
		}
		if strings.HasPrefix(line, `#`) {
			return nil
		}
		matches := regexConfigInclude.FindAllStringSubmatch(line, 1)
		if len(matches) == 0 || len(matches[0]) < 2 {
			return nil
		}
		line = matches[0][1]
		if strings.Contains(line, `sites-`) /*strings.Contains(line, `sites-enabled`) || strings.Contains(line, `sites-available`)*/ {
			includeSitesEnabled = line
			return echo.ErrExit
		}
		if strings.Contains(line, `conf.d`) {
			includeConfD = line
			return nil
		}
		includeDir = line
		return nil
	})
	if err != nil && err != echo.ErrExit {
		return ``, err
	}
	if len(includeSitesEnabled) > 0 {
		includeDir = includeSitesEnabled
	} else if len(includeConfD) > 0 {
		includeDir = includeConfD
	}
	includeDir = com.TrimFileName(includeDir)
	return includeDir, nil
}

func (c *Config) Reload(ctx context.Context) error {
	return c.sendSignal(ctx, `reload`)
}

func (c *Config) Stop(ctx context.Context) error {
	return c.sendSignal(ctx, `stop`)
}

func (c *Config) Quit(ctx context.Context) error {
	return c.sendSignal(ctx, `quit`)
}

func (c *Config) Reopen(ctx context.Context) error {
	return c.sendSignal(ctx, `reopen`)
}

// signal: stop, quit, reopen, reload
func (c *Config) sendSignal(ctx context.Context, signal string) error {
	_, err := c.exec(ctx, `-s`, signal)
	return err
}

func (c *Config) exec(ctx context.Context, args ...string) ([]byte, error) {
	if len(c.Command) == 0 {
		c.Command = `nginx`
		if com.IsWindows {
			c.Command += `.exe`
		}
	}
	return c.CommonConfig.Exec(ctx, args...)
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
	findString := `[\s]*include[\s]+["']?` + regexp.QuoteMeta(vhostDir) + `[\/]?\*(\.conf)?["']?[\s]*;`
	re, err := regexp.Compile(findString)
	if err != nil {
		return false, err
	}
	var httpBlockStart bool
	var seekedContent string
	var hasUpdate bool
	err = com.SeekFileLines(c.EngineConfigLocalFile, func(line string) error {
		if httpBlockStart && strings.TrimRight(line, "\t ") == `}` {
			if !delmode {
				dir := c.FixVhostDirPath(vhostDir)
				line = "\n\tinclude \"" + dir + "*.conf\";\n" + line
				hasUpdate = true
			}
			httpBlockStart = false
		}
		seekedContent += line + "\n"
		if hasUpdate {
			return nil
		}
		cleaned := strings.TrimSpace(line)
		if len(cleaned) == 0 {
			return nil
		}
		if strings.HasPrefix(cleaned, `#`) {
			return nil
		}
		if !httpBlockStart && regexHttpBlockStart.MatchString(cleaned) {
			httpBlockStart = true
			return nil
		}
		if httpBlockStart && re.MatchString(cleaned) {
			if delmode {
				seekedContent = strings.TrimSuffix(seekedContent, line+"\n")
				hasUpdate = true
				httpBlockStart = false
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

func (c *Config) RemoveVhostConfig(id uint) error {
	vhostDir, err := c.GetVhostConfigLocalDirAbs()
	if err != nil {
		return err
	}
	dir := filepath.Join(vhostDir, `htpasswd`)
	if !com.FileExists(dir) {
		return nil
	}
	if id > 0 {
		filePath := filepath.Join(dir, engine.NgingConfigPrefix+strconv.FormatUint(uint64(id), 10))
		if com.FileExists(filePath) {
			return os.Remove(filePath)
		}
		return nil
	}
	return c.CommonConfig.RemoveDir(`htpasswd(`+c.ID+`)`, dir, engine.NgingConfigPrefix)
}

func (c *Config) OnVhostConfigSaving(id uint, formData *form.Values) error {
	vhostDir, err := c.GetVhostConfigLocalDirAbs()
	if err != nil {
		return err
	}
	filePath := filepath.Join(vhostDir, `htpasswd`, engine.NgingConfigPrefix+strconv.FormatUint(uint64(id), 10))
	if formData.IsEnabled(`basicauth`) {
		username := strings.TrimSpace(formData.GetAttrVal("basicauth", "username"))
		password := strings.TrimSpace(formData.GetAttrVal("basicauth", "password"))
		if len(username) > 0 && len(password) > 0 {
			a := htpasswd.Accounts{}
			err = a.SetPassword(username, password, htpasswd.AlgoBCrypt)
			if err != nil {
				return err
			}
			if !com.FileExists(filePath) {
				com.MkdirAll(filepath.Dir(filePath), os.ModePerm)
			}
			return a.WriteToFile(filePath)
		}
	}
	if !com.FileExists(filePath) {
		return nil
	}
	return os.Remove(filePath)
}

func (c *Config) RemoveCertFile(id uint) error {
	if len(c.CertLocalDir) == 0 {
		return nil
	}
	certDir := filepath.Join(c.CertLocalDir, engine.NgingConfigPrefix+strconv.FormatUint(uint64(id), 10))
	return os.RemoveAll(certDir)
}

func (c *Config) RenewCert(ctx echo.Context, id uint, domains []string, email string, isObtain bool) error {
	command := strings.TrimSpace(c.Command)
	command = strings.TrimSuffix(command, `.exe`)
	command = strings.TrimSuffix(command, `nginx`)
	certUpdater := c.CertPathFormat.CertLocalUpdater()
	if len(certUpdater) == 0 {
		certUpdater = `certbot`
	}
	item := engine.CertUpdaters.GetItem(certUpdater)
	if item == nil {
		return ctx.NewError(code.Unsupported, `不支持证书更新工具: %v`, certUpdater)
	}

	up := item.X.(engine.CertUpdater)

	var stdCtx context.Context = ctx
	if len(c.CertLocalDir) > 0 && len(c.CertPathFormat.CertLocalUpdater()) > 0 {
		certDir := filepath.Join(c.CertLocalDir, engine.NgingConfigPrefix+strconv.FormatUint(uint64(id), 10), `well-known`)
		com.MkdirAll(certDir, os.ModePerm)
		if len(c.CertPathFormat.SaveDir) > 0 && c.Environ == engine.EnvironLocal {
			stdCtx = context.WithValue(ctx, engine.CtxCertDir, c.CertPathFormat.SaveDir)
		}
		return up.Update(stdCtx, ``, domains, email, certDir, isObtain)
	}
	if c.Environ == engine.EnvironContainer {
		if len(c.CertContainerDir) == 0 {
			return fmt.Errorf(`[%d][%s]%w`, id, strings.Join(domains, `,`), engine.ErrNotSetCertContainerDir)
		}
		if len(c.Endpoint) > 0 {
			err := fmt.Errorf(`[%s]Updating certificates is not supported in the API mode of the container`, c.GetIdent())
			return err
		}
		if len(c.CertPathFormat.CertContainerUpdater()) > 0 {
			certUpdater = c.CertPathFormat.CertContainerUpdater()
		}
		certDir := filepath.Join(c.CertContainerDir, engine.NgingConfigPrefix+strconv.FormatUint(uint64(id), 10), `well-known`)
		certDir = filepath.ToSlash(certDir)
		command = strings.TrimSpace(command)
		parts := com.ParseArgs(command)
		if len(parts) == 0 {
			return fmt.Errorf(`failed to parse command %q`, command)
		}
		executeableFile := parts[0]
		var args []string
		if len(parts) > 1 {
			args = append(args, parts[1:]...)
		}
		args = append(args, `mkdir`, `-p`, certDir)
		cmd := exec.CommandContext(ctx, executeableFile, args...)
		result, err := cmd.CombinedOutput()
		//log.Okay(cmd.String())
		if err != nil {
			err = fmt.Errorf(`%s: %w`, result, err)
			return err
		}
		if len(c.CertPathFormat.SaveDir) > 0 {
			stdCtx = context.WithValue(ctx, engine.CtxCertDir, c.CertPathFormat.SaveDir)
		}
		return up.Update(stdCtx, command, domains, email, certDir, isObtain)
	}

	if len(c.CertLocalDir) == 0 {
		return fmt.Errorf(`[%d][%s]%w`, id, strings.Join(domains, `,`), engine.ErrNotSetCertLocalDir)
	}
	return fmt.Errorf(`更新证书操作失败：没有在本机安装%q`, certUpdater)
}
