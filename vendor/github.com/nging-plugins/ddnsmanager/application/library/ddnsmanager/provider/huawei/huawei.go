package huawei

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/admpub/log"
	"github.com/webx-top/echo"

	"github.com/nging-plugins/ddnsmanager/application/library/ddnsmanager/domain/dnsdomain"
	"github.com/nging-plugins/ddnsmanager/application/library/ddnsmanager/provider"
)

const (
	huaweicloudEndpoint = "https://dns.myhuaweicloud.com"
	signUpURL           = `https://console.huaweicloud.com/iam/?locale=zh-cn#/mine/accessKey`
	//docLineType         = `https://support.huaweicloud.com/api-dns/zh-cn_topic_0085546214.html`
	docLineType          = `` //暂时留空(文档有疑问，获取zones列表的返回值里是否包含line字段(文档中没有))
	defaultTTL           = 300
	defaultTimeout int64 = 10
)

// https://support.huaweicloud.com/api-dns/dns_api_64001.html
// Huaweicloud Huaweicloud
type Huaweicloud struct {
	clientID     string
	clientSecret string
	Domains      []*dnsdomain.Domain
	TTL          int
	client       *http.Client
}

// HuaweicloudZonesResp zones response
type HuaweicloudZonesResp struct {
	Zones []struct {
		ID   string
		Name string
	}
}

// HuaweicloudRecordsResp 记录返回结果
type HuaweicloudRecordsResp struct {
	Recordsets []HuaweicloudRecordsets
}

// HuaweicloudRecordsets 记录
type HuaweicloudRecordsets struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	ZoneID  string   `json:"zone_id"`
	Status  string   `json:"status"`
	Type    string   `json:"type"`
	TTL     int      `json:"ttl"`
	Records []string `json:"records"`
	Line    string   `json:"line,omitempty"`
}

const Name = `HuaWei`

func (*Huaweicloud) Name() string {
	return Name
}

func (*Huaweicloud) Description() string {
	return `华为云DNS`
}

func (*Huaweicloud) SignUpURL() string {
	return signUpURL
}

func (*Huaweicloud) LineTypeURL() string {
	return docLineType
}

var configItems = echo.KVList{
	echo.NewKV(`clientId`, `AK`).SetHKV(`inputType`, `text`).SetHKV(`helpBlock`, `Access Key ID`).SetHKV(`required`, true),
	echo.NewKV(`clientSecret`, `SK`).SetHKV(`inputType`, `text`).SetHKV(`helpBlock`, `Secret Access Key`).SetHKV(`required`, true),
	echo.NewKV(`timeout`, `接口超时(秒)`).SetHKV(`inputType`, `number`).SetX(defaultTimeout),
	echo.NewKV(`ttl`, `域名TTL`).SetHKV(`inputType`, `number`).SetX(defaultTTL),
}

func (*Huaweicloud) ConfigItems() echo.KVList {
	return configItems
}

var support = dnsdomain.Support{
	A:    true,
	AAAA: true,
	Line: true,
}

func (*Huaweicloud) Support() dnsdomain.Support {
	return support
}

// Init 初始化
func (hw *Huaweicloud) Init(settings echo.H, domains []*dnsdomain.Domain) error {
	hw.TTL = settings.Int(`ttl`)
	hw.clientID = settings.String(`clientId`)
	hw.clientSecret = settings.String(`clientSecret`)
	if hw.TTL < 1 {
		hw.TTL = defaultTTL
	}
	hw.Domains = domains
	timeout := settings.Int64(`timeout`)
	if timeout < 1 {
		timeout = defaultTimeout
	}
	hw.client = &http.Client{Timeout: time.Duration(timeout) * time.Second}
	return nil
}

func (hw *Huaweicloud) Update(ctx context.Context, recordType string, ipAddr string) error {
	for _, domain := range hw.Domains {

		var records HuaweicloudRecordsResp

		err := hw.request(
			ctx,
			"GET",
			fmt.Sprintf(huaweicloudEndpoint+"/v2/recordsets?type=%s&name=%s", recordType, domain),
			nil,
			&records,
		)

		if err != nil {
			domain.UpdateStatus = dnsdomain.UpdatedFailed
			return err
		}

		find := false
		for _, record := range records.Recordsets {
			// 名称相同才更新。华为云默认是模糊搜索
			if record.Name == domain.String()+"." {
				var linesame bool
				if len(domain.Line) > 0 {
					linesame = record.Line == domain.Line
				} else {
					linesame = true
				}
				if linesame {
					// 更新
					hw.modify(ctx, record, domain, recordType, ipAddr)
					find = true
					break
				}
			}
		}

		if !find {
			// 新增
			hw.create(ctx, domain, recordType, ipAddr)
		}

	}
	return nil
}

// 创建
func (hw *Huaweicloud) create(ctx context.Context, domain *dnsdomain.Domain, recordType string, ipAddr string) {
	ipAddr = domain.IP(ipAddr)
	zone, err := hw.getZones(ctx, domain)
	if err != nil {
		return
	}
	if len(zone.Zones) == 0 {
		log.Info("未能找到公网域名, 请检查域名是否添加")
		domain.UpdateStatus = dnsdomain.UpdatedFailed
		return
	}

	zoneID := zone.Zones[0].ID
	for _, z := range zone.Zones {
		if z.Name == domain.DomainName+"." {
			zoneID = z.ID
			break
		}
	}

	record := &HuaweicloudRecordsets{
		Type:    recordType,
		Name:    domain.String() + ".",
		Records: []string{ipAddr},
		TTL:     hw.TTL,
	}
	var result HuaweicloudRecordsets
	apiVer := `2`
	if len(domain.Line) > 0 {
		apiVer = `2.1`
		record.Line = domain.Line
	}

	err = hw.request(
		ctx,
		"POST",
		fmt.Sprintf(huaweicloudEndpoint+"/v"+apiVer+"/zones/%s/recordsets", zoneID),
		record,
		&result,
	)
	if err == nil && (len(result.Records) > 0 && result.Records[0] == ipAddr) {
		log.Infof("新增域名解析 %s 成功！IP: %s", domain, ipAddr)
		domain.UpdateStatus = dnsdomain.UpdatedSuccess
	} else {
		log.Errorf("新增域名解析 %s 失败！Status: %s, Error: %v", domain, result.Status, err)
		domain.UpdateStatus = dnsdomain.UpdatedFailed
	}
}

// 修改
func (hw *Huaweicloud) modify(ctx context.Context, record HuaweicloudRecordsets, domain *dnsdomain.Domain, recordType string, ipAddr string) {
	ipAddr = domain.IP(ipAddr)

	// 相同不修改
	if len(record.Records) > 0 && record.Records[0] == ipAddr {
		log.Infof("你的IP %s 没有变化, 域名 %s", ipAddr, domain)
		domain.UpdateStatus = dnsdomain.UpdatedNothing
		return
	}

	var request map[string]interface{} = make(map[string]interface{})
	request["records"] = []string{ipAddr}
	request["ttl"] = hw.TTL

	var result HuaweicloudRecordsets

	apiVer := `2`
	if len(domain.Line) > 0 {
		apiVer = `2.1`
		request["line"] = domain.Line
	}

	err := hw.request(
		ctx,
		"PUT",
		fmt.Sprintf(huaweicloudEndpoint+"/v"+apiVer+"/zones/%s/recordsets/%s", record.ZoneID, record.ID),
		&request,
		&result,
	)

	if err == nil && (len(result.Records) > 0 && result.Records[0] == ipAddr) {
		log.Infof("更新域名解析 %s 成功！IP: %s, 状态: %s", domain, ipAddr, result.Status)
		domain.UpdateStatus = dnsdomain.UpdatedSuccess
	} else {
		log.Errorf("更新域名解析 %s 失败！Status: %s, Error: %v", domain, result.Status, err)
		domain.UpdateStatus = dnsdomain.UpdatedFailed
	}
}

// 获得域名记录列表
func (hw *Huaweicloud) getZones(ctx context.Context, domain *dnsdomain.Domain) (result HuaweicloudZonesResp, err error) {
	err = hw.request(
		ctx,
		"GET",
		fmt.Sprintf(huaweicloudEndpoint+"/v2/zones?name=%s", domain.DomainName),
		nil,
		&result,
	)

	return
}

// request 统一请求接口
func (hw *Huaweicloud) request(ctx context.Context, method string, url string, data interface{}, result interface{}) (err error) {
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
		log.Error("http.NewRequest失败. Error: ", err)
		return
	}

	s := Signer{
		Key:    hw.clientID,
		Secret: hw.clientSecret,
	}
	err = s.Sign(req)
	if err != nil {
		return
	}

	req.Header.Add("content-type", "application/json")

	var resp *http.Response
	resp, err = hw.client.Do(req)
	if err != nil {
		return
	}
	err = provider.UnmarshalHTTPResponse(resp, url, err, result)

	return
}
