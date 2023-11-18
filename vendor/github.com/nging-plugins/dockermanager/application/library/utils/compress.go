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

package utils

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/admpub/archiver/v4"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
)

func GetFilenames(srcPath string) (filenames map[string]string, err error) {
	srcPath, err = filepath.Abs(srcPath)
	if err != nil {
		return
	}
	var fi fs.FileInfo
	fi, err = os.Stat(srcPath)
	if err != nil {
		return
	}
	filenames = map[string]string{}
	if !fi.IsDir() {
		filenames[srcPath] = ``
		return
	}
	if !strings.HasSuffix(srcPath, echo.FilePathSeparator) {
		srcPath += echo.FilePathSeparator
	}
	err = filepath.Walk(srcPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		relPath := strings.TrimPrefix(path, srcPath)
		if len(relPath) == 0 {
			return nil
		}
		filenames[path] = relPath
		return nil
	})
	return
}

func CompressTar(ctx context.Context, filenames map[string]string) (*PipeRW, error) {
	files, err := archiver.FilesFromDisk(&archiver.FromDiskOptions{}, filenames)
	if err != nil {
		return nil, err
	}
	pipeRW := NewPipe()
	pipeRW.DoWrite(func(w io.Writer) error {
		tar := archiver.Tar{}
		return tar.Archive(ctx, w, files)
	})
	// defer pipeRW.Close()
	return pipeRW, err
}

func CompressTarAsync(ctx context.Context, files <-chan archiver.File) *PipeRW {
	pipeRW := NewPipe()
	pipeRW.DoWrite(func(w io.Writer) error {
		tar := archiver.Tar{}
		return tar.ArchiveAsync(ctx, w, files)
	})
	// defer pipeRW.Close()
	return pipeRW
}

func DecompressTar(ctx context.Context, reader io.Reader, dstPath string) error {
	tar := archiver.Tar{}
	err := com.MkdirAll(dstPath, os.ModePerm)
	if err != nil {
		return err
	}
	err = tar.Extract(ctx, reader, nil, func(ctx context.Context, f archiver.File) (err error) {
		saveFile := filepath.Join(dstPath, f.NameInArchive)
		if f.IsDir() {
			fmt.Printf("[DecompressTar] dir: %s\n", f.NameInArchive)
			err = com.MkdirAll(saveFile, f.Mode())
		} else {
			fmt.Printf("[DecompressTar] file: %s\n", f.NameInArchive)
			var fp *os.File
			fp, err = os.Create(saveFile)
			if err != nil {
				if !os.IsNotExist(err) {
					return
				}
				err = com.MkdirAll(filepath.Dir(saveFile), f.Mode())
				if err != nil {
					return
				}
				fp, err = os.Create(saveFile)
			}
			defer fp.Close()
			var rc io.ReadCloser
			rc, err = f.Open()
			if err != nil {
				return
			}
			defer rc.Close()
			_, err = io.Copy(fp, rc)
			if err != nil {
				return
			}
			err = fp.Sync()
			if err != nil {
				return
			}
			err = fp.Chmod(f.Mode())
		}
		return
	})
	return err
}
