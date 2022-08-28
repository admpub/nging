package model

import (
	"time"

	"github.com/admpub/nging/v4/application/handler"
	"github.com/admpub/nging/v4/application/library/common"
	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
)

func (u *User) CheckPasswd(username string, password string) (exists bool, err error) {
	exists = true
	err = u.Get(nil, `username`, username)
	if err != nil {
		if err == db.ErrNoMoreRows {
			exists = false
		}
		return
	}
	if u.NgingUser.Disabled == `Y` {
		err = u.Context().NewError(code.UserDisabled, `该用户已被禁用`).SetZone(`disabled`)
		return
	}
	if u.NgingUser.Password != com.MakePassword(password, u.NgingUser.Salt) {
		err = u.Context().NewError(code.InvalidParameter, `密码不正确`).SetZone(`password`)
	}
	return
}

func (u *User) FireLoginSuccess() error {
	c := u.Context()
	loginLogM := u.NewLoginLog(u.Username)
	u.NgingUser.LastLogin = uint(time.Now().Unix())
	u.NgingUser.LastIp = c.RealIP()
	set := echo.H{
		`last_login`:  u.NgingUser.LastLogin,
		`login_fails`: 0,
	}
	if !common.IsAnonymousMode(`user`) {
		set[`last_ip`] = u.NgingUser.LastIp
	}
	if len(u.NgingUser.SessionId) > 0 {
		if u.NgingUser.SessionId != loginLogM.SessionId {
			c.Session().RemoveID(u.NgingUser.SessionId)
			set.Set(`session_id`, loginLogM.SessionId)
		}
	} else {
		set.Set(`session_id`, loginLogM.SessionId)
	}

	// update user data
	u.NgingUser.UpdateFields(nil, set, `id`, u.NgingUser.Id)

	// loging
	loginLogM.OwnerId = uint64(u.Id)
	loginLogM.Success = `Y`
	loginLogM.AddAndSaveSession()

	// session
	u.SetSession()
	if u.NeedCheckU2F(u.NgingUser.Id, 2) {
		c.Session().Set(`auth2ndURL`, handler.URLFor(`/gauth_check`))
	}
	return nil
}

func (u *User) FireLoginFailure(pass string, err error) error {
	if !echo.IsErrorCode(err, code.UserDisabled) {
		// 仅记录密码不正确的情况
		loginLogM := u.NewLoginLog(u.Username)
		loginLogM.Errpwd = pass
		loginLogM.Failmsg = err.Error()
		loginLogM.Success = `N`
		loginLogM.Add()
		u.IncrLoginFails()
	}
	return nil
}
