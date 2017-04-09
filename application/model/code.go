/*

   Copyright 2016 Wenhui Shen <www.webx.top>

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.

*/
package model

import (
	"errors"
	"time"

	"github.com/admpub/nging/application/dbschema"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/middleware/tplfunc"
)

func NewCode(ctx echo.Context) *Code {
	return &Code{
		Verification: &dbschema.CodeVerification{},
		Invitation:   &dbschema.CodeInvitation{},
		Base:         &Base{Context: ctx},
	}
}

type Code struct {
	Verification *dbschema.CodeVerification
	Invitation   *dbschema.CodeInvitation
	*Base
}

func (c *Code) VerfyInvitationCode(code string) (err error) {
	err = c.Invitation.Get(nil, `code`, code)
	if err != nil {
		if err == db.ErrNoMoreRows {
			err = errors.New(c.T(`邀请码无效`))
		}
		return
	}
	if c.Invitation.Used > 0 {
		err = errors.New(c.T(`该邀请码已被使用过了`))
		return
	}
	if c.Invitation.Disabled == `Y` {
		err = errors.New(c.T(`该邀请码已被禁用`))
		return
	}
	now := uint(time.Now().Unix())
	if c.Invitation.Start > now {
		if c.Invitation.End > 0 {
			err = errors.New(c.T(`该邀请码只能在“%s - %s”这段时间内使用`,
				tplfunc.TsToDate(`2006/01/02 15:04:05`, c.Invitation.Start),
				tplfunc.TsToDate(`2006/01/02 15:04:05`, c.Invitation.End),
			))
		} else {
			err = errors.New(c.T(`该邀请码只能在“%s”之后使用`,
				tplfunc.TsToDate(`2006/01/02 15:04:05`, c.Invitation.Start),
			))
		}
		return
	}
	if c.Invitation.End > 0 && c.Invitation.End < now {
		err = errors.New(c.T(`该邀请码已过期`))
		return
	}
	return
}

func (c *Code) UseInvitationCode(m *dbschema.CodeInvitation, usedUid uint) (err error) {
	m.Used = uint(time.Now().Unix())
	m.RecvUid = usedUid
	err = m.Edit(nil, `id`, m.Id)
	return
}
