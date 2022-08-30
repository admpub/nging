package cloudflare

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/nging-plugins/ddnsmanager/application/library/ddnsmanager/domain/dnsdomain"
	"github.com/nging-plugins/ddnsmanager/application/library/ddnsmanager/provider"
	"github.com/webx-top/echo"
)

const (
	zonesAPI             = "https://api.cloudflare.com/client/v4/zones"
	signUpURL            = `https://dash.cloudflare.com/profile/api-tokens`
	docLineType          = `` // 文档网址留空表示不支持
	defaultTTL           = 1
	defaultTimeout int64 = 30
)

// Cloudflare Cloudflare实现
type Cloudflare struct {
	clientSecret string
	Domains      []*dnsdomain.Domain
	TTL          int
	client       *http.Client
}

// CloudflareZonesResp cloudflare zones返回结果
type CloudflareZonesResp struct {
	CloudflareStatus
	Result []struct {
		ID     string
		Name   string
		Status string
		Paused bool
	}
}

// CloudflareRecordsResp records
type CloudflareRecordsResp struct {
	CloudflareStatus
	Result []CloudflareRecord
}

// CloudflareRecord 记录实体
type CloudflareRecord struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Type    string `json:"type"`
	Content string `json:"content"`
	Proxied bool   `json:"proxied"`
	TTL     int    `json:"ttl"`
}

// CloudflareStatus 公共状态
type CloudflareStatus struct {
	Success  bool
	Messages []string
}

const Name = `Cloudflare`

func (*Cloudflare) Name() string {
	return Name
}

func (*Cloudflare) Description() string {
	return ``
}

func (*Cloudflare) SignUpURL() string {
	return signUpURL
}

func (*Cloudflare) LineTypeURL() string {
	return docLineType
}

var configItems = echo.KVList{
	echo.NewKV(`clientSecret`, `Token`).SetHKV(`inputType`, `text`).SetHKV(`required`, true),
	echo.NewKV(`timeout`, `接口超时(秒)`).SetHKV(`inputType`, `number`).SetX(defaultTimeout),
	echo.NewKV(`ttl`, `域名TTL`).SetHKV(`inputType`, `number`).SetX(defaultTTL),
}

func (*Cloudflare) ConfigItems() echo.KVList {
	return configItems
}

var support = dnsdomain.Support{
	A:    true,
	AAAA: true,
	Line: true,
}

func (*Cloudflare) Support() dnsdomain.Support {
	return support
}

// Init 初始化
func (cf *Cloudflare) Init(settings echo.H, domains []*dnsdomain.Domain) error {
	cf.TTL = settings.Int(`ttl`)
	cf.clientSecret = settings.String(`clientSecret`)
	if cf.TTL < 1 {
		cf.TTL = defaultTTL
	}
	cf.Domains = domains
	timeout := settings.Int64(`timeout`)
	if timeout < 1 {
		timeout = defaultTimeout
	}
	cf.client = &http.Client{Timeout: time.Duration(timeout) * time.Second}
	return nil
}

func (cf *Cloudflare) Update(ctx context.Context, recordType string, ipAddr string) error {
	for _, domain := range cf.Domains {
		// get zone
		result, err := cf.getZones(ctx, domain)
		if err != nil || len(result.Result) != 1 {
			domain.UpdateStatus = dnsdomain.UpdatedFailed
			return err
		}
		zoneID := result.Result[0].ID

		var records CloudflareRecordsResp
		// getDomains 最多更新前50条
		err = cf.request(
			ctx,
			"GET",
			fmt.Sprintf(zonesAPI+"/%s/dns_records?type=%s&name=%s&per_page=50", zoneID, recordType, domain),
			nil,
			&records,
		)

		if err != nil || !records.Success {
			domain.UpdateStatus = dnsdomain.UpdatedFailed
			return err
		}

		if len(records.Result) > 0 {
			// 更新
			cf.modify(ctx, records, zoneID, domain, recordType, ipAddr)
		} else {
			// 新增
			cf.create(ctx, zoneID, domain, recordType, ipAddr)
		}
	}
	return nil
}

// 创建
func (cf *Cloudflare) create(ctx context.Context, zoneID string, domain *dnsdomain.Domain, recordType string, ipAddr string) {
	ipAddr = domain.IP(ipAddr)
	record := &CloudflareRecord{
		Type:    recordType,
		Name:    domain.String(),
		Content: ipAddr,
		Proxied: false,
		TTL:     cf.TTL,
	}
	var status CloudflareStatus
	err := cf.request(
		ctx,
		"POST",
		fmt.Sprintf(zonesAPI+"/%s/dns_records", zoneID),
		record,
		&status,
	)
	if err == nil && status.Success {
		log.Printf("新增域名解析 %s 成功！IP: %s", domain, ipAddr)
		domain.UpdateStatus = dnsdomain.UpdatedSuccess
	} else {
		log.Printf("新增域名解析 %s 失败！Messages: %s", domain, status.Messages)
		domain.UpdateStatus = dnsdomain.UpdatedFailed
	}
}

// 修改
func (cf *Cloudflare) modify(ctx context.Context, result CloudflareRecordsResp, zoneID string, domain *dnsdomain.Domain, recordType string, ipAddr string) {
	ipAddr = domain.IP(ipAddr)
	for _, record := range result.Result {
		// 相同不修改
		if record.Content == ipAddr {
			log.Printf("你的IP %s 没有变化, 域名 %s", ipAddr, domain)
			domain.UpdateStatus = dnsdomain.UpdatedNothing
			continue
		}
		var status CloudflareStatus
		record.Content = ipAddr
		record.TTL = cf.TTL

		err := cf.request(
			ctx,
			"PUT",
			fmt.Sprintf(zonesAPI+"/%s/dns_records/%s", zoneID, record.ID),
			record,
			&status,
		)

		if err == nil && status.Success {
			log.Printf("更新域名解析 %s 成功！IP: %s", domain, ipAddr)
			domain.UpdateStatus = dnsdomain.UpdatedSuccess
		} else {
			log.Printf("更新域名解析 %s 失败！Messages: %s", domain, status.Messages)
			domain.UpdateStatus = dnsdomain.UpdatedFailed
		}
	}
}

// 获得域名记录列表
func (cf *Cloudflare) getZones(ctx context.Context, domain *dnsdomain.Domain) (result CloudflareZonesResp, err error) {
	err = cf.request(
		ctx,
		"GET",
		fmt.Sprintf(zonesAPI+"?name=%s&status=%s&per_page=%s", domain.DomainName, "active", "50"),
		nil,
		&result,
	)

	return
}

// request 统一请求接口
func (cf *Cloudflare) request(ctx context.Context, method string, url string, data interface{}, result interface{}) (err error) {
	jsonStr := make([]byte, 0)
	if data != nil {
		jsonStr, _ = json.Marshal(data)
	}
	var req *http.Request
	req, err = http.NewRequestWithContext(
		ctx,
		method,
		url,
		bytes.NewBuffer(jsonStr),
	)
	if err != nil {
		log.Println("http.NewRequest失败. Error: ", err)
		return
	}
	req.Header.Set("Authorization", "Bearer "+cf.clientSecret)
	req.Header.Set("Content-Type", "application/json")

	var resp *http.Response
	resp, err = cf.client.Do(req)
	if err != nil {
		return
	}
	err = provider.UnmarshalHTTPResponse(resp, url, err, result)

	return
}
