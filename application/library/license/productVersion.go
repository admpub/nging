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

	"github.com/admpub/log"
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
	executable       string
	isNew            bool
}

func (v *ProductVersion) IsNew() bool {
	return v.isNew
}

func (v *ProductVersion) Extract() error {
	if len(v.DownloadedPath) == 0 {
		return fmt.Errorf(`failed to download: %s`, v.DownloadURL)
	}
	v.extractedDir = filepath.Join(filepath.Dir(v.DownloadedPath), `extracted`)
	_, err := com.UnTarGz(v.DownloadedPath, v.extractedDir)
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

func (v *ProductVersion) Upgrade(ctx echo.Context, ngingDir string) error {
	if len(v.extractedDir) == 0 {
		if err := v.Extract(); err != nil {
			return err
		}
	}
	executable := filepath.Base(v.executable)
	backupDir := filepath.Join(filepath.Dir(v.extractedDir), `backup`)
	com.MkdirAll(backupDir, 0755)
	var backupFiles []string
	err := com.CopyDir(v.extractedDir, ngingDir, func(filePath string) bool {
		fmt.Println(filePath)
		oldFile := filepath.Join(ngingDir, filePath)
		if fi, err := os.Stat(oldFile); err == nil {
			if fi.IsDir() {
				dir := filepath.Join(backupDir, filePath)
				if !com.FileExists(dir) {
					com.MkdirAll(dir, fi.Mode())
				}
			} else {
				backupFile := filepath.Join(backupDir, filePath)
				err = com.Copy(oldFile, backupFile)
				if err != nil {
					log.Errorf(`failed to backup file %q: %v`, backupFile, err)
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
			err = com.Copy(backupFile, targetFile)
			if err != nil {
				log.Errorf(`failed to restore file %q: %v`, targetFile, err)
			}
		}
	}
	if len(v.executable) > 0 {
		fp, err := os.Open(v.executable)
		if err != nil {
			return fmt.Errorf(`%w: %v`, err, v.executable)
		}
		defer fp.Close()
		targetExecutable := filepath.Join(ngingDir, executable)
		err = selfupdate.Update(fp, targetExecutable)
		if err != nil {
			err = fmt.Errorf(`%w: %v`, err, targetExecutable)
			return err
		}
		err = selfupdate.Restart(func(err error) {
			if err == nil {
				return
			}
			restore()
		}, targetExecutable)
	}
	return err
}
