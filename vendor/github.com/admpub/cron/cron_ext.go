package cron

// RemoveCheckFunc 删除job的检查函数，返回true则删除
type RemoveCheckFunc func(e *Entry) bool

func (c *Cron) removeEntryByJob(cb RemoveCheckFunc) {
	var entries []*Entry
	for _, e := range c.entries {
		if !cb(e) {
			entries = append(entries, e)
			continue
		}
		c.logger.Info("removed", "entry", e.ID)
	}
	c.entries = entries
}

func (c *Cron) RemoveJob(cb RemoveCheckFunc) {
	c.runningMu.Lock()
	defer c.runningMu.Unlock()
	if c.running {
		c.removeJob <- cb
	} else {
		c.removeEntryByJob(cb)
	}
}

