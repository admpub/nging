package cookie

import (
	codec "github.com/admpub/securecookie"
	"github.com/admpub/sessions"
	ss "github.com/webx-top/echo/middleware/session/engine"
)

var defaultOptions = &CookieOptions{
	KeyPairs: [][]byte{
		[]byte(codec.GenerateRandomKey(32)),
		[]byte(codec.GenerateRandomKey(32)),
	},
}

func init() {
	RegWithOptions(defaultOptions)
}

func New(opts *CookieOptions) sessions.Store {
	if opts == nil {
		opts = defaultOptions
	}
	store := NewCookieStore(opts.KeyPairs...)
	return store
}

func Reg(store sessions.Store, args ...string) {
	name := `cookie`
	if len(args) > 0 {
		name = args[0]
	}
	ss.Reg(name, store)
}

func RegWithOptions(opts *CookieOptions, args ...string) sessions.Store {
	store := New(opts)
	Reg(store, args...)
	return store
}

func NewCookieOptions(keys ...string) *CookieOptions {
	options := &CookieOptions{
		KeyPairs: KeyPairs(keys...),
	}
	return options
}

type CookieOptions struct {
	KeyPairs [][]byte `json:"-"`
}

// Keys are defined in pairs to allow key rotation, but the common case is to set a single
// authentication key and optionally an encryption key.
//
// The first key in a pair is used for authentication and the second for encryption. The
// encryption key can be set to nil or omitted in the last pair, but the authentication key
// is required in all pairs.
//
// It is recommended to use an authentication key with 32 or 64 bytes. The encryption key,
// if set, must be either 16, 24, or 32 bytes to select AES-128, AES-192, or AES-256 modes.
func NewCookieStore(keyPairs ...[]byte) sessions.Store {
	return &cookieStore{sessions.NewCookieStore(keyPairs...)}
}

type cookieStore struct {
	*sessions.CookieStore
}
