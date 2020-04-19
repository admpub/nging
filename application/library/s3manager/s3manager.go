/*
   Nging is a toolbox for webmasters
   Copyright (C) 2018-present  Wenhui Shen <swh@admpub.com>

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU Affero General Public License as published
   by the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Affero General Public License for more details.

   You should have received a copy of the GNU Affero General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package s3manager

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/admpub/nging/application/library/charset"
	"github.com/admpub/nging/application/library/filemanager"
	"github.com/aws/aws-sdk-go/service/s3"
	minio "github.com/minio/minio-go"
	"github.com/pkg/errors"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
)

func New(client *minio.Client, bucketName string, editableMaxSize int64) *S3Manager {
	return &S3Manager{
		client:          client,
		bucketName:      bucketName,
		EditableMaxSize: editableMaxSize,
	}
}

type S3Manager struct {
	client          *minio.Client
	bucketName      string
	EditableMaxSize int64
}

func (s *S3Manager) Client() *minio.Client {
	return s.client
}

func (s *S3Manager) BucketName() string {
	return s.bucketName
}

func (s *S3Manager) SetBucketName(bucketName string) *S3Manager {
	s.bucketName = bucketName
	return s
}

func (s *S3Manager) Edit(ctx echo.Context, ppath string, content string, encoding string) (interface{}, error) {
	f, err := s.Get(ppath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if s.EditableMaxSize > 0 && fi.Size > s.EditableMaxSize {
		return nil, ctx.E(`很抱歉，不支持编辑超过%v的文件`, com.FormatBytes(s.EditableMaxSize))
	}
	encoding = strings.ToLower(encoding)
	isUTF8 := len(encoding) == 0 || encoding == `utf-8`
	if ctx.IsPost() {
		b := []byte(content)
		if !isUTF8 {
			b, err = charset.Convert(`utf-8`, encoding, b)
			if err != nil {
				return nil, err
			}
		}
		r := bytes.NewReader(b)
		err = s.Put(r, ppath, int64(len(b)))
		if err != nil {
			return nil, ctx.E(ppath + `:` + err.Error())
		}
		return nil, err
	}

	dat, err := ioutil.ReadAll(f)
	if err == nil && !isUTF8 {
		dat, err = charset.Convert(encoding, `utf-8`, dat)
	}
	return string(dat), err
}

func (s *S3Manager) Mkbucket(bucketName string, regions ...string) error {
	var region string
	if len(regions) > 0 {
		region = regions[0]
	}
	return s.client.MakeBucket(bucketName, region)
}

func (s *S3Manager) Mkdir(ppath, newName string) error {
	objectName := strings.TrimPrefix(ppath, `/`)
	objectName = path.Join(objectName, newName)
	if !strings.HasSuffix(objectName, `/`) {
		objectName += `/`
	}
	_, err := s.client.PutObject(s.bucketName, objectName, nil, 0, minio.PutObjectOptions{})
	return err
}

func (s *S3Manager) Rename(ppath, newName string) error {
	objectName := strings.TrimPrefix(ppath, `/`)
	// Source object
	src := minio.NewSourceInfo(s.bucketName, objectName, nil)
	newName = strings.TrimPrefix(newName, `/`)
	dst, err := minio.NewDestinationInfo(s.bucketName, newName, nil, nil)
	if err != nil {
		return err
	}

	// Initiate copy object.
	err = s.client.CopyObject(dst, src)
	if err != nil {
		return err
	}
	err = s.client.RemoveObject(s.bucketName, objectName)
	return err
}

func (s *S3Manager) Chown(ppath string, uid, gid int) error {
	return nil
}

func (s *S3Manager) Chmod(ppath string, mode os.FileMode) error {
	return nil
}

func (s *S3Manager) Search(ppath string, prefix string, num int) []string {
	var paths []string
	doneCh := make(chan struct{})
	defer close(doneCh)
	objectPrefix := path.Join(ppath, prefix)
	objectPrefix = strings.TrimPrefix(objectPrefix, `/`)
	objectCh := s.client.ListObjectsV2(s.bucketName, objectPrefix, false, doneCh)
	for object := range objectCh {
		if object.Err != nil {
			continue
		}
		paths = append(paths, object.Key)
	}
	return paths
}

func (s *S3Manager) Remove(ppath string) error {
	if len(ppath) == 0 {
		return errors.New("path invalid")
	}
	if strings.HasSuffix(ppath, `/`) {
		return s.RemoveDir(ppath)
	}
	objectName := strings.TrimPrefix(ppath, `/`)
	return s.client.RemoveObject(s.bucketName, objectName)
}

func (s *S3Manager) RemoveDir(ppath string) error {
	objectName := strings.TrimPrefix(ppath, `/`)
	if !strings.HasSuffix(objectName, `/`) {
		objectName += `/`
	}
	if objectName == `/` {
		return s.Clear()
	}
	s.client.RemoveObject(s.bucketName, objectName)
	doneCh := make(chan struct{})
	defer close(doneCh)
	objectCh := s.client.ListObjectsV2(s.bucketName, objectName, true, doneCh)
	for object := range objectCh {
		if object.Err != nil {
			continue
		}
		if len(object.Key) == 0 {
			continue
		}
		err := s.client.RemoveObject(s.bucketName, object.Key)
		if err != nil {
			return err
		}
	}
	return nil
}

// Clear 清空所有数据【慎用】
func (s *S3Manager) Clear() error {
	deleted := make(chan string)
	defer close(deleted)
	removeObjects := s.client.RemoveObjects(s.bucketName, deleted)
	for removeObject := range removeObjects {
		if removeObject.Err != nil {
			return removeObject.Err
		}
	}
	return nil
}

func (s *S3Manager) Upload(ctx echo.Context, ppath string) error {
	fileSrc, fileHdr, err := ctx.Request().FormFile(`file`)
	if err != nil {
		return err
	}
	defer fileSrc.Close()
	objectName := path.Join(ppath, fileHdr.Filename)
	return s.Put(fileSrc, objectName, fileHdr.Size)
}

// Put 提交数据
func (s *S3Manager) Put(reader io.Reader, objectName string, size int64) (err error) {
	opts := minio.PutObjectOptions{ContentType: "application/octet-stream"}
	objectName = strings.TrimPrefix(objectName, `/`)
	_, err = s.client.PutObject(s.bucketName, objectName, reader, size, opts)
	return
}

// Get 获取数据
func (s *S3Manager) Get(ppath string) (*minio.Object, error) {
	objectName := strings.TrimPrefix(ppath, `/`)
	f, err := s.client.GetObject(s.bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return f, errors.WithMessage(err, objectName)
	}
	return f, err
}

// Stat 获取对象信息
func (s *S3Manager) Stat(ppath string) (minio.ObjectInfo, error) {
	objectName := strings.TrimPrefix(ppath, `/`)
	f, err := s.client.StatObject(s.bucketName, objectName, minio.StatObjectOptions{})
	if err != nil {
		return f, errors.WithMessage(err, objectName)
	}
	return f, err
}

// Exists 对象是否存在
func (s *S3Manager) Exists(ppath string) (bool, error) {
	_, err := s.Stat(ppath)
	if s.ErrIsNotExist(err) == false {
		return false, err
	}
	return true, nil
}

func (s *S3Manager) ErrIsNotExist(err error) bool {
	if err == nil {
		return false
	}
	switch v := errors.Cause(err).(type) {
	case minio.ErrorResponse:
		return v.StatusCode == http.StatusNotFound || v.Code == s3.ErrCodeNoSuchKey
	default:
		rawErr := v.(error)
		if strings.Contains(rawErr.Error(), ` key does not exist`) {
			return true
		}
	}
	return false
}

func (s *S3Manager) Download(ctx echo.Context, ppath string) error {
	f, err := s.Get(ppath)
	if err != nil {
		return err
	}
	defer f.Close()
	fileName := path.Base(ppath)
	inline := ctx.Formx(`inline`).Bool()
	return ctx.Attachment(f, fileName, inline)
}

func (s *S3Manager) List(ctx echo.Context, ppath string, sortBy ...string) (err error, exit bool, dirs []os.FileInfo) {
	doneCh := make(chan struct{})
	defer close(doneCh)
	objectPrefix := strings.TrimPrefix(ppath, `/`)
	words := len(objectPrefix)
	var forceDir bool
	if words == 0 {
		forceDir = true
	} else {
		if strings.HasSuffix(objectPrefix, `/`) {
			forceDir = true
		} else {
			objectPrefix += `/`
		}
	}
	objectCh := s.client.ListObjectsV2(s.bucketName, objectPrefix, false, doneCh)
	for object := range objectCh {
		if object.Err != nil {
			continue
		}
		if len(objectPrefix) > 0 {
			object.Key = strings.TrimPrefix(object.Key, objectPrefix)
		}
		if len(object.Key) == 0 {
			continue
		}
		obj := NewFileInfo(object)
		dirs = append(dirs, obj)
	}
	if !forceDir && len(dirs) == 0 {
		return s.Download(ctx, ppath), true, nil
	}
	if len(sortBy) > 0 {
		switch sortBy[0] {
		case `time`:
			sort.Sort(filemanager.SortByModTime(dirs))
		case `-time`:
			sort.Sort(filemanager.SortByModTimeDesc(dirs))
		case `name`:
		case `-name`:
			sort.Sort(filemanager.SortByNameDesc(dirs))
		case `type`:
			fallthrough
		default:
			sort.Sort(filemanager.SortByFileType(dirs))
		}
	} else {
		sort.Sort(filemanager.SortByFileType(dirs))
	}
	if ctx.Format() == "json" {
		dirList, fileList := s.ListTransfer(dirs)
		data := ctx.Data()
		data.SetData(echo.H{
			`dirList`:  dirList,
			`fileList`: fileList,
		})
		return ctx.JSON(data), true, nil
	}
	return
}

func (s *S3Manager) ListTransfer(dirs []os.FileInfo) (dirList []echo.H, fileList []echo.H) {
	dirList = []echo.H{}
	fileList = []echo.H{}
	for _, d := range dirs {
		item := echo.H{
			`name`:  d.Name(),
			`size`:  d.Size(),
			`mode`:  d.Mode().String(),
			`mtime`: d.ModTime().Format(`2006-01-02 15:04:05`),
			//`sys`:   d.Sys(),
		}
		if d.IsDir() {
			dirList = append(dirList, item)
			continue
		}
		fileList = append(fileList, item)
	}
	return
}
