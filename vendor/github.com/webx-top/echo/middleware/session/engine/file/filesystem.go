package file

import (
	"os"
	"sync"

	"github.com/admpub/sessions"
	"github.com/webx-top/echo"
	ss "github.com/webx-top/echo/middleware/session/engine"
)

func New(opts *FileOptions) sessions.Store {
	store := NewFilesystemStore(opts.SavePath, opts.KeyPairs...)
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
	SavePath string   `json:"savePath"`
	KeyPairs [][]byte `json:"keyPairs"`
}

// NewFilesystemStore returns a new FilesystemStore.
//
// The path argument is the directory where sessions will be saved. If empty
// it will use os.TempDir().
//
// See NewCookieStore() for a description of the other parameters.
func NewFilesystemStore(path string, keyPairs ...[]byte) sessions.Store {
	if len(path) > 0 {
		fi, err := os.Stat(path)
		if os.IsNotExist(err) || !fi.IsDir() {
			err = os.MkdirAll(path, os.ModePerm)
			if err != nil {
				panic(err)
			}
		}
	}
	s := &filesystemStore{
		FilesystemStore: sessions.NewFilesystemStore(path, keyPairs...),
	}
	return s
}

type filesystemStore struct {
	*sessions.FilesystemStore
	quiteC chan<- struct{}
	doneC  <-chan struct{}
	once   sync.Once
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
	m.quiteC, m.doneC = m.Cleanup(0, 0)
}
