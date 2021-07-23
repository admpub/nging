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

package manager

import (
	"strings"

	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/param"

	"github.com/admpub/nging/v3/application/handler"
	"github.com/admpub/nging/v3/application/model"
	"github.com/admpub/nging/v3/application/registry/alert"
)

func AlertTopic(ctx echo.Context) error {
	m := model.NewAlertTopic(ctx)
	cond := db.Compounds{}
	topic := ctx.Formx(`q`).String()
	if len(topic) == 0 {
		topic = ctx.Formx(`topic`).String()
	} else {
		ctx.Request().Form().Set(`topic`, topic)
	}
	if len(topic) > 0 {
		cond.AddKV(`topic`, topic)
	}
	list := []*model.AlertTopicExt{}
	_, err := handler.PagingWithLister(ctx, handler.NewLister(m, &list, func(r db.Result) db.Result {
		return r.OrderBy(`-id`)
	}, cond.And()))
	ctx.Set(`listData`, list)
	ctx.Set(`title`, ctx.T(`所有专题账号`))
	ctx.Set(`topic`, topic)
	ctx.Set(`topicList`, alert.Topics.Slice())
	ctx.SetFunc(`topicName`, alert.Topics.Get)
	ctx.SetFunc(`platformName`, alert.RecipientPlatforms.Get)
	return ctx.Render(`/manager/alert_topic`, handler.Err(ctx, err))
}

func AlertTopicAdd(ctx echo.Context) error {
	var err error
	if ctx.IsPost() {
		m := model.NewAlertTopic(ctx)
		err = ctx.MustBind(m.NgingAlertTopic)
		recipientIds := ctx.Formx(`recipientIds`).String()
		if len(recipientIds) > 0 {
			recipientIds := param.StringSlice(strings.Split(recipientIds, `,`)).Uint()
			for _, recipientID := range recipientIds {
				row := *m.NgingAlertTopic
				row.RecipientId = recipientID
				_, err = m.Add(&row)
			}
		} else if err == nil {
			_, err = m.Add()
		}
		if err == nil {
			handler.SendOk(ctx, ctx.T(`操作成功`))
			return ctx.Redirect(handler.URLFor(`/manager/alert_recipient`))
		}
	}
	ctx.Set(`activeURL`, `/manager/alert_recipient`)
	ctx.Set(`title`, ctx.T(`添加警报接收人`))
	ctx.Set(`platforms`, alert.RecipientPlatforms.Slice())
	return ctx.Render(`/manager/alert_topic_edit`, handler.Err(ctx, err))
}

func AlertTopicEdit(ctx echo.Context) error {
	id := ctx.Formx(`id`).Uint()
	m := model.NewAlertTopic(ctx)
	err := m.Get(nil, `id`, id)
	if err != nil {
		handler.SendFail(ctx, err.Error())
		return ctx.Redirect(handler.URLFor(`/manager/alert_topic`))
	}
	if ctx.IsPost() {
		err = ctx.MustBind(m.NgingAlertTopic)
		if err == nil {
			m.Id = id
			err = m.Edit(nil, `id`, id)
		}
		if err == nil {
			handler.SendOk(ctx, ctx.T(`修改成功`))
			return ctx.Redirect(handler.URLFor(`/manager/alert_topic`))
		}
	} else if ctx.IsAjax() {
		disabled := ctx.Query(`disabled`)
		if len(disabled) > 0 {
			m.Disabled = disabled
			data := ctx.Data()
			err = m.SetField(nil, `disabled`, disabled, db.Cond{`id`: id})
			if err != nil {
				data.SetError(err)
				return ctx.JSON(data)
			}
			data.SetInfo(ctx.T(`操作成功`))
			return ctx.JSON(data)
		}
	} else {
		echo.StructToForm(ctx, m.NgingAlertTopic, ``, echo.LowerCaseFirstLetter)
	}

	ctx.Set(`activeURL`, `/manager/alert_topic`)
	ctx.Set(`title`, ctx.T(`修改警报接收人`))
	ctx.Set(`platforms`, alert.RecipientPlatforms.Slice())
	return ctx.Render(`/manager/alert_topic_edit`, handler.Err(ctx, err))
}

func AlertTopicDelete(ctx echo.Context) error {
	id := ctx.Formx(`id`).Uint()
	m := model.NewAlertTopic(ctx)
	err := m.Delete(nil, db.Cond{`id`: id})
	if err == nil {
		handler.SendOk(ctx, ctx.T(`操作成功`))
	} else {
		handler.SendFail(ctx, err.Error())
	}
	topic := ctx.Form("topic")
	return ctx.Redirect(handler.URLFor(`/manager/alert_topic?topic=` + topic))
}
