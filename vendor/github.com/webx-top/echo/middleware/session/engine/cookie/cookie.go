package cookie

import (
	"github.com/admpub/sessions"
	codec "github.com/gorilla/securecookie"
	"github.com/webx-top/echo"
	ss "github.com/webx-top/echo/middleware/session/engine"
)

func init() {
	RegWithOptions(&CookieOptions{
		KeyPairs: [][]byte{
			[]byte(codec.GenerateRandomKey(32)),
			[]byte(codec.GenerateRandomKey(32)),
		},
		SessionOptions: &echo.SessionOptions{
			Name:   `GOSESSIONID`,
			Engine: `cookie`,
			CookieOptions: &echo.CookieOptions{
				Path:     `/`,
				HttpOnly: true,
			},
		},
	})
}

func New(opts *CookieOptions) CookieStore {
	store := NewCookieStore(opts.KeyPairs...)
	store.Options(*opts.SessionOptions)
	return store
}

func Reg(store CookieStore, args ...string) {
	name := `cookie`
	if len(args) > 0 {
		name = args[0]
	}
	ss.Reg(name, store)
}

func RegWithOptions(opts *CookieOptions, args ...string) {
	Reg(New(opts), args...)
}

type CookieStore interface {
	ss.Store
}

type CookieOptions struct {
	KeyPairs             [][]byte `json:"keyPairs"`
	*echo.SessionOptions `json:"session"`
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
func NewCookieStore(keyPairs ...[]byte) CookieStore {
	return &cookieStore{sessions.NewCookieStore(keyPairs...)}
}

type cookieStore struct {
	*sessions.CookieStore
}

func (c *cookieStore) Options(options echo.SessionOptions) {
	c.CookieStore.Options = &sessions.Options{
		Path:     options.Path,
		Domain:   options.Domain,
		MaxAge:   options.MaxAge,
		Secure:   options.Secure,
		HttpOnly: options.HttpOnly,
	}
	c.CookieStore.MaxAge(c.CookieStore.Options.MaxAge)
}
