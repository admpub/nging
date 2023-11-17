package image

import (
	"context"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/admpub/nging/v5/application/handler"
	"github.com/admpub/nging/v5/application/library/notice"
	"github.com/docker/docker/api/types"
	"github.com/nging-plugins/dockermanager/application/library/dockerclient"
	"github.com/nging-plugins/dockermanager/application/library/utils"
	"github.com/nging-plugins/dockermanager/application/request"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
	"github.com/webx-top/echo/param"
)

// Import 导入容器快照
func Import(ctx echo.Context) error {
	user := handler.User(ctx)
	c, err := dockerclient.Client()
	if err != nil {
		return err
	}
	if ctx.IsPost() {
		req := echo.GetValidated(ctx).(*request.ImageImport)
		source := types.ImageImportSource{
			SourceName: req.SourceName,
		}
		var fp multipart.File
		var hd *multipart.FileHeader
		hasUploadedFile := len(req.SourceName) == 0 || req.SourceName == `-`
		if hasUploadedFile {
			fp, hd, err = ctx.Request().FormFile(`file`)
			if err != nil {
				return err
			}
		}
		op := ctx.Form(`op`, `2`)
		if op == `2` {
			// 方案2: 阻塞等待载入完成
			if hasUploadedFile {
				defer fp.Close()
				source.Source = fp
				source.SourceName = `-`
			}
			noticer := notice.NewP(ctx, `dockerImageImport`, user.Username, context.Background()).Add(100).AutoComplete(true)
			var reader io.ReadCloser
			reader, err = c.ImageImport(ctx, source, req.Ref, req.ImageImportOptions)
			if err != nil {
				noticer.Send(err.Error(), notice.StateFailure)
				return err
			}
			dockerclient.SyncReaderToNotice(noticer, reader)
			handler.SendOk(ctx, ctx.T(`导入成功`))
			return ctx.Redirect(handler.URLFor(`/docker/base/image/index`))
		}
		// 方案1: 先保存到服务器，然后在后台异步载入
		var cleanup = func() {}
		var saveFile string
		if hasUploadedFile {
			srcPath := filepath.Join(os.TempDir(), `nging/docker/`+param.AsString(user.Id)+`/import-`+time.Now().Format(`20060102150405.000`))
			err = com.MkdirAll(srcPath, os.ModePerm)
			if err != nil {
				fp.Close()
				return err
			}
			cleanup = func() {
				os.RemoveAll(srcPath)
			}
			saveFile, err = utils.SaveMultipartFile(fp, hd.Filename, srcPath)
			fp.Close()
			if err != nil {
				cleanup()
				return err
			}
		}
		err = dockerclient.StartBackgroundRun(ctx, user.Username, `dockerImageImport`, req.Ref, func(ctx context.Context) (io.ReadCloser, error) {
			defer cleanup()
			if hasUploadedFile {
				fp, err := os.Open(saveFile)
				if err != nil {
					return nil, err
				}
				source.Source = fp
				source.SourceName = `-`
				defer fp.Close()
			}
			return c.ImageImport(ctx, source, req.Ref, req.ImageImportOptions)
		})
		if err != nil {
			cleanup()
			return err
		}
		data := ctx.Data()
		data.SetInfo(ctx.T(`启动成功`), code.Success.Int())
		return ctx.JSON(data)
	}
	ctx.Set(`activeURL`, `/docker/base/image/index`)
	ctx.Set(`title`, ctx.T(`导入镜像`))
	return ctx.Render(`docker/base/image/import`, nil)
}
