package cloudbackup

import (
	"context"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/admpub/log"
	"github.com/admpub/nging/v5/application/library/s3manager"
)

func New(mgr *s3manager.S3Manager) *Cloudbackup {
	return &Cloudbackup{mgr: mgr}
}

type Cloudbackup struct {
	mgr        *s3manager.S3Manager
	SourcePath string
	DestPath   string
	Filter     func(string) bool

	WaitFillCompleted bool
	IgnoreWaitRegexp  *regexp.Regexp
}

func (c *Cloudbackup) OnCreate(file string) {
	fp, err := os.Open(file)
	if err != nil {
		log.Error(file + `: ` + err.Error())
		return
	}
	fi, err := fp.Stat()
	if err != nil {
		log.Error(file + `: ` + err.Error())
		return
	}
	if fi.IsDir() {
		fp.Close()
		err = filepath.Walk(file, func(ppath string, info os.FileInfo, werr error) error {
			if werr != nil {
				return werr
			}
			if info.IsDir() || !c.Filter(ppath) {
				return nil
			}
			_waitFillCompleted := c.WaitFillCompleted
			if _waitFillCompleted && c.IgnoreWaitRegexp != nil {
				_waitFillCompleted = c.IgnoreWaitRegexp.MatchString(ppath)
			}
			objectName := path.Join(c.DestPath, strings.TrimPrefix(ppath, c.SourcePath))
			FileChan() <- &PutFile{
				Manager:           c.mgr,
				ObjectName:        objectName,
				FilePath:          ppath,
				WaitFillCompleted: _waitFillCompleted,
			}
			return nil
		})
	} else {
		fp.Close()
		_waitFillCompleted := c.WaitFillCompleted
		if _waitFillCompleted && c.IgnoreWaitRegexp != nil {
			_waitFillCompleted = c.IgnoreWaitRegexp.MatchString(file)
		}
		objectName := path.Join(c.DestPath, strings.TrimPrefix(file, c.SourcePath))
		FileChan() <- &PutFile{
			Manager:           c.mgr,
			ObjectName:        objectName,
			FilePath:          file,
			WaitFillCompleted: _waitFillCompleted,
		}
	}
	if err != nil {
		log.Error(err)
	}
}

func (c *Cloudbackup) OnModify(file string) {
	objectName := path.Join(c.DestPath, strings.TrimPrefix(file, c.SourcePath))
	fp, err := os.Open(file)
	if err != nil {
		log.Error(file + `: ` + err.Error())
		return
	}
	fi, err := fp.Stat()
	if err != nil {
		log.Error(file + `: ` + err.Error())
		fp.Close()
		return
	}
	if fi.IsDir() {
		fp.Close()
		return
	}
	fp.Close()
	_waitFillCompleted := c.WaitFillCompleted
	if _waitFillCompleted && c.IgnoreWaitRegexp != nil {
		_waitFillCompleted = c.IgnoreWaitRegexp.MatchString(file)
	}
	FileChan() <- &PutFile{
		Manager:           c.mgr,
		ObjectName:        objectName,
		FilePath:          file,
		WaitFillCompleted: _waitFillCompleted,
	}
}

func (c *Cloudbackup) OnDelete(file string) {
	objectName := path.Join(c.DestPath, strings.TrimPrefix(file, c.SourcePath))
	err := c.mgr.RemoveDir(context.Background(), objectName)
	if err != nil {
		log.Error(file + `: ` + err.Error())
	}
	err = c.mgr.Remove(context.Background(), objectName)
	if err != nil {
		log.Error(file + `: ` + err.Error())
	}
}

func (c *Cloudbackup) OnRename(file string) {
	objectName := path.Join(c.DestPath, strings.TrimPrefix(file, c.SourcePath))
	err := c.mgr.RemoveDir(context.Background(), objectName)
	if err != nil {
		log.Error(file + `: ` + err.Error())
	}
	err = c.mgr.Remove(context.Background(), objectName)
	if err != nil {
		log.Error(file + `: ` + err.Error())
	}
}
