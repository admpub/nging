package upload

import (
	"errors"
	"fmt"
	"mime/multipart"
	"os"

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
	var body ReadCloserWithSize
	body, a.err = a.Body()
	if a.err != nil {
		return a
	}
	defer body.Close()

	if options.Result == nil {
		options.Result = a.Data
	} else {
		options.Result.CopyFrom(a.Data)
	}

	uploadMaxSize := options.MaxSize
	if uploadMaxSize == 0 {
		uploadMaxSize = a.UploadMaxSize()
	}
	if uploadMaxSize > 0 && body.Size() > uploadMaxSize {
		a.err = fmt.Errorf(`%w: %v>%v`, ErrFileTooLarge,
			com.FormatBytes(body.Size(), 2, true),
			com.FormatBytes(uploadMaxSize, 2, true),
		)
		return a
	}

	if err := a.fireReadBeforeHook(options, a.Data); err != nil {
		a.err = err
		return a
	}

	file, ok := body.(multipart.File)
	if !ok {
		file, a.err = AsFile(body)
		if a.err != nil {
			return a
		}
		defer file.Close()
	}
	if a.chunkUpload != nil {
		info := &ChunkInfo{
			Mapping:     a.fieldMapping,
			FileName:    options.Result.FileName,
			CurrentSize: uint64(options.Result.FileSize),
		}
		a.Context.Request().MultipartForm()
		info.Init(func(name string) string {
			return a.Form(name)
		}, a.Header)
		_, a.err = a.chunkUpload.ChunkUpload(a.Context, info, file)
		if a.err == nil { // 上传成功
			if a.chunkUpload.Merged() {
				var fp *os.File
				fp, a.err = os.Open(a.chunkUpload.GetSavePath())
				if a.err != nil {
					return a
				}
				defer fp.Close()
				a.err = a.saveFile(options.Result, fp, options)
				if a.err != nil {
					return a
				}
				// 上传到最终位置后删除合并后的文件
				os.Remove(a.chunkUpload.GetSavePath())
			}
			return a
		}
		if !errors.Is(a.err, ErrChunkUnsupported) { // 上传出错
			if errors.Is(a.err, ErrChunkUploadCompleted) ||
				errors.Is(a.err, ErrFileUploadCompleted) {
				a.err = nil
				return a
			}
			return a
		}
		// 不支持分片上传
	}

	a.err = a.saveFile(options.Result, file, options)
	return a
}

func (a *BaseClient) fireReadBeforeHook(options *Options, result *Result) error {
	for _, hook := range a.readBefore {
		err := hook(result)
		if err != nil {
			return err
		}
	}
	for _, hook := range options.ReadBefore {
		err := hook(result)
		if err != nil {
			return err
		}
	}
	return nil
}

// BatchUpload 批量上传
func (a *BaseClient) BatchUpload(opts ...OptionsSetter) Client {
	req := a.Request()
	if req == nil {
		a.err = ErrInvalidContent
		return a
	}
	if err := a.checkRequestBodySize(); err != nil {
		a.err = err
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
	uploadMaxSize := options.MaxSize
	if uploadMaxSize == 0 {
		uploadMaxSize = a.UploadMaxSize()
	}
	for _, fileHdr := range files {
		//for each fileheader, get a handle to the actual file
		if uploadMaxSize > 0 && fileHdr.Size > uploadMaxSize {
			a.err = fmt.Errorf(
				`%w: %v>%v`, ErrFileTooLarge,
				com.FormatBytes(fileHdr.Size, 2, true),
				com.FormatBytes(uploadMaxSize, 2, true),
			)
			return a
		}

		result := &Result{
			FileName: fileHdr.Filename,
			FileSize: fileHdr.Size,
		}

		if err := a.fireReadBeforeHook(options, result); err != nil {
			a.err = err
			return a
		}

		var file multipart.File
		file, a.err = fileHdr.Open()
		if a.err != nil {
			return a
		}
		if a.chunkUpload != nil {
			info := &ChunkInfo{
				Mapping:     a.fieldMapping,
				FileName:    fileHdr.Filename,
				CurrentSize: uint64(fileHdr.Size),
			}
			info.Init(func(name string) string {
				return a.Form(name)
			}, a.Header)
			_, a.err = a.chunkUpload.ChunkUpload(a.Context, info, file)
			if a.err == nil { // 上传成功
				file.Close()
				if a.chunkUpload.Merged() {
					file, a.err = os.Open(a.chunkUpload.GetSavePath())
					if a.err != nil {
						return a
					}
					a.err = a.saveFile(result, file, options)
					file.Close()
					if a.err != nil {
						return a
					}
					// 上传到最终位置后删除合并后的文件
					os.Remove(a.chunkUpload.GetSavePath())
					a.Results.Add(result)
				}
				continue
			}
			if !errors.Is(a.err, ErrChunkUnsupported) { // 上传出错
				file.Close()
				if errors.Is(a.err, ErrChunkUploadCompleted) ||
					errors.Is(a.err, ErrFileUploadCompleted) {
					a.err = nil
					return a
				}
				return a
			}
			// 不支持分片上传
		}
		a.err = a.saveFile(result, file, options)
		file.Close()
		if a.err != nil {
			return a
		}
		a.Results.Add(result)
	}
	return a
}

func (a *BaseClient) saveFile(result *Result, file multipart.File, options *Options) (err error) {
	if options.Checker != nil {
		err = options.Checker(result, file)
		file.Seek(0, 0)
		if err != nil {
			return
		}
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
	if len(result.Md5) == 0 {
		result.Md5, err = checksum.MD5sumReader(file)
		if err != nil {
			return
		}
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
