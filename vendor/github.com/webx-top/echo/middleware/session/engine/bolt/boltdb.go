package bolt

import (
	"runtime"
	"time"

	"github.com/admpub/boltstore/reaper"
	"github.com/admpub/boltstore/store"
	"github.com/admpub/sessions"
	"github.com/boltdb/bolt"
	"github.com/webx-top/echo"
	ss "github.com/webx-top/echo/middleware/session/engine"
)

func New(opts *BoltOptions) BoltStore {
	store, err := NewBoltStore(opts)
	if err != nil {
		panic(err.Error())
	}
	return store
}

func Reg(store BoltStore, args ...string) {
	name := `bolt`
	if len(args) > 0 {
		name = args[0]
	}
	ss.Reg(name, store)
}

func RegWithOptions(opts *BoltOptions, args ...string) {
	Reg(New(opts), args...)
}

type BoltStore interface {
	ss.Store
}

type BoltOptions struct {
	File           string               `json:"file"`
	KeyPairs       [][]byte             `json:"keyPairs"`
	BucketName     string               `json:"bucketName"`
	SessionOptions *echo.SessionOptions `json:"session"`
}

// NewBoltStore ./sessions.db
func NewBoltStore(opts *BoltOptions) (BoltStore, error) {
	config := store.Config{
		SessionOptions: sessions.Options{
			Path:     opts.SessionOptions.Path,
			Domain:   opts.SessionOptions.Domain,
			MaxAge:   opts.SessionOptions.MaxAge,
			Secure:   opts.SessionOptions.Secure,
			HttpOnly: opts.SessionOptions.HttpOnly,
		},
		DBOptions: store.Options{BucketName: []byte(opts.BucketName)},
	}
	b := &boltStore{
		config:   &config,
		keyPairs: opts.KeyPairs,
		dbFile:   opts.File,
		Storex: &Storex{
			Store: &store.Store{},
		},
	}
	b.Storex.b = b
	return b, nil
}

type Storex struct {
	*store.Store
	db          *bolt.DB
	b           *boltStore
	initialized bool
}

func (s *Storex) Get(ctx echo.Context, name string) (*sessions.Session, error) {
	if s.initialized == false {
		err := s.b.Init()
		if err != nil {
			return nil, err
		}
	}
	return s.Store.Get(ctx, name)
}

func (s *Storex) New(ctx echo.Context, name string) (*sessions.Session, error) {
	if s.initialized == false {
		err := s.b.Init()
		if err != nil {
			return nil, err
		}
	}
	return s.Store.New(ctx, name)
}

func (s *Storex) Save(ctx echo.Context, session *sessions.Session) error {
	if s.initialized == false {
		err := s.b.Init()
		if err != nil {
			return err
		}
	}
	return s.Store.Save(ctx, session)
}

type boltStore struct {
	*Storex
	config   *store.Config
	keyPairs [][]byte
	quiteC   chan<- struct{}
	doneC    <-chan struct{}
	dbFile   string
}

func (c *boltStore) Options(options echo.SessionOptions) {
	c.config.SessionOptions = sessions.Options{
		Path:     options.Path,
		Domain:   options.Domain,
		MaxAge:   options.MaxAge,
		Secure:   options.Secure,
		HttpOnly: options.HttpOnly,
	}
	stor, err := store.New(c.Storex.db, *c.config, c.keyPairs...)
	if err != nil {
		panic(err.Error())
	}
	c.Store = stor
}

func (c *boltStore) Close() error {
	// Invoke a reaper which checks and removes expired sessions periodically.
	if c.quiteC != nil && c.doneC != nil {
		reaper.Quit(c.quiteC, c.doneC)
	}

	if c.Storex.db != nil {
		c.Storex.db.Close()
	}

	return nil
}

func (b *boltStore) Init() error {
	if b.Storex.db == nil {
		var err error
		b.Storex.db, err = bolt.Open(b.dbFile, 0666, nil)
		if err != nil {
			return err
		}
		b.Storex.Store, err = store.New(b.Storex.db, *b.config, b.keyPairs...)
		if err != nil {
			return err
		}
		b.quiteC, b.doneC = reaper.Run(b.Storex.db, reaper.Options{
			BucketName:    b.config.DBOptions.BucketName,
			CheckInterval: time.Duration(int64(b.config.SessionOptions.MaxAge)) * time.Second,
		})
		runtime.SetFinalizer(b, func(b *boltStore) {
			b.Close()
		})
	}
	b.Storex.initialized = true
	return nil
}
