package gerberos

import "time"

type occurrences struct {
	registry map[string][]time.Time
	interval time.Duration
	count    int
}

func (r *occurrences) add(host string) bool {
	if _, f := r.registry[host]; !f {
		r.registry[host] = []time.Time{time.Now()}
		return false
	}

	r.registry[host] = append(r.registry[host], time.Now())
	if len(r.registry[host]) > r.count {
		r.registry[host] = r.registry[host][1:]
	}

	if len(r.registry[host]) == r.count {
		d := r.registry[host][r.count-1].Sub(r.registry[host][0])
		if d <= r.interval {
			delete(r.registry, host)
			return true
		}
	}

	return false
}

func newOccurrences(interval time.Duration, count int) *occurrences {
	return &occurrences{
		registry: make(map[string][]time.Time),
		interval: interval,
		count:    count,
	}
}
