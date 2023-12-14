package cloudbackup

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"path"
	"path/filepath"

	"github.com/admpub/nging/v5/application/dbschema"
	"github.com/admpub/nging/v5/application/library/common"
	"github.com/admpub/nging/v5/application/library/notice"
	"github.com/admpub/nging/v5/application/model"
	"github.com/studio-b12/gowebdav"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
)

func init() {
	Register(model.StorageEngineWebDAV, newStorageWebDAV, webDavForms, `WebDAV`)
}

func newStorageWebDAV(ctx echo.Context, cfg dbschema.NgingCloudBackup) (Storager, error) {
	if len(cfg.StorageConfig) == 0 {
		return nil, ErrEmptyConfig
	}
	conf := echo.H{}
	err := json.Unmarshal([]byte(cfg.StorageConfig), &conf)
	if err != nil {
		return nil, err
	}
	password := common.Crypto().Decode(conf.String(`password`))
	return NewStorageWebDAV(conf.String(`uri`), conf.String(`username`), password), nil
}

var webDavForms = []Form{
	{Type: `text`, Label: `网址`, Name: `storageConfig.uri`, Required: true, Placeholder: `http(s)://<IP或域名>:<端口>`, Pattern: `^http[s]?://`},
	{Type: `text`, Label: `用户名`, Name: `storageConfig.username`, Required: true},
	{Type: `password`, Label: `密码`, Name: `storageConfig.password`, Required: true},
}

func NewStorageWebDAV(uri, username, password string) Storager {
	return &StorageWebDAV{uri: uri, username: username, password: password}
}

type StorageWebDAV struct {
	uri      string // host:port
	username string
	password string
	conn     *gowebdav.Client
	prog     notice.Progressor
}

func (s *StorageWebDAV) Connect() (err error) {
	s.conn = gowebdav.NewClient(s.uri, s.username, s.password)
	err = s.conn.Connect()
	return
}

func (s *StorageWebDAV) Put(ctx context.Context, reader io.Reader, ppath string, size int64) (err error) {
	s.conn.MkdirAll(path.Dir(ppath), 0)
	err = s.conn.WriteStream(ppath, reader, 0)
	return err
}

func (s *StorageWebDAV) Download(ctx context.Context, ppath string, w io.Writer) error {
	resp, err := s.conn.ReadStream(ppath)
	if err != nil {
		return err
	}
	defer resp.Close()
	if s.prog != nil {
		stat, err := s.conn.Stat(ppath)
		if err != nil {
			return err
		}
		s.prog.Add(stat.Size())
		w = s.prog.ProxyWriter(w)
		defer s.prog.Reset()
	}
	_, err = io.Copy(w, resp)
	return err
}

func (s *StorageWebDAV) SetProgressor(prog notice.Progressor) {
	s.prog = prog
}

func (s *StorageWebDAV) Restore(ctx context.Context, ppath string, destpath string, callback func(from, to string)) error {
	stat, err := s.conn.Stat(ppath)
	if err != nil {
		return err
	}
	if !stat.IsDir() {
		if callback != nil {
			callback(ppath, destpath)
		}
		return DownloadFile(s, ctx, ppath, destpath)
	}
	return s.recursiveRestoreDir(ctx, ppath, destpath, callback)
}

func (s *StorageWebDAV) recursiveRestoreDir(ctx context.Context, ppath string, destpath string, callback func(from, to string)) error {
	dirs, err := s.conn.ReadDir(ppath)
	if err != nil {
		return err
	}
	for _, dir := range dirs {
		spath := path.Join(ppath, dir.Name())
		dest := filepath.Join(destpath, dir.Name())
		if dir.IsDir() {
			err = com.MkdirAll(dest, os.ModePerm)
			if err == nil {
				err = s.recursiveRestoreDir(ctx, spath, dest, callback)
			}
		} else {
			if callback != nil {
				callback(spath, dest)
			}
			err = DownloadFile(s, ctx, spath, dest)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *StorageWebDAV) RemoveDir(ctx context.Context, ppath string) error {
	return s.conn.RemoveAll(ppath)
}

func (s *StorageWebDAV) Remove(ctx context.Context, ppath string) error {
	return s.conn.Remove(ppath)
}

func (s *StorageWebDAV) Close() (err error) {
	if s.conn == nil {
		return
	}
	//err = s.conn.Close()
	return
}
