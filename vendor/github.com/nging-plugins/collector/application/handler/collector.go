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
	"github.com/admpub/gopiper"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/defaults"

	"github.com/admpub/nging/v5/application/handler"
	"github.com/admpub/nging/v5/application/library/notice"

	"github.com/nging-plugins/collector/application/dbschema"
	"github.com/nging-plugins/collector/application/library/collector"
	"github.com/nging-plugins/collector/application/library/collector/exec"
	"github.com/nging-plugins/collector/application/library/collector/export"
	"github.com/nging-plugins/collector/application/library/collector/sender"
	"github.com/nging-plugins/collector/application/model"
)

func Rule(c echo.Context) error {
	m := model.NewCollectorPage(c)
	groupID := c.Formx(`groupId`).Uint()
	cond := db.Compounds{db.Cond{`parent_id`: 0}}
	if groupID > 0 {
		cond.AddKV(`group_id`, groupID)
	}
	q := c.Formx(`q`).String()
	if len(q) > 0 {
		cond.AddKV(`name`, db.Like(`%`+q+`%`))
	}
	var rowAndGroup []*model.CollectorPageAndGroup
	page, size, totalRows, p := handler.PagingWithPagination(c)
	cnt, err := m.List(&rowAndGroup, func(r db.Result) db.Result {
		return r.OrderBy(`-id`)
	}, page, size, cond.And())
	ret := handler.Err(c, err)
	if totalRows <= 0 {
		totalRows = int(cnt())
		p.SetRows(totalRows)
	}
	mg := model.NewCollectorGroup(c)
	var groupList []*dbschema.NgingCollectorGroup
	mg.ListByOffset(&groupList, nil, 0, -1, db.Cond{`type`: `page`})
	c.Set(`pagination`, p)
	c.Set(`listData`, rowAndGroup)
	c.Set(`groupList`, groupList)
	c.Set(`groupId`, groupID)
	return c.Render(`collector/rule`, ret)
}

func RuleAdd(c echo.Context) error {
	user := handler.User(c)
	var err error
	pageM := model.NewCollectorPage(c)
	if c.IsPost() {
		result := c.Data()
		err = c.MustBind(pageM.NgingCollectorPage, func(key string, values []string) (string, []string) {
			if key == `Rule` || key == `Extra` {
				return ``, nil
			}
			return key, values
		})
		if err != nil {
			return c.JSON(result.SetError(err))
		}
		pageM.NgingCollectorPage.Uid = user.Id
		c.Begin()
		pageM.NgingCollectorPage.Id = 0
		_, err = parseFormToDb(c, pageM.NgingCollectorPage, `rule`, false)
		if err != nil {
			c.Rollback()
			return c.JSON(result.SetError(err))
		}
		pages := c.FormValues(`extra[index][]`)
		urls := c.FormValues(`extra[enterUrl][]`)
		urlCount := len(urls)
		//browsers := c.FormValues(`extra[browser][]`)
		//browserCount := len(browsers)
		types := c.FormValues(`extra[type][]`)
		typeCount := len(types)
		scopeRules := c.FormValues(`extra[scopeRule][]`)
		scopeRuleCount := len(scopeRules)
		contentTypes := c.FormValues(`extra[contentType][]`)
		contentTypeCount := len(contentTypes)
		charsets := c.FormValues(`extra[charset][]`)
		charsetCount := len(charsets)
		parentID := pageM.Id
		err = pageM.NgingCollectorPage.UpdateField(nil, `root_id`, pageM.Id, `id`, pageM.Id)
		if err != nil {
			c.Rollback()
			return c.JSON(result.SetError(err))
		}
		for key, index := range pages {
			pageData := &dbschema.NgingCollectorPage{
				Uid:      user.Id,
				ParentId: parentID,
				RootId:   pageM.Id,
				GroupId:  pageM.GroupId,
				Sort:     key,
				HasChild: `N`,
			}
			if key >= urlCount {
				break
			}
			pageData.EnterUrl = urls[key]
			/*
				if key >= browserCount {
					break
				}
				pageData.Browser = browsers[key]
			*/
			if key >= typeCount {
				break
			}
			pageData.Type = types[key]

			if key >= contentTypeCount {
				break
			}
			pageData.ContentType = contentTypes[key]

			if key >= scopeRuleCount {
				break
			}
			pageData.ScopeRule = scopeRules[key]

			if key >= charsetCount {
				break
			}
			pageData.Charset = charsets[key]
			pageData.Use(pageM.Trans())
			//extra[rule][{=idx=}]
			_, err = parseFormToDb(c, pageData, `extra[rule][`+index+`]`, false)
			if err == nil {
				err = pageM.NgingCollectorPage.UpdateField(nil, `has_child`, `Y`, `id`, pageData.ParentId)
			}
			if err != nil {
				c.Rollback()
				return c.JSON(result.SetError(err))
			}
			parentID = pageData.Id
		}
		c.End(err == nil)

		return c.JSON(result.SetInfo(c.T(`操作成功`), 1))
	}

	c.Set(`data`, exec.NewRules())
	id := c.Formx(`copyId`).Uint()
	if id > 0 {
		err = pageM.Get(nil, `id`, id)
		if err == nil {
			setFormData(c, pageM)
			c.Request().Form().Set(`id`, `0`)
		}
	}
	mg := model.NewCollectorGroup(c)
	if _, e := mg.ListByOffset(nil, nil, 0, -1, db.Cond{`type`: `page`}); e != nil {
		err = e
	}
	c.Set(`groupList`, mg.Objects())
	c.Set(`activeURL`, `/collector/rule`)
	c.Set(`dataTypes`, dataTypeList())
	c.Set(`browserList`, collector.BrowserKeys())
	c.Set(`allFilter`, gopiper.AllFilter())
	return c.Render(`collector/rule_edit`, handler.Err(c, err))
}

func RuleEdit(c echo.Context) error {
	user := handler.User(c)
	id := c.Formx(`id`).Uint()
	pageM := model.NewCollectorPage(c)
	err := pageM.Get(nil, `id`, id)
	if err != nil {
		if err == db.ErrNoMoreRows {
			err = c.E(`不存在id为%d的数据`, id)
		}
	}
	if c.IsPost() {
		result := c.Data()
		if err != nil {
			return c.JSON(result.SetError(err))
		}
		err = c.MustBind(pageM.NgingCollectorPage, func(key string, values []string) (string, []string) {
			if key == `Rule` || key == `Extra` {
				return ``, nil
			}
			return key, values
		})
		if err != nil {
			return c.JSON(result.SetError(err))
		}
		pageM.NgingCollectorPage.Uid = user.Id
		pageM.NgingCollectorPage.Id = id
		c.Begin()
		var rules []*dbschema.NgingCollectorRule
		//保存页面配置和规则
		rules, err = parseFormToDb(c, pageM.NgingCollectorPage, `rule`, true)
		if err != nil {
			c.Rollback()
			return c.JSON(result.SetError(err))
		}
		ruleIds := []uint{}
		for _, rule := range rules {
			ruleIds = append(ruleIds, rule.Id)
		}
		ruleM := model.NewCollectorRule(c)
		conds := []db.Compound{
			db.Cond{`page_id`: id},
		}
		if len(ruleIds) > 0 {
			conds = append(conds, db.Cond{`id`: db.NotIn(ruleIds)})
		}
		//删除已不再使用的规则
		err = ruleM.Delete(nil, db.And(conds...))
		if err != nil {
			c.Rollback()
			return c.JSON(result.SetError(err))
		}
		pages := c.FormValues(`extra[index][]`)
		pageIds := c.FormxValues(`extra[id][]`).Uint()
		pageIdCount := len(pageIds)
		urls := c.FormValues(`extra[enterUrl][]`)
		urlCount := len(urls)
		//browsers := c.FormValues(`extra[browser][]`)
		//browserCount := len(browsers)
		types := c.FormValues(`extra[type][]`)
		typeCount := len(types)
		scopeRules := c.FormValues(`extra[scopeRule][]`)
		scopeRuleCount := len(scopeRules)
		contentTypes := c.FormValues(`extra[contentType][]`)
		contentTypeCount := len(contentTypes)
		charsets := c.FormValues(`extra[charset][]`)
		charsetCount := len(charsets)
		parentID := pageM.Id
		postPageIds := []uint{}
		for key, index := range pages {
			pageData := &dbschema.NgingCollectorPage{
				Uid:      user.Id,
				ParentId: parentID,
				RootId:   pageM.Id,
				GroupId:  pageM.GroupId,
				Sort:     key,
				HasChild: `N`,
			}
			if key >= pageIdCount {
				break
			}
			pageData.Id = pageIds[key]

			if key >= urlCount {
				break
			}
			pageData.EnterUrl = urls[key]
			/*
				if key >= browserCount {
					break
				}
				pageData.Browser = browsers[key]
			*/
			if key >= typeCount {
				break
			}
			pageData.Type = types[key]

			if key >= contentTypeCount {
				break
			}
			pageData.ContentType = contentTypes[key]

			if key >= scopeRuleCount {
				break
			}
			pageData.ScopeRule = scopeRules[key]

			if key >= charsetCount {
				break
			}
			pageData.Charset = charsets[key]
			pageData.Use(pageM.Trans())
			//extra[rule][{=idx=}]
			//保存页面配置和规则
			rules, err = parseFormToDb(c, pageData, `extra[rule][`+index+`]`, true)
			if err == nil {
				err = pageM.NgingCollectorPage.UpdateField(nil, `has_child`, `Y`, `id`, pageData.ParentId)
			}
			if err != nil {
				c.Rollback()
				return c.JSON(result.SetError(err))
			}
			ruleIds = []uint{}
			for _, rule := range rules {
				ruleIds = append(ruleIds, rule.Id)
			}
			conds = []db.Compound{
				db.Cond{`page_id`: pageData.Id},
			}
			if len(ruleIds) > 0 {
				conds = append(conds, db.Cond{`id`: db.NotIn(ruleIds)})
			}
			//删除已不再使用的规则
			err = ruleM.Delete(nil, db.And(conds...))
			if err != nil {
				c.Rollback()
				return c.JSON(result.SetError(err))
			}
			parentID = pageData.Id
			postPageIds = append(postPageIds, pageData.Id)
		}
		conds = []db.Compound{
			db.Cond{`root_id`: id},
			db.Cond{`parent_id`: db.Gt(0)},
		}
		if len(postPageIds) > 0 {
			conds = append(conds, db.Cond{`id`: db.NotIn(postPageIds)})
		}
		var cnt func() int64
		cnt, err = pageM.ListByOffset(nil, nil, 0, -1, db.And(conds...))
		n := cnt()
		if n > 0 {
			ids := []uint{}
			for _, pageRow := range pageM.Objects() {
				ids = append(ids, pageRow.Id)
			}
			//删除已不再使用的规则
			err = ruleM.Delete(nil, db.Cond{`page_id`: db.In(ids)})
			if err == nil {
				//删除已不再使用的页面配置
				err = pageM.Delete(nil, db.And(conds...))
			}
			if err != nil {
				c.Rollback()
				return c.JSON(result.SetError(err))
			}
		}
		_ = rules
		c.End(err == nil)

		return c.JSON(result.SetInfo(c.T(`修改成功`), 1))
	}

	if err == nil {
		setFormData(c, pageM)
	}

	if err != nil {
		handler.SendFail(c, err.Error())
		return c.Redirect(handler.URLFor(`/collector/rule`))
	}

	mg := model.NewCollectorGroup(c)
	if _, e := mg.ListByOffset(nil, nil, 0, -1, db.Cond{`type`: `page`}); e != nil {
		err = e
	}
	c.Set(`groupList`, mg.Objects())
	c.Set(`activeURL`, `/collector/rule`)
	c.Set(`dataTypes`, dataTypeList())
	c.Set(`browserList`, collector.BrowserKeys())
	c.Set(`allFilter`, gopiper.AllFilter())
	return c.Render(`collector/rule_edit`, handler.Err(c, err))
}

func RuleDelete(c echo.Context) error {
	id := c.Formx(`id`).Uint()
	m := model.NewCollectorPage(c)
	c.Begin()
	err := m.Delete(nil, db.Cond{`id`: id})
	if err == nil {
		ruleM := model.NewCollectorRule(c)
		ruleM.Use(m.Trans())
		_, err = m.ListByOffset(nil, nil, 0, -1, db.Cond{`root_id`: id})
		if err != nil {
			c.Rollback()
			handler.SendFail(c, err.Error())
			return c.Redirect(handler.URLFor(`/collector/rule`))
		}
		ids := []uint{id}
		for _, row := range m.Objects() {
			ids = append(ids, row.Id)
		}
		err = ruleM.Delete(nil, db.Cond{`page_id`: db.In(ids)})
		if err == nil {
			err = m.Delete(nil, db.Cond{`root_id`: id})
		}
		if err != nil {
			c.Rollback()
			handler.SendFail(c, err.Error())
			return c.Redirect(handler.URLFor(`/collector/rule`))
		}
		handler.SendOk(c, c.T(`操作成功`))
	} else {
		handler.SendFail(c, err.Error())
	}

	c.End(err == nil)
	return c.Redirect(handler.URLFor(`/collector/rule`))
}

func RuleCollect(c echo.Context) error {
	var err error
	id := c.Formx(`id`).Int()
	if id < 1 {
		return c.E(`id值不正确`)
	}
	m := model.NewCollectorPage(c)
	err = m.Get(nil, db.Cond{`id`: id})
	if err != nil {
		return err
	}
	collected, err := m.FullData()
	if err != nil {
		return err
	}
	collected.SetExportFn(export.Export)
	user := handler.User(c)
	if c.Format() == `json` {
		data := c.Data()
		op := c.Form(`op`)
		if op == `stop` {
			_, err = Exit(m.Id)
			if err != nil {
				data.SetError(err)
			} else {
				data.SetInfo(c.T(`采集已终止`))
			}
			return c.JSON(data)
		}
		clientID := c.Formx(`clientID`).String()
		if len(clientID) == 0 {
			return c.E(`clientID值不正确`)
		}
		mockCtx := defaults.NewMockContext()
		mockCtx.SetTransaction(c.Transaction())
		mockCtx.SetTranslator(c.Object().Translator)
		err = Go(m.Id, collected, func() {
			var noticeSender sender.Notice
			progress := notice.NewProgress()
			if user != nil {
				notice.OpenMessage(user.Username, `collector`)
				defer notice.CloseMessage(user.Username, `collector`)
				noticeSender = func(message interface{}, statusCode int, progs ...*notice.Progress) error {
					msg := notice.NewMessageWithValue(
						`collector`,
						``,
						message,
						statusCode,
					).SetMode(`element`).SetID(id)
					if len(progs) > 0 && progs[0] != nil {
						progress = progs[0]
					}
					msg.SetProgress(progress).CalcPercent().SetClientID(clientID)
					sendErr := notice.Send(user.Username, msg)
					return sendErr
				}
			} else {
				noticeSender = sender.Default
			}
			_, err = collected.Collect(false, noticeSender, progress)
			if err != nil {
				if exec.ErrForcedExit == err {
					noticeSender(mockCtx.T(`[规则:%d] 采集结束`, id)+`: `+mockCtx.T(`强制退出`), 0)
				} else {
					noticeSender(mockCtx.T(`[规则:%d] 采集出错`, id)+`: `+err.Error(), 0)
				}
			} else {
				if progress.Total < 0 {
					progress.Total = 0
				}
				progress.Percent = 100
				progress.Finish = progress.Total
				progress.Complete = true
				noticeSender(mockCtx.T(`[规则:%d] 采集完毕(%d/%d)`, id, progress.Finish, progress.Total), 1, progress)
			}
		}, mockCtx)
		if err != nil {
			data.SetError(err)
		}
		data.SetInfo(c.T(`[规则:%d] 开始采集中...`, id))
		return c.JSON(data)
	}
	result, err := collected.Collect(true, nil, nil)
	if err != nil {
		return err
	}
	c.Set(`data`, m)
	c.Set(`result`, result)
	c.Set(`activeURL`, `/collector/rule`)
	return c.Render(`collector/rule_collect`, handler.Err(c, err))
}
