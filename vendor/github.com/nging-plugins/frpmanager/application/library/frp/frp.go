package frp

import (
	"math/rand"
	"os"
	"time"

	syncOnce "github.com/admpub/once"
	"github.com/fatedier/golib/crypto"
)

var (
	once      syncOnce.Once
	kcpDoneCh chan struct{}
)

func onceInit() {
	crypto.DefaultSalt = os.Getenv(`FRP_CRYPTO_SALT`)
	if len(crypto.DefaultSalt) == 0 {
		crypto.DefaultSalt = `frp`
	}
	rand.Seed(time.Now().UnixNano())
}
