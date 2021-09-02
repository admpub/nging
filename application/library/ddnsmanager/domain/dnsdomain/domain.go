package dnsdomain

import "fmt"

// Domain 域名实体
type Domain struct {
	Port         int
	DomainName   string
	SubDomain    string
	UpdateStatus UpdateStatusType // 更新状态
}

func (d Domain) String() string {
	if len(d.SubDomain) > 0 {
		return d.SubDomain + "." + d.DomainName
	}
	return d.DomainName
}

func (d Domain) IP(ip string) string {
	if d.Port > 0 {
		return fmt.Sprintf(`%s:%d`, ip, d.Port)
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
