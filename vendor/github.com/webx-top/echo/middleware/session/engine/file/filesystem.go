package file

import (
	"github.com/admpub/sessions"
	"github.com/webx-top/echo"
	ss "github.com/webx-top/echo/middleware/session/engine"
)

func New(opts *FileOptions) FilesystemStore {
	store := NewFilesystemStore(opts.SavePath, opts.KeyPairs...)
	store.Options(*opts.SessionOptions)
	return store
}

func Reg(store FilesystemStore, args ...string) {
	name := `file`
	if len(args) > 0 {
		name = args[0]
	}
	ss.Reg(name, store)
}

func RegWithOptions(opts *FileOptions, args ...string) {
	Reg(New(opts), args...)
}

type FilesystemStore interface {
	ss.Store
}

type FileOptions struct {
	SavePath             string   `json:"savePath"`
	KeyPairs             [][]byte `json:"keyPairs"`
	*echo.SessionOptions `json:"session"`
}

// NewFilesystemStore returns a new FilesystemStore.
//
// The path argument is the directory where sessions will be saved. If empty
// it will use os.TempDir().
//
// See NewCookieStore() for a description of the other parameters.
func NewFilesystemStore(path string, keyPairs ...[]byte) FilesystemStore {
	return &filesystemStore{sessions.NewFilesystemStore(path, keyPairs...)}
}

type filesystemStore struct {
	*sessions.FilesystemStore
}

func (c *filesystemStore) Options(options echo.SessionOptions) {
	c.FilesystemStore.Options = &sessions.Options{
		Path:     options.Path,
		Domain:   options.Domain,
		MaxAge:   options.MaxAge,
		Secure:   options.Secure,
		HttpOnly: options.HttpOnly,
	}
}
