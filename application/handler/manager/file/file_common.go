package file

import (
	"github.com/coscms/webcore/library/backend"
	"github.com/coscms/webcore/model/file"
	"github.com/coscms/webcore/registry/upload"
	"github.com/webx-top/echo"
)

func FileListWithOwner(ctx echo.Context, ownerType string, ownerID uint64) error {
	if err := setUploadURL(ctx); err != nil {
		return err
	}
	err := List(ctx, ownerType, ownerID)
	ctx.Set(`dialog`, false)
	ctx.Set(`multiple`, true)
	partial := ctx.Formx(`partial`).Bool()
	if partial {
		return ctx.Render(`manager/file/list.main.content`, err)
	}
	ctx.Set(`subdirList`, upload.Subdir.Slice())
	return ctx.Render(`manager/file/list`, err)
}

func FileDeleteWithOwner(ctx echo.Context, ownerType string, ownerID uint64) (err error) {
	id := ctx.Paramx("id").Uint64()
	fileM := file.NewFile(ctx)
	if id == 0 {
		ids := ctx.FormxValues(`id`).Uint64()
		for _, id := range ids {
			err = fileM.DeleteByID(id, ownerType, ownerID)
			if err != nil {
				return err
			}
		}
		goto END
	}
	err = fileM.DeleteByID(id, ownerType, ownerID)
	if err != nil {
		return err
	}

END:
	return ctx.Redirect(backend.URLFor(`/manager/file/list`))
}

func FinderWithOwner(ctx echo.Context, ownerType string, ownerID uint64) error {
	if err := setUploadURL(ctx); err != nil {
		return err
	}
	err := List(ctx, ownerType, ownerID)
	multiple := ctx.Formx(`multiple`).Bool()
	ctx.Set(`dialog`, true)
	ctx.Set(`multiple`, multiple)
	partial := ctx.Formx(`partial`).Bool()
	ctx.Set(`partial`, partial)
	if partial {
		return ctx.Render(`manager/file/list.main.content`, err)
	}
	ctx.Set(`subdirList`, upload.Subdir.Slice())
	return ctx.Render(`manager/file/finder`, err)
}
