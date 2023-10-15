package file

import (
	"os"
	"sync"
	"time"

	"github.com/admpub/sessions"
	"github.com/webx-top/echo"
	ss "github.com/webx-top/echo/middleware/session/engine"
)

func New(opts *FileOptions) sessions.Store {
	store := NewFilesystemStore(opts)
	return store
}

func Reg(store sessions.Store, args ...string) {
	name := `file`
	if len(args) > 0 {
		name = args[0]
	}
	ss.Reg(name, store)
}

func RegWithOptions(opts *FileOptions, args ...string) sessions.Store {
	store := New(opts)
	Reg(store, args...)
	return store
}

type FileOptions struct {
	SavePath      string        `json:"savePath"`
	KeyPairs      [][]byte      `json:"-"`
	CheckInterval time.Duration `json:"checkInterval"`
	MaxAge        int           `json:"maxAge"`
	MaxLength     int           `json:"maxLength"`
}

// NewFilesystemStore returns a new FilesystemStore.
//
// The path argument is the directory where sessions will be saved. If empty
// it will use os.TempDir().
//
// See NewCookieStore() for a description of the other parameters.
func NewFilesystemStore(opts *FileOptions) sessions.Store {
	if len(opts.SavePath) > 0 {
		fi, err := os.Stat(opts.SavePath)
		if os.IsNotExist(err) || !fi.IsDir() {
			err = os.MkdirAll(opts.SavePath, os.ModePerm)
			if err != nil {
				panic(err)
			}
		}
	}
	s := &filesystemStore{
		FilesystemStore: sessions.NewFilesystemStore(opts.SavePath, opts.KeyPairs...),
		options:         opts,
	}
	if opts.MaxLength > 0 {
		s.MaxLength(opts.MaxLength)
	}
	return s
}

type filesystemStore struct {
	*sessions.FilesystemStore
	options *FileOptions
	quiteC  chan<- struct{}
	doneC   <-chan struct{}
	once    sync.Once
}

func (m *filesystemStore) Get(ctx echo.Context, name string) (*sessions.Session, error) {
	m.Init()
	return m.FilesystemStore.Get(ctx, name)
}

func (m *filesystemStore) New(ctx echo.Context, name string) (*sessions.Session, error) {
	return m.FilesystemStore.New(ctx, name)
}

func (m *filesystemStore) Reload(ctx echo.Context, session *sessions.Session) error {
	return m.FilesystemStore.Reload(ctx, session)
}

func (m *filesystemStore) Save(ctx echo.Context, session *sessions.Session) error {
	return m.FilesystemStore.Save(ctx, session)
}

func (m *filesystemStore) Close() (err error) {
	// Invoke a reaper which checks and removes expired sessions periodically.
	if m.quiteC != nil && m.doneC != nil {
		m.StopCleanup(m.quiteC, m.doneC)
	}
	return
}

func (m *filesystemStore) Init() {
	m.once.Do(m.init)
}

func (m *filesystemStore) init() {
	m.Close()
	m.quiteC, m.doneC = m.Cleanup(m.options.CheckInterval, m.options.MaxAge)
}
