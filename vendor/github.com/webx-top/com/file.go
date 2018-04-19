// Copyright 2013 com authors
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package com

import (
	"archive/tar"
	"archive/zip"
	"bufio"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Storage unit constants.
const (
	Byte  = 1
	KByte = Byte * 1024
	MByte = KByte * 1024
	GByte = MByte * 1024
	TByte = GByte * 1024
	PByte = TByte * 1024
	EByte = PByte * 1024
)

func logn(n, b float64) float64 {
	return math.Log(n) / math.Log(b)
}

func humanateBytes(s uint64, base float64, sizes []string) string {
	if s < 10 {
		return fmt.Sprintf("%dB", s)
	}
	e := math.Floor(logn(float64(s), base))
	suffix := sizes[int(e)]
	val := float64(s) / math.Pow(base, math.Floor(e))
	f := "%.0f"
	if val < 10 {
		f = "%.1f"
	}

	return fmt.Sprintf(f+"%s", val, suffix)
}

// HumaneFileSize calculates the file size and generate user-friendly string.
func HumaneFileSize(s uint64) string {
	sizes := []string{"B", "KB", "MB", "GB", "TB", "PB", "EB"}
	return humanateBytes(s, 1024, sizes)
}

// FileMTime returns file modified time and possible error.
func FileMTime(file string) (int64, error) {
	f, err := os.Stat(file)
	if err != nil {
		return 0, err
	}
	return f.ModTime().Unix(), nil
}

// FileSize returns file size in bytes and possible error.
func FileSize(file string) (int64, error) {
	f, err := os.Stat(file)
	if err != nil {
		return 0, err
	}
	return f.Size(), nil
}

// Copy copies file from source to target path.
func Copy(src, dest string) error {
	// Gather file information to set back later.
	si, err := os.Lstat(src)
	if err != nil {
		return err
	}

	// Handle symbolic link.
	if si.Mode()&os.ModeSymlink != 0 {
		target, err := os.Readlink(src)
		if err != nil {
			return err
		}
		// NOTE: os.Chmod and os.Chtimes don't recoganize symbolic link,
		// which will lead "no such file or directory" error.
		return os.Symlink(target, dest)
	}

	sr, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sr.Close()

	dw, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer dw.Close()

	if _, err = io.Copy(dw, sr); err != nil {
		return err
	}

	// Set back file information.
	if err = os.Chtimes(dest, si.ModTime(), si.ModTime()); err != nil {
		return err
	}
	return os.Chmod(dest, si.Mode())
}

// WriteFile writes data to a file named by filename.
// If the file does not exist, WriteFile creates it
// and its upper level paths.
func WriteFile(filename string, data []byte) error {
	os.MkdirAll(filepath.Dir(filename), os.ModePerm)
	return ioutil.WriteFile(filename, data, 0655)
}

// IsFile returns true if given path is a file,
// or returns false when it's a directory or does not exist.
func IsFile(filePath string) bool {
	f, e := os.Stat(filePath)
	if e != nil {
		return false
	}
	return !f.IsDir()
}

// IsExist checks whether a file or directory exists.
// It returns false when the file or directory does not exist.
func IsExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

func Unlink(file string) bool {
	if err := os.Remove(file); err == nil {
		return true
	}
	return false
}

// SaveFile saves content type '[]byte' to file by given path.
// It returns error when fail to finish operation.
func SaveFile(filePath string, b []byte) (int, error) {
	os.MkdirAll(filepath.Dir(filePath), os.ModePerm)
	fw, err := os.Create(filePath)
	if err != nil {
		return 0, err
	}
	defer fw.Close()
	return fw.Write(b)
}

// SaveFileS saves content type 'string' to file by given path.
// It returns error when fail to finish operation.
func SaveFileS(filePath string, s string) (int, error) {
	return SaveFile(filePath, []byte(s))
}

// ReadFile reads data type '[]byte' from file by given path.
// It returns error when fail to finish operation.
func ReadFile(filePath string) ([]byte, error) {
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		return []byte(""), err
	}
	return b, nil
}

// ReadFileS reads data type 'string' from file by given path.
// It returns error when fail to finish operation.
func ReadFileS(filePath string) (string, error) {
	b, err := ReadFile(filePath)
	return string(b), err
}

// Zip 压缩为zip
func Zip(srcDirPath string, destFilePath string, args ...*regexp.Regexp) (n int64, err error) {
	root, err := filepath.Abs(srcDirPath)
	if err != nil {
		return 0, err
	}

	f, err := os.Create(destFilePath)
	if err != nil {
		return
	}
	defer f.Close()

	w := zip.NewWriter(f)
	var regexpIgnoreFile, regexpFileName *regexp.Regexp
	argLen := len(args)
	if argLen > 1 {
		regexpIgnoreFile = args[1]
		regexpFileName = args[0]
	} else if argLen == 1 {
		regexpFileName = args[0]
	}
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		name := info.Name()
		nameBytes := []byte(name)
		if regexpIgnoreFile != nil && regexpIgnoreFile.Match(nameBytes) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		} else if info.IsDir() {
			return nil
		}
		if regexpFileName != nil && !regexpFileName.Match(nameBytes) {
			return nil
		}
		relativePath := strings.TrimPrefix(path, root)
		relativePath = strings.Replace(relativePath, `\`, `/`, -1)
		relativePath = strings.TrimPrefix(relativePath, `/`)
		f, err := w.Create(relativePath)
		if err != nil {
			return err
		}
		sf, err := os.Open(path)
		if err != nil {
			return err
		}
		defer sf.Close()
		_, err = io.Copy(f, sf)
		return err
	})

	err = w.Close()
	if err != nil {
		return 0, err
	}

	fi, err := f.Stat()
	if err != nil {
		n = fi.Size()
	}
	return
}

// Unzip unzips .zip file to 'destPath'.
// It returns error when fail to finish operation.
func Unzip(srcPath, destPath string) error {
	// Open a zip archive for reading
	r, err := zip.OpenReader(srcPath)
	if err != nil {
		return err
	}
	defer r.Close()

	// Iterate through the files in the archive
	for _, f := range r.File {
		// Get files from archive
		rc, err := f.Open()
		if err != nil {
			return err
		}

		dir := filepath.Dir(f.Name)
		// Create directory before create file
		os.MkdirAll(destPath+"/"+dir, os.ModePerm)

		if f.FileInfo().IsDir() {
			continue
		}

		// Write data to file
		fw, _ := os.Create(filepath.Join(destPath, f.Name))
		if err != nil {
			return err
		}
		_, err = io.Copy(fw, rc)

		if fw != nil {
			fw.Close()
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func TarGz(srcDirPath string, destFilePath string) error {
	fw, err := os.Create(destFilePath)
	if err != nil {
		return err
	}
	defer fw.Close()

	// Gzip writer
	gw := gzip.NewWriter(fw)
	defer gw.Close()

	// Tar writer
	tw := tar.NewWriter(gw)
	defer tw.Close()

	// Check if it's a file or a directory
	f, err := os.Open(srcDirPath)
	if err != nil {
		return err
	}
	fi, err := f.Stat()
	if err != nil {
		return err
	}
	if fi.IsDir() {
		// handle source directory
		fmt.Println("Cerating tar.gz from directory...")
		if err := tarGzDir(srcDirPath, filepath.Base(srcDirPath), tw); err != nil {
			return err
		}
	} else {
		// handle file directly
		fmt.Println("Cerating tar.gz from " + fi.Name() + "...")
		if err := tarGzFile(srcDirPath, fi.Name(), tw, fi); err != nil {
			return err
		}
	}
	fmt.Println("Well done!")
	return err
}

// Deal with directories
// if find files, handle them with tarGzFile
// Every recurrence append the base path to the recPath
// recPath is the path inside of tar.gz
func tarGzDir(srcDirPath string, recPath string, tw *tar.Writer) error {
	// Open source diretory
	dir, err := os.Open(srcDirPath)
	if err != nil {
		return err
	}
	defer dir.Close()

	// Get file info slice
	fis, err := dir.Readdir(0)
	if err != nil {
		return err
	}
	for _, fi := range fis {
		// Append path
		curPath := srcDirPath + "/" + fi.Name()
		// Check it is directory or file
		if fi.IsDir() {
			// Directory
			// (Directory won't add unitl all subfiles are added)
			fmt.Printf("Adding path...%s\n", curPath)
			tarGzDir(curPath, recPath+"/"+fi.Name(), tw)
		} else {
			// File
			fmt.Printf("Adding file...%s\n", curPath)
		}

		tarGzFile(curPath, recPath+"/"+fi.Name(), tw, fi)
	}
	return err
}

// Deal with files
func tarGzFile(srcFile string, recPath string, tw *tar.Writer, fi os.FileInfo) error {
	if fi.IsDir() {
		// Create tar header
		hdr := new(tar.Header)
		// if last character of header name is '/' it also can be directory
		// but if you don't set Typeflag, error will occur when you untargz
		hdr.Name = recPath + "/"
		hdr.Typeflag = tar.TypeDir
		hdr.Size = 0
		//hdr.Mode = 0755 | c_ISDIR
		hdr.Mode = int64(fi.Mode())
		hdr.ModTime = fi.ModTime()

		// Write hander
		err := tw.WriteHeader(hdr)
		if err != nil {
			return err
		}
	} else {
		// File reader
		fr, err := os.Open(srcFile)
		if err != nil {
			return err
		}
		defer fr.Close()

		// Create tar header
		hdr := new(tar.Header)
		hdr.Name = recPath
		hdr.Size = fi.Size()
		hdr.Mode = int64(fi.Mode())
		hdr.ModTime = fi.ModTime()

		// Write hander
		err = tw.WriteHeader(hdr)
		if err != nil {
			return err
		}

		// Write file data
		_, err = io.Copy(tw, fr)
		if err != nil {
			return err
		}
	}
	return nil
}

// UnTarGz ungzips and untars .tar.gz file to 'destPath' and returns sub-directories.
// It returns error when fail to finish operation.
func UnTarGz(srcFilePath string, destDirPath string) ([]string, error) {
	// Create destination directory
	os.Mkdir(destDirPath, os.ModePerm)

	fr, err := os.Open(srcFilePath)
	if err != nil {
		return nil, err
	}
	defer fr.Close()

	// Gzip reader
	gr, err := gzip.NewReader(fr)
	if err != nil {
		return nil, err
	}
	defer gr.Close()

	// Tar reader
	tr := tar.NewReader(gr)

	dirs := make([]string, 0, 5)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			// End of tar archive
			break
		}

		// Check if it is directory or file
		if hdr.Typeflag != tar.TypeDir {
			// Get files from archive
			// Create directory before create file
			dir := filepath.Dir(hdr.Name)
			os.MkdirAll(destDirPath+"/"+dir, os.ModePerm)
			dirs = AppendStr(dirs, dir)

			// Write data to file
			fw, _ := os.Create(destDirPath + "/" + hdr.Name)
			if err != nil {
				return nil, err
			}
			_, err = io.Copy(fw, tr)
			if err != nil {
				return nil, err
			}
		}
	}
	return dirs, nil
}

var (
	selfPath string
	selfDir  string
)

func SelfPath() string {
	if len(selfPath) == 0 {
		selfPath, _ = filepath.Abs(os.Args[0])
	}
	return selfPath
}

func SelfDir() string {
	if len(selfDir) == 0 {
		selfDir = filepath.Dir(SelfPath())
	}
	return selfDir
}

// FileExists reports whether the named file or directory exists.
func FileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// SearchFile search a file in paths.
// this is offen used in search config file in /etc ~/
func SearchFile(filename string, paths ...string) (fullpath string, err error) {
	for _, path := range paths {
		if fullpath = filepath.Join(path, filename); FileExists(fullpath) {
			return
		}
	}
	err = errors.New(fullpath + " not found in paths")
	return
}

// GrepFile like command grep -E
// for example: GrepFile(`^hello`, "hello.txt")
// \n is striped while read
func GrepFile(patten string, filename string) (lines []string, err error) {
	re, err := regexp.Compile(patten)
	if err != nil {
		return
	}

	lines = make([]string, 0)
	err = SeekFileLines(filename, func(line string) error {
		if re.MatchString(line) {
			lines = append(lines, line)
		}
		return nil
	})
	return lines, err
}

func SeekFileLines(filename string, callback func(string) error) (err error) {
	fd, err := os.Open(filename)
	if err != nil {
		return
	}
	defer fd.Close()
	reader := bufio.NewReader(fd)
	prefix := ""
	for {
		byteLine, isPrefix, er := reader.ReadLine()
		if er != nil && er != io.EOF {
			return er
		}
		line := string(byteLine)
		if isPrefix {
			prefix += line
			continue
		}
		line = prefix + line
		if e := callback(line); e != nil {
			return e
		}
		prefix = ""
		if er == io.EOF {
			break
		}
	}
	return nil
}

func Readbuf(r io.Reader, length int) ([]byte, error) {
	buf := make([]byte, length)
	offset := 0

	for offset < length {
		read, err := r.Read(buf[offset:])
		if err != nil {
			return buf, err
		}

		offset += read
	}

	return buf, nil
}
