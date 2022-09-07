package alidns

import (
	"bytes"
	"context"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/admpub/log"
	"github.com/webx-top/echo"

	"github.com/nging-plugins/ddnsmanager/application/library/ddnsmanager/domain/dnsdomain"
	"github.com/nging-plugins/ddnsmanager/application/library/ddnsmanager/provider"
)

const (
	alidnsEndpoint       = "https://alidns.aliyuncs.com/"
	signUpURL            = `https://ram.console.aliyun.com/manage/ak`
	docLineType          = `https://help.aliyun.com/document_detail/29807.html`
	defaultTTL           = 600
	defaultTimeout int64 = 10
)

// https://help.aliyun.com/document_detail/29776.html?spm=a2c4g.11186623.6.672.715a45caji9dMA
// Alidns Alidns
type Alidns struct {
	clientID     string
	clientSecret string
	Domains      []*dnsdomain.Domain
	TTL          int
	client       *http.Client
}

// AlidnsSubDomainRecords 记录
type AlidnsSubDomainRecords struct {
	TotalCount    int
	DomainRecords struct {
		Record []struct {
			DomainName string
			RecordID   string
			Value      string
			Line       string
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
	return signUpURL
}

func (*Alidns) LineTypeURL() string {
	return docLineType
}

var configItems = echo.KVList{
	echo.NewKV(`clientId`, `AccessKey ID`).SetHKV(`inputType`, `text`).SetHKV(`required`, true),
	echo.NewKV(`clientSecret`, `AccessKey Secret`).SetHKV(`inputType`, `text`).SetHKV(`required`, true),
	echo.NewKV(`timeout`, `接口超时(秒)`).SetHKV(`inputType`, `number`).SetX(defaultTimeout),
	echo.NewKV(`ttl`, `域名TTL`).SetHKV(`inputType`, `number`).SetX(defaultTTL),
}

func (*Alidns) ConfigItems() echo.KVList {
	return configItems
}

var support = dnsdomain.Support{
	A:    true,
	AAAA: true,
	Line: true,
}

func (*Alidns) Support() dnsdomain.Support {
	return support
}

// Init 初始化
func (ali *Alidns) Init(settings echo.H, domains []*dnsdomain.Domain) error {
	ali.TTL = settings.Int(`ttl`)
	ali.clientID = settings.String(`clientId`)
	ali.clientSecret = settings.String(`clientSecret`)
	if ali.TTL < 1 { // 默认600s
		ali.TTL = defaultTTL
	}
	ali.Domains = domains
	timeout := settings.Int64(`timeout`)
	if timeout < 1 {
		timeout = defaultTimeout
	}
	ali.client = &http.Client{Timeout: time.Duration(timeout) * time.Second}
	return nil
}

func (ali *Alidns) Update(ctx context.Context, recordType string, ipAddr string) error {
	for _, domain := range ali.Domains {
		var record AlidnsSubDomainRecords
		// 获取当前域名信息
		params := url.Values{}
		params.Set("Action", "DescribeSubDomainRecords")
		params.Set("DomainName", domain.DomainName)
		params.Set("SubDomain", domain.GetFullDomain())
		params.Set("Type", recordType)
		if len(domain.Line) > 0 {
			params.Set("Line", domain.Line)
		}
		err := ali.request(ctx, params, &record)
		if err != nil {
			domain.UpdateStatus = dnsdomain.UpdatedFailed
			return err
		}

		if record.TotalCount > 0 {
			// 存在，更新
			ali.modify(ctx, record, domain, recordType, ipAddr)
		} else {
			// 不存在，创建
			ali.create(ctx, domain, recordType, ipAddr)
		}

	}
	return nil
}

// 创建
func (ali *Alidns) create(ctx context.Context, domain *dnsdomain.Domain, recordType string, ipAddr string) error {
	ipAddr = domain.IP(ipAddr)
	params := url.Values{}
	params.Set("Action", "AddDomainRecord")
	params.Set("DomainName", domain.DomainName)
	params.Set("RR", domain.GetSubDomain())
	params.Set("Type", recordType)
	params.Set("Value", ipAddr)
	params.Set("TTL", strconv.Itoa(ali.TTL))
	if len(domain.Line) > 0 {
		params.Set("Line", domain.Line)
	}

	var result AlidnsResp
	err := ali.request(ctx, params, &result)

	if err == nil && result.RecordID != "" {
		log.Infof("新增域名解析 %s 成功！IP: %s", domain, ipAddr)
		domain.UpdateStatus = dnsdomain.UpdatedSuccess
	} else {
		log.Errorf("新增域名解析 %s 失败: %v", domain, err)
		domain.UpdateStatus = dnsdomain.UpdatedFailed
	}
	return err
}

// 修改
func (ali *Alidns) modify(ctx context.Context, record AlidnsSubDomainRecords, domain *dnsdomain.Domain, recordType string, ipAddr string) error {
	ipAddr = domain.IP(ipAddr)

	// 相同不修改
	if len(record.DomainRecords.Record) > 0 && record.DomainRecords.Record[0].Value == ipAddr {
		log.Infof("你的IP %s 没有变化, 域名 %s", ipAddr, domain)
		domain.UpdateStatus = dnsdomain.UpdatedNothing
		return nil
	}

	params := url.Values{}
	params.Set("Action", "UpdateDomainRecord")
	params.Set("RR", domain.GetSubDomain())
	params.Set("RecordId", record.DomainRecords.Record[0].RecordID)
	params.Set("Type", recordType)
	params.Set("Value", ipAddr)
	params.Set("TTL", strconv.Itoa(ali.TTL))
	if len(domain.Line) > 0 {
		params.Set("Line", domain.Line)
	}

	var result AlidnsResp
	err := ali.request(ctx, params, &result)

	if err == nil && result.RecordID != "" {
		log.Infof("更新域名解析 %s 成功！IP: %s", domain, ipAddr)
		domain.UpdateStatus = dnsdomain.UpdatedSuccess
	} else {
		log.Errorf("更新域名解析 %s 失败: %v", domain, err)
		domain.UpdateStatus = dnsdomain.UpdatedFailed
	}
	return err
}

// request 统一请求接口
func (ali *Alidns) request(ctx context.Context, params url.Values, result interface{}) (err error) {

	AliyunSigner(ali.clientID, ali.clientSecret, &params)
	var req *http.Request
	req, err = http.NewRequestWithContext(
		ctx,
		"GET",
		alidnsEndpoint,
		bytes.NewBuffer(nil),
	)
	req.URL.RawQuery = params.Encode()

	if err != nil {
		log.Error("http.NewRequest失败. Error: ", err)
		return
	}

	var resp *http.Response
	resp, err = ali.client.Do(req)
	if err != nil {
		return
	}
	err = provider.UnmarshalHTTPResponse(resp, alidnsEndpoint, err, result)

	return
}
