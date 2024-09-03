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
	"io"
	"sync/atomic"
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

type MessageWithStatus struct {
	Message interface{}
	Status  int
}

type NoticeAndProgress struct {
	last atomic.Value
	send Noticer
	prog *Progress
}

// - Noticer -

func (a *NoticeAndProgress) Send(message interface{}, statusCode int) error {
	a.last.Store(MessageWithStatus{Message: message, Status: statusCode})
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
	if last, ok := a.last.Load().(MessageWithStatus); ok {
		a.send(last.Message, last.Status, a.prog)
	}
	return a
}

func (a *NoticeAndProgress) AutoComplete(on bool) NProgressor {
	a.prog.AutoComplete(on)
	return a
}

func (a *NoticeAndProgress) Complete() NProgressor {
	a.prog.SetComplete()
	if last, ok := a.last.Load().(MessageWithStatus); ok {
		if last.Status == StateSuccess {
			a.prog.SetPercent(100)
		}
		a.send(last.Message, last.Status, a.prog)
	}
	return a
}

func (a *NoticeAndProgress) Reset() {
	a.prog.Reset()
}

func (p *NoticeAndProgress) ProxyReader(r io.Reader) io.ReadCloser {
	return newProxyReader(r, p)
}

func (p *NoticeAndProgress) ProxyWriter(w io.Writer) io.WriteCloser {
	return newProxyWriter(w, p)
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
		noticeSender = MakeNoticer(progress, config.Type, config.Mode, config.ID, config.ClientID, config.User)
	} else {
		noticeSender = DefaultNoticer
	}
	return noticeSender
}

func MakeNoticer(progress *Progress, msgType string, mode string, id interface{}, clientID string, user string) Noticer {
	return func(message interface{}, statusCode int, progs ...*Progress) error {
		msg := acquireMessage()
		msg.Type = msgType
		msg.Title = ``
		msg.Status = statusCode
		msg.Content = message
		msg.Mode = mode
		msg.ID = id
		var prog *Progress
		if len(progs) > 0 {
			prog = progs[0]
		}
		if prog == nil {
			prog = progress
		}
		prog.CalcPercent()
		msg.SetProgress(prog).SetClientID(clientID)
		return Send(user, msg)
	}
}

func DownloadProxyFn(np NProgressor) func(name string, download int, size int64, r io.Reader) io.Reader {
	return func(name string, download int, size int64, r io.Reader) io.Reader {
		np.Add(size)
		np.Send(`downloading `+name, StateSuccess)
		return np.ProxyReader(r)
	}
}
