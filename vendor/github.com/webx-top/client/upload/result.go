package upload

import (
	"errors"
	"fmt"
	"path/filepath"
	"time"
)

// Results 批量上传时的结果数据记录
type Results []*Result

// Checker 上传合法性检查
type Checker func(r *Result) error

var (
	// ErrExistsFile 文件已存在
	ErrExistsFile = errors.New("This file already exists")
	// ErrInvalidContent 无效的上传内容
	ErrInvalidContent = errors.New("Invalid upload content")
	// ErrFileTooLarge 上传的文件太大
	ErrFileTooLarge = errors.New("The uploaded file is too large and exceeds the size limit")
)

func (r Results) FileURLs() (rs []string) {
	rs = make([]string, len(r))
	for k, v := range r {
		rs[k] = v.FileURL
	}
	return rs
}

func (r *Results) Add(result *Result) {
	*r = append(*r, result)
}

// FileNameGenerator 文件名称生成函数
type FileNameGenerator func(string) (string, error)

// Result 上传结果数据记录
type Result struct {
	FileID            int64
	FileName          string
	FileURL           string
	FileType          FileType
	FileSize          int64
	SavePath          string
	Md5               string
	Addon             interface{}
	fileNameGenerator FileNameGenerator
}

var DefaultNameGenerator FileNameGenerator = func(fileName string) (string, error) {
	return filepath.Join(time.Now().Format("2006/0102"), fileName), nil
}

func (r *Result) SetFileNameGenerator(generator FileNameGenerator) *Result {
	r.fileNameGenerator = generator
	return r
}

func (r *Result) CopyFrom(data *Result) *Result {
	if r == data {
		return r
	}
	r.FileID = data.FileID
	r.FileName = data.FileName
	r.FileURL = data.FileURL
	r.FileType = data.FileType
	r.FileSize = data.FileSize
	r.SavePath = data.SavePath
	r.Md5 = data.Md5
	r.Addon = data.Addon
	r.fileNameGenerator = data.fileNameGenerator
	return r
}

func (r *Result) FileNameGenerator() FileNameGenerator {
	if r.fileNameGenerator == nil {
		return DefaultNameGenerator
	}
	return r.fileNameGenerator
}

func (r *Result) GenFileName() (string, error) {
	return r.FileNameGenerator()(r.FileName)
}

func (r *Result) FileIdString() string {
	return fmt.Sprintf(`%d`, r.FileID)
}
