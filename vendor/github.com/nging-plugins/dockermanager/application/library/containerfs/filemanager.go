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

package containerfs

import (
	"archive/tar"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/admpub/log"
	"github.com/admpub/nging/v5/application/handler"
	"github.com/admpub/nging/v5/application/library/charset"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/nging-plugins/dockermanager/application/library/utils"
	uploadClient "github.com/webx-top/client/upload"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/param"
)

var (
	encodedSep   = com.URLEncode(`/`)
	encodedSlash = com.URLEncode(`\`)
)

func New(c *client.Client, containerID string, platform string, editableMaxSize int, ctx echo.Context) *fileManager {
	return &fileManager{
		Context:         ctx,
		client:          c,
		containerID:     containerID,
		platform:        platform,
		EditableMaxSize: editableMaxSize,
	}
}

type fileManager struct {
	echo.Context
	client          *client.Client
	containerID     string
	platform        string
	EditableMaxSize int
}

func (f *fileManager) RealPath(filePath string) string {
	if len(filePath) > 0 {
		filePath = filepath.Clean(filePath)
	}
	return filePath
}

func (f *fileManager) Seperator() string {
	switch f.platform {
	case `windows`:
		return `\`
	default:
		return `/`
	}
}

func (f *fileManager) URLEncodedSeperator() string {
	switch f.platform {
	case `windows`:
		return encodedSlash
	default:
		return encodedSep
	}
}

func (f *fileManager) RootDir() string {
	switch f.platform {
	case `windows`:
		return `c:\\`
	default:
		return `/`
	}
}

func (f *fileManager) Edit(absPath string, content string, encoding string) (interface{}, error) {
	fi, realPath, err := f.enterPath(absPath)
	if err != nil {
		return nil, err
	}
	if fi.Mode.IsDir() {
		return nil, errors.New(f.T(`不能编辑文件夹`))
	}
	if f.EditableMaxSize > 0 && fi.Size > int64(f.EditableMaxSize) {
		return nil, errors.New(f.T(`很抱歉，不支持编辑超过%v的文件`, com.FormatByte(f.EditableMaxSize, 2, true)))
	}
	absPath = realPath
	encoding = strings.ToLower(encoding)
	isUTF8 := encoding == `` || encoding == `utf-8`
	if f.IsPost() {
		b := []byte(content)
		if !isUTF8 {
			b, err = charset.Convert(`utf-8`, encoding, b)
			if err != nil {
				return ``, err
			}
		}
		var buffer bytes.Buffer
		writer := tar.NewWriter(&buffer)
		dir, name := filepath.Split(absPath)
		hdr := &tar.Header{
			Name:    name,
			Size:    int64(len(b)),
			Mode:    int64(fi.Mode),
			ModTime: time.Now(),
		}
		if err = writer.WriteHeader(hdr); err != nil {
			return nil, err
		}
		if _, err = writer.Write(b); err != nil {
			return nil, err
		}
		if err = writer.Close(); err != nil {
			return nil, err
		}
		opts := types.CopyToContainerOptions{}
		err = f.client.CopyToContainer(f, f.containerID, dir, &buffer, opts)
		return nil, err
	}
	var reader io.ReadCloser
	reader, _, err = f.client.CopyFromContainer(f, f.containerID, absPath)
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	tr := tar.NewReader(reader)
	if _, err = tr.Next(); err != nil {
		err = fmt.Errorf("Failed to find file in tar archive: %w", err)
		return nil, err
	}
	b, err := io.ReadAll(tr)
	if err == nil && !isUTF8 {
		b, err = charset.Convert(encoding, `utf-8`, b)
	}
	return string(b), err
}

func (f *fileManager) execLite(commands ...string) (err error) {
	var response *types.HijackedResponse
	response, err = f.exec(commands)
	if err != nil {
		return
	}
	defer response.Close()
	buf := bytes.NewBuffer(nil)
	_, err = utils.StdCopy(io.Discard, buf, response.Conn)
	if err != nil {
		return
	}
	errMsg := buf.String()
	if len(errMsg) > 0 {
		err = errors.New(errMsg)
	}
	return
}

func (f *fileManager) Remove(absPath string) error {
	fi, realPath, err := f.enterPath(absPath)
	if err != nil {
		return err
	}
	absPath = realPath
	commands := make([]string, 0, 1)
	if fi.Mode.IsDir() {
		commands = append(commands, `rm -rf `+strconv.Quote(absPath))
	} else {
		commands = append(commands, `rm -f `+strconv.Quote(absPath))
	}
	return f.execLite(commands...)
}

func (f *fileManager) Mkdir(absPath string, mode os.FileMode) error {
	return f.execLite(`mkdir -p ` + strconv.Quote(absPath))
}

func (f *fileManager) Rename(absPath string, newName string) (err error) {
	if len(newName) > 0 {
		commands := make([]string, 0, 1)
		newPath := filepath.Join(filepath.Dir(absPath), filepath.Base(newName))
		commands = append(commands, `mv `+strconv.Quote(absPath)+` `+strconv.Quote(newPath))
		err = f.execLite(commands...)
	} else {
		err = errors.New(f.T(`请输入有效的文件名称`))
	}
	return
}

func (f *fileManager) enterPath(absPath string) (stat types.ContainerPathStat, realPath string, err error) {
	if absPath != `/` && absPath != `\` {
		absPath = strings.TrimRight(absPath, `/`)
		absPath = strings.TrimRight(absPath, `\`)
	}
	stat, err = f.client.ContainerStatPath(f, f.containerID, absPath)
	if err != nil {
		return
	}
	if len(stat.LinkTarget) > 0 {
		realPath = stat.LinkTarget
		stat, err = f.client.ContainerStatPath(f, f.containerID, stat.LinkTarget)
	} else {
		realPath = absPath
	}
	return
}

func (f *fileManager) Upload(absPath string,
	chunkUpload *uploadClient.ChunkUpload,
	chunkOpts ...uploadClient.ChunkInfoOpter) (err error) {
	var fi types.ContainerPathStat
	var realPath string
	fi, realPath, err = f.enterPath(absPath)
	if err != nil {
		return
	}
	if !fi.Mode.IsDir() {
		return errors.New(f.T(`路径不正确: %s`, absPath))
	}
	absPath = realPath
	var filePath string
	var chunked bool // 是否支持分片
	if chunkUpload != nil {
		_, err := chunkUpload.Upload(f.Request().StdRequest(), chunkOpts...)
		if err != nil {
			if !errors.Is(err, uploadClient.ErrChunkUnsupported) {
				if errors.Is(err, uploadClient.ErrChunkUploadCompleted) ||
					errors.Is(err, uploadClient.ErrFileUploadCompleted) {
					return nil
				}
				return err
			}
		} else {
			if !chunkUpload.Merged() {
				return nil
			}
			chunked = true
			filePath = chunkUpload.GetSavePath()
		}
	}
	var filenames map[string]string
	if !chunked {
		user := handler.User(f.Context)
		srcPath := filepath.Join(os.TempDir(), `nging/docker/`+param.AsString(user.Id)+`/copy2container-`+time.Now().Format(`20060102150405.000`))
		err = com.MkdirAll(srcPath, os.ModePerm)
		if err != nil {
			return err
		}
		defer os.RemoveAll(srcPath)
		fileHdr, err := f.SaveUploadedFile(`file`, srcPath)
		if err != nil {
			return err
		}
		filePath = filepath.Join(absPath, fileHdr.Filename)
		filenames, err = utils.GetFilenames(srcPath)
		if err != nil {
			return err
		}
	}
	fileType := f.Form(`type`)
	switch fileType {
	case `tar`:
		var fp *os.File
		fp, err = os.Open(filePath)
		if err != nil {
			return err
		}
		defer fp.Close()
		opts := types.CopyToContainerOptions{}
		err = f.client.CopyToContainer(f, f.containerID, absPath, fp, opts)
		if err == nil {
			err := os.Remove(filePath)
			if err != nil {
				if !os.IsNotExist(err) {
					log.Error(err)
				}
			}
		}
		return err
	default:
		if chunked {
			filenames = map[string]string{
				filepath.Base(filePath): filePath,
			}
		}
	}
	prw, err := utils.CompressTar(f, filenames)
	if err != nil {
		return err
	}
	defer prw.Close()
	opts := types.CopyToContainerOptions{}
	err = prw.DoRead(func(r io.Reader) error {
		return f.client.CopyToContainer(f, f.containerID, absPath, r, opts)
	})
	return err
}

func (f *fileManager) Download(srcPath string) error {
	fi, realPath, err := f.enterPath(srcPath)
	if err != nil {
		return err
	}
	srcPath = realPath
	reader, pathStat, err := f.client.CopyFromContainer(f, f.containerID, srcPath)
	if err != nil {
		return err
	}
	defer reader.Close()
	if !fi.Mode.IsDir() {
		tr := tar.NewReader(reader)
		if _, err = tr.Next(); err != nil {
			err = fmt.Errorf("Failed to find file in tar archive: %w", err)
			return err
		}
		fileName := filepath.Base(srcPath)
		return f.Attachment(tr, fileName, pathStat.Mtime)
	}
	fileName := utils.ShortenID(f.containerID) + `-` + com.Md5(srcPath) + `.tar`
	echo.SetAttachmentHeader(f, fileName, false)
	return f.ServeContent(reader, fileName, pathStat.Mtime)
}

func (f *fileManager) exec(commands []string) (*types.HijackedResponse, error) {
	cfg := types.ExecConfig{
		Tty:          true,
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          []string{`sh`, `-c`, strings.Join(commands, ` && `)},
	}
	switch f.platform {
	case `windows`:
		return nil, fmt.Errorf(`Sorry, this operating system %q is not supported at this time`, f.platform)
		// cfg.Cmd[0] = `cmd.exe`
		// cfg.Cmd[1] = `/c`
		// TODO
	}
	idResp, err := f.client.ContainerExecCreate(f, f.containerID, cfg)
	if err != nil {
		return nil, err
	}
	var response types.HijackedResponse
	response, err = f.client.ContainerExecAttach(f, idResp.ID, types.ExecStartCheck{})
	//defer response.Close()
	return &response, err
}

func (f *fileManager) List(absPath string, sortBy ...string) (err error, exit bool, files []FileInfo) {
	var fi types.ContainerPathStat
	var realPath string
	fi, realPath, err = f.enterPath(absPath)
	if err != nil {
		return
	}
	absPath = realPath
	echo.Dump(echo.H{`isDir`: fi.Mode.IsDir(), `fi`: fi})
	if !fi.Mode.IsDir() {
		fileName := filepath.Base(absPath)
		inline := f.Formx(`inline`).Bool()
		var reader io.ReadCloser
		var pathStat types.ContainerPathStat
		reader, pathStat, err = f.client.CopyFromContainer(f, f.containerID, absPath)
		if err != nil {
			return
		}
		defer reader.Close()
		tr := tar.NewReader(reader)
		if _, err = tr.Next(); err != nil {
			err = fmt.Errorf("Failed to find file in tar archive: %w", err)
			return
		}
		return f.Attachment(tr, fileName, pathStat.Mtime, inline), true, nil
	}
	commands := make([]string, 0, 2)
	quotedPath := strconv.Quote(absPath)
	if absPath != `"."` {
		commands = append(commands, `cd `+quotedPath)
	}
	//commands = append(commands, `ls -al `+quotedPath)
	commands = append(commands, `ls -al --time-style long-iso `+quotedPath)
	var response *types.HijackedResponse
	response, err = f.exec(commands)
	if err != nil {
		return
	}
	defer response.Close()
	buf := bytes.NewBuffer(nil)
	_, err = utils.StdCopy(buf, buf, response.Conn)
	if err != nil {
		return
	}
	//var lines []string
	for {
		var message string
		message, err = buf.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				err = nil
				break
			}
			return
		}
		message = strings.TrimSuffix(message, com.StrLF)
		message = strings.TrimSuffix(message, com.StrR)

		// -rwxr-xr-x   1 root root    0 Nov 17 15:33 .dockerenv
		// lrwxrwxrwx   1 root root    7 Oct 30 00:00 bin -> usr/bin
		// drwxr-xr-x   1 root root    0 Sep 29 20:04 boot

		// --time-style long-iso
		// lrwxrwxrwx   1 root root    7 2023-10-30 00:00 bin -> usr/bin
		// drwxr-xr-x   1 root root    0 2023-09-29 20:04 boot
		// drwxr-xr-x   5 root root  320 2023-11-17 15:33 dev

		var perm, n, user, group, size, month, day, yearOrHourMinute, name string
		fields := strings.Fields(message)
		// if len(fields) < 9 {
		// 	continue
		// }
		// com.SliceExtract(fields, &perm, &n, &user, &group, &size, &month, &day, &yearOrHourMinute, &name)
		if len(fields) < 8 {
			continue
		}
		var dateStr, timeStr string
		com.SliceExtract(fields, &perm, &n, &user, &group, &size, &dateStr, &timeStr, &name)
		if name == `.` || name == `..` {
			continue
		}
		sizeN, _ := strconv.ParseInt(size, 10, 64)
		var modtime time.Time
		if len(dateStr) > 0 {
			t, _ := time.Parse(`2006-01-02 15:04`, dateStr+` `+timeStr)
			modtime = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, time.Local)
		} else {
			if len(day) == 1 {
				day = `0` + day
			}
			if strings.Contains(yearOrHourMinute, `:`) {
				t, _ := time.Parse(`Jan 02 15:04`, month+` `+day+` `+yearOrHourMinute)
				modtime = time.Date(time.Now().Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, time.Local)
			} else {
				t, _ := time.Parse(`Jan 02 2006`, month+` `+day+` `+yearOrHourMinute)
				modtime = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, time.Local)
			}
		}
		//lines = append(lines, message)
		files = append(files, FileInfo{
			FileMode: perm,
			Name:     name,
			User:     user,
			Group:    group,
			Size:     sizeN,
			ModTime:  modtime,
			IsDir:    strings.HasPrefix(perm, `d`),
		})
	}
	//echo.Dump(lines)
	if len(sortBy) > 0 {
		switch sortBy[0] {
		case `time`:
			sort.Sort(SortByModTime(files))
		case `-time`:
			sort.Sort(SortByModTimeDesc(files))
		case `name`:
		case `-name`:
			sort.Sort(SortByNameDesc(files))
		case `type`:
			fallthrough
		default:
			sort.Sort(SortByFileType(files))
		}
	} else {
		sort.Sort(SortByFileType(files))
	}
	if f.Format() == "json" {
		dirList, fileList := f.ListTransfer(files)
		data := f.Data()
		data.SetData(echo.H{
			`dirList`:  dirList,
			`fileList`: fileList,
		})
		return f.JSON(data), true, nil
	}
	return
}

func (f *fileManager) ListTransfer(dirs []FileInfo) (dirList []echo.H, fileList []echo.H) {
	dirList = []echo.H{}
	fileList = []echo.H{}
	for _, d := range dirs {
		item := echo.H{
			`name`:  d.Name,
			`size`:  d.Size,
			`mode`:  d.FileMode,
			`mtime`: d.ModTime.Format(`2006-01-02 15:04:05`),
		}
		if d.IsDir {
			dirList = append(dirList, item)
			continue
		}
		fileList = append(fileList, item)
	}
	return
}
