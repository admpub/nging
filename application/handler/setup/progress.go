package setup

import "sync"

type ProgressInfo struct {
	Finished  int64
	TotalSize int64
	Summary   string
	Timestamp int64
	mu        *sync.RWMutex
}

func (p *ProgressInfo) Clone() ProgressInfo {
	p.mu.RLock()
	r := *p
	p.mu.RUnlock()
	return r
}

func (p *ProgressInfo) GetTs() int64 {
	p.mu.RLock()
	r := p.Timestamp
	p.mu.RUnlock()
	return r
}

func (p *ProgressInfo) SetTs(ts int64) {
	p.mu.Lock()
	p.Timestamp = ts
	p.mu.Unlock()
}

func (p *ProgressInfo) Done(n int64) {
	p.mu.Lock()
	newVal := p.Finished + n
	if newVal > p.TotalSize {
		p.Finished = p.TotalSize
	} else {
		p.Finished = newVal
	}
	p.mu.Unlock()
}

func (p *ProgressInfo) Add(n int64) {
	p.mu.Lock()
	p.TotalSize += n
	p.mu.Unlock()
}

func (p *ProgressInfo) Reset() {
	p.mu.Lock()
	p.TotalSize = 0
	p.Finished = 0
	p.Summary = ``
	p.Timestamp = 0
	p.mu.Unlock()
}
