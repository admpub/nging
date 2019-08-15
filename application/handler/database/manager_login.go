package database

import (
	"github.com/admpub/nging/application/dbschema"
	"github.com/admpub/nging/application/library/dbmanager"
	"github.com/admpub/nging/application/library/dbmanager/driver"
	"github.com/admpub/nging/application/model"
	"github.com/webx-top/db"
)

func login(mgr dbmanager.Manager, accountID uint, m *model.DbAccount, user *dbschema.User) (err error) {
	ctx := mgr.Context()
	if !ctx.IsPost() {
		return nil
	}
	auth := mgr.Account()
	data := &driver.DbAuth{}
	ctx.Bind(data)
	if len(data.Username) == 0 {
		data.Username = `root`
	}
	if len(data.Host) == 0 {
		data.Host = `127.0.0.1`
	}
	auth.CopyFrom(data)
	if ctx.Form(`remember`) == `1` {
		m.Title = auth.Driver + `://` + auth.Username + `@` + auth.Host + `/` + auth.Db
		m.Engine = auth.Driver
		m.Host = auth.Host
		m.User = auth.Username
		m.Password = auth.Password
		m.Name = auth.Db
		err = m.SetOptions()
		if err != nil {
			return err
		}
		if accountID < 1 || err == db.ErrNoMoreRows {
			m.Uid = user.Id
			_, err = m.Add()
		} else {
			err = m.Edit(accountID, nil, db.Cond{`id`: accountID})
		}
	}
	ctx.Session().Set(`dbAuth`, auth)
	return nil
}
