/*
   Nging is a toolbox for webmasters
   Copyright (C) 2018-present Wenhui Shen <swh@admpub.com>

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
	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"

	"github.com/admpub/nging/v5/application/cmd/bootconfig"
	"github.com/admpub/nging/v5/application/handler"
	"github.com/admpub/nging/v5/application/model"
	"github.com/admpub/nging/v5/application/registry/alert"
)

func AlertRecipient(ctx echo.Context) error {
	m := model.NewAlertRecipient(ctx)
	cond := db.Compounds{}
	q := ctx.Formx(`q`).String()
	if len(q) > 0 {
		cond.AddKV(`name`, q)
	}
	excludeTopic := ctx.Formx(`excludeTopic`).String()
	if len(excludeTopic) > 0 {
		topicM := model.NewAlertTopic(ctx)
		cond.Add(db.Raw("NOT EXISTS (SELECT 1 FROM `"+topicM.Name_()+"` WHERE `topic`=? AND `recipient_id`=`"+m.Name_()+"`.`id`)", excludeTopic))
	}
	_, err := handler.PagingWithLister(ctx, handler.NewLister(m, nil, func(r db.Result) db.Result {
		return r.OrderBy(`-id`)
	}, cond.And()))
	ctx.Set(`listData`, m.Objects())
	ctx.Set(`title`, ctx.T(`警报收信账号`))
	ctx.SetFunc(`platformName`, alert.RecipientPlatforms.Get)
	ctx.Set(`topicList`, alert.Topics.Slice())
	ctx.SetFunc(`topicName`, alert.Topics.Get)
	return ctx.Render(`/manager/alert_recipient`, handler.Err(ctx, err))
}

func AlertRecipientAdd(ctx echo.Context) error {
	var err error
	m := model.NewAlertRecipient(ctx)
	if ctx.IsPost() {
		err = ctx.MustBind(m.NgingAlertRecipient)
		if err == nil {
			_, err = m.Add()
		}
		if err == nil {
			handler.SendOk(ctx, ctx.T(`操作成功`))
			return ctx.Redirect(handler.URLFor(`/manager/alert_recipient`))
		}
	} else {
		id := ctx.Formx(`copyId`).Uint()
		if id > 0 {
			err = m.Get(nil, db.Cond{`id`: id})
			if err == nil {
				echo.StructToForm(ctx, m.NgingAlertRecipient, ``, echo.LowerCaseFirstLetter)
				ctx.Request().Form().Set(`id`, `0`)
			}
		}
	}
	ctx.Set(`activeURL`, `/manager/alert_recipient`)
	ctx.Set(`title`, ctx.T(`添加收信账号`))
	platforms := alert.RecipientPlatforms.Slice()
	ctx.Set(`platforms`, platforms)
	var defaultPlatform string
	if len(platforms) > 0 {
		defaultPlatform = platforms[0].K
	}
	ctx.Set(`defaultPlatform`, defaultPlatform)
	webhookCustom := alert.NewWebhookCustom()
	b, _ := com.JSONEncode(webhookCustom, `  `)
	ctx.Set(`webhookCustomDefault`, string(b))
	ctx.Set(`webhookCustomDescriptions`, webhookCustom.Descriptions())
	return ctx.Render(`/manager/alert_recipient_edit`, handler.Err(ctx, err))
}

func AlertRecipientEdit(ctx echo.Context) error {
	id := ctx.Formx(`id`).Uint()
	m := model.NewAlertRecipient(ctx)
	err := m.Get(nil, `id`, id)
	if err != nil {
		handler.SendFail(ctx, err.Error())
		return ctx.Redirect(handler.URLFor(`/manager/alert_recipient`))
	}
	if ctx.IsPost() {
		err = ctx.MustBind(m.NgingAlertRecipient)
		if err == nil {
			m.Id = id
			err = m.Edit(nil, `id`, id)
		}
		if err == nil {
			handler.SendOk(ctx, ctx.T(`修改成功`))
			return ctx.Redirect(handler.URLFor(`/manager/alert_recipient`))
		}
	} else if ctx.IsAjax() {
		disabled := ctx.Query(`disabled`)
		if len(disabled) > 0 {
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
	} else {
		echo.StructToForm(ctx, m.NgingAlertRecipient, ``, echo.LowerCaseFirstLetter)
	}

	ctx.Set(`activeURL`, `/manager/alert_recipient`)
	ctx.Set(`title`, ctx.T(`修改收信账号`))
	platforms := alert.RecipientPlatforms.Slice()
	ctx.Set(`platforms`, platforms)
	var defaultPlatform string
	if len(platforms) > 0 {
		defaultPlatform = platforms[0].K
	}
	ctx.Set(`defaultPlatform`, defaultPlatform)
	webhookCustom := alert.NewWebhookCustom()
	b, _ := com.JSONEncode(webhookCustom, `  `)
	ctx.Set(`webhookCustomDefault`, string(b))
	ctx.Set(`webhookCustomDescriptions`, webhookCustom.Descriptions())
	return ctx.Render(`/manager/alert_recipient_edit`, handler.Err(ctx, err))
}

func AlertRecipientTest(ctx echo.Context) error {
	id := ctx.Formx(`id`).Uint()
	m := model.NewAlertRecipient(ctx)
	err := m.Get(nil, `id`, id)
	if err != nil {
		return err
	}
	user := handler.User(ctx)
	params := echo.H{
		`email-content`: []byte(ctx.T("您好，我是%s管理员`%s`，这是我发的测试信息，请忽略😊", bootconfig.SoftwareName, user.Username)),
	}
	params[`markdown-content`] = params[`email-content`]
	alertData := &alert.AlertData{
		Title:   ctx.T(`测试信息(%s)`, bootconfig.SoftwareName),
		Content: alert.DefaultTextContent,
		Data:    params,
	}
	data := ctx.Data()
	err = alertData.Send(m.NgingAlertRecipient)
	if err != nil {
		data.SetError(err)
	} else {
		data.SetInfo(ctx.T(`发送成功`))
	}
	return ctx.JSON(data)
}

func AlertRecipientDelete(ctx echo.Context) error {
	id := ctx.Formx(`id`).Uint()
	m := model.NewAlertRecipient(ctx)
	err := m.Delete(nil, db.Cond{`id`: id})
	if err == nil {
		handler.SendOk(ctx, ctx.T(`操作成功`))
	} else {
		handler.SendFail(ctx, err.Error())
	}

	return ctx.Redirect(handler.URLFor(`/manager/alert_recipient`))
}
