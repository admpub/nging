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

package notice

import (
	"context"
	"time"
)

func NewProgress() *Progress {
	return &Progress{
		Total:   -1,
		Finish:  -1,
		Percent: 0,
	}
}

type Progress struct {
	Total    int64   `json:"total" xml:"total"`
	Finish   int64   `json:"finish" xml:"finish"`
	Percent  float64 `json:"percent" xml:"percent"`
	Complete bool    `json:"complete" xml:"complete"`
	control  IsExited
}

type Control struct {
	exited bool
}

func (c *Control) IsExited() bool {
	return c.exited
}

func (c *Control) Exited() *Control {
	c.exited = true
	return c
}

func (c *Control) ListenContextAndTimeout(ctx context.Context, timeouts ...time.Duration) *Control {
	timeout := 24 * time.Hour
	if len(timeouts) > 0 && timeouts[0] != 0 {
		timeout = timeouts[0]
	}
	t := time.NewTicker(timeout)
	defer t.Stop()
	go func() {
		for {
			select {
			case <-ctx.Done():
				c.Exited()
				return
			case <-t.C:
				c.Exited()
				return
			}
		}
	}()
	return c
}

type IsExited interface {
	IsExited() bool
}

func (p *Progress) IsExited() bool {
	if p.control == nil {
		return false
	}
	return p.control.IsExited()
}

func (p *Progress) SetControl(control IsExited) *Progress {
	p.control = control
	return p
}

func (p *Progress) CalcPercent() *Progress {
	if p.Total > 0 {
		p.Percent = (float64(p.Finish) / float64(p.Total)) * 100
		if p.Percent < 0 {
			p.Percent = 0
		}
	} else if p.Total == 0 {
		p.Percent = 100
	} else {
		p.Percent = 0
	}
	return p
}
