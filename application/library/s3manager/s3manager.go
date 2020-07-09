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
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/admpub/nging/application/dbschema"
	"github.com/admpub/nging/application/library/charset"
	"github.com/admpub/nging/application/library/common"
	"github.com/admpub/nging/application/library/filemanager"
	"github.com/admpub/nging/application/library/s3manager/s3client/awsclient"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	minio "github.com/minio/minio-go"
	"github.com/pkg/errors"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
)

func New(client *minio.Client, config *dbschema.NgingCloudStorage, editableMaxSize int64) *S3Manager {
	return &S3Manager{
		client:          client,
		config:          config,
		bucketName:      config.Bucket,
		EditableMaxSize: editableMaxSize,
	}
}

type S3Manager struct {
	client          *minio.Client
	config          *dbschema.NgingCloudStorage
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

func (s *S3Manager) renameDirectory(ppath, newName string) error {
	dirName := strings.TrimPrefix(ppath, `/`)
	newName = strings.TrimPrefix(newName, `/`)
	if !strings.HasSuffix(newName, `/`) {
		newName += `/`
	}
	// 新建文件夹
	_, err := s.client.PutObject(s.bucketName, newName, nil, 0, minio.PutObjectOptions{})
	if err != nil {
		return err
	}
	doneCh := make(chan struct{})
	defer close(doneCh)
	objectCh := s.client.ListObjectsV2(s.bucketName, dirName, true, doneCh)
	for object := range objectCh {
		if object.Err != nil {
			continue
		}
		if len(object.Key) == 0 || object.Key == dirName {
			continue
		}
		dest := strings.TrimPrefix(object.Key, dirName)
		dest = path.Join(newName, dest)
		// println(object.Key, ` => `, dest)
		err = s.Move(object.Key, dest)
		if err != nil {
			return err
		}
	}
	return s.client.RemoveObject(s.bucketName, dirName)
}

func (s *S3Manager) Rename(ppath, newName string) error {
	if strings.HasSuffix(ppath, `/`) {
		return s.renameDirectory(ppath, newName)
	}
	objectName := strings.TrimPrefix(ppath, `/`)
	newName = strings.TrimPrefix(newName, `/`)
	return s.Move(objectName, newName)
}

func (s *S3Manager) Move(from, to string) error {
	err := s.Copy(from, to)
	if err != nil {
		return err
	}
	return s.client.RemoveObject(s.bucketName, from)
}

func (s *S3Manager) Copy(from, to string) error {
	// Source object
	src := minio.NewSourceInfo(s.bucketName, from, nil)
	dst, err := minio.NewDestinationInfo(s.bucketName, to, nil, nil)
	if err != nil {
		return err
	}
	return s.client.CopyObject(dst, src)
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

	return s.client.RemoveObject(s.bucketName, objectName)
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
	f, err := s.Stat(ppath)
	return s.StatIsExists(f, err)
}

// StatIsExists 对象是否存在
func (s *S3Manager) StatIsExists(f minio.ObjectInfo, err error) (bool, error) {
	if err == nil {
		if f.Err == nil {
			return len(f.Key) > 0, nil
		}
		err = f.Err
	}
	//echo.Dump(echo.H{`info`:f,`err`:err,`exists`:!s.ErrIsNotExist(err)})
	if s.ErrIsNotExist(err) { // 已经确定是不存在的状态，不需要返回err
		return false, nil
	}
	return true, err // 不知道存在状态，返回原始err
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
		if os.IsNotExist(rawErr) {
			return true
		}
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

func (s *S3Manager) listByMinio(ctx echo.Context, objectPrefix string) (err error, dirs []os.FileInfo) {
	doneCh := make(chan struct{})
	defer close(doneCh)
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
	return
}

func (s *S3Manager) listByAWS(ctx echo.Context, objectPrefix string) (err error, dirs []os.FileInfo) {
	var s3client *s3.S3
	s3client, err = awsclient.Connect(s.config)
	if err != nil {
		return
	}
	_, limit, _, pagination := common.PagingWithPagination(ctx)
	if limit < 1 {
		limit = 20
	}
	offset := ctx.Formx(`curr-offset`, `0`).Uint()
	endIndex := offset + uint(limit)
	prevOffset := ctx.Form(`prev-offset`, `0`)
	nextOffset := fmt.Sprint(offset)
	q := ctx.Request().URL().Query()
	q.Del(`curr-offset`)
	q.Del(`prev-offset`)
	q.Del(`_pjax`)
	pagination.SetURL(ctx.Request().URL().Path() + `?` + q.Encode() + `&curr-offset={curr}&prev-offset={prev}`)
	pagination.SetPosition(prevOffset, nextOffset, nextOffset)
	var seekNum uint
	err = s3client.ListObjectsPagesWithContext(ctx, &s3.ListObjectsInput{
		Bucket:    aws.String(s.bucketName),
		Prefix:    aws.String(objectPrefix),
		MaxKeys:   aws.Int64(int64(limit)),
		Delimiter: aws.String(`/`),
	}, func(p *s3.ListObjectsOutput, lastPage bool) bool {
		if seekNum < offset {
			return true
		}
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
			obj := NewStrFileInfo(*object.Prefix)
			dirs = append(dirs, obj)
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
			obj := NewS3FileInfo(object)
			dirs = append(dirs, obj)
		}
		seekNum += uint(len(dirs))
		return seekNum <= endIndex // continue paging
	})
	ctx.Set(`pagination`, pagination)
	return
}

func (s *S3Manager) List(ctx echo.Context, ppath string, sortBy ...string) (err error, exit bool, dirs []os.FileInfo) {
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
	engine := ctx.Form(`engine`, `minio`)
	switch engine {
	case `aws`:
		err, dirs = s.listByAWS(ctx, objectPrefix)
	default:
		err, dirs = s.listByMinio(ctx, objectPrefix)
	}
	if err != nil {
		return
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
