package image

import (
	"context"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/admpub/nging/v5/application/handler"
	"github.com/admpub/nging/v5/application/library/notice"
	"github.com/docker/docker/api/types"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
	"github.com/webx-top/echo/param"

	"github.com/nging-plugins/dockermanager/application/library/dockerclient"
	"github.com/nging-plugins/dockermanager/application/library/utils"
)

// Download 下载镜像(用于备份Image)
func Download(ctx echo.Context) error {
	c, err := dockerclient.Client()
	if err != nil {
		return err
	}
	var imageIDs []string
	imageID := ctx.Param(`id`)
	if len(imageID) == 0 || imageID == `0` {
		imageIDs = ctx.FormxValues(`id[]`).Unique().Filter().String()
	} else {
		imageIDs = append(imageIDs, imageID)
	}
	if len(imageIDs) == 0 {
		return ctx.NewError(code.InvalidParameter, `参数 %s 不能为空`, `id`).SetZone(`id`)
	}
	var reader io.ReadCloser
	reader, err = c.ImageSave(ctx, imageIDs)
	if err != nil {
		return err
	}
	defer reader.Close()
	imageShortIDs := make([]string, len(imageIDs))
	for index, imageID := range imageIDs {
		imageShortIDs[index] = utils.ShortenID(imageID)
	}
	fileName := `docker-image-` + strings.Join(imageShortIDs, `-`) + `.tar`
	echo.SetAttachmentHeader(ctx, fileName, false)
	return ctx.ServeContent(reader, fileName, time.Now())
}

// Load 载入镜像(用于恢复备份的Image)
func Load(ctx echo.Context) error {
	user := handler.User(ctx)
	c, err := dockerclient.Client()
	if err != nil {
		return err
	}
	if ctx.IsPost() {
		var fp multipart.File
		var hd *multipart.FileHeader
		fp, hd, err = ctx.Request().FormFile(`file`)
		if err != nil {
			return err
		}
		op := ctx.Form(`op`, `2`)
		if op == `2` {
			// 方案2: 阻塞等待载入完成
			defer fp.Close()
			noticer := notice.NewP(ctx, `dockerImageLoad`, user.Username, context.Background()).Add(100).AutoComplete(true)
			var result types.ImageLoadResponse
			result, err = c.ImageLoad(ctx, fp, true)
			if err != nil {
				noticer.Send(err.Error(), notice.StateFailure)
				return err
			}
			dockerclient.SyncReaderToNotice(noticer, result.Body)
			handler.SendOk(ctx, ctx.T(`载入成功`))
			return ctx.Redirect(handler.URLFor(`/docker/base/image/index`))
		}
		// 方案1: 先保存到服务器，然后在后台异步载入
		srcPath := filepath.Join(os.TempDir(), `nging/docker/`+param.AsString(user.Id)+`/load-`+time.Now().Format(`20060102150405.000`))
		err = com.MkdirAll(srcPath, os.ModePerm)
		if err != nil {
			fp.Close()
			return err
		}
		var saveFile string
		saveFile, err = utils.SaveMultipartFile(fp, hd.Filename, srcPath)
		fp.Close()
		if err != nil {
			os.RemoveAll(srcPath)
			return err
		}
		err = dockerclient.StartBackgroundRun(ctx, user.Username, `dockerImageLoad`, hd.Filename, func(ctx context.Context) (io.ReadCloser, error) {
			defer os.RemoveAll(srcPath)
			fp, err := os.Open(saveFile)
			if err != nil {
				return nil, err
			}
			defer fp.Close()
			result, err := c.ImageLoad(ctx, fp, true)
			if err != nil {
				return nil, err
			}
			return result.Body, err
		})
		if err != nil {
			os.RemoveAll(srcPath)
			return err
		}
		data := ctx.Data()
		data.SetInfo(ctx.T(`启动成功`), code.Success.Int())
		return ctx.JSON(data)
	}
	ctx.Set(`activeURL`, `/docker/base/image/index`)
	ctx.Set(`title`, ctx.T(`载入镜像`))
	return ctx.Render(`docker/base/image/load`, nil)
}
