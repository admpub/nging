package sender

import "github.com/admpub/nging/v4/application/registry/alert"

func init() {
	alert.Topics.Add(`ddnsUpdate`, `DDNS更新`)
}
