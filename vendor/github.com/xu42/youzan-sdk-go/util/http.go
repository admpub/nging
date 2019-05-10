package util

import (
	"fmt"
	"net/url"
	"strings"
)

// URLAPIBase API URL
const URLAPIBase string = "https://open.youzan.com/api/oauthentry/%s/%s/%s"

// BuildPostParams 组装HTTP POST参数
func BuildPostParams(data map[string]string) url.Values {

	params := make(url.Values)

	for key, value := range data {
		params.Set(key, value)
	}

	return params
}

// BuildURL 组装接口URL
func BuildURL(apiName, apiVersion string) (url string) {

	sl := strings.Split(apiName, ".")
	url = fmt.Sprintf(URLAPIBase, strings.Join(sl[0:len(sl)-1], "."), apiVersion, sl[len(sl)-1])

	return
}
