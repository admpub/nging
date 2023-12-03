package manager

import (
	"strings"

	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"

	"github.com/admpub/nging/v5/application/handler"
	"github.com/admpub/nging/v5/application/library/common"
	"github.com/admpub/nging/v5/application/model"
)

// oAuthIndex 应用列表
func oAuthIndex(ctx echo.Context) error {
	var err error
	m := model.NewOAuthApp(ctx)
	cond := db.Compounds{}
	name := ctx.Formx(`q`).String()
	if len(name) > 0 {
		cond.And(db.Or(
			db.Cond{`site_name`: db.Like(name + `%`)},
			db.Cond{`app_id`: name},
		))
	}
	sorts := common.Sorts(ctx, m.NgingOauthApp, `-id`)
	_, err = common.NewLister(m.NgingOauthApp, nil, func(r db.Result) db.Result {
		return r.OrderBy(sorts...)
	}, cond.And()).Paging(ctx)
	ctx.Set(`listData`, m.Objects())
	return ctx.Render(`manager/oauth_app/index`, common.Err(ctx, err))
}

// oAuthAdd 创建应用
func oAuthAdd(ctx echo.Context) error {
	var err error
	var id uint
	m := model.NewOAuthApp(ctx)
	if ctx.IsPost() {
		err = ctx.MustBind(
			m.NgingOauthApp,
			echo.FormatFieldValue(
				map[string]echo.FormDataFilter{
					`scopes`: echo.JoinValues(),
				},
			),
		)
		if err != nil {
			goto END
		}
		_, err = m.Add()
		if err != nil {
			goto END
		}
		common.SendOk(ctx, ctx.T(`添加成功`))
		return ctx.Redirect(handler.URLFor(`/manager/oauth_app/index`))
	}
	id = ctx.Formx(`copyId`).Uint()
	if id > 0 {
		err = m.Get(nil, `id`, id)
		if err == nil {
			echo.StructToForm(ctx, m.NgingOauthApp, ``, echo.LowerCaseFirstLetter)
			ctx.Request().Form().Set(`id`, `0`)
		}
	}

END:
	ctx.Set(`activeURL`, `/manager/oauth_app/index`)
	ctx.Set(`title`, ctx.T(`添加应用`))
	ctx.Set(`isEdit`, false)
	ctx.SetFunc(`isChecked`, func(ident string) bool {
		return false
	})
	return ctx.Render(`manager/oauth_app/edit`, common.Err(ctx, err))
}

// oAuthEdit 修改应用
func oAuthEdit(ctx echo.Context) error {
	id := ctx.Paramx(`id`).Uint64()
	if id < 1 {
		return ctx.E(`参数“%s”值无效`, `id`)
	}
	m := model.NewOAuthApp(ctx)
	err := m.Get(nil, `id`, id)
	if err != nil {
		if err == db.ErrNoMoreRows {
			err = ctx.E(`应用不存在`)
		}
		return err
	}
	if ctx.IsPost() {
		err = ctx.MustBind(
			m.NgingOauthApp,
			echo.ExcludeFieldName(`updated`, `appId`, `appSecret`),
			echo.FormatFieldValue(
				map[string]echo.FormDataFilter{
					`scopes`: echo.JoinValues(),
				},
			),
		)
		if err != nil {
			goto END
		}
		err = m.Edit(nil, `id`, m.Id)
		if err != nil {
			goto END
		}
		common.SendOk(ctx, ctx.T(`修改成功`))
		return ctx.Redirect(handler.URLFor(`/manager/oauth_app/index`))
	} else if ctx.IsAjax() {
		disabled := ctx.Query(`disabled`)
		if len(disabled) > 0 {
			if !common.IsBoolFlag(disabled) {
				return ctx.NewError(code.InvalidParameter, ``).SetZone(`disabled`)
			}
			m.Disabled = disabled
			data := ctx.Data()
			err = m.UpdateField(nil, `disabled`, disabled, db.Cond{`id`: id})
			if err != nil {
				data.SetError(err)
				return ctx.JSON(data)
			}
			data.SetInfo(ctx.T(`操作成功`))
			return ctx.JSON(data)
		}
	}
	echo.StructToForm(ctx, m.NgingOauthApp, ``, echo.LowerCaseFirstLetter)

END:
	ctx.Set(`activeURL`, `/manager/oauth_app/index`)
	ctx.Set(`title`, ctx.T(`修改应用`))
	ctx.Set(`isEdit`, true)
	var scopeIdents []string
	if len(m.Scopes) > 0 {
		scopeIdents = strings.Split(m.Scopes, `,`)
	}
	ctx.SetFunc(`isChecked`, func(ident string) bool {
		for _, idt := range scopeIdents {
			if idt == ident {
				return true
			}
		}
		return false
	})
	return ctx.Render(`manager/oauth_app/edit`, common.Err(ctx, err))
}

// oAuthDelete 删除应用
func oAuthDelete(ctx echo.Context) error {
	id := ctx.Paramx(`id`).Uint64()
	if id < 1 {
		return ctx.E(`参数“%s”值无效`, `id`)
	}
	m := model.NewOAuthApp(ctx)
	err := m.Get(nil, `id`, id)
	if err != nil {
		if err == db.ErrNoMoreRows {
			err = ctx.E(`应用不存在`)
		}
		return err
	}
	err = m.Delete(nil, `id`, id)
	if err != nil {
		return err
	}
	common.SendOk(ctx, ctx.T(`删除成功`))
	return ctx.Redirect(handler.URLFor(`/manager/oauth_app/index`))
}
