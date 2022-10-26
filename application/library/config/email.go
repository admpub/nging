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

package config

import (
	"github.com/admpub/mail"
	cronSend "github.com/admpub/nging/v5/application/library/cron/send"
	"github.com/admpub/nging/v5/application/library/email"
	"github.com/webx-top/echo"
)

type Email struct {
	*mail.SMTPConfig
	Timeout   int64  //超时时间(秒)，采用默认引擎发信时，此项无效
	Engine    string //值为email时采用github.com/jordan-wright/email包发送，否则采用默认的github.com/admpub/mail发送
	From      string //发信人Email地址
	QueueSize int    //允许同一时间发信的数量
}

func (c *Email) SetBy(r echo.H, defaults echo.H) *Email {
	if !r.Has(`smtp`) && defaults != nil {
		r.Set(`smtp`, defaults.GetStore(`smtp`))
	}
	smtp := r.GetStore(`smtp`)
	if c.SMTPConfig == nil {
		c.SMTPConfig = &mail.SMTPConfig{}
	}
	c.Username = smtp.String(`username`)
	c.Password = smtp.String(`password`)
	c.Host = smtp.String(`host`)
	c.Port = smtp.Int(`port`)
	c.Secure = smtp.String(`secure`)
	c.Identity = smtp.String(`identity`)
	c.Timeout = smtp.Int64(`timeout`)
	c.Engine = smtp.String(`engine`)
	c.From = smtp.String(`from`)
	c.QueueSize = smtp.Int(`queueSize`)
	return c
}

func (c *Email) Init() {
	if c.SMTPConfig == nil {
		c.SMTPConfig = &mail.SMTPConfig{}
	}
	cronSend.DefaultSMTPConfig = c.SMTPConfig
	cronSend.DefaultEmailConfig.Sender = c.From
	cronSend.DefaultEmailConfig.Engine = c.Engine
	if cronSend.DefaultEmailConfig.Timeout > 0 {
		cronSend.DefaultEmailConfig.Timeout = c.Timeout
	}
	if c.QueueSize > 0 {
		email.QueueSize = c.QueueSize
	}
}
