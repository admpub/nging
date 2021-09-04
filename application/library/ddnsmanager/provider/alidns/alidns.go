package alidns

import (
	"bytes"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/admpub/nging/v3/application/library/ddnsmanager/domain/dnsdomain"
	"github.com/admpub/nging/v3/application/library/ddnsmanager/provider"
	"github.com/webx-top/echo"
)

const (
	alidnsEndpoint string = "https://alidns.aliyuncs.com/"
)

// https://help.aliyun.com/document_detail/29776.html?spm=a2c4g.11186623.6.672.715a45caji9dMA
// Alidns Alidns
type Alidns struct {
	clientID     string
	clientSecret string
	Domains      []*dnsdomain.Domain
	TTL          int
}

// AlidnsSubDomainRecords 记录
type AlidnsSubDomainRecords struct {
	TotalCount    int
	DomainRecords struct {
		Record []struct {
			DomainName string
			RecordID   string
			Value      string
		}
	}
}

// AlidnsResp 修改/添加返回结果
type AlidnsResp struct {
	RecordID  string
	RequestID string
}

const Name = `AliDNS`

func (*Alidns) Name() string {
	return Name
}

func (*Alidns) Description() string {
	return `阿里云DNS`
}

func (*Alidns) SignUpURL() string {
	return ``
}

var configItems = echo.KVList{
	echo.NewKV(`ttl`, `TTL`).SetHKV(`inputType`, `number`),
	echo.NewKV(`clientId`, `clientId`).SetHKV(`inputType`, `text`),
	echo.NewKV(`clientSecret`, `clientSecret`).SetHKV(`inputType`, `text`),
}

func (*Alidns) ConfigItems() echo.KVList {
	return configItems
}

// Init 初始化
func (ali *Alidns) Init(settings echo.H, domains []*dnsdomain.Domain) error {
	ali.TTL = settings.Int(`ttl`)
	ali.clientID = settings.String(`clientId`)
	ali.clientSecret = settings.String(`clientSecret`)
	if ali.TTL <= 0 { // 默认600s
		ali.TTL = 600
	}
	return nil
}

func (ali *Alidns) Update(recordType string, ipAddr string) error {
	for _, domain := range ali.Domains {
		var record AlidnsSubDomainRecords
		// 获取当前域名信息
		params := url.Values{}
		params.Set("Action", "DescribeSubDomainRecords")
		params.Set("SubDomain", domain.GetFullDomain())
		params.Set("Type", recordType)
		err := ali.request(params, &record)
		if err != nil {
			return err
		}

		if record.TotalCount > 0 {
			// 存在，更新
			ali.modify(record, domain, recordType, ipAddr)
		} else {
			// 不存在，创建
			ali.create(domain, recordType, ipAddr)
		}

	}
	return nil
}

// 创建
func (ali *Alidns) create(domain *dnsdomain.Domain, recordType string, ipAddr string) {
	ipAddr = domain.IP(ipAddr)
	params := url.Values{}
	params.Set("Action", "AddDomainRecord")
	params.Set("DomainName", domain.DomainName)
	params.Set("RR", domain.GetSubDomain())
	params.Set("Type", recordType)
	params.Set("Value", ipAddr)
	params.Set("TTL", strconv.Itoa(ali.TTL))

	var result AlidnsResp
	err := ali.request(params, &result)

	if err == nil && result.RecordID != "" {
		log.Printf("新增域名解析 %s 成功！IP: %s", domain, ipAddr)
		domain.UpdateStatus = dnsdomain.UpdatedSuccess
	} else {
		log.Printf("新增域名解析 %s 失败！", domain)
		domain.UpdateStatus = dnsdomain.UpdatedFailed
	}
}

// 修改
func (ali *Alidns) modify(record AlidnsSubDomainRecords, domain *dnsdomain.Domain, recordType string, ipAddr string) {
	ipAddr = domain.IP(ipAddr)

	// 相同不修改
	if len(record.DomainRecords.Record) > 0 && record.DomainRecords.Record[0].Value == ipAddr {
		log.Printf("你的IP %s 没有变化, 域名 %s", ipAddr, domain)
		return
	}

	params := url.Values{}
	params.Set("Action", "UpdateDomainRecord")
	params.Set("RR", domain.GetSubDomain())
	params.Set("RecordId", record.DomainRecords.Record[0].RecordID)
	params.Set("Type", recordType)
	params.Set("Value", ipAddr)
	params.Set("TTL", strconv.Itoa(ali.TTL))

	var result AlidnsResp
	err := ali.request(params, &result)

	if err == nil && result.RecordID != "" {
		log.Printf("更新域名解析 %s 成功！IP: %s", domain, ipAddr)
		domain.UpdateStatus = dnsdomain.UpdatedSuccess
	} else {
		log.Printf("更新域名解析 %s 失败！", domain)
		domain.UpdateStatus = dnsdomain.UpdatedFailed
	}
}

// request 统一请求接口
func (ali *Alidns) request(params url.Values, result interface{}) (err error) {

	AliyunSigner(ali.clientID, ali.clientSecret, &params)
	var req *http.Request
	req, err = http.NewRequest(
		"GET",
		alidnsEndpoint,
		bytes.NewBuffer(nil),
	)
	req.URL.RawQuery = params.Encode()

	if err != nil {
		log.Println("http.NewRequest失败. Error: ", err)
		return
	}

	client := http.Client{Timeout: 10 * time.Second}
	var resp *http.Response
	resp, err = client.Do(req)
	if err != nil {
		return
	}
	err = provider.UnmarshalHTTPResponse(resp, alidnsEndpoint, err, result)

	return
}
