package model

import (
	"github.com/admpub/goth"
	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"

	"github.com/admpub/nging/v5/application/dbschema"
)

func NewUserOAuth(ctx echo.Context) *UserOAuth {
	m := &UserOAuth{
		NgingUserOauth: dbschema.NewNgingUserOauth(ctx),
	}
	return m
}

type UserOAuth struct {
	*dbschema.NgingUserOauth
}

func (f *UserOAuth) Add() (pk interface{}, err error) {
	old := dbschema.NewNgingUserOauth(f.Context())
	err = old.Get(nil, db.And(
		db.Cond{`uid`: f.Uid},
		db.Cond{`union_id`: f.UnionId},
		db.Cond{`open_id`: f.OpenId},
		db.Cond{`type`: f.Type},
	))
	if err == nil {
		pk = old.Id
		set := echo.H{}
		if len(f.Email) > 0 && old.Email != f.Email {
			set[`email`] = f.Email
		}
		if len(f.Mobile) > 0 && old.Mobile != f.Mobile {
			set[`mobile`] = f.Mobile
		}
		if len(f.Avatar) > 0 && old.Avatar != f.Avatar {
			set[`avatar`] = f.Avatar
		}
		if len(f.AccessToken) > 0 && old.AccessToken != f.AccessToken {
			set[`access_token`] = f.AccessToken
		}
		if len(f.RefreshToken) > 0 && old.RefreshToken != f.RefreshToken {
			set[`refresh_token`] = f.RefreshToken
		}
		if len(f.Name) > 0 && old.Name != f.Name {
			set[`name`] = f.Name
		}
		if len(f.NickName) > 0 && old.NickName != f.NickName {
			set[`nick_name`] = f.NickName
		}
		if f.Expired > 0 && old.Expired != f.Expired {
			set[`expired`] = f.Expired
		}
		if len(set) == 0 {
			return
		}
		err = f.UpdateFields(nil, set, `id`, old.Id)
		return
	}
	if err != db.ErrNoMoreRows {
		return
	}
	return f.NgingUserOauth.Insert()
}

func (f *UserOAuth) Upsert(mw func(db.Result) db.Result, args ...interface{}) (interface{}, error) {
	return f.NgingUserOauth.Upsert(mw, args...)
}

func (f *UserOAuth) Edit(mw func(db.Result) db.Result, args ...interface{}) error {
	return f.NgingUserOauth.Update(mw, args...)
}

func (f *UserOAuth) GetByOutUser(user *goth.User) (err error) {
	var unionID string
	if v, y := user.RawData[`unionid`]; y {
		unionID, _ = v.(string)
	}
	return f.Get(nil, db.And(
		db.Cond{`union_id`: unionID},
		db.Cond{`open_id`: user.UserID},
		db.Cond{`type`: user.Provider},
	))
}

func (f *UserOAuth) CopyFrom(user *goth.User) *UserOAuth {
	var unionID string
	if v, y := user.RawData[`unionid`]; y {
		unionID, _ = v.(string)
	}
	if v, y := user.RawData[`mobile`]; y {
		f.Mobile, _ = v.(string)
	}
	f.UnionId = unionID
	f.OpenId = user.UserID
	f.Type = user.Provider
	f.Avatar = user.AvatarURL
	f.Name = com.Substr(user.Name, ``, 30)
	f.NickName = com.Substr(user.NickName, ``, 30)
	f.Email = user.Email
	f.AccessToken = user.AccessToken
	f.RefreshToken = user.RefreshToken
	f.Expired = uint(user.ExpiresAt.Unix())
	return f
}

func (f *UserOAuth) Exists(uid uint64, unionID string, openID string, typ string) (bool, error) {
	return f.NgingUserOauth.Exists(nil, db.And(
		db.Cond{`uid`: uid},
		db.Cond{`union_id`: unionID},
		db.Cond{`open_id`: openID},
		db.Cond{`type`: typ},
	))
}

func (f *UserOAuth) ExistsOtherBinding(uid uint64, id uint64) (bool, error) {
	return f.NgingUserOauth.Exists(nil, db.And(
		db.Cond{`uid`: uid},
		db.Cond{`id`: db.NotEq(id)},
	))
}

func (f *UserOAuth) OAuthUserGender(ouser *goth.User) string {
	if v, y := ouser.RawData[`gender`]; y {
		gender := com.String(v)
		if len(gender) > 0 {
			switch gender[0] {
			case 'F', 'f', '0':
				return `female`
			case 'M', 'm', '1':
				return `male`
			default:
				return `secret` //保密
			}
		}
	}
	return ``
}
