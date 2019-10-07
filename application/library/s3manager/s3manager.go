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
	"io/ioutil"
	"os"
	"path"
	"sort"
	"strings"

	minio "github.com/minio/minio-go"
	"github.com/pkg/errors"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"

	"github.com/admpub/nging/application/library/charset"
	"github.com/admpub/nging/application/library/filemanager"
)

func New(client *minio.Client, bucketName string, editableMaxSize int64, ctx echo.Context) *s3Manager {
	return &s3Manager{
		Context:         ctx,
		client:          client,
		bucketName:      bucketName,
		EditableMaxSize: editableMaxSize,
	}
}

type s3Manager struct {
	echo.Context
	client          *minio.Client
	bucketName      string
	EditableMaxSize int64
}

func (s *s3Manager) Edit(ppath string, content string, encoding string) (interface{}, error) {
	objectName := strings.TrimPrefix(ppath, `/`)
	f, err := s.client.GetObject(s.bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if s.EditableMaxSize > 0 && fi.Size > s.EditableMaxSize {
		return nil, s.E(`很抱歉，不支持编辑超过%v的文件`, com.FormatBytes(s.EditableMaxSize))
	}
	encoding = strings.ToLower(encoding)
	isUTF8 := len(encoding) == 0 || encoding == `utf-8`
	if s.IsPost() {
		b := []byte(content)
		if !isUTF8 {
			b, err = charset.Convert(`utf-8`, encoding, b)
			if err != nil {
				return nil, err
			}
		}
		f.Close()
		r := bytes.NewReader(b)
		opts := minio.PutObjectOptions{ContentType: "application/octet-stream"}
		_, err = s.client.PutObject(s.bucketName, objectName, r, int64(len(b)), opts)
		if err != nil {
			return nil, s.E(ppath + `:` + err.Error())
		}
		return nil, err
	}

	dat, err := ioutil.ReadAll(f)
	if err == nil && !isUTF8 {
		dat, err = charset.Convert(encoding, `utf-8`, dat)
	}
	return string(dat), err
}

func (s *s3Manager) Mkbucket(bucketName string, regions ...string) error {
	var region string
	if len(regions) > 0 {
		region = regions[0]
	}
	return s.client.MakeBucket(bucketName, region)
}

func (s *s3Manager) Mkdir(ppath, newName string) error {
	objectName := strings.TrimPrefix(ppath, `/`)
	objectName = path.Join(objectName, newName)
	if !strings.HasSuffix(objectName, `/`) {
		objectName += `/`
	}
	_, err := s.client.PutObject(s.bucketName, objectName, nil, 0, minio.PutObjectOptions{})
	return err
}

func (s *s3Manager) Rename(ppath, newName string) error {
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

func (s *s3Manager) Chown(ppath string, uid, gid int) error {
	return nil
}

func (s *s3Manager) Chmod(ppath string, mode os.FileMode) error {
	return nil
}

func (s *s3Manager) Search(ppath string, prefix string, num int) []string {
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

func (s *s3Manager) Remove(ppath string) error {
	objectName := strings.TrimPrefix(ppath, `/`)
	return s.client.RemoveObject(s.bucketName, objectName)
}

func (s *s3Manager) Upload(ppath string) error {
	fileSrc, fileHdr, err := s.Request().FormFile(`file`)
	if err != nil {
		return err
	}
	defer fileSrc.Close()
	opts := minio.PutObjectOptions{ContentType: "application/octet-stream"}
	objectName := path.Join(ppath, fileHdr.Filename)
	objectName = strings.TrimPrefix(objectName, `/`)
	_, err = s.client.PutObject(s.bucketName, objectName, fileSrc, fileHdr.Size, opts)
	return err
}

func (s *s3Manager) Download(ppath string) error {
	objectName := strings.TrimPrefix(ppath, `/`)
	f, err := s.client.GetObject(s.bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return errors.WithMessage(err, objectName)
	}
	defer f.Close()
	fileName := path.Base(ppath)
	inline := s.Formx(`inline`).Bool()
	return s.Attachment(f, fileName, inline)
}

func (s *s3Manager) List(ppath string, sortBy ...string) (err error, exit bool, dirs []os.FileInfo) {
	doneCh := make(chan struct{})
	defer close(doneCh)
	objectPrefix := strings.TrimPrefix(ppath, `/`)
	forceDir := strings.HasSuffix(objectPrefix, "/")
	if !forceDir && len(objectPrefix) > 0 {
		objectPrefix += `/`
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
		return s.Download(ppath), true, nil
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
	if s.Format() == "json" {
		dirList, fileList := s.ListTransfer(dirs)
		data := s.Data()
		data.SetData(echo.H{
			`dirList`:  dirList,
			`fileList`: fileList,
		})
		return s.JSON(data), true, nil
	}
	return
}

func (s *s3Manager) ListTransfer(dirs []os.FileInfo) (dirList []echo.H, fileList []echo.H) {
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
