package middleware

import (
	"sync"
	"time"

	"github.com/admpub/timeago"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/middleware/tplfunc"
	"github.com/webx-top/echo/param"
	"github.com/webx-top/echo/subdomains"

	"github.com/admpub/nging/v5/application/library/codec"
	"github.com/admpub/nging/v5/application/library/common"
	"github.com/admpub/nging/v5/application/library/config"
	"github.com/admpub/nging/v5/application/library/license"
	uploadLibrary "github.com/admpub/nging/v5/application/library/upload"
	"github.com/admpub/nging/v5/application/registry/upload/checker"
)

var (
	tplFuncMap map[string]interface{}
	tplOnce    sync.Once
)

func initTplFuncMap() {
	tplFuncMap = tplfunc.New()
}

func TplFuncMap() map[string]interface{} {
	tplOnce.Do(initTplFuncMap)
	return tplFuncMap
}

func init() {
	timeago.Set(`language`, `zh-cn`)
	tplfunc.TplFuncMap[`Languages`] = languages
	tplfunc.TplFuncMap[`URLFor`] = subdomains.Default.URL
	tplfunc.TplFuncMap[`URLByName`] = subdomains.Default.URLByName
	tplfunc.TplFuncMap[`IsMessage`] = common.IsMessage
	tplfunc.TplFuncMap[`IsError`] = common.IsError
	tplfunc.TplFuncMap[`IsOk`] = common.IsOk
	tplfunc.TplFuncMap[`Message`] = common.Message
	tplfunc.TplFuncMap[`Ok`] = common.OkString
	tplfunc.TplFuncMap[`Version`] = func() *config.VersionInfo { return config.Version }
	tplfunc.TplFuncMap[`VersionNumber`] = func() string { return config.Version.Number }
	tplfunc.TplFuncMap[`CommitID`] = func() string { return config.Version.CommitID }
	tplfunc.TplFuncMap[`BuildTime`] = func() string { return config.Version.BuildTime }
	tplfunc.TplFuncMap[`TrackerURL`] = license.TrackerURL
	tplfunc.TplFuncMap[`TrackerHTML`] = license.TrackerHTML
	tplfunc.TplFuncMap[`Config`] = getConfig
	tplfunc.TplFuncMap[`WithURLParams`] = common.WithURLParams
	tplfunc.TplFuncMap[`FullURL`] = common.FullURL
	tplfunc.TplFuncMap[`MaxRequestBodySize`] = getMaxRequestBodySize
	tplfunc.TplFuncMap[`IndexStrSlice`] = indexStrSlice
	tplfunc.TplFuncMap[`HasString`] = hasString
	tplfunc.TplFuncMap[`Date`] = date
	tplfunc.TplFuncMap[`Token`] = checker.Token
	tplfunc.TplFuncMap[`BackendUploadURL`] = checker.BackendUploadURL
	tplfunc.TplFuncMap[`FrontendUploadURL`] = checker.FrontendUploadURL
	tplfunc.TplFuncMap[`Avatar`] = getAvatar
	tplfunc.TplFuncMap[`SM2PublicKey`] = codec.DefaultPublicKeyHex
	tplfunc.TplFuncMap[`FileTypeByName`] = uploadLibrary.FileTypeByName
	tplfunc.TplFuncMap[`FileTypeIcon`] = getFileTypeIcon
}

func getFileTypeIcon(typ string) string {
	return uploadLibrary.Get().FileIcon(typ)
}

func languages() []string {
	return config.FromFile().Language.AllList
}

func getConfig(args ...string) echo.H {
	if len(args) > 0 {
		return config.Setting(args...)
	}
	return config.Setting()
}

func getMaxRequestBodySize() int {
	return config.FromFile().GetMaxRequestBodySize()
}

func getAvatar(avatar string, defaults ...string) string {
	if len(avatar) > 0 {
		return tplfunc.AddSuffix(avatar, `_200_200`)
	}
	if len(defaults) > 0 && len(defaults[0]) > 0 {
		return defaults[0]
	}
	return DefaultAvatarURL
}

func indexStrSlice(slice []string, index int) string {
	if slice == nil {
		return ``
	}
	if index >= len(slice) {
		return ``
	}
	return slice[index]
}

func hasString(slice []string, str string) bool {
	if slice == nil {
		return false
	}
	for _, v := range slice {
		if v == str {
			return true
		}
	}
	return false
}

func date(timestamp interface{}) time.Time {
	v := param.AsInt64(timestamp)
	return time.Unix(v, 0)
}
