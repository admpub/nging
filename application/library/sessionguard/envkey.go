package sessionguard

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/admpub/nging/v5/application/library/common"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
)

func EnvKey(ctx echo.Context) string {
	return com.Md5(ctx.RealIP()+`|`+ctx.Request().UserAgent()) + `T` + strconv.FormatInt(time.Now().Unix(), 10)
}

func VerifyEnvKey(ctx echo.Context, envKey string, dontCheckTime ...bool) error {
	pos := strings.LastIndex(envKey, `T`)
	end := len(envKey) - 1
	if pos < 1 || pos >= end {
		return ctx.NewError(code.InvalidParameter, `凭证无效`).SetZone(`envKey`)
	}
	if len(dontCheckTime) == 0 || !dontCheckTime[0] {
		ts, err := strconv.ParseInt(envKey[pos+1:], 10, 64)
		if err != nil {
			return ctx.NewError(code.DataFormatIncorrect, `凭证解析失败`).SetZone(`envKey`)
		}
		if time.Now().Unix()-ts > 300 {
			return ctx.NewError(code.DataHasExpired, `凭证已经过期，请刷新页面后重新操作`).SetZone(`envKey`)
		}
	}
	if com.Md5(ctx.RealIP()+`|`+ctx.Request().UserAgent()) != envKey[0:pos] {
		return ctx.NewError(code.DataStatusIncorrect, `凭证来源不符合要求`).SetZone(`envKey`)
	}
	return nil
}

type PE struct {
	Password  string `json:"p"`
	EnvKey    string `json:"e"`
	Timestamp int64  `json:"t"`
}

func (p *PE) Verify(ctx echo.Context) error {
	sessionGuardCfg := common.ExtendConfig().Children(`sessionGuard`)
	encryptedPasswordExpires := sessionGuardCfg.Int64(`encryptedPasswordExpires`)
	if encryptedPasswordExpires <= 0 {
		encryptedPasswordExpires = 300
	}
	//fmt.Println(`encryptedPasswordExpires~~~~~~~~~~~~~~~~~~~~~~~~~>`, encryptedPasswordExpires)
	if time.Now().Unix()-p.Timestamp > encryptedPasswordExpires {
		return ctx.NewError(code.DataHasExpired, `凭证已经失效`).SetZone(`envKey`)
	}
	return VerifyEnvKey(ctx, p.EnvKey, true)
}

func Unpack(ctx echo.Context, encrypted string) (password string, err error) {
	if !strings.HasPrefix(encrypted, `{`) {
		return encrypted, nil
	}
	pe := &PE{}
	jerr := json.Unmarshal(com.Str2bytes(encrypted), pe)
	if jerr != nil {
		err = ctx.NewError(code.DataFormatIncorrect, `密码拆包失败`)
		return
	}
	return pe.Password, pe.Verify(ctx)
}
