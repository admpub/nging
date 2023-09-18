package cloudbackup

import (
	"context"
	"io"

	"github.com/admpub/nging/v5/application/dbschema"
	"github.com/admpub/nging/v5/application/library/common"
	"github.com/admpub/nging/v5/application/library/s3manager"
	"github.com/admpub/nging/v5/application/library/s3manager/s3client"
	"github.com/admpub/nging/v5/application/model"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
)

func init() {
	Register(model.StorageEngineS3, newStorageS3, s3Forms, `类S3对象存储`)
}

func newStorageS3(ctx echo.Context, cfg dbschema.NgingCloudBackup) (Storager, error) {
	m := dbschema.NewNgingCloudStorage(ctx)
	err := m.Get(nil, `id`, cfg.DestStorage)
	if err != nil {
		return nil, err
	}
	if len(m.Endpoint) == 0 {
		return nil, ctx.NewError(code.InvalidParameter, `Endpoint无效`)
	}
	m.Secret = common.Crypto().Decode(m.Secret)
	return NewStorageS3(*m), nil
}

var s3Forms = []Form{
	{Type: `text`, Label: `云存储账号`, Name: `destStorage`, Required: true},
}

func NewStorageS3(cfg dbschema.NgingCloudStorage) Storager {
	return &StorageS3{cfg: cfg}
}

type StorageS3 struct {
	cfg  dbschema.NgingCloudStorage
	conn *s3manager.S3Manager
}

func (s *StorageS3) Connect() (err error) {
	s.conn = s3client.New(&s.cfg, 0)
	s.conn.Client()
	if s.conn.ConnError() != nil {
		err = s.conn.ConnError()
	}
	return
}

func (s *StorageS3) Put(ctx context.Context, reader io.Reader, ppath string, size int64) (err error) {
	//s.conn.MkdirAll(ctx, path.Dir(ppath))
	err = s.conn.Put(ctx, reader, ppath, size)
	return
}

func (s *StorageS3) RemoveDir(ctx context.Context, ppath string) error {
	return s.conn.RemoveDir(ctx, ppath)
}

func (s *StorageS3) Remove(ctx context.Context, ppath string) error {
	return s.conn.Remove(ctx, ppath)
}

func (s *StorageS3) Close() (err error) {
	if s.conn == nil {
		return
	}
	// err = s.conn.Close()
	return
}
