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
	Cancel      = func() {}
	ErrInitFail = errors.New(`ddns boot failed`)
)

func Run(interval time.Duration) error {
	err := Domains().Update(Config)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithCancel(context.Background())
	Cancel = func() {
		cancel()
	}
	t := time.NewTicker(interval)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-t.C:
			d := Domains()
			if d == nil {
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
