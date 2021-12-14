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
	"github.com/admpub/nging/v4/application/handler"
	"github.com/admpub/nging/v4/application/model"
	"github.com/admpub/null"
	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/sqlbuilder"
	"github.com/webx-top/echo"
)

func Verification(ctx echo.Context) error {
	cond := db.Compounds{}
	sendStatus := ctx.Form(`sendStatus`)
	sendTo := ctx.Form(`sendTo`)
	usedStatus := ctx.Form(`usedStatus`)
	q := ctx.Formx(`q`).String()
	if len(q) > 0 {
		cond.AddKV(`code`, q)
	}
	if len(sendStatus) > 0 {
		cond.AddKV(`b.status`, sendStatus)
	}
	if len(sendTo) > 0 {
		cond.AddKV(`a.send_to`, sendTo)
	}
	if len(usedStatus) > 0 {
		if usedStatus == `1` {
			cond.AddKV(`a.used`, db.Gt(0))
		} else {
			cond.AddKV(`a.used`, 0)
		}
	}
	m := model.NewVerification(ctx)
	logM := model.NewSendingLog(ctx)
	recv := []null.StringMap{}
	pagination, err := handler.PagingWithSelectList(ctx, m.NewParam().SetRecv(&recv).SetArgs(cond.And()).SetAlias(`a`).SetMWSel(func(r sqlbuilder.Selector) sqlbuilder.Selector {
		return r.OrderBy(`-a.id`)
	}).SetCols(`a.*`, `a.created createdAt`, `b.sent_at`, `b.method`, `b.to`, `b.provider`, `b.result`, `b.status`, `b.retries`, `b.content`, `b.params`, `b.appointment_time`).AddJoin(`LEFT`, logM.Name_(), `b`, `b.source_type='code_verification' AND b.source_id=a.id`))
	if err != nil {
		return err
	}
	ctx.Set(`pagination`, pagination)
	ret := handler.Err(ctx, err)
	ctx.Set(`listData`, recv)
	return ctx.Render(`/manager/verification`, ret)
}

func VerificationDelete(ctx echo.Context) error {
	id := ctx.Formx(`id`).Uint64()
	m := model.NewVerification(ctx)
	err := m.Delete(nil, db.Cond{`id`: id})
	if err == nil {
		logM := model.NewSendingLog(ctx)
		logM.Delete(nil, db.And(db.Cond{`source_id`: id}, db.Cond{`source_type`: `code_verification`}))
		handler.SendOk(ctx, ctx.T(`操作成功`))
	} else {
		handler.SendFail(ctx, err.Error())
	}

	return ctx.Redirect(handler.URLFor(`/manager/verification`))
}
