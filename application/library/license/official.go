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
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/admpub/errors"
	godl "github.com/admpub/go-download/v2"

	"github.com/webx-top/com"
	"github.com/webx-top/echo"

	"github.com/admpub/nging/v5/application/handler"
	"github.com/admpub/nging/v5/application/library/config"
	"github.com/admpub/nging/v5/application/library/notice"
	"github.com/admpub/nging/v5/application/library/restclient"
)

var (
	ErrConnectionFailed       = errors.New(`连接授权服务器失败`)
	ErrOfficialDataUnexcepted = errors.New(`官方数据返回异常`)
	ErrLicenseDownloadFailed  = errors.New(`下载证书失败：官方数据返回异常`)
	ErrChecksumUnmatched      = errors.New("WARNING: Checksum don't match")
	ErrNoDownloadURL          = errors.New("暂无下载地址，请稍后再试")
)

type OfficialData struct {
	License   string
	Timestamp int64
}

type OfficialResponse struct {
	Code int
	Info string
	Zone string        `json:",omitempty" xml:",omitempty"`
	Data *OfficialData `json:",omitempty" xml:",omitempty"`
}

type ValidateResponse struct {
	Code int
	Info string
	Zone string    `json:",omitempty" xml:",omitempty"`
	Data Validator `json:",omitempty" xml:",omitempty"`
}

type Validator interface {
	Validate() error
}

var NewValidateResponse = func() *ValidateResponse {
	return &ValidateResponse{
		Data: ValidateResultInitor(),
	}
}

var ValidateResultInitor = func() Validator {
	return &ValidateResult{}
}

type ValidateResult struct {
}

func (v *ValidateResult) Validate() error {
	return nil
}

func validateFromOfficial(ctx echo.Context) error {
	client := restclient.RestyRetryable()
	client.SetHeader("Accept", "application/json")
	result := NewValidateResponse()
	client.SetResult(result)
	fullURL := FullLicenseURL(ctx)
	response, err := client.Get(fullURL)
	if err != nil {
		if strings.Contains(err.Error(), `connection refused`) {
			return ErrConnectionFailed
		}
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

type VersionResponse struct {
	Code int
	Info string
	Zone string          `json:",omitempty" xml:",omitempty"`
	Data *ProductVersion `json:",omitempty" xml:",omitempty"`
}

func LatestVersion(ctx echo.Context, version string, download bool) (*ProductVersion, error) {
	client := restclient.RestyRetryable()
	client.SetHeader("Accept", "application/json")
	result := &VersionResponse{}
	client.SetResult(result)
	surl := versionURL + `?` + URLValues(ctx, version).Encode()
	response, err := client.Get(surl)
	if err != nil {
		return nil, errors.Wrap(err, `check for the latest version failed`)
	}
	if response == nil {
		return nil, ErrConnectionFailed
	}
	switch response.StatusCode() {
	case http.StatusOK:
		if result.Code != 1 {
			return nil, errors.New(result.Info)
		}
		if result.Data == nil {
			return nil, ErrOfficialDataUnexcepted
		}
		var username string
		user := handler.User(ctx)
		if user != nil {
			username = user.Username
		}
		np := notice.NewP(ctx, `ngingDownloadNewVersion`, username, context.Background()).AutoComplete(true)
		defer result.Data.SetProgressor(np)
		result.Data.isNew = config.Version.IsNew(result.Data.Version, result.Data.Type)
		if !result.Data.isNew {
			np.Send(`no new version`, notice.StateSuccess)
			return result.Data, nil
		}
		np.Send(`new version found: v`+result.Data.Version, notice.StateSuccess)
		if !download {
			return result.Data, nil
		}
		if len(result.Data.DownloadURL) == 0 {
			return result.Data, ErrNoDownloadURL
		}
		np.Send(`automatically download the new version v`+result.Data.Version, notice.StateSuccess)
		saveTo := filepath.Join(echo.Wd(), `data/cache/nging-new-version`, result.Data.Version)
		err = com.MkdirAll(saveTo, os.ModePerm)
		if err != nil {
			return result.Data, err
		}
		saveTo += echo.FilePathSeparator + path.Base(result.Data.DownloadURL)
		result.Data.DownloadedPath = saveTo
		if com.FileExists(saveTo) {
			np.Send(`the file already exists: `+saveTo, notice.StateSuccess)
			return result.Data, nil
		}
		np.Send(`downloading `+result.Data.DownloadURL+` => `+saveTo, notice.StateSuccess)
		dlCfg := &godl.Options{
			Proxy: func(name string, download int, size int64, r io.Reader) io.Reader {
				np.Add(size)
				np.Send(`downloading `+name, notice.StateSuccess)
				return np.ProxyReader(r)
			},
		}
		_, err = godl.Download(result.Data.DownloadURL, saveTo, dlCfg)
		if err != nil {
			if len(result.Data.DownloadURLOther) > 0 {
				np.Send(`try to download from the mirror URL `+result.Data.DownloadURLOther, notice.StateSuccess)
				_, err = godl.Download(result.Data.DownloadURLOther, saveTo, dlCfg)
			}
		}
		np.Reset()
		if err != nil {
			np.Send(err.Error(), notice.StateFailure)
			np.Complete()
			return result.Data, err
		}
		np.Complete()
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
				return result.Data, com.ErrMd5Unmatched
			}
		} else {
			if resp, err := restclient.RestyRetryable().Get(result.Data.DownloadURL + `.sha256`); err == nil && resp.IsSuccess() {
				expectedSHA := resp.String()
				expectedSHA = strings.TrimSpace(expectedSHA)
				err = VerifyChecksum(saveTo, expectedSHA)
				if err != nil {
					return result.Data, err
				}
			}
		}
		//OK
		return result.Data, nil
	case http.StatusNotFound:
		return nil, fmt.Errorf(`%w: %s`, ErrConnectionFailed, surl)
	default:
		return nil, errors.New(response.Status())
	}
}

func VerifyChecksum(file string, expected string) error {
	f, err := os.OpenFile(file, os.O_RDONLY, 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	copyBuf := make([]byte, 1024*1024)

	h := sha256.New()
	_, err = io.CopyBuffer(h, f, copyBuf)
	if err != nil {
		return err
	}

	sha256Result := hex.EncodeToString(h.Sum(nil))
	if sha256Result == strings.SplitN(expected, ` `, 2)[0] {
		err = ErrChecksumUnmatched
	}
	return err
}
