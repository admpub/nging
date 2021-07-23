package alert

import (
	"github.com/admpub/nging/v3/application/library/imbot"
	_ "github.com/admpub/nging/v3/application/library/imbot/dingding"
	_ "github.com/admpub/nging/v3/application/library/imbot/workwx"
	"github.com/webx-top/echo"
)

var (
	// RecipientTypes 收信类型
	RecipientTypes = echo.NewKVData()

	// RecipientPlatforms 收信平台
	RecipientPlatforms = echo.NewKVData()

	// Topics 告警专题
	Topics = echo.NewKVData()
)

func init() {
	RecipientTypes.Add(`email`, `email`)
	RecipientTypes.Add(`webhook`, `webhook`)
	for name, mess := range imbot.Messagers() {
		RecipientPlatforms.Add(name, mess.Label)
	}

	//Topics.Add(`test`, `测试`)
}
