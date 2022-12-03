package common

import (
	"time"

	"github.com/admpub/log"
)

func NoRetry(err error) *errNoRetry {
	return &errNoRetry{error: err}
}

var _ error = NoRetry(nil)

type errNoRetry struct {
	error
}

type iNoRetry interface {
	NoRetry()
}

func (e *errNoRetry) NoRetry() {
}

func (e *errNoRetry) Unwrap() error {
	return e.error
}

func IsNoRetry(err error) bool {
	_, ok := err.(iNoRetry)
	return ok
}

func Retry(maxRetries int, fn func() error, stepDuration ...time.Duration) error {
	err := fn()
	if err == nil || IsNoRetry(err) {
		return err
	}
	log.Error(err.Error())
	var step time.Duration
	if len(stepDuration) > 0 {
		step = stepDuration[0]
	}
	if step == 0 {
		step = time.Second * 5
	}
	for i := 1; i <= maxRetries; i++ {
		wait := time.Duration(i) * step
		log.Infof(`retry(%d/%d) after %v seconds, waiting...`, i, maxRetries, wait.Seconds())
		time.Sleep(wait)
		err = fn()
		if err == nil || IsNoRetry(err) {
			break
		}
		log.Errorf(`retry(%d/%d) error: %v`, i, maxRetries, err)
	}
	return err
}
