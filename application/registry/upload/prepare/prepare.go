package prepare

import (
	"fmt"
	"io"
	"path"

	uploadClient "github.com/webx-top/client/upload"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
	"github.com/webx-top/echo/middleware/tplfunc"

	"github.com/admpub/nging/v4/application/library/common"
	uploadLibrary "github.com/admpub/nging/v4/application/library/upload"
	modelFile "github.com/admpub/nging/v4/application/model/file"
	storerUtils "github.com/admpub/nging/v4/application/model/file/storer"
	"github.com/admpub/nging/v4/application/registry/upload"
	"github.com/admpub/nging/v4/application/registry/upload/checker"
	"github.com/admpub/nging/v4/application/registry/upload/dbsaver"
	"github.com/admpub/nging/v4/application/registry/upload/driver"
	"github.com/admpub/nging/v4/application/registry/upload/thumb"
)

type PrepareData struct {
	ctx        echo.Context
	newStorer  driver.Constructor
	storer     driver.Storer
	StorerInfo storerUtils.Info
	DBSaver    dbsaver.DBSaver
	Checker    uploadClient.Checker
	Checkin    checker.Checker
	Subdir     string
	FileType   string
	multiple   bool // 是否为多文件上传
}

func (p *PrepareData) Storer() (driver.Storer, error) {
	var err error
	if p.storer == nil {
		p.storer, err = p.newStorer(p.ctx, p.Subdir)
	}
	return p.storer, err
}

func (p *PrepareData) NewStorer(subdir string) (driver.Storer, error) {
	return p.newStorer(p.ctx, subdir)
}

func (p *PrepareData) Close() error {
	if p.storer == nil {
		return nil
	}
	return p.storer.Close()
}

func (p *PrepareData) MakeModel(ownerType string, ownerID uint64) *modelFile.File {
	fileM := modelFile.NewFile(p.ctx)
	fileM.StorerName = p.StorerInfo.Name
	fileM.StorerId = p.StorerInfo.ID
	fileM.OwnerId = ownerID
	fileM.OwnerType = ownerType
	fileM.Type = p.FileType
	fileM.Subdir = p.Subdir
	return fileM
}

func (p *PrepareData) MakeCallback(fileM *modelFile.File, storer driver.Storer, subdir string) func(*uploadClient.Result, io.Reader, io.Reader) error {
	ctx := fileM.Context()
	callback := func(result *uploadClient.Result, originalReader io.Reader, _ io.Reader) error {
		fileM.Id = 0
		fileM.SetByUploadResult(result)
		if err := ctx.Begin(); err != nil {
			return err
		}
		fileM.Use(common.Tx(ctx))
		err := p.DBSaver(fileM, result, originalReader)
		if err != nil {
			ctx.Rollback()
			return err
		}
		if result.FileType.String() != `image` {
			ctx.Commit()
			return nil
		}
		thumbSizes := thumb.Registry.Get(subdir).AutoCrop()
		if len(thumbSizes) > 0 {
			thumbM := modelFile.NewThumb(ctx)
			thumbM.CPAFrom(fileM.NgingFile)
			for _, thumbSize := range thumbSizes {
				thumbM.Reset()
				if seek, ok := originalReader.(io.Seeker); ok {
					seek.Seek(0, 0)
				}
				thumbURL := tplfunc.AddSuffix(result.FileURL, fmt.Sprintf(`_%v_%v`, thumbSize.Width, thumbSize.Height))
				cropOpt := &modelFile.CropOptions{
					Options:          modelFile.ImageOptions(thumbSize.Width, thumbSize.Height),
					File:             fileM.NgingFile,
					SrcReader:        originalReader,
					Storer:           storer,
					DestFile:         storer.URLToFile(thumbURL),
					FileMD5:          ``,
					WatermarkOptions: storerUtils.GetWatermarkOptions(),
				}
				err = thumbM.Crop(cropOpt)
				if err != nil {
					ctx.Rollback()
					return err
				}
			}
		}
		ctx.Commit()
		return nil
	}
	return callback
}

func IsInvalidUpload(err error) bool {
	return uploadClient.ErrInvalidContent == err || echo.ErrNotFoundFileInput == err
}

func (p *PrepareData) IsInvalidUpload(err error) bool {
	return IsInvalidUpload(err)
}

func (p *PrepareData) MakeUploader(ownerType string, ownerID uint64, clientName string) func(fieldName string, fileType uploadClient.FileType, multiple bool) ([]*uploadClient.Result, error) {
	fileM := p.MakeModel(ownerType, ownerID)
	if len(clientName) == 0 {
		clientName = `default`
	}
	uploadCfg := uploadLibrary.Get()
	return func(fieldName string, fileType uploadClient.FileType, multiple bool) ([]*uploadClient.Result, error) {
		result := &uploadClient.Result{
			FileType: fileType,
		}
		client := NewClientWithResult(p.ctx, ownerType, ownerID, clientName, result)
		client.SetUploadMaxSize(int64(uploadCfg.MaxSizeBytes(fileType.String())))
		client.SetName(fieldName)
		_, err := p.SetMultiple(multiple).Save(fileM, clientName, client)
		if err != nil {
			return nil, err
		}
		if !multiple {
			return []*uploadClient.Result{client.GetUploadResult()}, nil
		}
		fileM.Reset()
		return client.GetBatchUploadResults(), nil
	}
}

func (p *PrepareData) SetMultiple(multiple bool) *PrepareData {
	p.multiple = multiple
	return p
}

func (p *PrepareData) Multiple() bool {
	return p.multiple
}

func (p *PrepareData) Save(fileM *modelFile.File, clientName string, clients ...uploadClient.Client) (client uploadClient.Client, err error) {
	ctx := fileM.Context()
	var result *uploadClient.Result
	if len(clients) == 0 || clients[0] == nil {
		result = &uploadClient.Result{
			FileType: uploadClient.FileType(p.FileType),
			FileName: ``,
		}
		client = NewClientWithModel(fileM, clientName, result)
	} else {
		client = clients[0]
		result = client.GetUploadResult()
		if len(result.FileType) == 0 {
			result.FileType = uploadClient.FileType(p.FileType)
		}
	}
	var (
		subdir string
		name   string
		storer driver.Storer
	)
	subdir, name, err = p.Checkin(ctx)
	if err != nil {
		client.SetError(err)
		return
	}
	result.SetFileNameGenerator(func(filename string) (string, error) {
		return storerUtils.SaveFilename(subdir, name, filename)
	})
	storer, err = p.Storer()
	if err != nil {
		client.SetError(err)
		return
	}

	callback := p.MakeCallback(fileM, storer, subdir)

	optionsSetters := []uploadClient.OptionsSetter{
		uploadClient.OptClientName(clientName),
		uploadClient.OptResult(result),
		uploadClient.OptStorer(storer),
		uploadClient.OptWatermarkOptions(storerUtils.GetWatermarkOptions()),
		uploadClient.OptChecker(p.Checker),
		uploadClient.OptCallback(callback),
	}
	if p.multiple {
		client.BatchUpload(optionsSetters...)
	} else {
		client.Upload(optionsSetters...)
	}
	err = client.GetError()
	return
}

func NoCheckTokenPrepare(ctx echo.Context, subdir string, fileType string, storerInfos ...storerUtils.Info) (*PrepareData, error) {
	p, err := Prepare(ctx, subdir, fileType, storerInfos...)
	if err != nil {
		return nil, err
	}
	p.Checkin = checker.DefaultNoCheck
	return p, nil
}

// Prepare 上传前的环境准备
func Prepare(ctx echo.Context, subdir string, fileType string, storerInfos ...storerUtils.Info) (*PrepareData, error) {
	if len(subdir) == 0 {
		subdir = `default`
	}
	if !upload.AllowedSubdir(subdir) {
		return nil, ctx.NewError(code.InvalidParameter, `subdir参数值“%s”未被登记`, subdir)
	}
	var storerInfo storerUtils.Info
	if len(storerInfos) > 0 {
		storerInfo = storerInfos[0]
	} else {
		storerInfo = storerUtils.Get()
	}
	//echo.Dump(ctx.Forms())
	newStore := driver.Get(storerInfo.Name)
	if newStore == nil {
		return nil, ctx.NewError(code.InvalidParameter, `存储引擎“%s”未被登记`, storerInfo.Name)
	}
	dbSaverFn := dbsaver.Get(subdir)
	checkerFn := func(rs *uploadClient.Result, rd io.Reader) error {
		var extension string
		cfg := uploadLibrary.Get()
		if len(fileType) > 0 {
			extension = path.Ext(rs.FileName)
			if !cfg.CheckTypeExtension(fileType, extension) {
				return ctx.NewError(code.InvalidParameter, ctx.T(`上传 %s 失败: 不支持的“%s”文件`, path.Base(rs.FileName), fileType))
			}
		}
		if len(rs.FileType) == 0 {
			if len(extension) == 0 {
				extension = path.Ext(rs.FileName)
			}
			rs.FileType = uploadClient.FileType(cfg.DetectType(extension))
		}
		if rd != nil {
			head, err := uploadClient.ReadHeadBytes(rd)
			if err != nil {
				return err
			}
			if !uploadClient.IsTypeString(head, string(rs.FileType)) {
				return ctx.NewError(code.InvalidParameter, `上传 %s 失败: 文件格式不正确`, path.Base(rs.FileName))
			}
		}
		return NopChecker(rs, rd)
	}
	data := &PrepareData{
		ctx:        ctx,
		newStorer:  newStore,
		StorerInfo: storerInfo,
		DBSaver:    dbSaverFn,
		Checker:    checkerFn,
		Checkin:    checker.DefaultWithVerify,
		Subdir:     subdir,
		FileType:   fileType,
	}
	return data, nil
}
