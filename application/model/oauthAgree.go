package model

import (
	"strings"

	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"

	"github.com/admpub/nging/v5/application/dbschema"
)

func NewOAuthAgree(ctx echo.Context) *OAuthAgree {
	m := &OAuthAgree{
		NgingOauthAgree: dbschema.NewNgingOauthAgree(ctx),
	}
	return m
}

type OAuthAgree struct {
	*dbschema.NgingOauthAgree
}

func (f *OAuthAgree) check() error {
	return nil
}

func (f *OAuthAgree) Add() (pk interface{}, err error) {
	pk, err = f.NgingOauthAgree.Insert()
	return
}

func (f *OAuthAgree) Get(uid uint, appID string) error {
	return f.NgingOauthAgree.Get(nil, db.And(
		db.Cond{`uid`: uid},
		db.Cond{`app_id`: appID},
	))
}

func (f *OAuthAgree) Save(uid uint, appID string, scopes []string) error {
	if len(scopes) == 0 {
		return nil
	}
	filteredScopes := []string{}
	existsScopes := map[string]struct{}{}
	filterScopes := func() {
		for _, scope := range scopes {
			scope = strings.TrimSpace(scope)
			if len(scope) == 0 {
				continue
			}
			if _, ok := existsScopes[scope]; ok {
				continue
			}
			existsScopes[scope] = struct{}{}
			filteredScopes = append(filteredScopes, scope)
		}
	}
	err := f.Get(uid, appID)
	if err != nil {
		if err != db.ErrNoMoreRows {
			return err
		}
		f.Uid = uid
		f.AppId = appID
		filterScopes()
		f.Scopes = strings.Join(filteredScopes, `,`)
		_, err = f.Add()
		return err
	}
	if len(f.Scopes) > 0 {
		oldScopes := strings.Split(f.Scopes, `,`)
		for _, scope := range oldScopes {
			scope = strings.TrimSpace(scope)
			if len(scope) == 0 {
				continue
			}
			filteredScopes = append(filteredScopes, scope)
			existsScopes[scope] = struct{}{}
		}
	}
	filterScopes()
	return f.UpdateField(nil, `scopes`, strings.Join(filteredScopes, `,`), db.And(
		db.Cond{`uid`: uid},
		db.Cond{`app_id`: appID},
	))
}

func (f *OAuthAgree) IsAgreed(uid uint, appID string, scopes []string) (bool, error) {
	if len(scopes) == 0 {
		return true, nil
	}
	err := f.Get(uid, appID)
	if err != nil {
		if err != db.ErrNoMoreRows {
			return false, err
		}
		return false, nil
	}
	if len(f.Scopes) == 0 {
		return false, nil
	}
	agreedScopes := strings.Split(f.Scopes, `,`)
	for _, scope := range scopes {
		if !com.InSlice(scope, agreedScopes) {
			return false, nil
		}
	}
	return true, nil
}
