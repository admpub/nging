package boot

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/admpub/confl"
	"github.com/admpub/log"
	"github.com/admpub/nging/v4/application/library/config/startup"
	syncOnce "github.com/admpub/once"
	"github.com/nging-plugins/ddnsmanager/application/library/ddnsmanager/config"
	"github.com/nging-plugins/ddnsmanager/application/library/ddnsmanager/domain"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
)

var (
	dflt            = config.New()
	domains         *domain.Domains
	once            syncOnce.Once
	mutex           sync.RWMutex
	cancel          context.CancelFunc
	defaultInterval = 5 * time.Minute
	waitingDuration = 500 * time.Millisecond
	ErrInitFail     = errors.New(`ddns boot failed`)
)

func IsRunning() bool {
	return cancel != nil
}

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
	if dflt.Closed {
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
	err = os.WriteFile(saveFile, b, os.ModePerm)
	if err != nil {
		return err
	}
	mutex.Lock()
	*dflt = *c
	mutex.Unlock()
	return nil
}

func Run(ctx context.Context, intervals ...time.Duration) (err error) {
	cfg := Config()
	if !cfg.IsValid() {
		log.Warn(`[DDNS] Exit task: The task does not meet the startup conditions`)
		return nil
	}
	d := Domains()
	if d == nil {
		return ErrInitFail
	}
	err = d.Update(ctx, cfg)
	if err != nil {
		log.Error(`[DDNS] Exit task`)
		return err
	}
	mutex.Lock()
	if cancel != nil {
		cancel()
		cancel = nil
		time.Sleep(waitingDuration)
	}
	var c context.Context
	c, cancel = context.WithCancel(ctx)
	mutex.Unlock()
	go func() {
		interval := cfg.Interval
		if len(intervals) > 0 {
			interval = intervals[0]
		}
		if interval < time.Second {
			interval = defaultInterval
		}
		t := time.NewTicker(interval)
		defer t.Stop()
		log.Okay(`[DDNS] Starting task. Interval: `, interval.String())
		for {
			select {
			case <-c.Done():
				log.Warn(`[DDNS] Forced exit task`)
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
					log.Error(`[DDNS] Exit task. Error: `, err.Error())
					return
				}
				log.Debug(`[DDNS] Checking network ip`)
				err := d.Update(ctx, Config())
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
	cfg := Config() // 含锁，小心使用
	mutex.Lock()
	once.Reset()
	if cancel != nil {
		cancel()
		cancel = nil
		time.Sleep(waitingDuration)
		if cfg.Closed {
			log.Warn(`[DDNS] Stopping task`)
		}
	}
	mutex.Unlock()
	if cfg.Closed {
		return nil
	}
	log.Warn(`[DDNS] Starting reboot task`)
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
