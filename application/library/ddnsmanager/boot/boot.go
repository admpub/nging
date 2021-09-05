package boot

import (
	"context"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/admpub/confl"
	"github.com/admpub/log"
	"github.com/admpub/nging/v3/application/library/config/startup"
	"github.com/admpub/nging/v3/application/library/ddnsmanager/config"
	"github.com/admpub/nging/v3/application/library/ddnsmanager/domain"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
)

var (
	dflt        = config.New()
	domains     *domain.Domains
	once        sync.Once
	mutex       sync.RWMutex
	cancel      context.CancelFunc
	ErrInitFail = errors.New(`ddns boot failed`)
)

func Config() *config.Config {
	mutex.RLock()
	c := *dflt
	mutex.RUnlock()
	return &c
}

func init() {
	startup.OnAfter(`web.installed`, start)
}

func start() {
	saveFile := filepath.Join(echo.Wd(), `config/ddns.yaml`)
	if !com.FileExists(saveFile) {
		return
	}
	_, err := confl.DecodeFile(saveFile, dflt)
	if err != nil {
		log.Error(saveFile+`: `, err)
		return
	}
	err = Run(context.Background())
	if err != nil {
		log.Error(err)
	}
}

func SetConfig(c *config.Config) error {
	saveFile := filepath.Join(echo.Wd(), `config/ddns.yaml`)
	b, err := confl.Marshal(c)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(saveFile, b, os.ModePerm)
	if err != nil {
		return err
	}
	mutex.Lock()
	*dflt = *c
	mutex.Unlock()
	return nil
}

func Run(ctx context.Context, intervals ...time.Duration) error {
	cfg := Config()
	if !cfg.IsValid() {
		return nil
	}
	d := Domains()
	if d == nil {
		return ErrInitFail
	}
	err := d.Update(cfg)
	if err != nil {
		return err
	}
	interval := cfg.Interval
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
	go func() {
		for {
			select {
			case <-c.Done():
				return
			case <-t.C:
				d := Domains()
				if d == nil {
					mutex.Lock()
					if cancel != nil {
						cancel()
						cancel = nil
					}
					mutex.Unlock()
					err = ErrInitFail
					return
				}
				log.Debug(`[DDNS] checking network ip`)
				err := d.Update(Config())
				if err != nil {
					log.Error(err)
				}
			}
		}
	}()
	return err
}

func Domains() *domain.Domains {
	once.Do(initDomains)
	return domains
}

func Reset(ctx context.Context) error {
	mutex.Lock()
	once = sync.Once{}
	if cancel != nil {
		cancel()
		cancel = nil
	}
	mutex.Unlock()
	return Run(ctx)
}

func initDomains() {
	err := commit()
	if err != nil {
		log.Error(err)
	}
}

func commit() error {
	err := dflt.Commit()
	if err != nil {
		return err
	}
	domains, err = domain.ParseDomain(dflt)
	return err
}
