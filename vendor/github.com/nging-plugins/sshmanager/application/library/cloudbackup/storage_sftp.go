package cloudbackup

import (
	"context"
	"io"
	"path"

	nd "github.com/admpub/nging/v5/application/dbschema"
	"github.com/admpub/nging/v5/application/library/cloudbackup"
	"github.com/admpub/nging/v5/application/library/sftpmanager"
	"github.com/admpub/nging/v5/application/model"
	"github.com/nging-plugins/sshmanager/application/dbschema"
	sshconf "github.com/nging-plugins/sshmanager/application/library/config"
	"github.com/webx-top/echo"
)

func init() {
	cloudbackup.Register(model.StorageEngineSFTP, newStorageSFTP, sftpForms, `SFTP`)
}

func newStorageSFTP(ctx echo.Context, cfg nd.NgingCloudBackup) (cloudbackup.Storager, error) {
	m := dbschema.NewNgingSshUser(ctx)
	err := m.Get(nil, `id`, cfg.DestStorage)
	if err != nil {
		return nil, err
	}
	conf := sshconf.ToSFTPConfig(m)
	return NewStorageSFTP(conf), nil
}

var sftpForms = []cloudbackup.Form{
	{Type: `text`, Label: `SSH账号`, Name: `destStorage`, Required: true},
}

func NewStorageSFTP(cfg sftpmanager.Config) cloudbackup.Storager {
	return &StorageSFTP{cfg: cfg}
}

type StorageSFTP struct {
	cfg  sftpmanager.Config
	conn *sftpmanager.SftpManager
}

func (s *StorageSFTP) Connect() (err error) {
	s.conn = sftpmanager.New(sftpmanager.DefaultConnector, &s.cfg, 0)
	s.conn.Client()
	if s.conn.ConnError() != nil {
		err = s.conn.ConnError()
	}
	return
}

func (s *StorageSFTP) Put(ctx context.Context, reader io.Reader, ppath string, size int64) (err error) {
	s.conn.MkdirAll(ctx, path.Dir(ppath))
	err = s.conn.Put(ctx, reader, ppath, size)
	return
}

func (s *StorageSFTP) RemoveDir(ctx context.Context, ppath string) error {
	return s.conn.RemoveDir(ppath)
}

func (s *StorageSFTP) Remove(ctx context.Context, ppath string) error {
	return s.conn.Remove(ppath)
}

func (s *StorageSFTP) Close() (err error) {
	if s.conn == nil {
		return
	}
	err = s.conn.Close()
	return
}
