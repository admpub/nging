package upload

import (
	"fmt"
	"mime/multipart"

	"github.com/admpub/checksum"
	"github.com/admpub/log"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
)

// Upload 单个文件上传
func (a *BaseClient) Upload(opts ...OptionsSetter) Client {
	options := &Options{}
	for _, opt := range opts {
		opt(options)
	}
	if options.Result == nil {
		options.Result = a.Data
	} else {
		options.Result.CopyFrom(a.Data)
	}
	var body ReadCloserWithSize
	body, a.err = a.Body()
	if a.err != nil {
		return a
	}
	defer body.Close()
	if body.Size() > a.uploadMaxSize {
		a.err = fmt.Errorf(`%w: %v`, ErrFileTooLarge, com.FormatBytes(a.uploadMaxSize))
		return a
	}
	file, ok := body.(multipart.File)
	if !ok {
		file, a.err = AsFile(body)
		if a.err != nil {
			return a
		}
	}

	a.err = a.saveFile(options.Result, file, options)
	return a
}

// BatchUpload 批量上传
func (a *BaseClient) BatchUpload(opts ...OptionsSetter) Client {
	req := a.Request()
	if req == nil {
		a.err = ErrInvalidContent
		return a
	}
	m := req.MultipartForm()
	if m == nil || m.File == nil {
		a.err = ErrInvalidContent
		return a
	}
	options := &Options{}
	for _, opt := range opts {
		opt(options)
	}
	files, ok := m.File[a.Name()]
	if !ok {
		a.err = echo.ErrNotFoundFileInput
		return a
	}
	for _, fileHdr := range files {
		//for each fileheader, get a handle to the actual file
		var file multipart.File
		file, a.err = fileHdr.Open()
		if a.err != nil {
			if file != nil {
				file.Close()
			}
			return a
		}
		if fileHdr.Size > a.uploadMaxSize {
			a.err = fmt.Errorf(`%w: %v`, ErrFileTooLarge, com.FormatBytes(a.uploadMaxSize))
			file.Close()
			return a
		}
		result := &Result{
			FileName: fileHdr.Filename,
			FileSize: fileHdr.Size,
		}
		err := a.saveFile(result, file, options)
		file.Close()
		if err != nil {
			a.err = err
			return a
		}
		a.Results.Add(result)
	}
	return a
}

func (a *BaseClient) saveFile(result *Result, file multipart.File, options *Options) (err error) {
	if options.Checker != nil {
		if err = options.Checker(result); err != nil {
			return
		}
	}
	result.Md5, err = checksum.MD5sumReader(file)
	if err != nil {
		return
	}
	var dstFile string
	dstFile, err = options.Result.FileNameGenerator()(result.FileName)
	if err != nil {
		if err == ErrExistsFile {
			log.Warn(result.FileName, `:`, ErrExistsFile)
			err = nil
		}
		return
	}
	if len(dstFile) == 0 {
		return
	}
	if len(result.SavePath) > 0 {
		return
	}
	originalFile := file
	file.Seek(0, 0)
	for _, hook := range options.SaveBefore {
		newFile, size, err := hook(file, result, options)
		if err != nil {
			return err
		}
		file = newFile
		if size > 0 {
			result.FileSize = size
		}
	}
	result.SavePath, result.FileURL, err = options.Storer.Put(dstFile, file, result.FileSize)
	if err != nil {
		return
	}
	file.Seek(0, 0)
	if err = options.Callback(result, originalFile, file); err != nil {
		options.Storer.Delete(dstFile)
		return
	}
	return
}
