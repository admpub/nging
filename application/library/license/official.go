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

package license

import (
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/admpub/errors"
	"github.com/admpub/license_gen/lib"

	"github.com/webx-top/com"
	"github.com/webx-top/echo"

	"github.com/admpub/nging/v3/application/library/config"
	"github.com/admpub/nging/v3/application/library/restclient"
)

var (
	ErrConnectionFailed       = errors.New(`连接授权服务器失败`)
	ErrOfficialDataUnexcepted = errors.New(`官方数据返回异常`)
)

type OfficialData struct {
	lib.LicenseData
	Timestamp int64
}

type OfficialResp struct {
	Code int
	Info string
	Zone string        `json:",omitempty" xml:",omitempty"`
	Data *OfficialData `json:",omitempty" xml:",omitempty"`
}

type ValidResp struct {
	Code int
	Info string
	Zone string    `json:",omitempty" xml:",omitempty"`
	Data Validator `json:",omitempty" xml:",omitempty"`
}

type Validator interface {
	Validate() error
}

var NewValidResp = func() *ValidResp {
	return &ValidResp{Data: &ValidResult{}}
}

type ValidResult struct {
}

func (v *ValidResult) Validate() error {
	return nil
}

func validateFromOfficial(machineID string, ctx echo.Context) error {
	client := restclient.Resty()
	client.SetHeader("Accept", "application/json")
	result := NewValidResp()
	client.SetResult(result)
	fullURL := FullLicenseURL(machineID, ctx)
	response, err := client.Get(fullURL)
	if err != nil {
		return errors.Wrap(err, `Connection to the license server failed`)
	}
	if response == nil {
		return ErrConnectionFailed
	}
	switch response.StatusCode() {
	case http.StatusOK:
		if result.Code != 1 {
			return errors.New(result.Info)
		}
		if result.Data == nil {
			return ErrOfficialDataUnexcepted
		}
		return result.Data.Validate()
	case http.StatusNotFound:
		return ErrConnectionFailed
	default:
		return errors.New(response.Status())
	}
}

type VersionResp struct {
	Code int
	Info string
	Zone string          `json:",omitempty" xml:",omitempty"`
	Data *ProductVersion `json:",omitempty" xml:",omitempty"`
}

func latestVersion() error {
	client := restclient.Resty()
	client.SetHeader("Accept", "application/json")
	result := &VersionResp{}
	client.SetResult(result)
	response, err := client.Get(versionURL)
	if err != nil {
		return errors.Wrap(err, `Check for the latest version failed`)
	}
	if response == nil {
		return ErrConnectionFailed
	}
	switch response.StatusCode() {
	case http.StatusOK:
		if result.Code != 1 {
			return errors.New(result.Info)
		}
		if result.Data == nil {
			return ErrOfficialDataUnexcepted
		}
		hasNew := config.Version.IsNew(result.Data.Version, result.Data.Type)
		if hasNew {
			if result.Data.ForceUpgrade == `Y` {
				if len(result.Data.DownloadUrl) > 0 {
					//TODO: download
					saveTo := filepath.Join(echo.Wd(), `data/cache/nging-new-version`)
					err = com.MkdirAll(saveTo, os.ModePerm)
					if err != nil {
						return err
					}
					saveTo += echo.FilePathSeparator + path.Base(result.Data.DownloadUrl)
					err = com.RangeDownload(result.Data.DownloadUrl, saveTo)
					if err != nil {
						if len(result.Data.DownloadUrlOther) > 0 {
							err = com.RangeDownload(result.Data.DownloadUrlOther, saveTo)
						}
					}
					if err != nil {
						return err
					}
					//TODO: verify sign
					var signList []string
					if len(result.Data.Sign) > 0 {
						signList = strings.Split(result.Data.Sign, `,`)
					}
					if len(signList) > 0 {
						fileMd5 := com.Md5file(saveTo)
						var matched bool
						for _, sign := range signList {
							if sign == fileMd5 {
								matched = true
								break
							}
						}
						if !matched {
							return com.ErrMd5Unmatched
						}
					}
					//OK
				}
			}
		}
		return nil
	case http.StatusNotFound:
		return ErrConnectionFailed
	default:
		return errors.New(response.Status())
	}
}
