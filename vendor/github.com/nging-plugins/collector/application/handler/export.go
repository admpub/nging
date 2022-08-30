/*
   Nging is a toolbox for webmasters
   Copyright (C) 2018-present  Wenhui Shen <swh@admpub.com>

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU Affero General Public License as published
   by the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Affero General Public License for more details.

   You should have received a copy of the GNU Affero General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package handler

import (
	"strings"

	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"

	"github.com/admpub/nging/v4/application/handler"

	"github.com/nging-plugins/collector/application/dbschema"
	"github.com/nging-plugins/collector/application/library/collector/export"
	"github.com/nging-plugins/collector/application/model"
	dbmgrmodel "github.com/nging-plugins/dbmanager/application/model"
)

var destTypeInputField = map[string]string{
	`API`:         `api`,
	`DSN`:         `dsn`,
	`dbAccountID`: `dbAccountId`,
}

func Export(c echo.Context) error {
	m := model.NewCollectorExport(c)
	cond := db.Compounds{}
	groupID := c.Formx(`groupId`).Uint()
	if groupID > 0 {
		cond.AddKV(`group_id`, groupID)
	}
	q := c.Formx(`q`).String()
	if len(q) > 0 {
		cond.AddKV(`name`, db.Like(`%`+q+`%`))
	}
	_, err := handler.PagingWithLister(c, handler.NewLister(m, nil, func(r db.Result) db.Result {
		return r.OrderBy(`-id`)
	}, cond.And()))
	ret := handler.Err(c, err)
	rows := m.Objects()
	gIds := []uint{}
	rowAndGroup := make([]*model.CollectorExportAndGroup, len(rows))
	for k, u := range rows {
		rowAndGroup[k] = &model.CollectorExportAndGroup{
			NgingCollectorExport: u,
		}
		if u.GroupId < 1 {
			continue
		}
		if !com.InUintSlice(u.GroupId, gIds) {
			gIds = append(gIds, u.GroupId)
		}
	}

	mg := model.NewCollectorGroup(c)
	var groupList []*dbschema.NgingCollectorGroup
	if len(gIds) > 0 {
		_, err = mg.List(&groupList, nil, 1, 1000, db.And(
			db.Cond{`id IN`: gIds},
			db.Cond{`type`: `export`},
		))
		if err != nil {
			if ret == nil {
				ret = err
			}
		} else {
			for k, v := range rowAndGroup {
				for _, g := range groupList {
					if g.Id == v.GroupId {
						rowAndGroup[k].Group = g
						break
					}
				}
			}
		}
	}
	c.Set(`listData`, rowAndGroup)
	mg.ListByOffset(&groupList, nil, 0, -1, db.Cond{`type`: `export`})
	c.Set(`groupList`, groupList)
	c.Set(`groupId`, groupID)
	return c.Render(`collector/export`, ret)
}

func setDest(ctx echo.Context, m *dbschema.NgingCollectorExport) error {
	if inputName, ok := destTypeInputField[m.DestType]; ok {
		m.Dest = ctx.Form(inputName)
	}
	names := ctx.FormValues(`mapping[name][]`)
	dests := ctx.FormValues(`mapping[dest][]`)
	destCount := len(dests)
	mappings := export.NewMappings()
	for index, name := range names {
		name = strings.TrimSpace(name)
		if len(name) == 0 {
			continue
		}
		if index >= destCount {
			break
		}
		mapping, err := export.NewMapping(name, dests[index])
		if err != nil {
			if err == export.ErrSameTableName {
				return ctx.E(`来源字段中所指定的目地表名称不能与当前字段的目地表名称相同`)
			}
			return err
		}
		mappings.Add(mapping)
	}
	var err error
	m.Mapping, err = mappings.String()
	return err
}

func setExportFormData(ctx echo.Context) {
	dbaM := dbmgrmodel.NewDbAccount(ctx)
	cond := []db.Compound{db.Cond{`engine`: `mysql`}}
	dbaM.ListByOffset(nil, nil, 0, -1, db.And(cond...))
	dbAccountList := dbaM.Objects()
	ctx.Set(`dbAccountList`, dbAccountList)

	pageM := model.NewCollectorPage(ctx)
	cond = []db.Compound{db.Cond{`parent_id`: 0}}
	pageM.ListByOffset(nil, nil, 0, -1, db.And(cond...))
	pageList := pageM.Objects()
	ctx.Set(`pageList`, pageList)
}

func ajaxPageRuleFieldList(ctx echo.Context) error {
	pageID := ctx.Formx(`pageId`).Uint()
	pageM := model.NewCollectorPage(ctx)
	cond := []db.Compound{
		db.Cond{`id`: pageID},
		//db.Cond{`root_id`: pageID},
		//db.Cond{`has_child`: `N`},
	}
	err := pageM.Get(nil, db.And(cond...))
	if err != nil && err != db.ErrNoMoreRows {
		return err
	}
	ruleM := model.NewCollectorRule(ctx)
	cond = []db.Compound{db.Cond{`page_id`: pageM.Id}}
	_, err = ruleM.ListByOffset(nil, nil, 0, -1, db.And(cond...))
	if err != nil && err != db.ErrNoMoreRows {
		return err
	}
	ruleList := ruleM.Objects()
	data := ctx.Data()
	data.SetData(ruleList)
	return ctx.JSON(data)
}

func ajaxChildrenPageList(ctx echo.Context) error {
	pageID := ctx.Formx(`pageId`).Uint()
	pageM := model.NewCollectorPage(ctx)
	cond := []db.Compound{
		db.Cond{`parent_id`: pageID},
	}
	_, err := pageM.ListByOffset(nil, func(r db.Result) db.Result {
		return r.OrderBy(`sort`, `id`)
	}, 0, -1, db.And(cond...))
	if err != nil && err != db.ErrNoMoreRows {
		return err
	}
	pageList := pageM.Objects()
	data := ctx.Data()
	data.SetData(pageList)
	return ctx.JSON(data)
}

func ajaxOperate(ctx echo.Context, operate string) error {
	switch operate {
	case `ruleList`:
		return ajaxPageRuleFieldList(ctx)
	case `childrenPageList`:
		return ajaxChildrenPageList(ctx)
	default:
		return nil
	}
}

func ExportAdd(ctx echo.Context) error {
	operate := ctx.Form(`op`)
	if len(operate) > 0 {
		return ajaxOperate(ctx, operate)
	}
	var err error
	m := model.NewCollectorExport(ctx)
	if ctx.IsPost() {
		err = ctx.MustBind(m.NgingCollectorExport)
		if err == nil {
			err = setDest(ctx, m.NgingCollectorExport)
		}
		if err == nil {
			_, err = m.Add()
		}
		if err == nil {
			handler.SendOk(ctx, ctx.T(`操作成功`))
			return ctx.Redirect(handler.URLFor(`/collector/export`))
		}
	} else {
		id := ctx.Formx(`copyId`).Uint()
		if id > 0 {
			err = m.Get(nil, `id`, id)
			if err == nil {
				echo.StructToForm(ctx, m.NgingCollectorExport, ``, func(topName, fieldName string) string {
					return echo.LowerCaseFirstLetter(topName, fieldName)
				})
				if inputName, ok := destTypeInputField[m.DestType]; ok {
					ctx.Request().Form().Set(inputName, m.NgingCollectorExport.Dest)
				}
				ctx.Request().Form().Set(`id`, `0`)
			}
		}
	}
	ctx.Set(`activeURL`, `/collector/export`)
	mg := model.NewCollectorGroup(ctx)
	if _, e := mg.ListByOffset(nil, nil, 0, -1, db.Cond{`type`: `export`}); e != nil {
		err = e
	}
	ctx.Set(`groupList`, mg.Objects())
	ctx.Set(`data`, dbschema.NewNgingCollectorExport(ctx))
	mappings := export.NewMappings()
	if len(m.Mapping) > 0 {
		err = com.JSONDecode([]byte(m.Mapping), mappings)
	}
	ctx.Set(`fieldList`, mappings)
	setExportFormData(ctx)
	return ctx.Render(`collector/export_edit`, handler.Err(ctx, err))
}

func ExportEdit(ctx echo.Context) error {
	operate := ctx.Form(`op`)
	if len(operate) > 0 {
		return ajaxOperate(ctx, operate)
	}
	id := ctx.Formx(`id`).Uint()
	m := model.NewCollectorExport(ctx)
	//user := handler.User(ctx)
	err := m.Get(nil, `id`, id)
	if err != nil {
		handler.SendFail(ctx, err.Error())
		return ctx.Redirect(handler.URLFor(`/collector/export`))
	}
	if ctx.IsPost() {
		err = ctx.MustBind(m.NgingCollectorExport)
		if err == nil {
			err = setDest(ctx, m.NgingCollectorExport)
		}
		if err == nil {
			m.Id = id
			err = m.Edit(nil, `id`, id)
		}
		if err == nil {
			handler.SendOk(ctx, ctx.T(`修改成功`))
			return ctx.Redirect(handler.URLFor(`/collector/export`))
		}
	} else if ctx.IsAjax() {
		return ExportEditStatus(ctx)
	}
	echo.StructToForm(ctx, m.NgingCollectorExport, ``, echo.LowerCaseFirstLetter)
	if inputName, ok := destTypeInputField[m.DestType]; ok {
		ctx.Request().Form().Set(inputName, m.NgingCollectorExport.Dest)
	}
	mappings := export.NewMappings()
	if len(m.Mapping) > 0 {
		err = com.JSONDecode([]byte(m.Mapping), mappings)
	}
	ctx.Set(`fieldList`, mappings)
	ctx.Set(`activeURL`, `/collector/export`)
	mg := model.NewCollectorGroup(ctx)
	if _, e := mg.ListByOffset(nil, nil, 0, -1, db.Cond{`type`: `export`}); e != nil {
		err = e
	}
	var childrenPageList []*dbschema.NgingCollectorPage
	if m.PageRoot > 0 {
		childM := dbschema.NewNgingCollectorPage(ctx)
		cond := []db.Compound{
			db.Cond{`parent_id`: m.PageRoot},
		}
		childM.ListByOffset(nil, func(r db.Result) db.Result {
			return r.OrderBy(`sort`, `id`)
		}, 0, -1, db.And(cond...))
		childrenPageList = childM.Objects()
	}
	ctx.Set(`groupList`, mg.Objects())
	ctx.Set(`data`, m)
	ctx.Set(`childrenPageList`, childrenPageList)
	setExportFormData(ctx)
	return ctx.Render(`collector/export_edit`, handler.Err(ctx, err))
}

func ExportEditStatus(ctx echo.Context) error {
	disabled := ctx.Form(`disabled`)
	id := ctx.Formx(`id`).Uint()
	m := dbschema.NewNgingCollectorExport(ctx)
	err := m.UpdateField(nil, `disabled`, disabled, db.Cond{`id`: id})
	if err != nil {
		return ctx.JSON(ctx.Data().SetError(err))
	}
	return ctx.JSON(ctx.Data().SetInfo(ctx.T(`操作成功`)))
}

func ExportDelete(ctx echo.Context) error {
	id := ctx.Formx(`id`).Uint()
	m := model.NewCollectorExport(ctx)
	//user := handler.User(ctx)
	err := m.Delete(nil, db.And(
		db.Cond{`id`: id},
		//db.Cond{`uid`: user.Id},
	))
	if err == nil {
		handler.SendOk(ctx, ctx.T(`操作成功`))
	} else {
		handler.SendFail(ctx, err.Error())
	}
	return ctx.Redirect(handler.URLFor(`/collector/export`))
}
