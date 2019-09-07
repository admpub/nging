package caddy

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/admpub/nging/application/handler"
	"github.com/admpub/nging/application/library/config"
	"github.com/admpub/nging/application/library/filemanager"
	"github.com/admpub/nging/application/library/notice"
	"github.com/admpub/nging/application/model"
	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
)

func VhostFile(ctx echo.Context) error {
	var err error
	id := ctx.Formx(`id`).Uint()
	filePath := ctx.Form(`path`)
	do := ctx.Form(`do`)
	m := model.NewVhost(ctx)
	err = m.Get(nil, db.Cond{`id`: id})
	mgr := filemanager.New(m.Root, config.DefaultConfig.Sys.EditableFileMaxBytes, ctx)
	absPath := m.Root
	if err == nil && len(m.Root) > 0 {
		var exit bool

		if len(filePath) > 0 {
			filePath = filepath.Clean(filePath)
			absPath = filepath.Join(m.Root, filePath)
		}

		switch do {
		case `edit`:
			data := ctx.Data()
			if _, ok := Editable(absPath); !ok {
				data.SetInfo(ctx.T(`此文件不能在线编辑`), 0)
			} else {
				content := ctx.Form(`content`)
				encoding := ctx.Form(`encoding`)
				dat, err := mgr.Edit(absPath, content, encoding)
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
			err = mgr.Rename(absPath, newName)
			if err != nil {
				data.SetInfo(err.Error(), 0)
			} else {
				data.SetCode(1)
			}
			return ctx.JSON(data)
		case `mkdir`:
			data := ctx.Data()
			newName := ctx.Form(`name`)
			err = mgr.Mkdir(filepath.Join(absPath, newName), os.ModePerm)
			if err != nil {
				data.SetInfo(err.Error(), 0)
			} else {
				data.SetCode(1)
			}
			return ctx.JSON(data)
		case `delete`:
			err = mgr.Remove(absPath)
			if err != nil {
				handler.SendFail(ctx, err.Error())
			}
			return ctx.Redirect(ctx.Referer())
		case `upload`:
			err = mgr.Upload(absPath)
			if err != nil {
				user := handler.User(ctx)
				if user != nil {
					notice.OpenMessage(user.Username, `upload`)
					notice.Send(user.Username, notice.NewMessageWithValue(`upload`, ctx.T(`文件上传出错`), err.Error()))
				}
				return ctx.JSON(echo.H{`error`: err.Error()}, 500)
			}
			return ctx.String(`OK`)
		default:
			var dirs []os.FileInfo
			err, exit, dirs = mgr.List(absPath)
			ctx.Set(`dirs`, dirs)
		}
		if exit {
			return err
		}
	}
	ctx.Set(`data`, m)
	if filePath == `.` {
		filePath = ``
	}
	pathSlice := strings.Split(strings.Trim(filePath, echo.FilePathSeparator), echo.FilePathSeparator)
	pathLinks := make(echo.KVList, len(pathSlice))
	encodedSep := filemanager.EncodedSepa
	urlPrefix := fmt.Sprintf(`/caddy/vhost_file?id=%d&path=`, id) + encodedSep
	for k, v := range pathSlice {
		urlPrefix += com.URLEncode(v)
		pathLinks[k] = &echo.KV{K: v, V: urlPrefix}
		urlPrefix += encodedSep
	}
	ctx.Set(`pathLinks`, pathLinks)
	ctx.Set(`rootPath`, strings.TrimSuffix(m.Root, echo.FilePathSeparator))
	ctx.Set(`path`, filePath)
	ctx.Set(`absPath`, absPath)
	ctx.SetFunc(`Editable`, func(fileName string) bool {
		_, ok := Editable(fileName)
		return ok
	})
	ctx.SetFunc(`Playable`, func(fileName string) string {
		mime, _ := Playable(fileName)
		return mime
	})
	ctx.Set(`activeURL`, `/caddy/vhost`)
	return ctx.Render(`caddy/file`, err)
}

func Editable(fileName string) (string, bool) {
	if config.DefaultConfig.Sys.EditableFileExtensions == nil {
		return "", false
	}
	ext := strings.TrimPrefix(filepath.Ext(fileName), `.`)
	ext = strings.ToLower(ext)
	typ, ok := config.DefaultConfig.Sys.EditableFileExtensions[ext]
	return typ, ok
}

func Playable(fileName string) (string, bool) {
	if config.DefaultConfig.Sys.PlayableFileExtensions == nil {
		config.DefaultConfig.Sys.PlayableFileExtensions = map[string]string{
			`mp4`:  `video/mp4`,
			`m3u8`: `application/x-mpegURL`,
		}
	}
	ext := strings.TrimPrefix(filepath.Ext(fileName), `.`)
	ext = strings.ToLower(ext)
	typ, ok := config.DefaultConfig.Sys.PlayableFileExtensions[ext]
	return typ, ok
}
