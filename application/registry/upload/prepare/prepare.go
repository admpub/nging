package prepare

import (
	"path"

	"github.com/admpub/nging/v3/application/model/file/storer"
	"github.com/admpub/nging/v3/application/registry/upload"
	"github.com/admpub/nging/v3/application/registry/upload/checker"
	"github.com/admpub/nging/v3/application/registry/upload/dbsaver"
	"github.com/admpub/nging/v3/application/registry/upload/driver"
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
	FileType   string
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

// Prepare 上传前的环境准备
func Prepare(ctx echo.Context, subdir string, fileType string, storerInfo storer.Info) (*PrepareData, error) {
	if len(subdir) == 0 {
		subdir = `default`
	}
	if !upload.Subdir.Has(subdir) {
		return nil, ctx.NewError(code.InvalidParameter, `subdir参数值“%s”未被登记`, subdir)
	}
	//echo.Dump(ctx.Forms())
	newStore := driver.Get(storerInfo.Name)
	if newStore == nil {
		return nil, ctx.NewError(code.InvalidParameter, `存储引擎“%s”未被登记`, storerInfo.Name)
	}
	dbSaverFn := dbsaver.Get(subdir)
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
		Checkin:    checker.Default,
		Subdir:     subdir,
		FileType:   fileType,
	}
	return data, nil
}
