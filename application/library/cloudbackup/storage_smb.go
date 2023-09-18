package cloudbackup

import (
	"context"
	"encoding/json"
	"io"
	"net"
	"os"
	"path"
	"strings"

	"github.com/admpub/nging/v5/application/dbschema"
	"github.com/admpub/nging/v5/application/library/common"
	"github.com/admpub/nging/v5/application/model"
	smb "github.com/hirochachacha/go-smb2"
	"github.com/webx-top/echo"
)

func init() {
	Register(model.StorageEngineSMB, newStorageSMB, smbForms, `SMB`)
}

func newStorageSMB(ctx echo.Context, cfg dbschema.NgingCloudBackup) (Storager, error) {
	if len(cfg.StorageConfig) == 0 {
		return nil, ErrEmptyConfig
	}
	conf := echo.H{}
	err := json.Unmarshal([]byte(cfg.StorageConfig), &conf)
	if err != nil {
		return nil, err
	}
	password := common.Crypto().Decode(conf.String(`password`))
	return NewStorageSMB(conf.String(`addr`), conf.String(`username`), password, conf.String(`sharename`)), nil
}

var smbForms = []Form{
	{Type: `text`, Label: `主机地址`, Name: `storageConfig.addr`, Required: true, Placeholder: `<IP或域名>:<端口>`},
	{Type: `text`, Label: `用户名`, Name: `storageConfig.username`, Required: true},
	{Type: `password`, Label: `密码`, Name: `storageConfig.password`, Required: true},
	{Type: `text`, Label: `共享名称`, Name: `storageConfig.sharename`, Required: true},
}

func NewStorageSMB(addr, username, password, sharename string) Storager {
	return &StorageSMB{addr: addr, username: username, password: password, sharename: sharename}
}

type StorageSMB struct {
	addr      string // host:port
	username  string
	password  string
	sharename string
	conn      net.Conn
	session   *smb.Session
	share     *smb.Share
}

func (s *StorageSMB) Connect() (err error) {
	if !strings.Contains(s.addr, `:`) {
		s.addr += `:445`
	}
	s.conn, err = net.Dial("tcp", s.addr)
	if err != nil {
		return
	}

	d := &smb.Dialer{
		Initiator: &smb.NTLMInitiator{
			User:     s.username,
			Password: s.password,
		},
	}

	s.session, err = d.Dial(s.conn)
	if err != nil {
		s.Close()
		return
	}

	s.share, err = s.session.Mount(s.sharename)
	if err != nil {
		s.Close()
		return
	}
	return
}

func (s *StorageSMB) Put(ctx context.Context, reader io.Reader, ppath string, size int64) (err error) {
	s.share.MkdirAll(path.Dir(ppath), os.ModePerm)
	var fp *smb.File
	fp, err = s.share.Create(ppath)
	if err != nil {
		return
	}
	defer fp.Close()
	_, err = io.Copy(fp, reader)
	return
}

func (s *StorageSMB) RemoveDir(ctx context.Context, ppath string) error {
	return s.share.RemoveAll(ppath)
}

func (s *StorageSMB) Remove(ctx context.Context, ppath string) error {
	return s.share.Remove(ppath)
}

func (s *StorageSMB) Close() (err error) {
	if s.share != nil {
		err = s.share.Umount()
	}
	if s.session != nil {
		err = s.session.Logoff()
	}
	if s.conn != nil {
		err = s.conn.Close()
	}
	return
}
