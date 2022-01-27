package ssystem

import (
	"regexp"
	"strconv"
	"time"

	"github.com/admpub/log"
	"github.com/webx-top/echo/middleware/bytes"
)

var (
	reNumeric                  = regexp.MustCompile(`^[0-9]+$`)
	defaultMaxRequestBodyBytes = 2 << 20 // 2M
)

func ParseTimeDuration(timeout string) time.Duration {
	var timeoutDuration time.Duration
	if len(timeout) > 0 {
		if reNumeric.MatchString(timeout) {
			if val, err := strconv.ParseUint(timeout, 10, 64); err != nil {
				log.Error(err)
			} else {
				timeoutDuration = time.Second * time.Duration(val)
			}
		} else {
			timeoutDuration, _ = time.ParseDuration(timeout)
		}
	}
	return timeoutDuration
}

func ParseBytes(size string) (int, error) {
	bz, err := bytes.Parse(size)
	return int(bz), err
}
