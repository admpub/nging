package cloudbackup

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/admpub/nging/v5/application/dbschema"
	"github.com/admpub/nging/v5/application/model"
	"github.com/webx-top/echo"
)

type Storager interface {
	Connect() (err error)
	Put(ctx context.Context, reader io.Reader, ppath string, size int64) (err error)
	RemoveDir(ctx context.Context, ppath string) error
	Remove(ctx context.Context, ppath string) error
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

var storages = map[string]Constructor{}

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
