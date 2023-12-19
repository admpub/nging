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

package license

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/admpub/nging/v5/application/cmd/bootconfig"
	"github.com/admpub/nging/v5/application/library/notice"
	"github.com/admpub/nging/v5/application/library/selfupdate"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
)

type ProductVersion struct {
	Version          string `comment:"版本号(格式1.0.1)" json:"version" xml:"version"`
	Type             string `comment:"版本类型(stable-稳定版;beta-公测版;alpha-内测版)" json:"type" xml:"type"`
	OsArch           string `comment:"支持的操作系统(多个用逗号分隔)，留空表示不限制" json:"os_arch" xml:"os_arch"`
	ReleasedAt       uint   `comment:"发布时间" json:"released_at" xml:"released_at"`
	ForceUpgrade     string `comment:"是否强行升级为此版本" json:"force_upgrade" xml:"force_upgrade"`
	Description      string `comment:"发布说明" json:"description" xml:"description"`
	Remark           string `comment:"备注" json:"remark" xml:"remark"`
	DownloadURL      string `comment:"下载网址" json:"download_url" xml:"download_url"`
	Sign             string `comment:"下载后验证签名(多个签名之间用逗号分隔)" json:"sign" xml:"sign"`
	DownloadURLOther string `comment:"备用下载网址" json:"download_url_other" xml:"download_url_other"`
	DownloadedPath   string `comment:"自动下载保存路径" json:"-" xml:"-"`
	extractedDir     string
	backupDir        string
	executable       string
	isNew            bool
	prog             notice.NProgressor
}

func (v *ProductVersion) SetProgressor(prog notice.NProgressor) {
	prog.Reset()
	v.prog = prog
}

func (v *ProductVersion) IsNew() bool {
	return v.isNew
}

func (v *ProductVersion) clean(downloadDir string, newVersionDir string) {
	if len(newVersionDir) == 0 {
		return
	}
	dirEntries, _ := os.ReadDir(newVersionDir)
	downloadFolder := filepath.Base(downloadDir)
	for _, dirEntry := range dirEntries {
		if dirEntry.Name() == downloadFolder {
			continue
		}
		ppath := filepath.Join(newVersionDir, dirEntry.Name())
		os.RemoveAll(ppath)
		v.prog.Send(fmt.Sprintf(`clean up old files %q`, ppath), notice.StateSuccess)
	}
}

func (v *ProductVersion) Extract() error {
	if len(v.DownloadedPath) == 0 {
		return fmt.Errorf(`failed to download: %s`, v.DownloadURL)
	}
	downloadDir := filepath.Dir(v.DownloadedPath)
	v.backupDir = filepath.Join(downloadDir, `backup`)
	newVersionDir := filepath.Join(echo.Wd(), `data/cache/nging-new-version`)
	v.extractedDir = filepath.Join(newVersionDir, `latest`)
	v.clean(downloadDir, newVersionDir)
	com.MkdirAll(v.extractedDir, os.ModePerm)
	ddp := filepath.Join(v.extractedDir, `download_dir.txt`)
	if err := os.WriteFile(ddp, com.Str2bytes(downloadDir), 0666); err != nil {
		return fmt.Errorf(`%w: %s`, err, ddp)
	}
	v.prog.Send(fmt.Sprintf(`extract the file %q to %q`, v.DownloadedPath, v.extractedDir), notice.StateSuccess)
	_, err := com.UnTarGz(v.DownloadedPath, v.extractedDir)
	if err != nil {
		v.prog.Send(fmt.Sprintf(`failed to extract %q to %q`, v.DownloadedPath, v.extractedDir), notice.StateFailure)
		return fmt.Errorf(`%w: %s`, err, v.DownloadedPath)
	}
	subDir := strings.SplitN(filepath.Base(v.DownloadedPath), `.`, 2)[0]
	_extractedDir := filepath.Join(v.extractedDir, subDir)
	if com.FileExists(_extractedDir) {
		v.extractedDir = _extractedDir
	}
	files, err := filepath.Glob(v.extractedDir + echo.FilePathSeparator + `*.sha256`)
	if err != nil {
		return err
	}
	for _, sha256file := range files {
		executable := strings.TrimSuffix(sha256file, `.sha256`)
		if com.FileExists(executable) {
			os.Chmod(executable, 0755)
			v.executable = executable
		}
	}
	return err
}

func (v *ProductVersion) Upgrade(ctx echo.Context, ngingDir string, restartMode ...string) error {
	if len(v.extractedDir) == 0 {
		if err := v.Extract(); err != nil {
			return err
		}
	}
	executable := filepath.Base(v.executable)
	backupDir := v.backupDir
	com.MkdirAll(backupDir, 0755)
	v.prog.Send(fmt.Sprintf(`copy the files from %q to %q`, v.extractedDir, ngingDir), notice.StateSuccess)
	var backupFiles []string
	var extension string
	if com.IsWindows {
		extension = `.exe`
	}
	err := com.CopyDir(v.extractedDir, ngingDir, func(filePath string) bool {
		//fmt.Println(filePath)
		oldFile := filepath.Join(ngingDir, filePath)
		if fi, err := os.Stat(oldFile); err == nil {
			if fi.IsDir() {
				dir := filepath.Join(backupDir, filePath)
				if !com.FileExists(dir) {
					com.MkdirAll(dir, fi.Mode())
				}
			} else {
				if filePath == `startup`+extension {
					return true // 跳过此处复制。如果需要升级 startup，需要手动升级
				}
				backupFile := filepath.Join(backupDir, filePath)
				err = com.Copy(oldFile, backupFile)
				if err != nil {
					v.prog.Send(fmt.Sprintf(`failed to back up file %q: %v`, backupFile, err), notice.StateFailure)
				} else {
					v.prog.Send(fmt.Sprintf(`back up file %q to %q`, oldFile, backupFile), notice.StateSuccess)
				}
				backupFiles = append(backupFiles, backupFile)
			}
		}
		if executable == filePath {
			return true // 跳过此处复制，采用单独的替换逻辑来处理
		}
		return false
	})
	if err != nil {
		return err
	}
	restore := func() {
		for _, backupFile := range backupFiles {
			targetFile := strings.TrimPrefix(backupFile, backupDir)
			targetFile = filepath.Join(ngingDir, targetFile)
			err := com.Copy(backupFile, targetFile)
			if err != nil {
				v.prog.Send(fmt.Sprintf(`failed to restore file %q: %v`, targetFile, err), notice.StateFailure)
			} else {
				v.prog.Send(fmt.Sprintf(`restore file %q to %q`, backupFile, targetFile), notice.StateSuccess)
			}
		}
	}
	if len(v.executable) > 0 {
		var fp *os.File
		fp, err = os.Open(v.executable)
		if err != nil {
			return fmt.Errorf(`%w: %v`, err, v.executable)
		}
		defer fp.Close()
		targetExecutable := filepath.Join(ngingDir, executable)
		v.prog.Send(fmt.Sprintf(`update file %q`, targetExecutable), notice.StateSuccess)
		err = selfupdate.Update(fp, targetExecutable)
		if err != nil {
			v.prog.Send(fmt.Sprintf(`failed to update file %q: %v`, targetExecutable, err), notice.StateFailure)
			err = fmt.Errorf(`%w: %v`, err, targetExecutable)
			return err
		}
		v.prog.Send(fmt.Sprintf(`restart file %q`, targetExecutable), notice.StateSuccess)
		err = selfupdate.Restart(func(err error) {
			if err == nil {
				// v.prog.Send(`exit the current process`, notice.StateSuccess)
				// os.Exit(0)
				return
			}
			v.prog.Send(fmt.Sprintf(`failed to restart file %q: %v`, targetExecutable, err), notice.StateFailure)
			v.prog.Send(`start restoring files`, notice.StateSuccess)
			restore()
		}, targetExecutable, restartMode...)
		if err == nil {
			v.prog.Send(`successfully upgrade `+bootconfig.SoftwareName+` to version `+v.Version, notice.StateSuccess)
		}
	}
	v.prog.Complete()
	return err
}
