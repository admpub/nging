package frp

import (
	"math/rand"
	"os"
	"time"

	"github.com/admpub/nging/v3/application/library/hook"
	syncOnce "github.com/admpub/once"
	"github.com/fatedier/golib/crypto"
)

var (
	once      syncOnce.Once
	kcpDoneCh chan struct{}
	Hook      = hook.New()
)

func onceInit() {
	crypto.DefaultSalt = os.Getenv(`FRP_CRYPTO_SALT`)
	if len(crypto.DefaultSalt) == 0 {
		crypto.DefaultSalt = `frp`
	}
	rand.Seed(time.Now().UnixNano())
}
