package session

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	codec "github.com/admpub/securecookie"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/middleware/session/engine/cookie"
	"github.com/webx-top/echo/middleware/session/engine/file"
)

var (
	DefaultSessionEngine = `file`
	DefaultSessionName   = `SID`
	DefaultTempDir       = os.TempDir()
	DefaultCookieOptions = &echo.CookieOptions{
		Domain:   ``,
		Path:     `/`,
		MaxAge:   0,
		HttpOnly: true,
	}
	once                  = sync.Once{}
	defaultSessionOptions *echo.SessionOptions
)

func DefaultSessionOptions() *echo.SessionOptions {
	once.Do(initDefaultSessionOptions)
	return defaultSessionOptions
}

func ResetDefaultSessionOptions() {
	once = sync.Once{}
}

func initDefaultSessionOptions() {
	defaultSessionOptions = echo.NewSessionOptions(DefaultSessionEngine, DefaultSessionName, DefaultCookieOptions)
	if err := registerSessionStore(); err != nil {
		panic(err)
	}
}

type cookieSecretKey struct {
	HashKey  string `json:"hashKey"`
	BlockKey string `json:"blockKey"`
}

func (c *cookieSecretKey) Generate() *cookieSecretKey {
	c.HashKey = string(codec.GenerateRandomKey(32))
	c.BlockKey = string(codec.GenerateRandomKey(32))
	return c
}

func registerSessionStore() error {
	tempDir := `./`
	if len(os.Args) > 0 {
		tempDir = filepath.Dir(os.Args[0])
	}
	r := &cookieSecretKey{}
	cookieSecretFile := filepath.Join(tempDir, `cookiesecret`)
	b, err := ioutil.ReadFile(cookieSecretFile)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	} else {
		err = json.Unmarshal(b, r)
	}
	if err != nil {
		r.Generate()
		b, err = json.Marshal(r)
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(cookieSecretFile, b, os.ModePerm)
	}

	//==================================
	// 注册session存储引擎
	//==================================

	//1. 注册默认引擎：cookie
	cookieStoreOptions := cookie.NewCookieOptions(r.HashKey, r.BlockKey)
	cookie.RegWithOptions(cookieStoreOptions)

	//2. 注册文件引擎：file
	file.RegWithOptions(&file.FileOptions{
		SavePath: filepath.Join(DefaultTempDir, `sessions`),
		KeyPairs: cookieStoreOptions.KeyPairs,
	})
	return err
}
