package boot

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/admpub/log"
	"github.com/admpub/nging/v3/application/library/ddnsmanager/config"
	"github.com/admpub/nging/v3/application/library/ddnsmanager/domain"
)

var (
	Config      = config.New()
	domains     *domain.Domains
	once        sync.Once
	mutex       sync.RWMutex
	cancel      context.CancelFunc
	ErrInitFail = errors.New(`ddns boot failed`)
)

func Run(ctx context.Context, intervals ...time.Duration) error {
	if Config.Closed {
		return nil
	}
	d := Domains()
	if d == nil {
		return ErrInitFail
	}
	err := d.Update(Config)
	if err != nil {
		return err
	}
	interval := Config.Interval
	if len(intervals) > 0 {
		interval = intervals[0]
	}
	if interval < time.Second {
		interval = 5 * time.Minute
	}
	mutex.Lock()
	if cancel != nil {
		cancel()
		cancel = nil
	}
	var c context.Context
	c, cancel = context.WithCancel(ctx)
	mutex.Unlock()
	t := time.NewTicker(interval)
	defer t.Stop()
	for {
		select {
		case <-c.Done():
			return nil
		case <-t.C:
			d := Domains()
			if d == nil {
				mutex.Lock()
				if cancel != nil {
					cancel()
					cancel = nil
				}
				mutex.Unlock()
				return ErrInitFail
			}
			log.Debug(`[DDNS] checking network ip`)
			err := d.Update(Config)
			if err != nil {
				log.Error(err)
			}
		}
	}
}

func Domains() *domain.Domains {
	once.Do(initDomains)
	return domains
}

func Reset(ctx context.Context) {
	mutex.Lock()
	once = sync.Once{}
	if cancel != nil {
		cancel()
		cancel = nil
	}
	mutex.Unlock()
	Run(ctx)
}

func initDomains() {
	err := Commit()
	if err != nil {
		log.Error(err)
	}
}

func Commit() error {
	err := Config.Commit()
	if err != nil {
		return err
	}
	domains, err = domain.ParseDomain(Config)
	return err
}
