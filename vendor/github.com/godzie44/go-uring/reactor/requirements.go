//go:build linux

package reactor

import (
	"errors"
	"fmt"
	"github.com/godzie44/go-uring/uring"
)

var ErrRequirements = errors.New("ring does not meet the requirements")

func checkRingReq(r *uring.Ring, net bool) error {
	if !r.Params.ExtArgFeature() {
		return fmt.Errorf("%w: ext arg feature must exists", ErrRequirements)
	}

	if net && !r.Params.FastPollFeature() {
		return fmt.Errorf("%w: fast poll feature must exists for NetReactor", ErrRequirements)
	}

	return nil
}
