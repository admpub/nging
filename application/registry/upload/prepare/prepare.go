package prepare

import (
	"path"

	"github.com/admpub/nging/application/model/file/storer"
	"github.com/admpub/nging/application/registry/upload/checker"
	"github.com/admpub/nging/application/registry/upload/dbsaver"
	"github.com/admpub/nging/application/registry/upload/driver"
	uploadSubdir "github.com/admpub/nging/application/registry/upload/subdir"
	uploadClient "github.com/webx-top/client/upload"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
)

var NopChecker uploadClient.Checker = func(r *uploadClient.Result) error {
	return nil
}

type PrepareData struct {
	newStorer  driver.Constructor
	storer     driver.Storer
	StorerInfo storer.Info
	DBSaver    dbsaver.DBSaver
	Checker    uploadClient.Checker
	Checkin    checker.Checker
	Subdir     string
	TableName  string
	FieldName  string
	FileType   string
	SubdirInfo *uploadSubdir.SubdirInfo
}

func (p *PrepareData) Storer(ctx echo.Context) (driver.Storer, error) {
	var err error
	if p.storer == nil {
		p.storer, err = p.newStorer(ctx, p.Subdir)
	}
	return p.storer, err
}

func (p *PrepareData) Close() error {
	if p.storer == nil {
		return nil
	}
	return p.storer.Close()
}

func (p *PrepareData) GetFieldInfo() *uploadSubdir.FieldInfo {
	return p.SubdirInfo.GetField(p.FieldName)
}

func (p *PrepareData) AutoCropThumbSize() []uploadSubdir.ThumbSize {
	return p.SubdirInfo.AutoCropThumbSize(p.FieldName)
}

// Prepare 上传前的环境准备
func Prepare(ctx echo.Context, uploadType string, fileType string, storerInfo storer.Info) (*PrepareData, error) {
	if len(uploadType) == 0 {
		return nil, ctx.NewError(code.InvalidParameter, `请提供参数“%s”`, ctx.Path())
	}
	params := uploadSubdir.ParseUploadType(uploadType)
	if !params.IsAllowed() {
		return nil, ctx.NewError(code.InvalidParameter, `参数“%s”未被登记`, uploadType)
	}
	//echo.Dump(ctx.Forms())
	newStore := driver.Get(storerInfo.Name)
	if newStore == nil {
		return nil, ctx.NewError(code.InvalidParameter, `存储引擎“%s”未被登记`, storerInfo.Name)
	}
	dbSaverKey := params.MustGetTable()
	if len(params.Field) > 0 {
		dbSaverKey += `.` + params.Field
	}
	dbSaverFn := dbsaver.Get(dbSaverKey)
	checkerFn := func(r *uploadClient.Result) error {
		extension := path.Ext(r.FileName)
		if len(r.FileType) > 0 {
			if !uploadClient.CheckTypeExtension(fileType, extension) {
				return ctx.NewError(code.InvalidParameter, `不支持将扩展名为“%v”的文件作为“%v”类型的文件来进行上传`, extension, fileType)
			}
		} else {
			r.FileType = uploadClient.FileType(uploadClient.DetectType(extension))
		}
		return NopChecker(r)
	}
	data := &PrepareData{
		newStorer:  newStore,
		StorerInfo: storerInfo,
		DBSaver:    dbSaverFn,
		Checker:    checkerFn,
		Checkin:    uploadSubdir.CheckerGet(params.MustGetSubdir()),
		Subdir:     params.MustGetSubdir(),
		TableName:  params.MustGetTable(),
		FieldName:  params.Field,
		FileType:   fileType,
		SubdirInfo: params.SubdirInfo,
	}
	return data, nil
}
