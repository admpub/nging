package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/admpub/archiver"
	"github.com/admpub/log"
	"github.com/webx-top/com"

	dcommon "github.com/nging-plugins/dbmanager/application/library/common"
)

type ImportFile struct {
	delDirs     []string
	StructFiles []string
	DataFiles   []string
}

func (a *ImportFile) Close() error {
	for _, delDir := range a.delDirs {
		os.RemoveAll(delDir)
	}
	for _, sqlFile := range a.StructFiles {
		if !com.FileExists(sqlFile) {
			continue
		}
		os.Remove(sqlFile)
	}
	return nil
}

func (a *ImportFile) AllSqlFiles() []string {
	files := make([]string, 0, len(a.StructFiles)+len(a.DataFiles))
	files = append(files, a.StructFiles...)
	files = append(files, a.DataFiles...)
	return files
}

func (a *ImportFile) AllSqlFileNum() int {
	return len(a.StructFiles) + len(a.DataFiles)
}

func ParseImportFile(cacheDir string, files []string) (*ImportFile, error) {
	var (
		delDirs        []string
		sqlStructFiles []string
		sqlDataFiles   []string
	)
	nowTime := com.String(time.Now().Unix())
	for index, sqlFile := range files {
		switch strings.ToLower(filepath.Ext(sqlFile)) {
		case `.sql`:
			if strings.Contains(filepath.Base(sqlFile), `struct`) {
				sqlStructFiles = append(sqlStructFiles, sqlFile)
			} else {
				sqlDataFiles = append(sqlDataFiles, sqlFile)
			}
		case `.zip`:
			dir := filepath.Join(cacheDir, fmt.Sprintf("upload-"+nowTime+"-%d", index))
			err := com.MkdirAll(dir, os.ModePerm)
			if err != nil {
				return nil, err
			}
			err = archiver.Unarchive(sqlFile, dir)
			if err != nil {
				log.Error(err)
				continue
			}
			delDirs = append(delDirs, dir)
			err = os.Remove(sqlFile)
			if err != nil {
				log.Error(err)
			}
			err = filepath.Walk(dir, func(fpath string, info os.FileInfo, err error) error {
				if err != nil || info.IsDir() {
					return err
				}
				if strings.ToLower(filepath.Ext(fpath)) != `.sql` {
					return nil
				}
				if strings.Contains(info.Name(), `struct`) {
					sqlStructFiles = append(sqlStructFiles, fpath)
					return nil
				}
				sqlDataFiles = append(sqlDataFiles, fpath)
				return nil
			})
			if err != nil {
				return nil, err
			}
		}
	}
	dcommon.SortStrings(sqlStructFiles)
	dcommon.SortStrings(sqlDataFiles)
	return &ImportFile{
		delDirs:     delDirs,
		StructFiles: sqlStructFiles,
		DataFiles:   sqlDataFiles,
	}, nil
}
