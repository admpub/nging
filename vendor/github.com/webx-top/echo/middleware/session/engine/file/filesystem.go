package file

import (
	"os"

	"github.com/admpub/sessions"
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
	return &filesystemStore{sessions.NewFilesystemStore(path, keyPairs...)}
}

type filesystemStore struct {
	*sessions.FilesystemStore
}
