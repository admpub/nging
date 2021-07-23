/*
   Nging is a toolbox for webmasters
   Copyright (C) 2018-present  Wenhui Shen <swh@admpub.com>

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

package term

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/pkg/sftp"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"

	"github.com/admpub/nging/v3/application/dbschema"
	"github.com/admpub/nging/v3/application/handler"
	"github.com/admpub/nging/v3/application/handler/caddy"
	"github.com/admpub/nging/v3/application/library/config"
	"github.com/admpub/nging/v3/application/library/filemanager"
	"github.com/admpub/nging/v3/application/library/notice"
	"github.com/admpub/nging/v3/application/library/respond"
	"github.com/admpub/nging/v3/application/library/sftpmanager"
	"github.com/admpub/nging/v3/application/model"
	"github.com/admpub/web-terminal/library/ssh"

	uploadChunk "github.com/admpub/nging/v3/application/registry/upload/chunk"
	uploadClient "github.com/webx-top/client/upload"
	uploadDropzone "github.com/webx-top/client/upload/driver/dropzone"
)

func sftpConnect(m *dbschema.NgingSshUser) (*sftp.Client, error) {
	account := &ssh.AccountConfig{
		User:     m.Username,
		Password: config.DefaultConfig.Decode(m.Password),
	}
	if len(m.PrivateKey) > 0 {
		account.PrivateKey = []byte(m.PrivateKey)
	}
	if len(m.Passphrase) > 0 {
		account.Passphrase = []byte(config.DefaultConfig.Decode(m.Passphrase))
	}
	config, err := ssh.NewSSHConfig(nil, nil, account)
	if err != nil {
		return nil, err
	}
	sshClient := ssh.New(config)
	err = sshClient.Connect(m.Host, m.Port)
	if err != nil {
		return nil, err
	}
	return sftp.NewClient(sshClient.Client)
}

func SftpSearch(ctx echo.Context, id uint) error {
	m := model.NewSshUser(ctx)
	err := m.Get(nil, `id`, id)
	if err != nil {
		return err
	}
	client, err := sftpConnect(m.NgingSshUser)
	if err != nil {
		return err
	}
	defer client.Close()
	query := ctx.Form(`query`)
	num := ctx.Formx(`size`, `10`).Int()
	if num <= 0 {
		num = 10
	}
	paths := sftpmanager.Search(client, query, ctx.Form(`type`), num)
	data := ctx.Data().SetData(paths)
	return ctx.JSON(data)
}

func Sftp(ctx echo.Context) error {
	ctx.Set(`activeURL`, `/term/account`)
	id := ctx.Formx(`id`).Uint()
	m := model.NewSshUser(ctx)
	err := m.Get(nil, `id`, id)
	if err != nil {
		return err
	}
	client, err := sftpConnect(m.NgingSshUser)
	if err != nil {
		return err
	}
	defer client.Close()

	mgr := sftpmanager.New(client, config.DefaultConfig.Sys.EditableFileMaxBytes, ctx)

	ppath := ctx.Form(`path`)
	do := ctx.Form(`do`)
	parentPath := ppath
	if len(ppath) == 0 {
		ppath = `/`
	} else {
		parentPath = path.Dir(ppath)
	}
	user := handler.User(ctx)
	switch do {
	case `edit`:
		data := ctx.Data()
		if _, ok := caddy.Editable(ppath); !ok {
			data.SetInfo(ctx.T(`此文件不能在线编辑`), 0)
		} else {
			content := ctx.Form(`content`)
			encoding := ctx.Form(`encoding`)
			dat, err := mgr.Edit(ppath, content, encoding)
			if err != nil {
				data.SetInfo(err.Error(), 0)
			} else {
				if ctx.IsPost() {
					data.SetInfo(ctx.T(`保存成功`), 1)
				}
				data.SetData(dat, 1)
			}
		}
		return ctx.JSON(data)
	case `mkdir`:
		data := ctx.Data()
		newName := ctx.Form(`name`)
		if len(newName) == 0 {
			data.SetInfo(ctx.T(`请输入文件夹名`), 0)
		} else {
			err = mgr.Mkdir(ppath, newName)
			if err != nil {
				data.SetError(err)
			}
			if data.GetCode() == 1 {
				data.SetInfo(ctx.T(`创建成功`))
			}
		}
		return ctx.JSON(data)
	case `rename`:
		data := ctx.Data()
		newName := ctx.Form(`name`)
		err = mgr.Rename(ppath, newName)
		if err != nil {
			data.SetInfo(err.Error(), 0)
		} else {
			data.SetInfo(ctx.T(`重命名成功`), 1)
		}
		return ctx.JSON(data)
	case `chown`:
		data := ctx.Data()
		uid := ctx.Formx(`uid`).Int()
		gid := ctx.Formx(`gid`).Int()
		err = mgr.Chown(ppath, uid, gid)
		if err != nil {
			data.SetInfo(err.Error(), 0)
		} else {
			data.SetInfo(ctx.T(`操作成功`), 1)
		}
		return ctx.JSON(data)
	case `chmod`:
		data := ctx.Data()
		mode := ctx.Formx(`mode`).Uint32() //0777 etc...
		err = mgr.Chmod(ppath, os.FileMode(mode))
		if err != nil {
			data.SetInfo(err.Error(), 0)
		} else {
			data.SetInfo(ctx.T(`操作成功`), 1)
		}
		return ctx.JSON(data)
	case `search`:
		prefix := ctx.Form(`query`)
		num := ctx.Formx(`size`, `10`).Int()
		if num <= 0 {
			num = 10
		}
		paths := mgr.Search(ppath, prefix, num)
		data := ctx.Data().SetData(paths)
		return ctx.JSON(data)
	case `delete`:
		err = mgr.Remove(ppath)
		if err != nil {
			handler.SendFail(ctx, err.Error())
		}
		return ctx.Redirect(ctx.Referer())
	case `upload`:
		var cu *uploadClient.ChunkUpload
		var opts []uploadClient.ChunkInfoOpter
		if user != nil {
			_cu := uploadChunk.ChunkUploader()
			_cu.UID = fmt.Sprintf(`user/%d`, user.Id)
			cu = &_cu
			opts = append(opts, uploadClient.OptChunkInfoMapping(uploadDropzone.MappingChunkInfo))
		}
		err = mgr.Upload(ppath, cu, opts...)
		if err != nil {
			user := handler.User(ctx)
			if user != nil {
				notice.OpenMessage(user.Username, `upload`)
				notice.Send(user.Username, notice.NewMessageWithValue(`upload`, ctx.T(`文件上传出错`), err.Error()))
			}
		}
		return respond.Dropzone(ctx, err, nil)
	default:
		var dirs []os.FileInfo
		var exit bool
		err, exit, dirs = mgr.List(ppath)
		if exit {
			return err
		}
		ctx.Set(`dirs`, dirs)
	}
	ctx.Set(`parentPath`, parentPath)
	ctx.Set(`path`, ppath)
	pathPrefix := ppath
	if ppath != `/` {
		pathPrefix = ppath + `/`
	}
	pathSlice := strings.Split(strings.Trim(pathPrefix, `/`), `/`)
	pathLinks := make(echo.KVList, len(pathSlice))
	encodedSep := filemanager.EncodedSep
	urlPrefix := ctx.Request().URL().Path() + fmt.Sprintf(`?id=%d&path=`, id) + encodedSep
	for k, v := range pathSlice {
		urlPrefix += com.URLEncode(v)
		pathLinks[k] = &echo.KV{K: v, V: urlPrefix}
		urlPrefix += encodedSep
	}
	ctx.Set(`pathLinks`, pathLinks)
	ctx.Set(`pathPrefix`, pathPrefix)
	ctx.SetFunc(`Editable`, func(fileName string) bool {
		_, ok := caddy.Editable(fileName)
		return ok
	})
	ctx.SetFunc(`Playable`, func(fileName string) string {
		mime, _ := caddy.Playable(fileName)
		return mime
	})
	ctx.Set(`data`, m.NgingSshUser)
	return ctx.Render(`term/sftp`, err)
}
