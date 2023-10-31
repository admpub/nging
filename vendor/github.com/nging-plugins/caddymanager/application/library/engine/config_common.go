package engine

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/admpub/log"
	"github.com/webx-top/com"

	"github.com/admpub/nging/v5/application/library/common"
	"github.com/nging-plugins/caddymanager/application/dbschema"
)

func NewConfig(engineName, templateFile string) *CommonConfig {
	return &CommonConfig{engineName: engineName, templateFile: templateFile}
}

type ParsedCommand struct {
	Command         string
	Args            []string
	ContainerEngine string
	ContainerName   string
}

type CommonConfig struct {
	ID                        string
	Command                   string
	Endpoint                  string
	CmdWithConfig             bool
	WorkDir                   string
	EnvVars                   []string
	Environ                   string
	EngineConfigLocalFile     string
	EngineConfigContainerFile string
	VhostConfigLocalDir       string
	VhostConfigContainerDir   string
	CertLocalDir              string
	CertContainerDir          string

	endpointTLSCert []byte
	endpointTLSKey  []byte
	engineName      string
	templateFile    string
	parsedCommand   *ParsedCommand
}

func (c *CommonConfig) GetIdent() string {
	return c.ID
}

func (c *CommonConfig) GetTemplateFile() string {
	return c.templateFile
}

func (c *CommonConfig) GetEngine() string {
	return c.engineName
}

func (c *CommonConfig) GetEnviron() string {
	return c.Environ
}

func (c *CommonConfig) GetCertLocalDir() string {
	return c.CertLocalDir
}

func (c *CommonConfig) GetCertContainerDir() string {
	return c.CertContainerDir
}

func (c *CommonConfig) GetVhostConfigLocalDir() string {
	return c.VhostConfigLocalDir
}

func (c *CommonConfig) GetVhostConfigContainerDir() string {
	return c.VhostConfigContainerDir
}

func (c *CommonConfig) GetEngineConfigLocalFile() string {
	return c.EngineConfigLocalFile
}

func (c *CommonConfig) GetEngineConfigContainerFile() string {
	return c.EngineConfigContainerFile
}

func (c *CommonConfig) EngineConfigFile() string {
	if c.Environ == EnvironContainer {
		return c.EngineConfigContainerFile
	}
	return c.EngineConfigLocalFile
}

func (c *CommonConfig) VhostConfigDir() string {
	if c.Environ == EnvironContainer {
		return c.VhostConfigContainerDir
	}
	return c.VhostConfigLocalDir
}

func (c *CommonConfig) CopyFrom(m *dbschema.NgingVhostServer) {
	c.ID = m.Ident
	c.Command = m.ExecutableFile
	c.CmdWithConfig = m.CmdWithConfig == common.BoolY
	c.WorkDir = m.WorkDir
	kvs := strings.Split(m.Env, com.StrLF)
	c.EnvVars = make([]string, 0, len(kvs))
	for _, kv := range kvs {
		kv = strings.TrimSpace(kv)
		if len(kv) == 0 {
			continue
		}
		c.EnvVars = append(c.EnvVars, kv)
	}
	c.Environ = m.Environ
	c.EngineConfigLocalFile, _ = filepath.Abs(m.ConfigLocalFile)
	c.EngineConfigContainerFile = m.ConfigContainerFile
	c.VhostConfigLocalDir, _ = filepath.Abs(m.VhostConfigLocalDir)
	c.VhostConfigContainerDir = m.VhostConfigContainerDir
	c.CertLocalDir, _ = filepath.Abs(m.CertLocalDir)
	c.CertContainerDir = m.CertContainerDir
	c.Endpoint = m.Endpoint
	c.endpointTLSCert = com.Str2bytes(m.EndpointTlsCert)
	c.endpointTLSKey = com.Str2bytes(m.EndpointTlsKey)
}

func ParseContainerInfo(parts []string) (string, string) {
	return parts[0], parts[len(parts)-1]
}

func (c *CommonConfig) Exec(ctx context.Context, args ...string) ([]byte, error) {
	if c.Environ == EnvironContainer && len(c.Endpoint) > 0 {
		return c.execEndpoint(ctx, args...)
	}
	return c.execCommand(ctx, args...)
}

func (c *CommonConfig) execEndpoint(ctx context.Context, args ...string) ([]byte, error) {
	client, err := NewAPIClient(c.endpointTLSCert, c.endpointTLSKey)
	if err != nil {
		return nil, err
	}
	if c.parsedCommand == nil {
		c.parsedCommand = &ParsedCommand{
			Command: c.Command,
		}
		rootArgs := com.ParseArgs(c.Command)
		if len(rootArgs) > 1 {
			c.parsedCommand.ContainerEngine = ``
			c.parsedCommand.ContainerName = path.Base(strings.SplitN(c.Endpoint, `/exec`, 2)[0])
			c.parsedCommand.Command = rootArgs[0]
			c.parsedCommand.Args = rootArgs[1:]
		}
	}
	data := RequestDockerExec{
		Cmd: append([]string{c.parsedCommand.Command}, c.parsedCommand.Args...),
		Env: c.EnvVars,
	}
	err = client.Post(c.Endpoint, data)
	return nil, err
}

func (c *CommonConfig) execCommand(ctx context.Context, args ...string) ([]byte, error) {
	command := c.Command
	if c.Environ == EnvironContainer {
		if c.parsedCommand == nil {
			c.parsedCommand = &ParsedCommand{
				Command: c.Command,
			}
			rootArgs := com.ParseArgs(command)
			if len(rootArgs) > 1 {
				c.parsedCommand.ContainerEngine, c.parsedCommand.ContainerName = ParseContainerInfo(rootArgs)
				c.parsedCommand.Command = rootArgs[0]
				c.parsedCommand.Args = rootArgs[1:]
			}
		}
		if len(c.parsedCommand.Command) > 0 {
			command = c.parsedCommand.Command
			rootArgs := append([]string{}, c.parsedCommand.Args...)
			args = append(rootArgs, args...)
		}
	}
	cmd := exec.CommandContext(ctx, command, args...)
	cmd.Dir = c.WorkDir
	cmd.Env = append(cmd.Env, c.EnvVars...)
	cmd.Env = append(cmd.Env, os.Environ()...)
	if stderr := GetCtxStderr(ctx); stderr != nil {
		cmd.Stderr = stderr
	}
	if stdout := GetCtxStdout(ctx); stdout != nil {
		cmd.Stdout = stdout
	}
	if cmd.Stderr == nil && cmd.Stdout == nil {
		result, err := cmd.CombinedOutput()
		if err != nil {
			err = fmt.Errorf(`%s: %w`, result, err)
		}
		return result, err
	}
	err := cmd.Run()
	return nil, err
}

func (c *CommonConfig) RemoveDir(typeName string, rootDir string, prefix string, extensions ...string) error {
	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if len(extensions) > 0 {
			ext := filepath.Ext(path)
			if !com.InSlice(ext, extensions) {
				return nil
			}
		}
		if len(prefix) > 0 && !strings.HasPrefix(info.Name(), prefix) {
			return nil
		}
		log.Info(`Delete the `+typeName+` file: `, path)
		return os.Remove(path)
	})
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
	} else {
		os.Remove(rootDir)
	}
	return err
}

func (c *CommonConfig) FixVhostDirPath(vhostDir string) string {
	dir := vhostDir
	var sep string
	if strings.Contains(dir, `\`) {
		sep = `\`
	} else {
		sep = `/`
	}
	if !strings.HasSuffix(dir, sep) {
		dir += sep
	}
	return dir
}
