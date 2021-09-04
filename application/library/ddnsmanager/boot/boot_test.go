package boot

import (
	"context"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	select {
	case <-ctx.Done():
		Cancel()
	default:
		err := Run(time.Second * 5)
		if err != nil {
			panic(err)
		}
	}
}
