package dnspod

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/nging-plugins/ddnsmanager/application/library/ddnsmanager/domain/dnsdomain"
	"github.com/nging-plugins/ddnsmanager/application/library/ddnsmanager/provider"
	"github.com/webx-top/echo"
)

const (
	recordListAPI         = "https://dnsapi.cn/Record.List"
	recordModifyURL       = "https://dnsapi.cn/Record.Modify"
	recordCreateAPI       = "https://dnsapi.cn/Record.Create"
	signUpURL             = `https://console.dnspod.cn/account/token`
	docLineType           = `https://docs.dnspod.cn/api/5f5623f9e75cf42d25bf6776/`
	defaultTTL            = 600
	defaultTimeout  int64 = 10
)

// https://docs.dnspod.cn/api/5f562ae4e75cf42d25bf689e/
// Dnspod 腾讯云dns实现
type Dnspod struct {
	clientID     string
	clientSecret string
	Domains      []*dnsdomain.Domain
	TTL          int
	client       *http.Client
}

// DnspodRecordListResp recordListAPI结果
type DnspodRecordListResp struct {
	DnspodStatus
	Records []struct {
		ID      string
		Name    string
		Type    string
		Value   string
		Enabled string
	}
}

// DnspodStatus DnspodStatus
type DnspodStatus struct {
	Status struct {
		Code    string
		Message string
	}
}

const Name = `DNSPod`

func (*Dnspod) Name() string {
	return Name
}

func (*Dnspod) Description() string {
	return ``
}

func (*Dnspod) SignUpURL() string {
	return signUpURL
}

func (*Dnspod) LineTypeURL() string {
	return docLineType
}

var configItems = echo.KVList{
	echo.NewKV(`clientId`, `ID`).SetHKV(`inputType`, `text`).SetHKV(`required`, true),
	echo.NewKV(`clientSecret`, `Token`).SetHKV(`inputType`, `text`).SetHKV(`required`, true),
	echo.NewKV(`timeout`, `接口超时(秒)`).SetHKV(`inputType`, `number`).SetX(defaultTimeout),
	echo.NewKV(`ttl`, `域名TTL`).SetHKV(`inputType`, `number`).SetX(defaultTTL),
}

func (*Dnspod) ConfigItems() echo.KVList {
	return configItems
}

var support = dnsdomain.Support{
	A:    true,
	AAAA: true,
	Line: true,
}

func (*Dnspod) Support() dnsdomain.Support {
	return support
}

// Init 初始化
func (dnspod *Dnspod) Init(settings echo.H, domains []*dnsdomain.Domain) error {
	dnspod.TTL = settings.Int(`ttl`)
	dnspod.clientID = settings.String(`clientId`)
	dnspod.clientSecret = settings.String(`clientSecret`)
	if dnspod.TTL < 1 {
		dnspod.TTL = defaultTTL
	}
	dnspod.Domains = domains
	timeout := settings.Int64(`timeout`)
	if timeout < 1 {
		timeout = defaultTimeout
	}
	dnspod.client = &http.Client{Timeout: time.Duration(timeout) * time.Second}
	return nil
}

func (dnspod *Dnspod) Update(ctx context.Context, recordType string, ipAddr string) error {

	for _, domain := range dnspod.Domains {
		result, err := dnspod.getRecordList(domain, recordType)
		if err != nil {
			domain.UpdateStatus = dnsdomain.UpdatedFailed
			return err
		}

		if len(result.Records) > 0 { // 更新
			err = dnspod.modify(ctx, result, domain, recordType, ipAddr)
		} else { // 新增
			err = dnspod.create(ctx, result, domain, recordType, ipAddr)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

// 创建
func (dnspod *Dnspod) create(ctx context.Context, result DnspodRecordListResp, domain *dnsdomain.Domain, recordType string, ipAddr string) error {
	ipAddr = domain.IP(ipAddr)
	values := url.Values{
		"login_token": {dnspod.clientID + "," + dnspod.clientSecret},
		"domain":      {domain.DomainName},
		"sub_domain":  {domain.GetSubDomain()},
		"record_type": {recordType},
		"record_line": {"默认"},
		"value":       {ipAddr},
		"ttl":         {strconv.Itoa(dnspod.TTL)},
		"format":      {"json"},
	}
	if len(domain.Line) > 0 {
		values.Set(`record_line`, domain.Line)
	}
	status, err := dnspod.commonRequest(
		ctx,
		recordCreateAPI,
		values,
		domain,
	)
	if err == nil && status.Status.Code == "1" {
		log.Printf("新增域名解析 %s 成功！IP: %s", domain, ipAddr)
		domain.UpdateStatus = dnsdomain.UpdatedSuccess
	} else {
		log.Printf("新增域名解析 %s 失败！Code: %s, Message: %s", domain, status.Status.Code, status.Status.Message)
		domain.UpdateStatus = dnsdomain.UpdatedFailed
	}
	return err
}

// 修改
func (dnspod *Dnspod) modify(ctx context.Context, result DnspodRecordListResp, domain *dnsdomain.Domain, recordType string, ipAddr string) error {
	ipAddr = domain.IP(ipAddr)
	for _, record := range result.Records {
		// 相同不修改
		if record.Value == ipAddr {
			log.Printf("你的IP %s 没有变化, 域名 %s", ipAddr, domain)
			domain.UpdateStatus = dnsdomain.UpdatedNothing
			continue
		}
		values := url.Values{
			"login_token": {dnspod.clientID + "," + dnspod.clientSecret},
			"domain":      {domain.DomainName},
			"sub_domain":  {domain.GetSubDomain()},
			"record_type": {recordType},
			"record_line": {"默认"},
			"record_id":   {record.ID},
			"value":       {ipAddr},
			"ttl":         {strconv.Itoa(dnspod.TTL)},
			"format":      {"json"},
		}
		if len(domain.Line) > 0 {
			values.Set(`record_line`, domain.Line)
		}
		status, err := dnspod.commonRequest(
			ctx,
			recordModifyURL,
			values,
			domain,
		)
		if err == nil && status.Status.Code == "1" {
			log.Printf("更新域名解析 %s 成功！IP: %s", domain, ipAddr)
			domain.UpdateStatus = dnsdomain.UpdatedSuccess
		} else {
			log.Printf("更新域名解析 %s 失败！Code: %s, Message: %s", domain, status.Status.Code, status.Status.Message)
			domain.UpdateStatus = dnsdomain.UpdatedFailed
		}
	}
	return nil
}

// 公共
func (dnspod *Dnspod) commonRequest(ctx context.Context, apiAddr string, values url.Values, domain *dnsdomain.Domain) (status DnspodStatus, err error) {

	var resp *http.Response
	resp, err = http.PostForm(
		apiAddr,
		values,
	)
	if err != nil {
		return
	}

	err = provider.UnmarshalHTTPResponse(resp, apiAddr, err, &status)

	return
}

// 获得域名记录列表
func (dnspod *Dnspod) getRecordList(domain *dnsdomain.Domain, typ string) (result DnspodRecordListResp, err error) {
	values := url.Values{
		"login_token": {dnspod.clientID + "," + dnspod.clientSecret},
		"domain":      {domain.DomainName},
		"record_type": {typ},
		"sub_domain":  {domain.GetSubDomain()},
		"format":      {"json"},
	}
	if len(domain.Line) > 0 {
		values.Set(`record_line`, domain.Line)
	}

	var resp *http.Response
	resp, err = dnspod.client.PostForm(
		recordListAPI,
		values,
	)
	if err != nil {
		return
	}

	err = provider.UnmarshalHTTPResponse(resp, recordListAPI, err, &result)

	return
}
