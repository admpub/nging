package provider

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// UnmarshalHTTPResponse 处理HTTP结果，返回序列化的json
func UnmarshalHTTPResponse(resp *http.Response, url string, err error, result interface{}) error {
	body, err := GetHTTPResponse(resp, url, err)
	if err != nil {
		return err
	}

	// log.Println(string(body))
	err = json.Unmarshal(body, &result)
	if err != nil {
		err = fmt.Errorf("请求接口%s解析json结果失败! ERROR: %s", url, err)
	}

	return err
}

var MaxReadBodySize = int64(2 << 20) // 2M

// GetHTTPResponse 处理HTTP结果，返回byte
func GetHTTPResponse(resp *http.Response, url string, err error) ([]byte, error) {
	if err != nil {
		log.Printf("请求接口%s失败! ERROR: %s\n", url, err)
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, MaxReadBodySize))

	if err != nil {
		log.Printf("请求接口%s失败! ERROR: %s\n", url, err)
	}

	// 300及以上状态码都算异常
	if resp.StatusCode >= 300 {
		err = fmt.Errorf("请求接口 %s 失败! 返回内容: %s , 返回状态码: %d", url, string(body), resp.StatusCode)
	}

	return body, err
}
