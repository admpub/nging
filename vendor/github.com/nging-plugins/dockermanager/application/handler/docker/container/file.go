package container

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/webx-top/com"
	"github.com/webx-top/echo"

	"github.com/admpub/nging/v5/application/handler"
	"github.com/admpub/nging/v5/application/library/config"
	"github.com/admpub/nging/v5/application/library/notice"
	"github.com/admpub/nging/v5/application/library/respond"
	uploadChunk "github.com/admpub/nging/v5/application/registry/upload/chunk"
	uploadClient "github.com/webx-top/client/upload"
	uploadDropzone "github.com/webx-top/client/upload/driver/dropzone"

	"github.com/nging-plugins/dockermanager/application/library/containerfs"
	"github.com/nging-plugins/dockermanager/application/library/dockerclient"
)

func File(ctx echo.Context) error {
	user := handler.User(ctx)
	c, err := dockerclient.Client()
	if err != nil {
		return err
	}
	containerID := ctx.Param(`id`)
	data, err := c.ContainerInspect(ctx, containerID)
	if err != nil {
		return err
	}
	data.ID = containerID
	filePath := ctx.Form(`path`, data.Config.WorkingDir)
	if len(filePath) > 0 {
		filePath, err = filepath.Abs(filePath)
		if err != nil {
			return err
		}
	}
	do := ctx.Form(`do`)
	mgr := containerfs.New(c, containerID, data.Platform, config.FromFile().Sys.EditableFileMaxBytes(), ctx)
	if len(filePath) == 0 {
		filePath = mgr.RootDir()
	}
	switch do {
	case `edit`:
		data := ctx.Data()
		if _, ok := Editable(filePath); !ok {
			data.SetInfo(ctx.T(`此文件不能在线编辑`), 0)
		} else {
			content := ctx.Form(`content`)
			encoding := ctx.Form(`encoding`)
			dat, err := mgr.Edit(filePath, content, encoding)
			if err != nil {
				data.SetInfo(err.Error(), 0)
			} else {
				data.SetData(dat, 1)
			}
		}
		return ctx.JSON(data)
	case `rename`:
		data := ctx.Data()
		newName := ctx.Form(`name`)
		err = mgr.Rename(filePath, newName)
		if err != nil {
			data.SetInfo(err.Error(), 0)
		} else {
			data.SetCode(1)
		}
		return ctx.JSON(data)
	case `mkdir`:
		data := ctx.Data()
		newName := ctx.Form(`name`)
		err = mgr.Mkdir(filepath.Join(filePath, newName), os.ModePerm)
		if err != nil {
			data.SetInfo(err.Error(), 0)
		} else {
			data.SetCode(1)
		}
		return ctx.JSON(data)
	case `delete`:
		err = mgr.Remove(filePath)
		if err != nil {
			handler.SendFail(ctx, err.Error())
		}
		next := ctx.Referer()
		if len(next) == 0 {
			next = ctx.Request().URL().Path() + fmt.Sprintf(`?path=%s`, com.URLEncode(filepath.Dir(filePath)))
		}
		return ctx.Redirect(next)
	case `upload`:
		var cu *uploadClient.ChunkUpload
		var opts []uploadClient.ChunkInfoOpter
		if user != nil {
			cu = uploadChunk.NewUploader(fmt.Sprintf(`user/%d`, user.Id))
			opts = append(opts, uploadClient.OptChunkInfoMapping(uploadDropzone.MappingChunkInfo))
		}
		err = mgr.Upload(filePath, cu, opts...)
		if err != nil {
			user := handler.User(ctx)
			if user != nil {
				notice.OpenMessage(user.Username, `upload`)
				notice.Send(user.Username, notice.NewMessageWithValue(`upload`, ctx.T(`文件上传出错`), err.Error()))
			}
		}
		return respond.Dropzone(ctx, err, nil)
	case `download`:
		return mgr.Download(filePath)
	default:
		var files []containerfs.FileInfo
		var exit bool
		err, exit, files = mgr.List(filePath)
		if exit {
			return err
		}
		ctx.Set(`files`, files)
	}
	//echo.Dump(lines)
	pathSeperator := mgr.Seperator()
	pathSlice := strings.Split(strings.Trim(filePath, pathSeperator), pathSeperator)
	pathLinks := make(echo.KVList, 0, len(pathSlice))
	encodedSep := mgr.URLEncodedSeperator()
	urlPrefix := ctx.Request().URL().Path() + `?path=` + encodedSep
	for _, v := range pathSlice {
		if len(v) == 0 {
			continue
		}
		urlPrefix += com.URLEncode(v)
		pathLinks = append(pathLinks, &echo.KV{K: v, V: urlPrefix})
		urlPrefix += encodedSep
	}
	ctx.Set(`pathLinks`, pathLinks)
	ctx.Set(`pathSeperator`, pathSeperator)
	if !strings.HasSuffix(filePath, pathSeperator) {
		filePath += pathSeperator
	}
	ctx.Set(`path`, filePath)
	ctx.Set(`data`, data)
	ctx.Set(`rootDir`, mgr.RootDir())
	ctx.SetFunc(`Editable`, func(fileName string) bool {
		_, ok := Editable(fileName)
		return ok
	})
	ctx.SetFunc(`Playable`, func(fileName string) string {
		mime, _ := Playable(fileName)
		return mime
	})
	ctx.Set(`activeURL`, `/docker/base/container/index`)
	ctx.Set(`title`, ctx.T(`容器文件管理`))
	return ctx.Render(`docker/base/container/file`, handler.Err(ctx, err))
}

func Editable(fileName string) (string, bool) {
	return config.FromFile().Sys.Editable(fileName)
}

func Playable(fileName string) (string, bool) {
	return config.FromFile().Sys.Playable(fileName)
}
