package gerberos

import (
	"net"
	"regexp"
	"sync"
	"time"
)

type aggregate struct {
	registry      map[string]net.IP
	registryMutex sync.Mutex
	interval      time.Duration
	regexp        []*regexp.Regexp
}

func newAggregate(interval time.Duration, res []*regexp.Regexp) *aggregate {
	return &aggregate{
		registry: make(map[string]net.IP),
		interval: interval,
		regexp:   res,
	}
}
