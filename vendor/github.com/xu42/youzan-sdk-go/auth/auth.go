package auth

import (
	"encoding/json"
	"github.com/xu42/youzan-sdk-go/util"
	"io/ioutil"
	"net/http"
)

// URLOauthToken 认证Token
const URLOauthToken string = "https://open.youzan.com/oauth/token"

// GenSelfToken 获取自用型AccessToken
func GenSelfToken(request GenSelfTokenRequest) (response GenSelfTokenResponse, err error) {
	body, err := get(request.toMap())
	err = json.Unmarshal(body, &response)
	return
}

// GenToolToken 生成工具型Token
func GenToolToken(request GenToolTokenRequest) (response GenToolTokenResponse, err error) {
	body, err := get(request.toMap())
	err = json.Unmarshal(body, &response)
	return
}

// get get
func get(data map[string]string) (body []byte, err error) {

	resp, err := http.DefaultClient.PostForm(URLOauthToken, util.BuildPostParams(data))
	if err != nil {
		return
	}

	body, err = ioutil.ReadAll(resp.Body)
	return
}

// GenTokenBaseRequest 获取AccessToken基本请求参数结构体
type GenTokenBaseRequest struct {
	GrantType    string `json:"grant_type"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

// GenToolTokenRequest  获取工具型AccessToken请求参数结构体
type GenToolTokenRequest struct {
	GenTokenBaseRequest
	Code        string `json:"code"`
	RedirectURI string `json:"redirect_uri"`
}

// GenSelfTokenRequest 获取自用型AccessToken请求参数结构体
type GenSelfTokenRequest struct {
	GenTokenBaseRequest
	KdtID string `json:"kdt_id"`
}

// GenToolTokenResponse 获取工具型AccessToken响应参数结构体
type GenToolTokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
}

// GenSelfTokenResponse 获取自用型AccessToken响应参数结构体
type GenSelfTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
}

func (req *GenTokenBaseRequest) toMap(grantType string) (m map[string]string) {
	m = make(map[string]string)
	m["client_secret"] = req.ClientSecret
	m["client_id"] = req.ClientID
	m["grant_type"] = grantType
	return
}

func (req *GenSelfTokenRequest) toMap() (m map[string]string) {
	m = make(map[string]string)
	m = req.GenTokenBaseRequest.toMap("silent")
	m["kdt_id"] = req.KdtID
	return
}

func (req *GenToolTokenRequest) toMap() (m map[string]string) {
	m = make(map[string]string)
	m = req.GenTokenBaseRequest.toMap("authorization_code")
	m["code"] = req.Code
	m["redirect_uri"] = req.RedirectURI
	return
}
