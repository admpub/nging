package cloudbackup

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/admpub/nging/v5/application/dbschema"
	"github.com/admpub/nging/v5/application/library/common"
	"github.com/admpub/nging/v5/application/library/notice"
	"github.com/admpub/nging/v5/application/model"
	"github.com/jlaffaye/ftp"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
)

func init() {
	Register(model.StorageEngineFTP, newStorageFTP, ftpForms, `FTP`)
}

func newStorageFTP(ctx echo.Context, cfg dbschema.NgingCloudBackup) (Storager, error) {
	if len(cfg.StorageConfig) == 0 {
		return nil, ErrEmptyConfig
	}
	conf := echo.H{}
	err := json.Unmarshal([]byte(cfg.StorageConfig), &conf)
	if err != nil {
		return nil, err
	}
	password := common.Crypto().Decode(conf.String(`password`))
	return NewStorageFTP(conf.String(`addr`), conf.String(`username`), password), nil
}

var ftpForms = []Form{
	{Type: `text`, Label: `主机地址`, Name: `storageConfig.addr`, Required: true, Placeholder: `<IP或域名>:<端口>`},
	{Type: `text`, Label: `用户名`, Name: `storageConfig.username`, Required: true},
	{Type: `password`, Label: `密码`, Name: `storageConfig.password`, Required: true},
}

func NewStorageFTP(addr, username, password string) Storager {
	return &StorageFTP{addr: addr, username: username, password: password}
}

type StorageFTP struct {
	addr     string // host:port
	username string
	password string
	conn     *ftp.ServerConn
	prog     notice.Progressor
}

func (s *StorageFTP) Connect() (err error) {
	if !strings.Contains(s.addr, `:`) {
		s.addr += `:21`
	}
	s.conn, err = ftp.Dial(s.addr, ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		return
	}

	err = s.conn.Login(s.username, s.password) // anonymous
	return
}

func (s *StorageFTP) MkdirAll(dir string) (err error) {
	err = s.conn.ChangeDir(dir)
	if err == nil {
		return
	}
	var notExistDirs []string

LOOP:
	notExistDirs = append(notExistDirs, path.Base(dir))
	dir = path.Dir(dir)
	if err = s.conn.ChangeDir(dir); err != nil {
		if len(dir) == 0 || dir == `/` || dir == `.` {
			return
		}
		goto LOOP
	}

	for j := len(notExistDirs) - 1; j >= 0; j-- {
		if dir != `/` {
			dir += `/`
		}
		dir += notExistDirs[j]
		//println(`mkdir:`, dir)
		err = s.conn.MakeDir(dir)
		if err != nil {
			break
		}
	}
	return
}

func (s *StorageFTP) Put(ctx context.Context, reader io.Reader, ppath string, size int64) (err error) {
	dir := path.Dir(ppath)
	s.MkdirAll(dir)
	err = s.conn.Stor(ppath, reader)
	return err
}

func (s *StorageFTP) Download(ctx context.Context, ppath string, w io.Writer) error {
	resp, err := s.conn.Retr(ppath)
	if err != nil {
		return err
	}
	defer resp.Close()
	// if s.prog != nil {
	// 	stat, err := resp.Stat()
	// 	if err != nil {
	// 		return err
	// 	}
	// 	s.prog.Add(stat.Size)
	// 	w = s.prog.ProxyWriter(w)
	//	defer s.prog.Reset()
	// }
	_, err = io.Copy(w, resp)
	return err
}

func (s *StorageFTP) SetProgressor(prog notice.Progressor) {
	s.prog = prog
}

func (s *StorageFTP) Restore(ctx context.Context, ppath string, destpath string, callback func(from, to string)) error {
	walker := s.conn.Walk(ppath)
	var err error
	for walker.Next() {
		spath := walker.Path()
		subdir := strings.TrimPrefix(spath, ppath)
		dest := filepath.Join(destpath, subdir)
		if walker.Stat().Type == ftp.EntryTypeFolder {
			err = com.MkdirAll(dest, os.ModePerm)
		} else {
			err = DownloadFile(s, ctx, spath, dest)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *StorageFTP) RemoveDir(ctx context.Context, ppath string) error {
	return s.conn.RemoveDir(ppath)
}

func (s *StorageFTP) Remove(ctx context.Context, ppath string) error {
	return s.conn.Delete(ppath)
}

func (s *StorageFTP) Close() (err error) {
	if s.conn == nil {
		return
	}
	err = s.conn.Quit()
	return
}
