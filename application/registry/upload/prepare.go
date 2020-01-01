package upload

import (
	"path"

	uploadClient "github.com/webx-top/client/upload"
	"github.com/webx-top/echo"
)

var NopChecker uploadClient.Checker = func(r *uploadClient.Result) error {
	return nil
}

type PrepareData struct {
	newStorer    Constructor
	storer       Storer
	StorerEngine string
	DBSaver      DBSaver
	Checker      uploadClient.Checker
	Checkin      Checker
	TableName    string
	FieldName    string
}

func (p *PrepareData) Storer(ctx echo.Context) Storer {
	if p.storer == nil {
		p.storer = p.newStorer(ctx, p.TableName) // 使用表名称作为文件夹名
	}
	return p.storer
}

func (p *PrepareData) Close() error {
	if p.storer == nil {
		return nil
	}
	return p.storer.Close()
}

// Prepare 上传前的环境准备
func Prepare(ctx echo.Context, uploadType string, fileType string, storerEngine string) (*PrepareData, error) {
	if len(uploadType) == 0 {
		return nil, ctx.E(`请提供参数“%s”`, ctx.Path())
	}
	tableName, fieldName, defaults := GetTableInfo(uploadType)
	if !SubdirIsAllowed(uploadType, defaults...) {
		return nil, ctx.E(`参数“%s”未被登记`, uploadType)
	}
	//echo.Dump(ctx.Forms())
	newStore := StorerGet(storerEngine)
	if newStore == nil {
		return nil, ctx.E(`存储引擎“%s”未被登记`, storerEngine)
	}
	dbsaver := DBSaverGet(uploadType, defaults...)
	checker := func(r *uploadClient.Result) error {
		extension := path.Ext(r.FileName)
		if len(r.FileType) > 0 {
			if !uploadClient.CheckTypeExtension(fileType, extension) {
				return ctx.E(`不支持将扩展名为“%v”的文件作为“%v”类型的文件来进行上传`, extension, fileType)
			}
		} else {
			r.FileType = uploadClient.FileType(uploadClient.DetectType(extension))
		}
		return NopChecker(r)
	}
	data := &PrepareData{
		newStorer:    newStore,
		StorerEngine: storerEngine,
		DBSaver:      dbsaver,
		Checker:      checker,
		Checkin:      CheckerGet(uploadType, defaults...),
		TableName:    tableName,
		FieldName:    fieldName,
	}
	return data, nil
}
