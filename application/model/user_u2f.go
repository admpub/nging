package model

import (
	"github.com/admpub/nging/v5/application/dbschema"
	"github.com/webx-top/db"
)

func (u *User) NeedCheckU2F(uid uint, step uint) bool {
	u2f := dbschema.NewNgingUserU2f(u.Context())
	n, _ := u2f.Count(nil, db.And(
		db.Cond{`uid`: uid},
		db.Cond{`step`: GetU2FStepCondValue(step)},
	))
	return n > 0
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
