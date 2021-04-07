package frp

import (
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/fatedier/golib/crypto"
)

var (
	once      sync.Once
	kcpDoneCh chan struct{}
)

func onceInit() {
	crypto.DefaultSalt = os.Getenv(`FRP_CRYPTO_SALT`)
	if len(crypto.DefaultSalt) == 0 {
		crypto.DefaultSalt = `frp`
	}
	rand.Seed(time.Now().UnixNano())
}
