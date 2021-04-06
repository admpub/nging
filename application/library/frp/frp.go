package frp

import (
	"math/rand"
	"sync"
	"time"

	"github.com/fatedier/golib/crypto"
)

var once sync.Once

func onceInit() {
	crypto.DefaultSalt = `frp`
	rand.Seed(time.Now().UnixNano())
}
