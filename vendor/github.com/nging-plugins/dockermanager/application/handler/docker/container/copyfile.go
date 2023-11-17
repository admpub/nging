package container

import (
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
	"github.com/webx-top/echo/param"

	"github.com/admpub/nging/v5/application/handler"

	"github.com/nging-plugins/dockermanager/application/library/dockerclient"
	"github.com/nging-plugins/dockermanager/application/library/utils"
)

func FileImport(ctx echo.Context) error {
	var err error
	if ctx.IsPost() {
		_, err = ctx.Request().MultipartForm()
		if err != nil {
			goto END
		}
		err = CopyToContainer(ctx)
		if err != nil {
			goto END
		}
		return err
	}

END:
	ctx.Set(`activeURL`, `/docker/base/container/index`)
	ctx.Set(`title`, ctx.T(`导入文件`))
	return ctx.Render(`docker/base/container/file_import`, handler.Err(ctx, err))
}

func FileExport(ctx echo.Context) error {
	var err error
	if ctx.IsPost() {
		err = CopyFromContainer(ctx)
		if err != nil {
			goto END
		}
		return err
	}

END:
	ctx.Set(`activeURL`, `/docker/base/container/index`)
	ctx.Set(`title`, ctx.T(`导出文件`))
	return ctx.Render(`docker/base/container/file_export`, handler.Err(ctx, err))
}

func CopyToContainer(ctx echo.Context) error {
	c, err := dockerclient.Client()
	if err != nil {
		return err
	}
	user := handler.User(ctx)
	containerID := ctx.Param(`id`)
	srcPath := ctx.Form(`srcPath`)
	dstPath := ctx.Form(`dstPath`)
	if len(dstPath) == 0 {
		return ctx.NewError(code.InvalidParameter, `请填写容器目录`).SetZone(`dstPath`)
	}
	if len(srcPath) == 0 {
		mf, err := ctx.Request().MultipartForm()
		if err != nil {
			return err
		}
		files, ok := mf.File[`files`]
		if !ok || len(files) == 0 {
			return ctx.NewError(code.InvalidParameter, `上传文件无效`).SetZone(`files`)
		}
		defer mf.RemoveAll()
		srcPath = filepath.Join(os.TempDir(), `nging/docker/`+param.AsString(user.Id)+`/copy2container-`+time.Now().Format(`20060102150405.000`))
		err = com.MkdirAll(srcPath, os.ModePerm)
		if err != nil {
			return err
		}
		defer os.RemoveAll(srcPath)
		for _, file := range files {
			_, err = utils.SaveUploadedFile(file, srcPath)
			if err != nil {
				return err
			}
		}
	}
	opts := types.CopyToContainerOptions{}
	filenames, err := utils.GetFilenames(srcPath)
	if err != nil {
		return err
	}
	prw, err := utils.CompressTar(ctx, filenames)
	if err != nil {
		return err
	}
	defer prw.Close()
	err = prw.DoRead(func(r io.Reader) error {
		return c.CopyToContainer(ctx, containerID, dstPath, r, opts)
	})
	if err == nil {
		handler.SendOk(ctx, ctx.T(`操作成功`))
	} else {
		handler.SendFail(ctx, err.Error())
	}

	return ctx.Redirect(handler.URLFor(`/docker/base/container/index`))
}

func StatPath(ctx echo.Context) error {
	c, err := dockerclient.Client()
	if err != nil {
		return err
	}
	containerID := ctx.Param(`id`)
	path := ctx.Form(`path`)
	var stat types.ContainerPathStat
	stat, err = c.ContainerStatPath(ctx, containerID, path)
	if err != nil {
		return err
	}
	ctx.Logger().Debugf(`ContainerStatPath: %+v`, stat)
	return err
}

func CopyFromContainer(ctx echo.Context) error {
	c, err := dockerclient.Client()
	if err != nil {
		return err
	}
	containerID := ctx.Param(`id`)
	srcPath := ctx.Form(`srcPath`)
	if len(srcPath) == 0 {
		return ctx.NewError(code.InvalidParameter, `请填写容器文件路径`).SetZone(`srcPath`)
	}
	dstPath := ctx.Form(`dstPath`)
	reader, pathStat, err := c.CopyFromContainer(ctx, containerID, srcPath)
	if err != nil {
		return err
	}
	defer reader.Close()
	if len(dstPath) == 0 {
		fileName := utils.ShortenID(containerID) + `-` + com.Md5(srcPath) + `.tar`
		echo.SetAttachmentHeader(ctx, fileName, false)
		return ctx.ServeContent(reader, fileName, pathStat.Mtime)
	}
	err = utils.DecompressTar(ctx, reader, dstPath)
	if err == nil {
		handler.SendOk(ctx, ctx.T(`操作成功`))
	} else {
		handler.SendFail(ctx, err.Error())
	}

	return ctx.Redirect(handler.URLFor(`/docker/base/container/index`))
}
