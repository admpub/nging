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
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/admpub/nging/application/dbschema"
	"github.com/admpub/nging/application/handler"
	"github.com/admpub/nging/application/handler/caddy"
	"github.com/admpub/nging/application/library/charset"
	"github.com/admpub/nging/application/library/config"
	"github.com/admpub/nging/application/library/filemanager"
	"github.com/admpub/nging/application/library/notice"
	"github.com/admpub/nging/application/model"
	"github.com/admpub/web-terminal/library/ssh"
	"github.com/pkg/sftp"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
)

func sftpConnect(m *dbschema.SshUser) (*sftp.Client, error) {
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
	mgr, err := sftpConnect(m.SshUser)
	if err != nil {
		return err
	}
	defer mgr.Close()
	var (
		paths  []string
		prefix string
		ppath  string
	)
	query := ctx.Form(`query`)
	if strings.HasSuffix(query, `/`) {
		ppath = query
	} else {
		prefix = path.Base(query)
		ppath = path.Dir(query)
	}
	num := ctx.Formx(`size`, `10`).Int()
	if num <= 0 {
		num = 10
	}
	if len(ppath) == 0 {
		ppath = `/`
	}
	var onlyDir bool
	switch ctx.Form(`type`) {
	case `dir`:
		onlyDir = true
	case `file`:
		onlyDir = false
	default:
		onlyDir = true
	}
	dirs, _ := mgr.ReadDir(ppath)
	for _, d := range dirs {
		if onlyDir && d.IsDir() == false {
			continue
		}
		if len(paths) >= num {
			break
		}
		name := d.Name()
		if len(prefix) == 0 || strings.HasPrefix(name, prefix) {
			paths = append(paths, path.Join(ppath, name)+`/`)
			continue
		}
	}
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
	mgr, err := sftpConnect(m.SshUser)
	if err != nil {
		return err
	}
	defer mgr.Close()
	ppath := ctx.Form(`path`)
	do := ctx.Form(`do`)
	parentPath := ppath
	if len(ppath) == 0 {
		ppath = `/`
	} else {
		parentPath = path.Dir(ppath)
	}
	switch do {
	case `edit`:
		data := ctx.Data()
		if _, ok := caddy.Editable(ppath); !ok {
			data.SetInfo(ctx.T(`此文件不能在线编辑`), 0)
		} else {
			content := ctx.Form(`content`)
			encoding := ctx.Form(`encoding`)

			f, err := mgr.Open(ppath)
			if err != nil {
				return ctx.JSON(data.SetError(err))
			}
			defer f.Close()
			fi, err := f.Stat()
			if err != nil {
				return ctx.JSON(data.SetError(err))
			}
			if fi.IsDir() {
				return ctx.JSON(data.SetInfo(ctx.T(`不能编辑文件夹`), 0))
			}
			if config.DefaultConfig.Sys.EditableFileMaxBytes > 0 && fi.Size() > config.DefaultConfig.Sys.EditableFileMaxBytes {
				return ctx.JSON(data.SetInfo(ctx.T(`很抱歉，不支持编辑超过%v的文件`, com.FormatByte(config.DefaultConfig.Sys.EditableFileMaxBytes), 0)))
			}
			encoding = strings.ToLower(encoding)
			isUTF8 := len(encoding) == 0 || encoding == `utf-8`
			if ctx.IsPost() {
				b := []byte(content)
				if !isUTF8 {
					b, err = charset.Convert(`utf-8`, encoding, b)
					if err != nil {
						return ctx.JSON(data.SetError(err))
					}
				}
				f.Close()
				r := bytes.NewReader(b)
				f, err = mgr.OpenFile(ppath, os.O_CREATE|os.O_RDWR|os.O_TRUNC)
				if err != nil {
					return ctx.JSON(data.SetError(err))
				}
				_, err = io.Copy(f, r)
				if err != nil {
					data.SetInfo(ppath+`:`+err.Error(), 0)
				} else {
					data.SetInfo(ctx.T(`保存成功`), 1)
				}
				return ctx.JSON(data)
			}

			dat, err := ioutil.ReadAll(f)
			if err == nil && !isUTF8 {
				dat, err = charset.Convert(encoding, `utf-8`, dat)
			}
			if err != nil {
				data.SetInfo(err.Error(), 0)
			} else {
				data.SetData(string(dat), 1)
			}
		}
		return ctx.JSON(data)
	case `mkdir`:
		data := ctx.Data()
		newName := ctx.Form(`name`)
		if len(newName) == 0 {
			data.SetInfo(ctx.T(`请输入文件夹名`), 0)
		} else {
			dirPath := path.Join(ppath, newName)
			if f, err := mgr.Open(dirPath); err == nil {
				if finfo, err := f.Stat(); err != nil {
					data.SetError(err)
				} else if finfo.IsDir() {
					data.SetInfo(ctx.T(`已经存在相同名称的文件夹`), 0)
				} else {
					data.SetInfo(ctx.T(`已经存在相同名称的文件`), 0)
				}
			} else if !os.IsNotExist(err) {
				data.SetError(err)
			} else {
				err = mgr.Mkdir(dirPath)
				if err != nil {
					data.SetError(err)
				}
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
		var paths []string
		prefix := ctx.Form(`query`)
		num := ctx.Formx(`size`, `10`).Int()
		if num <= 0 {
			num = 10
		}
		dirs, _ := mgr.ReadDir(ppath)
		for _, d := range dirs {
			if len(paths) >= num {
				break
			}
			name := d.Name()
			if strings.HasPrefix(name, prefix) {
				paths = append(paths, name)
				continue
			}
		}
		data := ctx.Data().SetData(paths)
		return ctx.JSON(data)
	case `delete`:
		err = mgr.Remove(ppath)
		if err != nil {
			handler.SendFail(ctx, err.Error())
		}
		return ctx.Redirect(ctx.Referer())
	case `upload`:
		d, err := mgr.Open(ppath)
		if err != nil {
			return err
		}
		defer d.Close()
		fi, err := d.Stat()
		if !fi.IsDir() {
			return ctx.E(`路径不正确`)
		}
		fileSrc, fileHdr, err := ctx.Request().FormFile(`file`)
		if err != nil {
			return err
		}
		defer fileSrc.Close()

		// Destination
		fileName := fileHdr.Filename
		fileDst, err := mgr.Create(path.Join(ppath, fileName))
		if err != nil {
			return err
		}
		defer fileDst.Close()

		// Copy
		_, err = io.Copy(fileDst, fileSrc)
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
		d, err := mgr.Open(ppath)
		if err != nil {
			return err
		}
		defer d.Close()
		fi, err := d.Stat()
		if !fi.IsDir() {
			fileName := path.Base(ppath)
			return ctx.Attachment(d, fileName)
		}

		dirs, err := mgr.ReadDir(ppath)
		sortBy := ctx.Form(`sortBy`)
		switch sortBy {
		case `time`:
			sort.Sort(filemanager.SortByModTime(dirs))
		case `-time`:
			sort.Sort(filemanager.SortByModTimeDesc(dirs))
		case `name`:
		case `-name`:
			sort.Sort(filemanager.SortByNameDesc(dirs))
		case `type`:
			fallthrough
		default:
			sort.Sort(filemanager.SortByFileType(dirs))
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
	urlPrefix := fmt.Sprintf(`/term/sftp?id=%d&path=`, id) + encodedSep
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
	ctx.Set(`data`, m.SshUser)
	return ctx.Render(`term/sftp`, err)
}
