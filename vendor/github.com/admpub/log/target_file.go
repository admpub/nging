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
	"sync"
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
	mutex        sync.Mutex
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
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.Filter.Init()
	if len(t.FileName) == 0 {
		return errors.New("FileTarget.FileName must be set")
	}
	t.timeFormat = ``
	t.openedFile = ``
	t.filePrefix, t.timeFormat, t.FileName, err = DateFormatFilename(t.FileName)
	if err != nil {
		return
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
	t.openedFile = t.getFileName()
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

func (t *FileTarget) getFileName() string {
	if len(t.timeFormat) > 0 {
		return fmt.Sprintf(t.FileName, time.Now().Format(t.timeFormat))
	}
	return t.FileName
}

func (t *FileTarget) rotate(bytes int64) {
	fileName := t.getFileName()
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

func DateFormatFilename(dfile string) (prefix string, dateformat string, filename string, err error) {
	p := strings.Index(dfile, `{date:`)
	prefix = dfile
	if p > -1 {
		fileName := dfile[0:p]
		if p == 0 {
			fileName = "./"
		}
		prefix = fileName
		placeholder := dfile[p+6:]
		p2 := strings.Index(placeholder, `}`)
		var hs bool
		switch fileName[len(fileName)-1] {
		case '/', '\\':
			hs = true
		}
		if fileName, err = filepath.Abs(fileName); err != nil {
			return "", "", "", err
		}
		if p2 > -1 {
			dateformat = placeholder[0:p2]
			fileSuffix := placeholder[p2+1:]
			switch filepath.Separator {
			case '/':
				dateformat = strings.Replace(dateformat, "\\", "/", -1)
				fileSuffix = strings.Replace(fileSuffix, "\\", "/", -1)
			case '\\':
				dateformat = strings.Replace(dateformat, "/", "\\", -1)
				fileSuffix = strings.Replace(fileSuffix, "/", "\\", -1)
			}
			if hs {
				fileName = filepath.Join(fileName, `%v`+fileSuffix)
			} else {
				fileName += `%v` + fileSuffix
			}
		}
		filename = fileName
		if len(filename) == 0 {
			err = errors.New("FileTarget.FileName must be set")
			return
		}
		if prefix, err = filepath.Abs(prefix); err != nil {
			return "", "", "", err
		}
	} else {
		if prefix, err = filepath.Abs(prefix); err != nil {
			return "", "", "", err
		}
		filename = prefix
	}
	return
}
