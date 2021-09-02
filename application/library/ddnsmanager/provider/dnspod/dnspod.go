package dnspod

import (
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/admpub/nging/v3/application/library/ddnsmanager/domain/dnsdomain"
	"github.com/admpub/nging/v3/application/library/ddnsmanager/provider"
	"github.com/webx-top/echo"
)

const (
	recordListAPI   string = "https://dnsapi.cn/Record.List"
	recordModifyURL string = "https://dnsapi.cn/Record.Modify"
	recordCreateAPI string = "https://dnsapi.cn/Record.Create"
)

// https://cloud.tencent.com/document/api/302/8516
// Dnspod 腾讯云dns实现
type Dnspod struct {
	clientID     string
	clientSecret string
	Domains      []*dnsdomain.Domain
	TTL          string
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

var configItems = echo.KVList{
	echo.NewKV(`ttl`, `TTL`).SetHKV(`inputType`, `number`),
	echo.NewKV(`clientId`, `clientId`).SetHKV(`inputType`, `text`),
	echo.NewKV(`clientSecret`, `clientSecret`).SetHKV(`inputType`, `text`),
}

func (*Dnspod) ConfigItems() echo.KVList {
	return configItems
}

// Init 初始化
func (dnspod *Dnspod) Init(settings echo.H, domains []*dnsdomain.Domain) error {
	dnspod.TTL = settings.String(`ttl`)
	dnspod.clientID = settings.String(`clientId`)
	dnspod.clientSecret = settings.String(`clientSecret`)
	if dnspod.TTL == "" { // 默认600s
		dnspod.TTL = "600"
	}
	return nil
}

func (dnspod *Dnspod) Update(recordType string, ipAddr string) error {

	for _, domain := range dnspod.Domains {
		result, err := dnspod.getRecordList(domain, recordType)
		if err != nil {
			return err
		}

		if len(result.Records) > 0 { // 更新
			err = dnspod.modify(result, domain, recordType, ipAddr)
		} else { // 新增
			err = dnspod.create(result, domain, recordType, ipAddr)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

// 创建
func (dnspod *Dnspod) create(result DnspodRecordListResp, domain *dnsdomain.Domain, recordType string, ipAddr string) error {
	ipAddr = domain.IP(ipAddr)
	status, err := dnspod.commonRequest(
		recordCreateAPI,
		url.Values{
			"login_token": {dnspod.clientID + "," + dnspod.clientSecret},
			"domain":      {domain.DomainName},
			"sub_domain":  {domain.GetSubDomain()},
			"record_type": {recordType},
			"record_line": {"默认"},
			"value":       {ipAddr},
			"ttl":         {dnspod.TTL},
			"format":      {"json"},
		},
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
func (dnspod *Dnspod) modify(result DnspodRecordListResp, domain *dnsdomain.Domain, recordType string, ipAddr string) error {
	ipAddr = domain.IP(ipAddr)
	for _, record := range result.Records {
		// 相同不修改
		if record.Value == ipAddr {
			log.Printf("你的IP %s 没有变化, 域名 %s", ipAddr, domain)
			continue
		}
		status, err := dnspod.commonRequest(
			recordModifyURL,
			url.Values{
				"login_token": {dnspod.clientID + "," + dnspod.clientSecret},
				"domain":      {domain.DomainName},
				"sub_domain":  {domain.GetSubDomain()},
				"record_type": {recordType},
				"record_line": {"默认"},
				"record_id":   {record.ID},
				"value":       {ipAddr},
				"ttl":         {dnspod.TTL},
				"format":      {"json"},
			},
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
func (dnspod *Dnspod) commonRequest(apiAddr string, values url.Values, domain *dnsdomain.Domain) (status DnspodStatus, err error) {

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

	client := http.Client{Timeout: 10 * time.Second}
	var resp *http.Response
	resp, err = client.PostForm(
		recordListAPI,
		values,
	)
	if err != nil {
		return
	}

	err = provider.UnmarshalHTTPResponse(resp, recordListAPI, err, &result)

	return
}
