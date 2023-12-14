package cloudbackup

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/admpub/nging/v5/application/dbschema"
	"github.com/admpub/nging/v5/application/library/common"
	"github.com/admpub/nging/v5/application/library/s3manager"
	"github.com/admpub/nging/v5/application/library/s3manager/fileinfo"
	"github.com/admpub/nging/v5/application/library/s3manager/s3client"
	"github.com/admpub/nging/v5/application/model"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
)

func init() {
	Register(model.StorageEngineS3, newStorageS3, s3Forms, `类S3对象存储`)
}

var storageS3Test bool

func newStorageS3(ctx echo.Context, cfg dbschema.NgingCloudBackup) (Storager, error) {
	m := dbschema.NewNgingCloudStorage(ctx)
	err := m.Get(nil, `id`, cfg.DestStorage)
	if err != nil {
		return nil, err
	}
	if len(m.Endpoint) == 0 {
		return nil, ctx.NewError(code.InvalidParameter, `Endpoint无效`)
	}
	if !storageS3Test {
		m.Secret = common.Crypto().Decode(m.Secret)
	}
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

func (s *StorageS3) Download(ctx context.Context, ppath string, w io.Writer) error {
	resp, err := s.conn.Get(ctx, ppath)
	if err != nil {
		return err
	}
	defer resp.Close()
	_, err = io.Copy(w, resp)
	return err
}

func (s *StorageS3) Restore(ctx context.Context, ppath string, destpath string) error {
	objectPrefix := strings.TrimPrefix(ppath, `/`)
	if !strings.HasSuffix(objectPrefix, `/`) {
		_, err := s.conn.Stat(ctx, objectPrefix)
		if err == nil {
			return DownloadFile(s, ctx, ppath, destpath)
		}
		if !s.conn.ErrIsNotExist(err) {
			return err
		}
		objectPrefix += `/`
	}
	awsclient, err := s.conn.AWSClient()
	if err != nil {
		return err
	}
	input := &s3.ListObjectsInput{
		Bucket:  aws.String(s.conn.BucketName()),
		Prefix:  aws.String(objectPrefix),
		MaxKeys: aws.Int64(200),
	}
	var _err error
	err = awsclient.ListObjectsPagesWithContext(ctx, input, func(p *s3.ListObjectsOutput, lastPage bool) bool {
		for _, object := range p.CommonPrefixes {
			if object.Prefix == nil {
				continue
			}
			if len(objectPrefix) > 0 {
				key := strings.TrimPrefix(*object.Prefix, objectPrefix)
				object.Prefix = &key
			}
			if len(*object.Prefix) == 0 {
				continue
			}
			obj := fileinfo.NewStr(*object.Prefix)
			dest := filepath.Join(destpath, obj.Name())
			spath := filepath.Join(objectPrefix, obj.Name())
			if obj.IsDir() {
				_err = com.MkdirAll(dest, os.ModePerm)
			} else {
				_err = DownloadFile(s, ctx, spath, dest)
			}
			if _err != nil {
				return false
			}
		}
		for _, object := range p.Contents {
			if object.Key == nil {
				continue
			}
			if len(objectPrefix) > 0 {
				key := strings.TrimPrefix(*object.Key, objectPrefix)
				object.Key = &key
			}
			if len(*object.Key) == 0 {
				continue
			}
			obj := fileinfo.NewS3(object)
			dest := filepath.Join(destpath, obj.Name())
			spath := filepath.Join(objectPrefix, obj.Name())
			if obj.IsDir() {
				_err = com.MkdirAll(dest, os.ModePerm)
			} else {
				_err = DownloadFile(s, ctx, spath, dest)
			}
			if _err != nil {
				return false
			}
		}
		return true // continue paging
	})
	if _err != nil {
		err = _err
	}
	return err
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
