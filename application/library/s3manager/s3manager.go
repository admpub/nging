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
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/admpub/nging/v4/application/dbschema"
	"github.com/admpub/nging/v4/application/library/charset"
	"github.com/admpub/nging/v4/application/library/common"
	"github.com/admpub/nging/v4/application/library/filemanager"
	"github.com/admpub/nging/v4/application/library/s3manager/s3client/awsclient"

	"github.com/admpub/errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	minio "github.com/minio/minio-go/v7"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"

	uploadClient "github.com/webx-top/client/upload"
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
	f, err := s.Get(ctx, ppath)
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
		err = s.Put(ctx, r, ppath, int64(len(b)))
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

func (s *S3Manager) Mkbucket(ctx context.Context, bucketName string, regions ...string) error {
	var region string
	if len(regions) > 0 {
		region = regions[0]
	}
	return s.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: region})
}

func (s *S3Manager) Mkdir(ctx context.Context, ppath, newName string) error {
	objectName := strings.TrimPrefix(ppath, `/`)
	if len(newName) > 0 {
		objectName = path.Join(objectName, newName)
	}
	if !strings.HasSuffix(objectName, `/`) {
		objectName += `/`
	}
	_, err := s.client.PutObject(ctx, s.bucketName, objectName, nil, 0, minio.PutObjectOptions{})
	return err
}

func (s *S3Manager) renameDirectory(ctx context.Context, ppath, newName string) error {
	dirName := strings.TrimPrefix(ppath, `/`)
	newName = strings.TrimPrefix(newName, `/`)
	if !strings.HasSuffix(newName, `/`) {
		newName += `/`
	}
	// 新建文件夹
	_, err := s.client.PutObject(ctx, s.bucketName, newName, nil, 0, minio.PutObjectOptions{})
	if err != nil {
		return err
	}
	objectCh := s.client.ListObjects(ctx, s.bucketName, minio.ListObjectsOptions{Prefix: dirName, Recursive: true})
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
		err = s.Move(ctx, object.Key, dest)
		if err != nil {
			return err
		}
	}
	return s.client.RemoveObject(ctx, s.bucketName, dirName, minio.RemoveObjectOptions{ForceDelete: true})
}

func (s *S3Manager) Rename(ctx context.Context, ppath, newName string) error {
	if strings.HasSuffix(ppath, `/`) {
		return s.renameDirectory(ctx, ppath, newName)
	}
	objectName := strings.TrimPrefix(ppath, `/`)
	newName = strings.TrimPrefix(newName, `/`)
	return s.Move(ctx, objectName, newName)
}

func (s *S3Manager) Move(ctx context.Context, from, to string) error {
	err := s.Copy(ctx, from, to)
	if err != nil {
		return err
	}
	return s.client.RemoveObject(ctx, s.bucketName, from, minio.RemoveObjectOptions{ForceDelete: true})
}

func (s *S3Manager) Copy(ctx context.Context, from, to string) error {
	// Source object
	src := minio.CopySrcOptions{Bucket: s.bucketName, Object: from}
	dst := minio.CopyDestOptions{Bucket: s.bucketName, Object: to}
	_, err := s.client.CopyObject(ctx, dst, src)
	return err
}

func (s *S3Manager) Chown(ppath string, uid, gid int) error {
	return nil
}

func (s *S3Manager) Chmod(ppath string, mode os.FileMode) error {
	return nil
}

func (s *S3Manager) Search(ctx context.Context, ppath string, prefix string, num int) []string {
	var paths []string
	objectPrefix := path.Join(ppath, prefix)
	objectPrefix = strings.TrimPrefix(objectPrefix, `/`)
	objectCh := s.client.ListObjects(ctx, s.bucketName, minio.ListObjectsOptions{Prefix: objectPrefix})
	for object := range objectCh {
		if object.Err != nil {
			continue
		}
		paths = append(paths, object.Key)
	}
	return paths
}

func (s *S3Manager) Remove(ctx context.Context, ppath string) error {
	if len(ppath) == 0 {
		return errors.New("path invalid")
	}
	if strings.HasSuffix(ppath, `/`) {
		return s.RemoveDir(ctx, ppath)
	}
	objectName := strings.TrimPrefix(ppath, `/`)
	return s.client.RemoveObject(ctx, s.bucketName, objectName, minio.RemoveObjectOptions{ForceDelete: true})
}

func (s *S3Manager) RemoveDir(ctx context.Context, ppath string) error {
	objectName := strings.TrimPrefix(ppath, `/`)
	if !strings.HasSuffix(objectName, `/`) {
		objectName += `/`
	}
	if objectName == `/` {
		return s.Clear(ctx)
	}
	objectCh := s.client.ListObjects(ctx, s.bucketName, minio.ListObjectsOptions{Prefix: objectName, Recursive: true})
	for object := range objectCh {
		if object.Err != nil {
			continue
		}
		if len(object.Key) == 0 {
			continue
		}
		err := s.client.RemoveObject(ctx, s.bucketName, object.Key, minio.RemoveObjectOptions{ForceDelete: true})
		if err != nil {
			return err
		}
	}

	return s.client.RemoveObject(ctx, s.bucketName, objectName, minio.RemoveObjectOptions{ForceDelete: true})
}

// Clear 清空所有数据【慎用】
func (s *S3Manager) Clear(ctx context.Context) error {
	deleted := make(chan minio.ObjectInfo)
	defer close(deleted)
	removeObjects := s.client.RemoveObjects(ctx, s.bucketName, deleted, minio.RemoveObjectsOptions{})
	for removeObject := range removeObjects {
		if removeObject.Err != nil {
			return removeObject.Err
		}
	}
	return nil
}

func (s *S3Manager) Upload(ctx echo.Context, ppath string,
	chunkUpload *uploadClient.ChunkUpload,
	chunkOpts ...uploadClient.ChunkInfoOpter) error {
	var fileSrc io.Reader
	var objectName string
	var objectSize int64
	var chunked bool // 是否支持分片
	if chunkUpload != nil {
		_, err := chunkUpload.Upload(ctx.Request().StdRequest(), chunkOpts...)
		if err != nil {
			if !errors.Is(err, uploadClient.ErrChunkUnsupported) {
				if errors.Is(err, uploadClient.ErrChunkUploadCompleted) ||
					errors.Is(err, uploadClient.ErrFileUploadCompleted) {
					return nil
				}
				return err
			}
		} else {
			if !chunkUpload.Merged() {
				return nil
			}
			_fp, err := os.Open(chunkUpload.GetSavePath())
			if err != nil {
				return err
			}
			fileSrc = _fp
			defer func() {
				_fp.Close()
				os.Remove(chunkUpload.GetSavePath())
			}()
			chunked = true
			objectName = path.Join(ppath, filepath.Base(chunkUpload.GetSavePath()))
			objectSize = chunkUpload.GetSaveSize()
		}
	}
	if !chunked {
		_fileSrc, _fileHdr, err := ctx.Request().FormFile(`file`)
		if err != nil {
			return err
		}
		defer _fileSrc.Close()
		fileSrc = _fileSrc
		objectName = path.Join(ppath, _fileHdr.Filename)
		objectSize = _fileHdr.Size
	}
	//return s.uploadByAWS(ctx,fileSrc, objectName)
	return s.Put(ctx, fileSrc, objectName, objectSize)
}

// Put 提交数据
func (s *S3Manager) Put(ctx context.Context, reader io.Reader, objectName string, size int64) (err error) {
	_, err = s.PutObject(ctx, reader, objectName, size)
	return
}

func (s *S3Manager) PutObject(ctx context.Context, reader io.Reader, objectName string, size int64) (int64, error) {
	opts := minio.PutObjectOptions{ContentType: "application/octet-stream"}
	objectName = strings.TrimPrefix(objectName, `/`)
	info, err := s.client.PutObject(ctx, s.bucketName, objectName, reader, size, opts)
	return info.Size, err
}

func (s *S3Manager) FPutObject(ctx context.Context, filePath string, objectName string) (int64, error) {
	opts := minio.PutObjectOptions{ContentType: "application/octet-stream"}
	objectName = strings.TrimPrefix(objectName, `/`)
	info, err := s.client.FPutObject(ctx, s.bucketName, objectName, filePath, opts)
	return info.Size, err
}

// Get 获取数据
func (s *S3Manager) Get(ctx context.Context, ppath string) (*minio.Object, error) {
	objectName := strings.TrimPrefix(ppath, `/`)
	f, err := s.client.GetObject(ctx, s.bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return f, errors.WithMessage(err, objectName)
	}
	return f, err
}

// Stat 获取对象信息
func (s *S3Manager) Stat(ctx context.Context, ppath string) (minio.ObjectInfo, error) {
	objectName := strings.TrimPrefix(ppath, `/`)
	f, err := s.client.StatObject(ctx, s.bucketName, objectName, minio.StatObjectOptions{})
	if err != nil {
		return f, errors.WithMessage(err, objectName)
	}
	return f, err
}

// Exists 对象是否存在
func (s *S3Manager) Exists(ctx context.Context, ppath string) (bool, error) {
	return s.StatIsExists(s.Stat(ctx, ppath))
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
	if err != nil && s.ErrIsNotExist(err) { // 已经确定是不存在的状态，不需要返回err
		return false, nil
	}
	return true, err // 不知道存在状态，返回原始err
}

func (s *S3Manager) ErrIsNotExist(err error) bool {
	if err == nil {
		return false
	}
	if rawErr, ok := err.(minio.ErrorResponse); ok {
		return rawErr.StatusCode == http.StatusNotFound || rawErr.Code == s3.ErrCodeNoSuchKey
	}
	if os.IsNotExist(err) {
		return true
	}
	switch v := errors.Unwrap(err).(type) {
	case minio.ErrorResponse:
		return v.StatusCode == http.StatusNotFound || v.Code == s3.ErrCodeNoSuchKey
	case nil:
		if strings.Contains(err.Error(), ` key does not exist`) {
			return true
		}
	default:
		if strings.Contains(v.Error(), ` key does not exist`) {
			return true
		}
	}
	return false
}

func (s *S3Manager) Download(ctx echo.Context, ppath string) error {
	f, err := s.Get(ctx, ppath)
	if err != nil {
		return err
	}
	defer f.Close()
	fileName := path.Base(ppath)
	inline := ctx.Formx(`inline`).Bool()
	return ctx.Attachment(f, fileName, inline)
}

func (s *S3Manager) listByMinio(ctx context.Context, objectPrefix string) (dirs []os.FileInfo, err error) {
	objectCh := s.client.ListObjects(ctx, s.bucketName, minio.ListObjectsOptions{Prefix: objectPrefix})
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

func (s *S3Manager) uploadByAWS(ctx context.Context, reader io.Reader, objectName string) error {
	// Set up a new s3manager client
	sess, err := awsclient.NewSession(s.config)
	if err != nil {
		return err
	}
	uploader := s3manager.NewUploader(sess)
	result, err := uploader.Upload(&s3manager.UploadInput{
		Body:   reader,
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(objectName),
	})
	//com.Dump(result)
	_ = result
	return err
}

func (s *S3Manager) listByAWS(ctx echo.Context, objectPrefix string) (dirs []os.FileInfo, err error) {
	var s3client *s3.S3
	s3client, err = awsclient.Connect(s.config)
	if err != nil {
		return
	}
	_, limit, _, pagination := common.PagingWithPagination(ctx)
	if limit < 1 {
		limit = 20
	}
	offset := ctx.Form(`offset`)
	prevOffset := ctx.Form(`prev`)
	var nextOffset string
	q := ctx.Request().URL().Query()
	q.Del(`offset`)
	q.Del(`prev`)
	q.Del(`_pjax`)
	pagination.SetURL(ctx.Request().URL().Path() + `?` + q.Encode() + `&offset={curr}&prev={prev}`)
	input := &s3.ListObjectsInput{
		Bucket:    aws.String(s.bucketName),
		Prefix:    aws.String(objectPrefix),
		MaxKeys:   aws.Int64(int64(limit)),
		Delimiter: aws.String(`/`),
		Marker:    aws.String(offset),
	}
	var n int
	err = s3client.ListObjectsPagesWithContext(ctx, input, func(p *s3.ListObjectsOutput, lastPage bool) bool {
		if p.NextMarker != nil {
			nextOffset = *p.NextMarker
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
		n += len(dirs)
		return n <= limit // continue paging
	})
	pagination.SetPosition(prevOffset, nextOffset, offset)
	ctx.Set(`pagination`, pagination)
	return
}

func (s *S3Manager) List(ctx echo.Context, ppath string, sortBy ...string) (dirs []os.FileInfo, exit bool, err error) {
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
	engine := ctx.Form(`engine`, `aws`)
	switch engine {
	case `aws`:
		dirs, err = s.listByAWS(ctx, objectPrefix)
	default:
		dirs, err = s.listByMinio(ctx, objectPrefix)
	}
	if err != nil {
		return
	}
	if !forceDir && len(dirs) == 0 {
		return nil, true, s.Download(ctx, ppath)
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
		return nil, true, ctx.JSON(data)
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
