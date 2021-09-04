package dnsdomain

import (
	"encoding/json"
	"strings"

	"github.com/webx-top/echo"
)

// Domain 域名实体
type Domain struct {
	IPFormat     string
	DomainName   string
	SubDomain    string
	UpdateStatus UpdateStatusType // 更新状态
	Extra        echo.H
}

func (d Domain) String() string {
	if len(d.SubDomain) > 0 {
		return d.SubDomain + "." + d.DomainName
	}
	return d.DomainName
}

func (d Domain) IP(ip string) string {
	if len(d.IPFormat) > 0 {
		return strings.ReplaceAll(d.IPFormat, Tag(`ip`), ip)
	}
	return ip
}

// GetFullDomain 获得全部的，子域名
func (d Domain) GetFullDomain() string {
	if len(d.SubDomain) > 0 {
		return d.SubDomain + "." + d.DomainName
	}
	return "@." + d.DomainName
}

// GetSubDomain 获得子域名，为空返回@
// 阿里云，dnspod需要
func (d Domain) GetSubDomain() string {
	if len(d.SubDomain) > 0 {
		return d.SubDomain
	}
	return "@"
}

type Result struct {
	Domain string
	Status string
}

func (d Domain) Result() *Result {
	return &Result{
		Domain: d.String(),
		Status: string(d.UpdateStatus),
	}
}

type Results map[string][]*Result

func (r *Results) Add(provider string, result *Result) {
	if _, ok := (*r)[provider]; !ok {
		(*r)[provider] = []*Result{}
	}
	(*r)[provider] = append((*r)[provider], result)
}

func (r *Results) String() string {
	b, _ := json.MarshalIndent(r, ``, `  `)
	return string(b)
}
