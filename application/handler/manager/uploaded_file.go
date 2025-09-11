/*
   Nging is a toolbox for webmasters
   Copyright (C) 2018-present Wenhui Shen <swh@admpub.com>

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU Affero General Public License as published
   by the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Affero General Public License for more details.

   You should have received a copy of the GNU Affero General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package manager

import (
	"fmt"

	"github.com/webx-top/echo"

	"github.com/coscms/webcore/library/backend"
	"github.com/coscms/webcore/library/filemanager"
	"github.com/coscms/webcore/library/filemanager/filemanagerhandler"
	uploadLibrary "github.com/coscms/webcore/library/upload"
	"github.com/coscms/webcore/registry/upload/chunk"
)

// UploadedFile 本地附件文件管理
func UploadedFile(ctx echo.Context) error {
	ctx.Set(`activeURL`, `/manager/uploaded/file`)
	return Uploaded(ctx, `file`)
}

func UploadedChunk(ctx echo.Context) error {
	ctx.Set(`activeURL`, `/manager/uploaded/file`)
	return Uploaded(ctx, `chunk`)
}

func UploadedMerged(ctx echo.Context) error {
	ctx.Set(`activeURL`, `/manager/uploaded/file`)
	return Uploaded(ctx, `merged`)
}

// Uploaded 本地附件文件管理
func Uploaded(ctx echo.Context, uploadType string) error {
	var (
		root      string
		canUpload bool
		canEdit   bool
		canDelete bool
	)
	switch uploadType {
	case `chunk`:
		root = chunk.ChunkTempDir
		canDelete = true
	case `merged`:
		root = chunk.MergeSaveDir
		canDelete = true
	case `file`:
		canUpload = true
		canEdit = true
		canDelete = true
		root = uploadLibrary.UploadDir
	default:
		return echo.ErrNotFound
	}
	id := ctx.Formx(`id`).Uint()
	gallery := ctx.Formx(`gallery`).Bool()
	var urlParam string
	if gallery {
		urlParam += `&gallery=1`
	}
	globalPrefix := backend.URLFor(`/manager/uploaded/` + uploadType)
	urlPrefix := globalPrefix + fmt.Sprintf(`?id=%d`+urlParam+`&path=`, id) + filemanager.EncodedSepa
	h := filemanagerhandler.New(root, urlPrefix)
	h.SetCanUpload(canUpload)
	h.SetCanEdit(canEdit)
	h.SetCanDelete(canDelete)
	h.SetCanChmod(false)
	h.SetCanChown(false)
	err := h.Handle(ctx)
	if err != nil || ctx.Response().Committed() {
		return err
	}
	ctx.Set(`uploadType`, uploadType)
	ctx.SetFunc(`URLPrefix`, func() string {
		return globalPrefix
	})
	if gallery {
		return ctx.Render(`manager/uploaded_photo`, err)
	}
	return ctx.Render(`manager/uploaded_file`, err)
}
