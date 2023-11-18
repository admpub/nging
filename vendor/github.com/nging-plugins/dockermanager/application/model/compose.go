package model

import (
	"github.com/nging-plugins/dockermanager/application/dbschema"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
)

func NewStack(ctx echo.Context) *Compose {
	m := dbschema.NewNgingDockerCompose(ctx)
	m.Type = TypeStack
	return &Compose{
		NgingDockerCompose: m,
		_type:              m.Type,
	}
}

func NewCompose(ctx echo.Context) *Compose {
	m := dbschema.NewNgingDockerCompose(ctx)
	m.Type = TypeCompose
	return &Compose{
		NgingDockerCompose: m,
		_type:              m.Type,
	}
}

type Compose struct {
	*dbschema.NgingDockerCompose
	_type string
}

func (c *Compose) Exists(excludeID ...uint) (*dbschema.NgingDockerCompose, error) {
	cond := db.NewCompounds()
	cond.AddKV(`type`, c._type)
	cond.AddKV(`name`, c.Name)
	if len(excludeID) > 0 {
		cond.AddKV(`id`, db.NotEq(excludeID[0]))
	}
	bean := dbschema.NewNgingDockerCompose(c.Context())
	err := bean.Get(nil, cond.And())
	if err != nil {
		if err == db.ErrNoMoreRows {
			return nil, nil
		}
		return nil, err
	}
	return bean, err
}

func (c *Compose) check() error {
	if len(c.Name) == 0 {
		return c.Context().NewError(code.InvalidParameter, `参数 %s 无效`, `name`).SetZone(`name`)
	}
	return nil
}

func (c *Compose) Add() (interface{}, error) {
	if err := c.check(); err != nil {
		return nil, err
	}
	old, err := c.Exists()
	if err != nil {
		return nil, err
	}
	if old != nil {
		cond := db.NewCompounds()
		cond.AddKV(`type`, c._type)
		cond.AddKV(`name`, c.Name)
		c.Id = old.Id
		c.Created = old.Created
		return nil, c.NgingDockerCompose.Update(nil, `id`, c.Id)
	}
	return c.NgingDockerCompose.Insert()
}

func (c *Compose) Edit(mw func(db.Result) db.Result, args ...interface{}) error {
	if err := c.check(); err != nil {
		return err
	}
	return c.NgingDockerCompose.Update(mw, args...)
}

func (c *Compose) ListPage(cond *db.Compounds, sorts ...interface{}) error {
	cond.AddKV(`type`, c.Type)
	return c.NgingDockerCompose.ListPage(cond, sorts...)
}

func (c *Compose) Reset() *Compose {
	c.NgingDockerCompose.Reset()
	c.Type = c._type
	return c
}
