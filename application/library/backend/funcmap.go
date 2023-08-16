package backend

import (
	"sync"
	"time"

	"github.com/admpub/timeago"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/middleware/tplfunc"
	"github.com/webx-top/echo/param"
	"github.com/webx-top/echo/subdomains"

	"github.com/admpub/nging/v5/application/handler"
	"github.com/admpub/nging/v5/application/library/codec"
	"github.com/admpub/nging/v5/application/library/common"
	"github.com/admpub/nging/v5/application/library/config"
	"github.com/admpub/nging/v5/application/library/config/cmder"
	"github.com/admpub/nging/v5/application/library/license"
	uploadLibrary "github.com/admpub/nging/v5/application/library/upload"
	"github.com/admpub/nging/v5/application/registry/navigate"
	"github.com/admpub/nging/v5/application/registry/upload/checker"
)

var (
	DefaultAvatarURL = `/public/assets/backend/images/user_128.png`
	AssetsURLPath    = `/public/assets/backend`
	tplFuncMap       map[string]interface{}
	tplOnce          sync.Once
)

func initTplFuncMap() {
	tplFuncMap = addGlobalFuncMap(tplfunc.New())
}

func GlobalFuncMap() map[string]interface{} {
	tplOnce.Do(initTplFuncMap)
	return tplFuncMap
}

func init() {
	timeago.Set(`language`, `zh-cn`)
	tplfunc.TplFuncMap[`Languages`] = languages
	tplfunc.TplFuncMap[`URLFor`] = subdomains.Default.URL
	tplfunc.TplFuncMap[`URLByName`] = subdomains.Default.URLByName
	tplfunc.TplFuncMap[`BackendURLByName`] = getBackendURLByName
	tplfunc.TplFuncMap[`FrontendURLByName`] = getFrontendURLByName
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
	tplfunc.TplFuncMap[`TemplateTags`] = common.TemplateTags
	tplfunc.TplFuncMap[`CmdIsRunning`] = cmdIsRunning
	tplfunc.TplFuncMap[`CmdHasGroup`] = cmdHasGroup
	tplfunc.TplFuncMap[`CmdExists`] = cmdExists
	tplfunc.TplFuncMap[`HasService`] = hasService
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

func addGlobalFuncMap(fm map[string]interface{}) map[string]interface{} {
	fm[`AssetsURL`] = getAssetsURL
	fm[`BackendURL`] = getBackendURL
	fm[`FrontendURL`] = getFrontendURL
	fm[`Project`] = navigate.ProjectGet
	fm[`ProjectSearchIdent`] = navigate.ProjectSearchIdent
	fm[`Projects`] = navigate.ProjectListAll
	return fm
}

func getAssetsURL(paths ...string) (r string) {
	r = AssetsURLPath
	for _, ppath := range paths {
		r += ppath
	}
	return r
}

func getBackendURL(paths ...string) (r string) {
	r = handler.BackendPrefix
	for _, ppath := range paths {
		r += ppath
	}
	return r
	//return subdomains.Default.URL(r, `backend`)
}

func getFrontendURL(paths ...string) (r string) {
	r = handler.FrontendPrefix
	for _, ppath := range paths {
		r += ppath
	}
	return subdomains.Default.URL(r, `frontend`)
}

func getBackendURLByName(name string, params ...interface{}) string {
	info := subdomains.Default.Get(`backend`)
	if info == nil {
		return `/not-found:` + name
	}
	return info.URLByName(subdomains.Default, name, params...)
}

func getFrontendURLByName(name string, params ...interface{}) string {
	info := subdomains.Default.Get(`frontend`)
	if info == nil {
		return `/not-found:` + name
	}
	return info.URLByName(subdomains.Default, name, params...)
}

func cmdIsRunning(name string) bool {
	return config.FromCLI().IsRunning(name)
}

func cmdHasGroup(group string) bool {
	return config.FromCLI().CmdHasGroup(group)
}

func cmdExists(name string) bool {
	return config.FromCLI().CmdGet(name) != nil
}

func hasService(name string) bool {
	return cmder.Has(name)
}
