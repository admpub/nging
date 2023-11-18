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

package notice

import (
	"context"
)

func NewWithProgress(noticer Noticer, progresses ...*Progress) *NoticeAndProgress {
	var progress *Progress
	if len(progresses) > 0 {
		progress = progresses[0]
	}
	if progress == nil {
		progress = NewProgress()
	}
	return &NoticeAndProgress{
		send: noticer,
		prog: progress,
	}
}

type NoticeAndProgress struct {
	send Noticer
	prog *Progress
}

// - Noticer -

func (a *NoticeAndProgress) Send(message interface{}, statusCode int) error {
	return a.send(message, statusCode, a.prog)
}

func (a *NoticeAndProgress) Success(message interface{}) error {
	return a.Send(message, StateSuccess)
}

func (a *NoticeAndProgress) Failure(message interface{}) error {
	return a.Send(message, StateFailure)
}

// - Progress -

func (a *NoticeAndProgress) Add(n int64) NProgressor {
	a.prog.Add(n)
	return a
}

func (a *NoticeAndProgress) Done(n int64) NProgressor {
	a.prog.Done(n)
	return a
}

func (a *NoticeAndProgress) AutoComplete(on bool) NProgressor {
	a.prog.AutoComplete(on)
	return a
}

func (a *NoticeAndProgress) Complete() NProgressor {
	a.prog.SetComplete()
	return a
}

func (a *NoticeAndProgress) Callback(total int64, exec func(callback func(strLen int)) error) error {
	return a.prog.Callback(total, exec)
}

func NewNoticer(ctx context.Context, config *HTTPNoticerConfig) Noticer {
	var noticeSender Noticer
	if config.IsExited == nil && config.Timeout != 0 {
		config.IsExited = NewControlWithContext(ctx, config.Timeout)
	}
	if len(config.Mode) == 0 {
		if config.ID != nil {
			config.Mode = ModeElement
		} else {
			config.Mode = ModeNotify
		}
	}
	progress := NewProgress().SetControl(config.IsExited)
	if len(config.User) > 0 {
		OpenMessage(config.User, config.Type)
		//defer CloseMessage(config.User, config.Type)
		noticeSender = func(message interface{}, statusCode int, progs ...*Progress) error {
			msg := NewMessageWithValue(
				config.Type,
				``,
				message,
				statusCode,
			).SetMode(config.Mode).SetID(config.ID)
			var prog *Progress
			if len(progs) > 0 {
				prog = progs[0]
			}
			if prog == nil {
				prog = progress
			}
			msg.SetProgress(prog).CalcPercent().SetClientID(config.ClientID)
			return Send(config.User, msg)
		}
	} else {
		noticeSender = DefaultNoticer
	}
	return noticeSender
}
