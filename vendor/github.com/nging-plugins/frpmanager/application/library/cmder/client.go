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

func NewClient() cmder.Cmder {
	return &FRPClient{
		Base: NewBase(),
	}
}

type FRPClient struct {
	*Base
}

func (c *FRPClient) Init() error {
	id := c.CLIConfig.GenerateIDFromConfigFileName(c.CLIConfig.Confx)
	return frp.StartClientByConfigFile(c.CLIConfig.Confx, c.PidFile(id, false))
}

func (c *FRPClient) StopHistory(ids ...string) error {
	if len(ids) > 0 {
		for _, id := range ids {
			pidPath := c.PidFile(id, false)
			err := com.CloseProcessFromPidFile(pidPath)
			if err != nil {
				log.Error(err.Error() + `: ` + pidPath)
			}
		}
		return nil
	}
	pidFilePath := filepath.Join(echo.Wd(), `data/pid/frp/client`)
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

func (c *FRPClient) StartBy(id uint, writer ...io.Writer) (err error) {
	configFile := c.ConfigFile(id, false)
	params := []string{os.Args[0], `--config`, c.CLIConfig.Conf, `--subconfig`, configFile, `--type`, `frpclient`}
	cmd, iErr, _, rErr := com.RunCmdWithWriterx(params, time.Millisecond*500, writer...)
	key := fmt.Sprintf("frpclient.%d", id)
	c.CLIConfig.CmdSet(key, cmd)
	if iErr != nil {
		err = fmt.Errorf(iErr.Error()+`: %s`, rErr.Buffer().String())
		return
	}
	return
}

func (c *FRPClient) Start(writer ...io.Writer) (err error) {
	err = c.StopHistory()
	if err != nil {
		log.Error(err.Error())
	}
	md := dbschema.NewNgingFrpClient(nil)
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

func (c *FRPClient) Stop() error {
	return c.CLIConfig.CmdGroupStop("frpclient")
}

func (c *FRPClient) Restart(writer ...io.Writer) error {
	err := c.Stop()
	if err != nil {
		return err
	}
	return c.Start(writer...)
}

func (c *FRPClient) RestartBy(id string, writer ...io.Writer) error {
	err := c.CLIConfig.CmdStop("frpclient." + id)
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

func (c *FRPClient) StopBy(id string) error {
	err := c.CLIConfig.CmdStop("frpclient." + id)
	if err != nil && !strings.Contains(err.Error(), `finished`) {
		return err
	}
	pidPath := c.PidFile(id, false)
	err = com.CloseProcessFromPidFile(pidPath)
	if err != nil && !strings.Contains(err.Error(), `finished`) {
		log.Error(err.Error() + `: ` + pidPath)
	}
	return nil
}
