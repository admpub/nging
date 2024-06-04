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
package model

import (
	"errors"
	"time"

	"github.com/admpub/nging/v5/application/dbschema"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
	"github.com/webx-top/echo/middleware/tplfunc"
)

var (
	// SMSWaitingSeconds 短信发送后等待秒数
	SMSWaitingSeconds int64 = 60
	// SMSMaxPerDay 短信每人每天发送上限
	SMSMaxPerDay int64 = 10
)

func NewCode(ctx echo.Context) *Code {
	return &Code{
		Verification: dbschema.NewNgingCodeVerification(ctx),
		Invitation:   dbschema.NewNgingCodeInvitation(ctx),
	}
}

type Code struct {
	Verification *dbschema.NgingCodeVerification
	Invitation   *dbschema.NgingCodeInvitation
}

func (c *Code) AddVerificationCode() (interface{}, error) {
	if len(c.Verification.Disabled) == 0 {
		c.Verification.Disabled = `N`
	}
	return c.Verification.Insert()
}

func (c *Code) UseVerificationCode(m *dbschema.NgingCodeVerification) (err error) {
	m.Used = uint(time.Now().Unix())
	err = m.UpdateField(nil, `used`, m.Used, `id`, m.Id)
	return
}

func (c *Code) LastVerificationCode(ownerID uint64, ownerType string, sendMethod string) (err error) {
	err = c.Verification.Get(func(r db.Result) db.Result {
		return r.OrderBy(`-id`)
	}, db.And(
		db.Cond{`disabled`: `N`},
		db.Cond{`owner_type`: ownerType},
		db.Cond{`owner_id`: ownerID},
		db.Cond{`send_method`: sendMethod},
	))
	return
}

func (c *Code) CountTodayVerificationCode(ownerID uint64, ownerType string, sendMethod string) (int64, error) {
	now := time.Now()
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	end := start.AddDate(0, 0, 1)
	return c.Verification.Count(nil, db.And(
		db.Cond{`disabled`: `N`},
		db.Cond{`owner_type`: ownerType},
		db.Cond{`owner_id`: ownerID},
		db.Cond{`send_method`: sendMethod},
		db.Cond{`created`: db.Between(start.Unix(), end.Unix())},
	))
}

func (c *Code) CheckFrequency(ownerID uint64, ownerType string, sendMethod string, frequencyCfg echo.H) error {
	ctx := c.Verification.Context()
	if err := c.LastVerificationCode(ownerID, ownerType, sendMethod); err != nil {
		if err != db.ErrNoMoreRows {
			return err
		}
	} else {
		interval := frequencyCfg.Int64(`interval`, SMSWaitingSeconds)
		waitingSeconds := time.Now().Unix() - int64(c.Verification.Created)
		if waitingSeconds < interval {
			return ctx.NewError(code.FrequencyTooFast, `请等待%d秒之后再发送`, interval-waitingSeconds)
		}
		maxPerDay := frequencyCfg.Int64(`maxPerDay`, SMSMaxPerDay)
		if count, err := c.CountTodayVerificationCode(ownerID, ownerType, sendMethod); err != nil {
			return err
		} else if count >= maxPerDay {
			return ctx.NewError(code.FrequencyTooFast, `您今天的发送次数已达上限: %d`, maxPerDay)
		}
	}
	return nil
}

func (c *Code) CheckVerificationCode(vcode string, purpose string, ownerID uint64, ownerType string, sendMethod string, sendTo string) (err error) {
	ctx := c.Verification.Context()
	err = c.Verification.Get(nil, db.And(
		db.Cond{`disabled`: `N`},
		db.Cond{`owner_type`: ownerType},
		db.Cond{`owner_id`: ownerID},
		db.Cond{`send_method`: sendMethod},
		db.Cond{`send_to`: sendTo},
		db.Cond{`code`: vcode},
		db.Cond{`purpose`: purpose},
	))
	var objectName string
	switch sendMethod {
	case `email`:
		objectName = ctx.T(`邮件`)
	case `mobile`:
		objectName = ctx.T(`短信`)
	}
	if err != nil {
		if err == db.ErrNoMoreRows {
			return ctx.NewError(code.InvalidParameter, `%s验证码无效`, objectName)
		}
		return err
	}
	now := uint(time.Now().Unix())
	if !(c.Verification.Start <= now && c.Verification.End >= now) {
		return ctx.NewError(code.DataHasExpired, `%s验证码已经过期`, objectName)
	}
	if c.Verification.Used > 0 {
		return ctx.NewError(code.DataUnavailable, `%s验证码已经使用过了`, objectName)
	}
	return
}

func (c *Code) VerfyInvitationCode(code string) (err error) {
	ctx := c.Invitation.Context()
	err = c.Invitation.Get(nil, `code`, code)
	if err != nil {
		if err == db.ErrNoMoreRows {
			err = ctx.E(`邀请码无效`)
		}
		return
	}
	if c.Invitation.Used > 0 {
		err = ctx.E(`该邀请码已被使用过了`)
		return
	}
	if c.Invitation.Disabled == `Y` {
		err = ctx.E(`该邀请码已被禁用`)
		return
	}
	now := uint(time.Now().Unix())
	if c.Invitation.Start > now {
		if c.Invitation.End > 0 {
			err = errors.New(ctx.T(`该邀请码只能在“%s - %s”这段时间内使用`,
				tplfunc.TsToDate(`2006/01/02 15:04:05`, c.Invitation.Start),
				tplfunc.TsToDate(`2006/01/02 15:04:05`, c.Invitation.End),
			))
		} else {
			err = errors.New(ctx.T(`该邀请码只能在“%s”之后使用`,
				tplfunc.TsToDate(`2006/01/02 15:04:05`, c.Invitation.Start),
			))
		}
		return
	}
	if c.Invitation.End > 0 && c.Invitation.End < now {
		err = ctx.E(`该邀请码已过期`)
		return
	}
	return
}

func (c *Code) UseInvitationCode(m *dbschema.NgingCodeInvitation, usedUid uint) (err error) {
	m.Used = uint(time.Now().Unix())
	m.RecvUid = usedUid
	err = m.Update(nil, `id`, m.Id)
	return
}
