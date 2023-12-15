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
	"io"
	"sync"
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
	return &Progress{
		total:   -1,
		finish:  -1,
		percent: 0,
		mu:      &sync.RWMutex{},
	}
}

type Progress struct {
	total        int64
	finish       int64
	percent      float64
	complete     bool
	mu           *sync.RWMutex
	control      IsExited
	autoComplete bool
}

func (p *Progress) IsExited() bool {
	if p.control == nil {
		return false
	}
	return p.control.IsExited()
}

func (p *Progress) CloneInfo() ProgressInfo {
	p.mu.RLock()
	cloned := ProgressInfo{
		Total:    p.total,
		Finish:   p.finish,
		Percent:  p.percent,
		Complete: p.complete,
	}
	p.mu.RUnlock()
	return cloned
}

func (p *Progress) CopyToInfo(to *ProgressInfo) {
	p.mu.RLock()
	to.Total = p.total
	to.Finish = p.finish
	to.Percent = p.percent
	to.Complete = p.complete
	p.mu.RUnlock()
}

func (p *Progress) SetControl(control IsExited) *Progress {
	p.control = control
	return p
}

func (p *Progress) Percent() float64 {
	p.mu.RLock()
	percent := p.percent
	p.mu.RUnlock()
	return percent
}

func (p *Progress) Total() int64 {
	p.mu.RLock()
	total := p.total
	p.mu.RUnlock()
	return total
}

func (p *Progress) Finish() int64 {
	p.mu.RLock()
	finish := p.finish
	p.mu.RUnlock()
	return finish
}

func (p *Progress) Complete() bool {
	p.mu.RLock()
	complete := p.complete
	p.mu.RUnlock()
	return complete
}

func (p *Progress) Reset() {
	p.mu.Lock()
	p.total = -1
	p.finish = -1
	p.percent = 0
	p.complete = false
	p.mu.Unlock()
}

func (p *Progress) CalcPercent() *Progress {
	p.mu.Lock()
	if p.total > 0 {
		p.percent = (float64(p.finish) / float64(p.total)) * 100
		if p.percent < 0 {
			p.percent = 0
		}
	} else if p.total == 0 {
		p.percent = 100
	} else {
		p.percent = 0
	}
	p.mu.Unlock()
	return p
}

func (p *Progress) SetPercent(percent float64) *Progress {
	p.mu.Lock()
	p.percent = percent
	p.mu.Unlock()
	return p
}

func (p *Progress) Add(n int64) *Progress {
	p.mu.Lock()
	if p.finish > 0 {
		p.finish = 0
	}
	if p.total == -1 {
		p.total++
	}
	p.total += n
	p.mu.Unlock()
	return p
}

func (p *Progress) Done(n int64) int64 {
	p.mu.Lock()
	if p.finish == -1 {
		p.finish++
	}
	p.finish += n
	newN := p.finish
	if p.autoComplete && newN >= p.total {
		p.complete = true
	}
	p.mu.Unlock()
	return newN
}

func (p *Progress) AutoComplete(on bool) *Progress {
	p.autoComplete = on
	return p
}

func (p *Progress) SetComplete() *Progress {
	p.mu.Lock()
	p.complete = true
	p.mu.Unlock()
	return p
}

func (p *Progress) ProxyReader(r io.Reader) io.ReadCloser {
	return newProxyReader(r, p)
}

func (p *Progress) ProxyWriter(w io.Writer) io.WriteCloser {
	return newProxyWriter(w, p)
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
