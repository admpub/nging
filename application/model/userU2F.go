package model

import (
	"github.com/admpub/nging/v4/application/dbschema"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
)

func NewUserU2F(ctx echo.Context) *UserU2F {
	m := &UserU2F{
		NgingUserU2f: dbschema.NewNgingUserU2f(ctx),
	}
	return m
}

type UserU2F struct {
	*dbschema.NgingUserU2f
}

func (u *UserU2F) check() error {
	exists, err := u.Exists(nil, db.And(
		db.Cond{`uid`: u.Uid},
		db.Cond{`type`: u.Type},
		db.Cond{`step`: u.Step},
		db.Cond{`token`: u.Token},
	))
	if err != nil {
		return err
	}
	if exists {
		err = u.Context().NewError(code.DataAlreadyExists, `Token已经存在`).SetZone(`token`)
	}
	return err
}

func (u *UserU2F) Add() (interface{}, error) {
	if err := u.check(); err != nil {
		return nil, err
	}
	return u.NgingUserU2f.Insert()
}

func (u *UserU2F) HasType(uid uint, authType string, step uint) (bool, error) {
	return u.NgingUserU2f.Exists(nil, db.And(
		db.Cond{`uid`: uid},
		db.Cond{`type`: authType},
		db.Cond{`step`: GetU2FStepCondValue(step)},
	))
}

func (u *UserU2F) Unbind(uid uint, typ string, step uint) error {
	return u.NgingUserU2f.Delete(nil, db.And(
		db.Cond{`uid`: uid},
		db.Cond{`type`: typ},
		db.Cond{`step`: GetU2FStepCondValue(step)},
	))
}

func (u *UserU2F) UnbindByToken(uid uint, typ string, step uint, token string) error {
	return u.NgingUserU2f.Delete(nil, db.And(
		db.Cond{`uid`: uid},
		db.Cond{`type`: typ},
		db.Cond{`step`: GetU2FStepCondValue(step)},
		db.Cond{`token`: token},
	))
}

func (u *UserU2F) ListPageByType(uid uint, typ string, step uint, sorts ...interface{}) error {
	cond := db.NewCompounds()
	cond.AddKV(`uid`, uid)
	cond.AddKV(`type`, typ)
	cond.AddKV(`step`, GetU2FStepCondValue(step))
	return u.ListPage(cond, sorts...)
}
