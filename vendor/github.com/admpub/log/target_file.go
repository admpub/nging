// Copyright 2015 Qiang Xue. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package log

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"
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
	fileSuffix   string
	mutex        sync.Mutex
	logFiles     logFiles
}

// NewFileTarget creates a FileTarget.
// The new FileTarget takes these default options:
// MaxLevel: LevelDebug, Rotate: true, BackupCount: 10, MaxBytes: 1 << 20
// You must specify the FileName field.
func NewFileTarget() *FileTarget {
	return &FileTarget{
		Filter:      &Filter{MaxLevel: LevelDebug},
		Rotate:      true,
		BackupCount: DefaultFileBackupCount,
		MaxBytes:    DefaultFileMaxBytes,
		close:       make(chan bool),
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
	t.filePrefix, t.fileSuffix, t.timeFormat, t.FileName, err = DateFormatFilename(t.FileName)
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
		t.recordOldLogs()
	}
	return nil
}

type logFiles []*logFileInfo

type logFileInfo struct {
	Path  string
	MTime int64
}

func (l *logFiles) Add(fpath string, mtime time.Time) {
	*l = append(*l, &logFileInfo{
		Path:  fpath,
		MTime: mtime.UnixNano(),
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

func Dump(v interface{}, returnOnly ...bool) string {
	b, _ := json.MarshalIndent(v, ``, `  `)
	if len(returnOnly) > 0 && returnOnly[0] {
		return string(b)
	}
	fmt.Println(string(b))
	return ``
}

var fileTimeSuffix = regexp.MustCompile(`\.[0-9]{14,}\.[0-9]{5}$`)

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
		if info.IsDir() {
			return nil
		}
		if f == t.filePrefix || strings.HasPrefix(info.Name(), `.`) {
			return nil
		}
		if len(t.fileSuffix) > 0 && !strings.HasSuffix(fileTimeSuffix.ReplaceAllString(f, ``), t.fileSuffix) {
			return nil
		}

		if strings.HasPrefix(f, t.filePrefix) {
			files.Add(f, info.ModTime())
		}
		return nil
	})
	sort.Sort(files)
	t.logFiles = *files
	if err != nil {
		fmt.Fprintf(t.errWriter, "%v\n", err)
	} else if t.BackupCount > 0 {
		for len(t.logFiles) > t.BackupCount {
			if err = os.Remove(t.popFile().Path); err != nil {
				fmt.Fprintf(t.errWriter, "%v\n", err)
				break
			}
		}
	}
}

func (t *FileTarget) popFile() *logFileInfo {
	pathInfo := t.logFiles[0]
	if len(t.logFiles) > 1 {
		t.logFiles = t.logFiles[1:]
	} else {
		t.logFiles = t.logFiles[0:0]
	}
	return pathInfo
}

func (t *FileTarget) CountFiles() int {
	return t.logFiles.Len()
}

func (t *FileTarget) ClearFiles() {
	var old logFiles
	for _, file := range t.logFiles {
		if err := os.Remove(file.Path); err != nil {
			fmt.Fprintf(t.errWriter, "%v\n", err)
			old = append(old, file)
		}
	}
	if old.Len() > 0 {
		t.logFiles = old
	} else {
		t.logFiles = t.logFiles[0:0]
	}
}

// Process saves an allowed log message into the log file.
func (t *FileTarget) Process(e *Entry) {
	if e == nil {
		t.closeFile()
		t.close <- true
		return
	}
	if !t.Allow(e) {
		return
	}
	_, err := t.Write([]byte(e.String() + "\n"))
	if err != nil {
		fmt.Fprintf(t.errWriter, "FileTarge write error: %v\n", err)
	}
	if t.Rotate {
		t.rotate()
	}
}

func (t *FileTarget) Write(b []byte) (int, error) {
	if t.fd == nil {
		if err := t.createLogFile(t.getFileName(), true); err != nil {
			return 0, err
		}
	}
	n, err := t.fd.Write(b)
	t.mutex.Lock()
	t.currentBytes += int64(n)
	t.mutex.Unlock()
	return n, err
}

// Close closes the file target.
func (t *FileTarget) Close() {
	<-t.close
	t.closeFile()
}

func (t *FileTarget) getFileName() string {
	if len(t.timeFormat) > 0 {
		return fmt.Sprintf(t.FileName, time.Now().Format(t.timeFormat))
	}
	return t.FileName
}

func (t *FileTarget) rotate() {
	fileName := t.getFileName()
	if t.openedFile == fileName && t.currentBytes <= t.MaxBytes {
		return
	}
	var err error
	if t.BackupCount > 0 {
		for i := len(t.logFiles) - t.BackupCount; i >= 0; i-- {
			pathInfo := t.popFile()
			if pathInfo.Path == fileName {
				t.logFiles = append(t.logFiles, pathInfo)
				if len(t.logFiles) < 2 {
					break
				}
				pathInfo = t.popFile()
			}
			if err = os.Remove(pathInfo.Path); err != nil {
				fmt.Fprintf(t.errWriter, "%v\n", err)
			}
		}
	}
	now := time.Now().Local()
	newPath := fileName
	if t.openedFile == fileName { // 文件名没变但尺寸超过设定值
		newPath = fileName + `.` + now.Format(`20060102150405.00000`)
		err = os.Rename(t.openedFile, newPath)
		if err != nil {
			fmt.Fprintf(t.errWriter, "%v\n", err)
		}
	}
	//println(`newPath:`, newPath)
	t.logFiles.Add(newPath, now)
	t.createLogFile(fileName)
}

func (t *FileTarget) closeFile() {
	if t.fd != nil {
		t.fd.Close()
		t.fd = nil
	}
}

func (t *FileTarget) createLogFile(fileName string, recordFile ...bool) (err error) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.closeFile()
	t.currentBytes = 0
	t.createDir(fileName)
	t.fd, err = os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		t.fd = nil
		fmt.Fprintf(t.errWriter, "FileTarget was unable to create a log file: %v\n", err)
	}
	t.openedFile = fileName
	if len(recordFile) > 0 && recordFile[0] {
		t.logFiles.Add(fileName, time.Now().Local())
	}
	return
}

func (t *FileTarget) createDir(fileName string) {
	fdir := filepath.Dir(fileName)
	if finf, err := os.Stat(fdir); err != nil || !finf.IsDir() {
		err := os.MkdirAll(fdir, os.ModePerm)
		if err != nil {
			fmt.Fprintf(t.errWriter, "%v\n", err)
		}
		os.Chmod(fdir, os.ModePerm)
	}
}

func DateFormatFilename(dfile string) (prefix string, suffix, dateformat string, filename string, err error) {
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
			return
		}
		if p2 > -1 {
			dateformat = placeholder[0:p2]
			suffix = placeholder[p2+1:]
			switch filepath.Separator {
			case '/':
				dateformat = strings.Replace(dateformat, "\\", "/", -1)
				suffix = strings.Replace(suffix, "\\", "/", -1)
			case '\\':
				dateformat = strings.Replace(dateformat, "/", "\\", -1)
				suffix = strings.Replace(suffix, "/", "\\", -1)
			}
			if hs {
				fileName = filepath.Join(fileName, `%v`+suffix)
			} else {
				fileName += `%v` + suffix
			}
		}
		filename = fileName
		if len(filename) == 0 {
			err = errors.New("FileTarget.FileName must be set")
			return
		}
		if prefix, err = filepath.Abs(prefix); err != nil {
			return
		}
	} else {
		if prefix, err = filepath.Abs(prefix); err != nil {
			return
		}
		filename = prefix
	}
	return
}
