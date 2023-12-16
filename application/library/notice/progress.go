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
	"math"
	"sync/atomic"
)

type ProgressInfo struct {
	Total    int64   `json:"total" xml:"total"`
	Finish   int64   `json:"finish" xml:"finish"`
	Percent  float64 `json:"percent" xml:"percent"`
	Complete bool    `json:"complete" xml:"complete"`
}

func (p *ProgressInfo) reset() {
	p.Total = 0
	p.Finish = 0
	p.Percent = 0
	p.Complete = false
}

func NewProgress() *Progress {
	p := &Progress{}
	p.total.Store(-1)
	p.finish.Store(-1)
	return p
}

type Progress struct {
	total        atomic.Int64
	finish       atomic.Int64
	percent      atomic.Int64
	complete     atomic.Bool
	control      IsExited
	autoComplete atomic.Bool
}

func (p *Progress) IsExited() bool {
	if p.control == nil {
		return false
	}
	return p.control.IsExited()
}

func (p *Progress) CloneInfo() ProgressInfo {
	cloned := ProgressInfo{
		Total:    p.total.Load(),
		Finish:   p.finish.Load(),
		Percent:  p.Percent(),
		Complete: p.complete.Load(),
	}
	return cloned
}

func (p *Progress) CopyToInfo(to *ProgressInfo) {
	to.Total = p.total.Load()
	to.Finish = p.finish.Load()
	to.Percent = p.Percent()
	to.Complete = p.complete.Load()
}

func (p *Progress) SetControl(control IsExited) *Progress {
	p.control = control
	return p
}

func (p *Progress) Percent() float64 {
	percent := float64(p.percent.Load()) / 10000
	return percent
}

func (p *Progress) Total() int64 {
	total := p.total.Load()
	return total
}

func (p *Progress) Finish() int64 {
	finish := p.finish.Load()
	return finish
}

func (p *Progress) Complete() bool {
	complete := p.complete.Load()
	return complete
}

func (p *Progress) Reset() {
	p.total.Store(-1)
	p.finish.Store(-1)
	p.percent.Store(0)
	p.complete.Store(false)
}

func (p *Progress) CalcPercent() *Progress {
	total := p.Total()
	var percent float64
	if total > 0 {
		percent = (float64(p.Finish()) / float64(total)) * 100
		if percent < 0 {
			percent = 0
		}
	} else if total == 0 {
		percent = 100
	} else {
		percent = 0
	}
	p.SetPercent(percent)
	return p
}

func (p *Progress) SetPercent(percent float64) *Progress {
	p.percent.Store(int64(math.Floor(percent * 10000)))
	return p
}

func (p *Progress) Add(n int64) *Progress {
	if p.Finish() > 0 {
		p.finish.Store(0)
	}
	p.total.CompareAndSwap(-1, 0)
	p.total.Add(n)
	return p
}

func (p *Progress) Done(n int64) int64 {
	p.finish.CompareAndSwap(-1, 0)
	newN := p.finish.Add(n)
	if p.autoComplete.Load() && newN >= p.total.Load() {
		p.SetComplete()
	}
	return newN
}

func (p *Progress) AutoComplete(on bool) *Progress {
	p.autoComplete.Store(on)
	return p
}

func (p *Progress) SetComplete() *Progress {
	p.complete.Store(true)
	return p
}

func (p *Progress) Callback(total int64, exec func(callback func(strLen int)) error) error {
	var remains int64 = 100
	var partPercent float64
	var perByteVal float64
	if total > 0 {
		perByteVal = float64(remains) / float64(total)
	}
	var callback = func(strLen int) {
		if perByteVal <= 0 {
			return
		}
		partPercent += perByteVal * float64(strLen)
		if partPercent < 1 {
			return
		}
		percent := int64(partPercent)
		remains -= percent
		p.Done(percent)
		partPercent = partPercent - float64(percent)
	}
	err := exec(callback)
	if remains > 0 {
		p.Done(remains)
	}
	return err
}
