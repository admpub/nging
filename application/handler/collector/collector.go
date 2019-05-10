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

package collector

import (
	"context"
	"strings"

	"github.com/admpub/gopiper"
	"github.com/admpub/nging/application/dbschema"
	"github.com/admpub/nging/application/handler"
	"github.com/admpub/nging/application/library/collector"
	"github.com/admpub/nging/application/library/collector/exec"
	"github.com/admpub/nging/application/library/collector/export"
	"github.com/admpub/nging/application/library/collector/sender"
	"github.com/admpub/nging/application/library/cron"
	"github.com/admpub/nging/application/library/notice"
	"github.com/admpub/nging/application/model"
	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
)

func init() {
	handler.RegisterToGroup(`/collector`, func(g *echo.Group) {
		g.Route(`GET,POST`, `/export`, g.MetaHandler(echo.H{`name`: `导出管理`}, Export))
		g.Route(`GET,POST`, `/export_log`, g.MetaHandler(echo.H{`name`: `日子列表`}, ExportLog))
		g.Route(`GET,POST`, `/export_log_view/:id`, g.MetaHandler(echo.H{`name`: `日志详情`}, ExportLogView))
		g.Route(`GET,POST`, `/export_log_delete`, g.MetaHandler(echo.H{`name`: `删除日志`}, ExportLogDelete))
		g.Route(`GET,POST`, `/export_add`, g.MetaHandler(echo.H{`name`: `添加导出规则`}, ExportAdd))
		g.Route(`GET,POST`, `/export_edit`, g.MetaHandler(echo.H{`name`: `修改导出规则`}, ExportEdit))
		g.Route(`GET,POST`, `/export_edit_status`, g.MetaHandler(echo.H{`name`: `修改导出规则`}, ExportEditStatus))
		g.Route(`GET,POST`, `/export_delete`, g.MetaHandler(echo.H{`name`: `删除导出规则`}, ExportDelete))
		g.Route(`GET,POST`, `/history`, g.MetaHandler(echo.H{`name`: `历史记录`}, History))
		g.Route(`GET,POST`, `/history_view`, g.MetaHandler(echo.H{`name`: `查看历史内容`}, HistoryView))
		g.Route(`GET,POST`, `/history_delete`, g.MetaHandler(echo.H{`name`: `删除历史记录`}, HistoryDelete))
		g.Route(`GET,POST`, `/rule`, g.MetaHandler(echo.H{`name`: `规则列表`}, Rule))
		g.Route(`GET,POST`, `/rule_add`, g.MetaHandler(echo.H{`name`: `添加规则`}, RuleAdd))
		g.Route(`GET,POST`, `/rule_edit`, g.MetaHandler(echo.H{`name`: `修改规则`}, RuleEdit))
		g.Route(`GET,POST`, `/rule_delete`, g.MetaHandler(echo.H{`name`: `删除规则`}, RuleDelete))
		g.Route(`GET,POST`, `/rule_collect`, g.MetaHandler(echo.H{`name`: `采集`}, RuleCollect))
		g.Route(`GET,POST`, `/group`, g.MetaHandler(echo.H{`name`: `任务分组列表`}, Group))
		g.Route(`GET,POST`, `/group_add`, g.MetaHandler(echo.H{`name`: `添加分组`}, GroupAdd))
		g.Route(`GET,POST`, `/group_edit`, g.MetaHandler(echo.H{`name`: `修改分组`}, GroupEdit))
		g.Route(`GET,POST`, `/group_delete`, g.MetaHandler(echo.H{`name`: `删除分组`}, GroupDelete))
		g.Route(`GET,POST`, `/regexp_test`, g.MetaHandler(echo.H{`name`: `测试正则表达式`}, RegexpTest))
	})
	cron.AddSYSJob(`collect_page`, CollectPageJob, `>collect_page:1`, `网页采集`)
}

func Rule(c echo.Context) error {
	m := model.NewCollectorPage(c)
	groupID := c.Formx(`groupId`).Uint()
	cond := []db.Compound{db.Cond{`parent_id`: 0}}
	if groupID > 0 {
		cond = append(cond, db.Cond{`group_id`: groupID})
	}
	page, size, totalRows, p := handler.PagingWithPagination(c)
	cnt, err := m.List(nil, func(r db.Result) db.Result {
		return r.OrderBy(`-id`)
	}, page, size, db.And(cond...))
	ret := handler.Err(c, err)
	if totalRows <= 0 {
		totalRows = int(cnt())
		p.SetRows(totalRows)
	}
	rows := m.Objects()
	gIds := []uint{}
	rowAndGroup := make([]*model.CollectorPageAndGroup, len(rows))
	for k, u := range rows {
		rowAndGroup[k] = &model.CollectorPageAndGroup{
			CollectorPage: u,
		}
		if u.GroupId < 1 {
			continue
		}
		if !com.InUintSlice(u.GroupId, gIds) {
			gIds = append(gIds, u.GroupId)
		}
	}

	mg := model.NewCollectorGroup(c)
	var groupList []*dbschema.CollectorGroup
	if len(gIds) > 0 {
		_, err = mg.List(&groupList, nil, 1, 1000, db.And(
			db.Cond{`id IN`: gIds},
			db.Cond{`type`: `page`},
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
	c.Set(`pagination`, p)
	c.Set(`listData`, rowAndGroup)
	mg.ListByOffset(&groupList, nil, 0, -1, db.Cond{`type`: `page`})
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
		err = c.MustBind(pageM.CollectorPage, func(key string, values []string) (string, []string) {
			if strings.HasPrefix(key, `rule[`) {
				return ``, nil
			}
			if strings.HasPrefix(key, `extra[`) {
				return ``, nil
			}
			return key, values
		})
		if err != nil {
			return c.JSON(result.SetError(err))
		}
		pageM.CollectorPage.Uid = user.Id
		c.Begin()
		pageM.CollectorPage.Id = 0
		_, err = parseFormToDb(c, pageM.CollectorPage, `rule`, false)
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
		err = pageM.CollectorPage.SetField(nil, `root_id`, pageM.Id, `id`, pageM.Id)
		if err != nil {
			c.Rollback()
			return c.JSON(result.SetError(err))
		}
		for key, index := range pages {
			pageData := &dbschema.CollectorPage{
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
				err = pageM.CollectorPage.SetField(nil, `has_child`, `Y`, `id`, pageData.ParentId)
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
		err = c.MustBind(pageM.CollectorPage, func(key string, values []string) (string, []string) {
			if strings.HasPrefix(key, `rule[`) {
				return ``, nil
			}
			if strings.HasPrefix(key, `extra[`) {
				return ``, nil
			}
			return key, values
		})
		if err != nil {
			return c.JSON(result.SetError(err))
		}
		pageM.CollectorPage.Uid = user.Id
		pageM.CollectorPage.Id = id
		c.Begin()
		var rules []*dbschema.CollectorRule
		//保存页面配置和规则
		rules, err = parseFormToDb(c, pageM.CollectorPage, `rule`, true)
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
			pageData := &dbschema.CollectorPage{
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
				err = pageM.CollectorPage.SetField(nil, `has_child`, `Y`, `id`, pageData.ParentId)
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
	clientID := c.Formx(`clientID`).Uint()
	if clientID < 0 {
		return c.E(`clientID值不正确`)
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
		ctx := context.Background()
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
					noticeSender(c.T(`[规则:%d] 采集结束`, id)+`: `+c.T(`强制退出`), 0)
				} else {
					noticeSender(c.T(`[规则:%d] 采集出错`, id)+`: `+err.Error(), 0)
				}
			} else {
				if progress.Total < 0 {
					progress.Total = 0
				}
				progress.Percent = 100
				progress.Finish = progress.Total
				progress.Complete = true
				noticeSender(c.T(`[规则:%d] 采集完毕(%d/%d)`, id, progress.Finish, progress.Total), 1, progress)
			}
		}, ctx)
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
	return c.Render(`/collector/rule_collect`, handler.Err(c, err))
}
