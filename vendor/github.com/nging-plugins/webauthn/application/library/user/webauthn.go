package user

import (
	"encoding/json"

	cw "github.com/coscms/webauthn"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"

	"github.com/admpub/nging/v5/application/dbschema"
	"github.com/admpub/nging/v5/application/handler"
	"github.com/admpub/nging/v5/application/model"
	"github.com/nging-plugins/webauthn/application/library/common"
)

var handle cw.UserHandler = &UserHandle{}

type UserHandle struct {
}

func (u *UserHandle) GetUser(ctx echo.Context, username string, opType cw.Type, stage cw.Stage) (webauthn.User, error) {
	if opType == cw.TypeRegister || opType == cw.TypeUnbind {
		user := handler.User(ctx)
		if user == nil {
			return nil, ctx.NewError(code.Unauthenticated, `请先登录`)
		}
		if username != user.Username {
			return nil, ctx.NewError(code.NonPrivileged, `用户名不匹配`)
		}
	}
	userM := model.NewUser(ctx)
	err := userM.Get(func(r db.Result) db.Result {
		return r.Select(`id`, `username`, `avatar`, `disabled`)
	}, `username`, username)
	if err != nil {
		if err == db.ErrNoMoreRows {
			err = ctx.NewError(code.UserNotFound, `用户不存在`).SetZone(`username`)
		}
		return nil, err
	}
	if userM.Disabled == `Y` {
		err = ctx.NewError(code.UserDisabled, `该用户已被禁用`).SetZone(`disabled`)
		return nil, err
	}
	user := &cw.User{
		ID:          uint64(userM.Id),
		Name:        userM.Username,
		DisplayName: userM.Username,
		Icon:        userM.Avatar,
	}
	u2f := dbschema.NewNgingUserU2f(ctx)
	_, err = u2f.ListByOffset(nil, nil, 0, -1, db.And(
		db.Cond{`uid`: userM.Id},
		db.Cond{`type`: `webauthn`},
		db.Cond{`step`: 1},
	))
	if err != nil {
		return nil, err
	}
	u2fList := u2f.Objects()
	if len(u2fList) == 0 {
		err = ctx.NewError(code.Unsupported, `该用户不支持免密登录`)
		return nil, err
	}
	user.Credentials = make([]webauthn.Credential, len(u2fList))
	for index, row := range u2fList {
		cred := webauthn.Credential{}
		err = json.Unmarshal([]byte(row.Extra), &cred)
		if err != nil {
			return nil, err
		}
		user.Credentials[index] = cred
	}
	if opType == cw.TypeUnbind && stage == cw.StageBegin {
		unbind := ctx.Form(`unbind`)
		ctx.Session().Set(common.SessionKeyUnbindToken, unbind)
	}
	return user, nil
}

func (u *UserHandle) Register(ctx echo.Context, user webauthn.User, cred *webauthn.Credential) error {
	userM := model.NewUser(ctx)
	err := userM.Get(func(r db.Result) db.Result {
		return r.Select(`id`, `disabled`)
	}, `username`, user.WebAuthnName())
	if err != nil {
		return err
	}
	if userM.Disabled == `Y` {
		err = ctx.NewError(code.UserDisabled, `该用户已被禁用`).SetZone(`disabled`)
		return err
	}
	u2fM := model.NewUserU2F(ctx)
	u2fM.Uid = userM.Id
	u2fM.Token = com.ByteMd5(cred.ID)
	u2fM.Name = common.GetOS(ctx.Request().UserAgent())
	b, err := json.Marshal(cred)
	if err != nil {
		return err
	}
	u2fM.Extra = string(b)
	u2fM.Type = `webauthn`
	u2fM.Step = 1
	_, err = u2fM.Add()
	return err
}

func (u *UserHandle) Login(ctx echo.Context, user webauthn.User, cred *webauthn.Credential) error {
	userM := model.NewUser(ctx)
	err := userM.Get(nil, `username`, user.WebAuthnName())
	if err != nil {
		return err
	}
	err = userM.FireLoginSuccess(`webauthn`)
	//userM.SetSession()
	return err
}

func (u *UserHandle) Unbind(ctx echo.Context, user webauthn.User, cred *webauthn.Credential) error {
	userM := model.NewUser(ctx)
	err := userM.Get(nil, `username`, user.WebAuthnName())
	if err != nil {
		return err
	}
	u2fM := model.NewUserU2F(ctx)
	unbind, _ := ctx.Session().Get(common.SessionKeyUnbindToken).(string)
	err = u2fM.UnbindByToken(userM.Id, `webauthn`, 1, unbind)
	if err == nil {
		ctx.Session().Delete(common.SessionKeyUnbindToken)
	}
	return err
}
