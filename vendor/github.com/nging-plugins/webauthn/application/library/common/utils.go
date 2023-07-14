package common

import (
	ua "github.com/admpub/useragent"
)

func GetOS(userAgent string) string {
	infoUA := ua.Parse(userAgent)
	return infoUA.OS
}
