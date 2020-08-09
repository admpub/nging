package helper

import (
	"fmt"
	"strings"

	"github.com/webx-top/db"
	"github.com/webx-top/echo"

	"github.com/admpub/nging/application/dbschema"
	fileModel "github.com/admpub/nging/application/model/file"
	uploadSubdir "github.com/admpub/nging/application/registry/upload/subdir"
)

func OnRemoveOwnerFile(ctx echo.Context, typ string, id interface{}, ownerDir string) error {
	fileM := fileModel.NewFile(ctx)
	err := fileM.DeleteBy(db.And(
		db.Cond{`table_id`: id},
		db.Cond{`table_name`: typ},
	))
	return err
}

func OnUpdateOwnerFilePath(ctx echo.Context,
	src string, typ string, id interface{},
	newSavePath string, newViewURL string) error {
	fileM := &dbschema.NgingFile{}
	//embedM := &dbschema.NgingFileEmbedded{}
	params := uploadSubdir.ParseUploadType(typ)
	thumbM := &dbschema.NgingFileThumb{}
	cond := db.NewCompounds()
	cond.Add(db.Cond{`table_id`: id})
	cond.Add(db.Cond{`table_name`: params.MustGetTable()})
	cond.Add(db.Cond{`field_name`: params.Field})
	cond.Add(db.Cond{`view_url`: src})
	_, err := fileM.ListByOffset(nil, nil, 0, -1, cond.And())
	if err != nil {
		return err
	}
	replaceFrom := `/0/`
	replaceTo := `/` + fmt.Sprint(id) + `/`
	for _, row := range fileM.Objects() {
		err = row.SetFields(nil, echo.H{
			`save_path`:  newSavePath,
			`view_url`:   newViewURL,
			`used_times`: 1,
		}, db.Cond{`id`: row.Id})
		if err != nil {
			return err
		}
		_, err = thumbM.ListByOffset(nil, nil, 0, -1, db.Cond{`file_id`: row.Id})
		if err != nil {
			return err
		}
		for _, thumb := range thumbM.Objects() {
			thumb.SavePath = strings.Replace(thumb.SavePath, replaceFrom, replaceTo, -1)
			thumb.ViewUrl = strings.Replace(thumb.ViewUrl, replaceFrom, replaceTo, -1)
			err = thumb.SetFields(nil, echo.H{
				`save_path`: thumb.SavePath,
				`view_url`:  thumb.ViewUrl,
			}, db.Cond{`id`: thumb.Id})
			if err != nil {
				return err
			}
		}
	}
	return err
}
