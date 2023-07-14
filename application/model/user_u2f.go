package model

import (
	"errors"
	"strings"

	"github.com/admpub/nging/v5/application/dbschema"
	"github.com/webx-top/com"
	"github.com/webx-top/db"
)

func (u *User) NeedCheckU2F(authType string, uid uint, step uint) (need bool, err error) {
	u2f := dbschema.NewNgingUserU2f(u.Context())
	if authType == AuthTypePassword {
		var n int64
		n, err = u2f.Count(nil, db.And(
			db.Cond{`uid`: uid},
			db.Cond{`step`: GetU2FStepCondValue(step)},
		))
		if err != nil {
			if errors.Is(err, db.ErrNoMoreRows) {
				err = nil
			}
			return
		}
		need = n > 0
		return
	}
	err = u2f.Get(func(r db.Result) db.Result {
		return r.Select(`precondition`)
	}, db.And(
		db.Cond{`uid`: uid},
		db.Cond{`step`: GetU2FStepCondValue(step)},
	))
	if err != nil {
		if errors.Is(err, db.ErrNoMoreRows) {
			err = nil
		}
		return
	}
	if len(u2f.Precondition) == 0 {
		return
	}
	parts := strings.Split(u2f.Precondition, `,`)
	need = com.InSlice(authType, parts)
	return
}

func (u *User) GetUserAllU2F(uid uint) ([]*dbschema.NgingUserU2f, error) {
	u2f := dbschema.NewNgingUserU2f(u.Context())
	all := []*dbschema.NgingUserU2f{}
	_, err := u2f.ListByOffset(&all, nil, 0, -1, `uid`, uid)
	return all, err
}

func GetU2FStepCondValue(step uint) interface{} {
	var stepValue interface{}
	if step == 2 {
		stepValue = db.In([]uint{0, 2})
	} else {
		stepValue = step
	}
	return stepValue
}

func (u *User) U2F(uid uint, typ string, step uint) (u2f *dbschema.NgingUserU2f, err error) {
	u2f = dbschema.NewNgingUserU2f(u.Context())
	err = u2f.Get(nil, db.And(
		db.Cond{`uid`: uid},
		db.Cond{`type`: typ},
		db.Cond{`step`: GetU2FStepCondValue(step)},
	))
	return
}
