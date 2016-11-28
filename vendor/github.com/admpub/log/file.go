// Copyright 2015 Qiang Xue. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package log

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/admpub/queueChan"
)

// FileTarget writes filtered log messages to a file.
// FileTarget supports file rotation by keeping certain number of backup log files.
type FileTarget struct {
	*Filter
	// the log file name. When Rotate is true, log file name will be suffixed
	// to differentiate different backup copies (e.g. app.log.1)
	FileName string
	// whether to enable file rotating at specific time interval or when maximum file size is reached.
	Rotate bool
	// how many log files should be kept when Rotate is true (the current log file is not included).
	// This field is ignored when Rotate is false.
	BackupCount int
	// maximum number of bytes allowed for a log file. Zero means no limit.
	// This field is ignored when Rotate is false.
	MaxBytes int64

	fd           *os.File
	currentBytes int64
	errWriter    io.Writer
	close        chan bool
	timeFormat   string
	openedFile   string
	scaned       bool
	filePrefix   string
	queue        queueChan.QueueChan
}

// NewFileTarget creates a FileTarget.
// The new FileTarget takes these default options:
// MaxLevel: LevelDebug, Rotate: true, BackupCount: 10, MaxBytes: 1 << 20
// You must specify the FileName field.
func NewFileTarget() *FileTarget {
	return &FileTarget{
		Filter:      &Filter{MaxLevel: LevelDebug},
		Rotate:      true,
		BackupCount: 10,
		MaxBytes:    1 << 20, // 1MB
		close:       make(chan bool, 0),
	}
}

// Open prepares FileTarget for processing log messages.
func (t *FileTarget) Open(errWriter io.Writer) (err error) {
	t.Filter.Init()
	if t.FileName == `` {
		return errors.New("FileTarget.FileName must be set")
	}
	t.timeFormat = ``
	t.openedFile = ``
	p := strings.Index(t.FileName, `{date:`)
	t.filePrefix = t.FileName
	if p > -1 {
		fileName := t.FileName[0:p]
		if p == 0 {
			fileName = "./"
		}
		t.filePrefix = fileName
		placeholder := t.FileName[p+6:]
		p2 := strings.Index(placeholder, `}`)
		var hs bool
		switch fileName[len(fileName)-1] {
		case '/', '\\':
			hs = true
		}
		if fileName, err = filepath.Abs(fileName); err != nil {
			return err
		}
		if p2 > -1 {
			t.timeFormat = placeholder[0:p2]
			fileSuffix := placeholder[p2+1:]
			switch filepath.Separator {
			case '/':
				t.timeFormat = strings.Replace(t.timeFormat, "\\", "/", -1)
				fileSuffix = strings.Replace(fileSuffix, "\\", "/", -1)
			case '\\':
				t.timeFormat = strings.Replace(t.timeFormat, "/", "\\", -1)
				fileSuffix = strings.Replace(fileSuffix, "/", "\\", -1)
			}
			if hs {
				fileName = filepath.Join(fileName, `%v`+fileSuffix)
			} else {
				fileName += `%v` + fileSuffix
			}

		}
		t.FileName = fileName
		if t.FileName == `` {
			return errors.New("FileTarget.FileName must be set")
		}
		if t.filePrefix, err = filepath.Abs(t.filePrefix); err != nil {
			return err
		}
	} else {
		if t.filePrefix, err = filepath.Abs(t.filePrefix); err != nil {
			return err
		}
		t.FileName = t.filePrefix
	}

	t.errWriter = errWriter

	if t.Rotate {
		if t.BackupCount < 0 {
			return errors.New("FileTarget.BackupCount must be no less than 0")
		}
		if t.MaxBytes <= 0 {
			return errors.New("FileTarget.MaxBytes must be no less than 0")
		}
		t.queue = queueChan.New(t.BackupCount)
		t.queue.Dynamic()
	}
	t.openedFile = t.fileName()
	t.createDir(t.openedFile)
	t.fd, err = os.OpenFile(t.openedFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0660)
	if err != nil {
		return fmt.Errorf("FileTarget was unable to create a log file: %v", err)
	}
	if t.Rotate {
		t.recordOldLogs()
	}
	return nil
}

type logFiles []*logFileInfo

type logFileInfo struct {
	Path  string
	MTime int64
}

func (l *logFiles) Add(fpath string, mtime int64) {
	*l = append(*l, &logFileInfo{
		Path:  fpath,
		MTime: mtime,
	})
}

func (l *logFiles) Len() int {
	return len(*l)
}

func (l *logFiles) Less(i, j int) bool {
	return (*l)[i].MTime < (*l)[j].MTime
}

func (l *logFiles) Swap(i, j int) {
	(*l)[i], (*l)[j] = (*l)[j], (*l)[i]
}

func (t *FileTarget) recordOldLogs() {

	if t.scaned {
		return
	}

	t.scaned = true

	files := &logFiles{}
	err := filepath.Walk(filepath.Dir(t.filePrefix), func(f string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || f == t.filePrefix {
			return nil
		}

		if strings.HasPrefix(f, t.filePrefix) {
			files.Add(f, info.ModTime().UnixNano())
		}
		return nil
	})
	sort.Sort(files)
	for _, f := range *files {
		t.queue.PushTS(f.Path)
	}
	if err != nil {
		fmt.Fprintf(t.errWriter, "%v\n", err)
	} else if t.BackupCount > 0 {
		for t.queue.Length() > t.BackupCount {
			path, ok := t.queue.PopTS().(string)
			if !ok {
				continue
			}

			if err = os.Remove(path); err != nil {
				fmt.Fprintf(t.errWriter, "%v\n", err)
				break
			}
		}
	}
}

// Process saves an allowed log message into the log file.
func (t *FileTarget) Process(e *Entry) {
	if e == nil {
		t.fd.Close()
		t.close <- true
		return
	}
	if t.fd != nil && t.Allow(e) {
		if t.Rotate {
			t.rotate(int64(len(e.String()) + 1))
		}
		n, err := t.fd.Write([]byte(e.String() + "\n"))
		t.currentBytes += int64(n)
		if err != nil {
			fmt.Fprintf(t.errWriter, "FileTarge write error: %v\n", err)
		}
	}
}

// Close closes the file target.
func (t *FileTarget) Close() {
	<-t.close
	if t.fd != nil {
		t.fd.Close()
		t.fd = nil
	}
}

func (t *FileTarget) fileName() string {
	if t.timeFormat != `` {
		return fmt.Sprintf(t.FileName, time.Now().Format(t.timeFormat))
	}
	return t.FileName
}

func (t *FileTarget) rotate(bytes int64) {
	fileName := t.fileName()
	if t.openedFile == fileName && (t.currentBytes+bytes <= t.MaxBytes || bytes > t.MaxBytes) {
		return
	}
	t.fd.Close()
	t.currentBytes = 0
	var err error
	if t.BackupCount > 0 {
		for i := t.queue.Length() - t.BackupCount; i >= 0; i-- {
			path, ok := t.queue.PopTS().(string)
			if !ok {
				continue
			}
			if path == fileName {
				t.queue.PushTS(path)
				if t.queue.Length() > 1 {
					path, ok = t.queue.PopTS().(string)
					if !ok {
						continue
					}
				} else {
					break
				}
			}
			if err = os.Remove(path); err != nil {
				fmt.Fprintf(t.errWriter, "%v\n", err)
			}
		}
	}
	newPath := fileName
	if t.openedFile == fileName {
		newPath = fileName + `.` + time.Now().Format(`20060102150405`)
		err = os.Rename(t.openedFile, newPath)
		if err != nil {
			fmt.Fprintf(t.errWriter, "%v\n", err)
		}
	}
	t.queue.PushTS(newPath)
	/*
		for i := t.BackupCount; i >= 0; i-- {
			path := fileName
			if i > 0 {
				path = fmt.Sprintf("%v.%v", fileName, i)
			}
			if _, err = os.Lstat(path); err != nil {
				// file not exists
				continue
			}
			if i == t.BackupCount {
				os.Remove(path)
			} else {
				os.Rename(path, fmt.Sprintf("%v.%v", fileName, i+1))
			}
		}
	*/
	t.createDir(fileName)
	t.fd, err = os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0660)
	if err != nil {
		t.fd = nil
		fmt.Fprintf(t.errWriter, "FileTarget was unable to create a log file: %v\n", err)
	}
	t.openedFile = fileName
}

func (t *FileTarget) createDir(fileName string) {
	fdir := filepath.Dir(fileName)
	if finf, err := os.Stat(fdir); err != nil || !finf.IsDir() {
		err := os.MkdirAll(fdir, os.ModePerm)
		if err != nil {
			fmt.Fprintf(t.errWriter, "%v\n", err)
		}
	}
}
