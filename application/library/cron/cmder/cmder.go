package cmder

import (
	"context"
	"io"

	"github.com/admpub/log"

	"github.com/admpub/nging/v4/application/library/config/cmder"
	"github.com/admpub/nging/v4/application/library/cron"
)

func Get() cmder.Cmder {
	return cmder.Get(`task`)
}

func New() cmder.Cmder {
	return &taskCmd{}
}

type taskCmd struct {
}

func (c *taskCmd) Init() error {
	return nil
}

func (c *taskCmd) StopHistory(_ ...string) error {
	return nil
}

func (c *taskCmd) Start(writer ...io.Writer) error {
	if err := cron.InitJobs(context.Background()); err != nil {
		log.Error(err)
	}
	return nil
}

func (c *taskCmd) Stop() error {
	cron.Close()
	return nil
}

func (c *taskCmd) Reload() error {
	return nil
}

func (c *taskCmd) Restart(writer ...io.Writer) error {
	c.Stop()
	return c.Start(writer...)
}
