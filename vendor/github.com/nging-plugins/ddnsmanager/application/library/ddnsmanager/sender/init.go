package sender

import "github.com/admpub/nging/v5/application/registry/alert"

func init() {
	alert.Topics.Add(`ddnsUpdate`, `DDNS更新`)
}
