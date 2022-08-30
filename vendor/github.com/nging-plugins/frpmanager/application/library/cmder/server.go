package cmder

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/admpub/log"
	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"

	"github.com/admpub/nging/v4/application/library/config/cmder"

	"github.com/nging-plugins/frpmanager/application/dbschema"
	"github.com/nging-plugins/frpmanager/application/library/frp"
)

func NewServer() cmder.Cmder {
	return &FRPServer{
		Base: NewBase(),
	}
}

type FRPServer struct {
	*Base
}

func (c *FRPServer) Init() error {
	id := c.CLIConfig.GenerateIDFromConfigFileName(c.CLIConfig.Confx)
	return frp.StartServerByConfigFile(c.CLIConfig.Confx, c.PidFile(id, true))
}

func (c *FRPServer) StopHistory(ids ...string) error {
	if len(ids) > 0 {
		for _, id := range ids {
			pidPath := c.PidFile(id, true)
			err := com.CloseProcessFromPidFile(pidPath)
			if err != nil {
				log.Error(err.Error() + `: ` + pidPath)
			}
		}
		return nil
	}
	pidFilePath := filepath.Join(echo.Wd(), `data/pid/frp/server`)
	err := filepath.Walk(pidFilePath, func(pidPath string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}
		err = com.CloseProcessFromPidFile(pidPath)
		if err != nil {
			log.Error(err.Error() + `: ` + pidPath)
		}
		return os.Remove(pidPath)
	})
	return err
}

func (c *FRPServer) StartBy(id uint, writer ...io.Writer) (err error) {
	configFile := c.ConfigFile(id, true)
	params := []string{os.Args[0], `--config`, c.CLIConfig.Conf, `--subconfig`, configFile, `--type`, `frpserver`}
	cmd, iErr, _, rErr := com.RunCmdWithWriterx(params, time.Millisecond*500, writer...)
	key := fmt.Sprintf("frpserver.%d", id)
	c.CLIConfig.CmdSet(key, cmd)
	if iErr != nil {
		err = fmt.Errorf(iErr.Error()+`: %s`, rErr.Buffer().String())
		return
	}
	return
}

func (c *FRPServer) Start(writer ...io.Writer) (err error) {
	err = c.StopHistory()
	if err != nil {
		log.Error(err.Error())
	}
	md := dbschema.NewNgingFrpServer(nil)
	cd := db.And(
		db.Cond{`disabled`: `N`},
	)
	_, err = md.ListByOffset(nil, nil, 0, -1, cd)
	if err != nil {
		if err == db.ErrNoMoreRows {
			return nil
		}
		return
	}
	for _, row := range md.Objects() {
		err = c.StartBy(row.Id, writer...)
		if err != nil {
			return
		}
	}
	return nil
}

func (c *FRPServer) Stop() error {
	return c.CLIConfig.CmdGroupStop("frpserver")
}

func (c *FRPServer) Restart(writer ...io.Writer) error {
	err := c.Stop()
	if err != nil {
		return err
	}
	return c.Start(writer...)
}

func (c *FRPServer) RestartBy(id string, writer ...io.Writer) error {
	err := c.CLIConfig.CmdStop("frpserver." + id)
	if err != nil {
		return err
	}
	pidPath := c.PidFile(id, true)
	err = com.CloseProcessFromPidFile(pidPath)
	if err != nil {
		log.Error(err.Error() + `: ` + pidPath)
	}
	idv, _ := strconv.ParseUint(id, 10, 32)
	return c.StartBy(uint(idv), writer...)
}

func (c *FRPServer) StopBy(id string) error {
	err := c.CLIConfig.CmdStop("frpserver." + id)
	if err != nil && !strings.Contains(err.Error(), `finished`) {
		return err
	}
	pidPath := c.PidFile(id, true)
	err = com.CloseProcessFromPidFile(pidPath)
	if err != nil && !strings.Contains(err.Error(), `finished`) {
		log.Error(err.Error() + `: ` + pidPath)
	}
	return nil
}
