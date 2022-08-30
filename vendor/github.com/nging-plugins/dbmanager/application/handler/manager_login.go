package handler

import (
	"github.com/webx-top/db"
	"github.com/webx-top/echo"

	ngingdbschema "github.com/admpub/nging/v4/application/dbschema"

	"github.com/nging-plugins/dbmanager/application/library/dbmanager"
	"github.com/nging-plugins/dbmanager/application/library/dbmanager/driver"
	"github.com/nging-plugins/dbmanager/application/model"
)

func addAuth(ctx echo.Context, auth *driver.DbAuth) {
	accounts := getAccounts(ctx)
	accounts.Add(auth)
	ctx.Session().Set(`dbAccounts`, accounts)
}

func getAccounts(ctx echo.Context) driver.AuthAccounts {
	if accounts, ok := ctx.Internal().Get(`dbAccounts`).(*driver.AuthAccounts); ok {
		return *accounts
	}
	accounts, ok := ctx.Session().Get(`dbAccounts`).(driver.AuthAccounts)
	if !ok {
		accounts = driver.AuthAccounts{}
	} else {
		ctx.Internal().Set(`dbAccounts`, &accounts)
	}
	return accounts
}

func deleteAuth(ctx echo.Context, auth *driver.DbAuth) {
	accounts := getAccounts(ctx)
	accounts.Delete(auth)
	ctx.Session().Set(`dbAccounts`, accounts)
}

func clearAuth(ctx echo.Context) {
	accounts := getAccounts(ctx)
	for key := range accounts {
		accounts.DeleteByKey(key)
	}
	ctx.Session().Delete(`dbAccounts`)
}

func getLoginInfo(mgr dbmanager.Manager, accountID uint, m *model.DbAccount, user *ngingdbschema.NgingUser) (err error) {
	ctx := mgr.Context()
	if !ctx.IsPost() {
		return nil
	}
	auth := mgr.Account()
	data := &driver.DbAuth{AccountID: accountID, AccountTitle: m.Title}
	ctx.Bind(data)
	if len(data.Username) == 0 {
		data.Username = `root`
	}
	if len(data.Host) == 0 {
		data.Host = `127.0.0.1`
	}
	auth.CopyFrom(data)
	if ctx.Form(`remember`) == `1` {
		if len(m.Title) == 0 {
			m.Title = auth.Driver + `://` + auth.Username + `@` + auth.Host + `/` + auth.Db
		}
		m.Engine = auth.Driver
		m.Host = auth.Host
		m.User = auth.Username
		m.Password = auth.Password
		m.Name = auth.Db
		err = m.SetOptions()
		if err != nil {
			return err
		}
		if accountID < 1 || m.Id < 1 {
			m.Uid = user.Id
			_, err = m.Add()
		} else {
			err = m.Edit(accountID, nil, db.Cond{`id`: accountID})
		}
	}
	return
}
