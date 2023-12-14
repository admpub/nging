package cloudbackup

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/admpub/nging/v5/application/dbschema"
	"github.com/admpub/nging/v5/application/model"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
)

type Storager interface {
	Connect() (err error)
	Put(ctx context.Context, reader io.Reader, ppath string, size int64) (err error)
	Download(ctx context.Context, ppath string, w io.Writer) error
	RemoveDir(ctx context.Context, ppath string) error
	Remove(ctx context.Context, ppath string) error
	Restore(ctx context.Context, ppath string, destpath string) error
	Close() (err error)
}

type Form struct {
	Type        string
	Label       string
	Name        string
	Required    bool
	Pattern     string
	Placeholder string
}

var Forms = map[string][]Form{}

func HasForm(engineName string, formName string) bool {
	for _, f := range Forms[engineName] {
		if f.Name == formName {
			return true
		}
	}
	return false
}

type Constructor func(ctx echo.Context, cfg dbschema.NgingCloudBackup) (Storager, error)

var storages = map[string]Constructor{
	`mock`: newStorageMock,
}

func Register(name string, constructor Constructor, forms []Form, label string) {
	storages[name] = constructor
	Forms[name] = forms
	model.CloudBackupStorageEngines.Add(name, label)
}

var ErrUnsupported = errors.New(`unsupported storage engine`)
var ErrEmptyConfig = errors.New(`empty config`)

func NewStorage(ctx echo.Context, cfg dbschema.NgingCloudBackup) (Storager, error) {
	cr, ok := storages[cfg.StorageEngine]
	if !ok {
		return nil, fmt.Errorf(`%w: %s`, ErrUnsupported, cfg.StorageEngine)
	}
	return cr(ctx, cfg)
}

func DownloadFile(s Storager, ctx context.Context, ppath string, dest string) error {
	dir := filepath.Dir(dest)
	com.MkdirAll(dir, os.ModePerm)
	fi, err := os.Create(dest)
	if err != nil {
		return err
	}
	err = s.Download(ctx, ppath, fi)
	fi.Close()
	return err
}
