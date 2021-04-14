package ipsimplefilter

import (
	"bytes"
	"net"
	"strings"

	"github.com/admpub/log"
)

func New(ip string) IPR {
	r := IPR{}
	r.Parse(ip)
	return r
}

type IPR [2]net.IP

func (i *IPR) Parse(ip string) {
	ips := strings.SplitN(ip, `-`, 2)
	ips[0] = strings.TrimSpace(ips[0])
	ipStart := net.ParseIP(ips[0])
	var ipEnd net.IP
	if len(ips) != 2 {
		ipEnd = ipStart
	} else {
		ips[1] = strings.TrimSpace(ips[1])
		ipEnd = net.ParseIP(ips[1])
	}
	(*i)[0] = ipStart
	(*i)[1] = ipEnd
}

func (i *IPR) Contains(trial net.IP) bool {
	return bytes.Compare(trial, i[0]) >= 0 && bytes.Compare(trial, i[1]) <= 0
}

func NewFilter() *Filter {
	return &Filter{}
}

type Filter struct {
	allowList []IPR
	blockList []IPR
}

func (f *Filter) Allowed(ip string) error {
	return f.add(ip, true)
}

func (f *Filter) Blocked(ip string) error {
	return f.add(ip, false)
}

func (f *Filter) add(ip string, isAllowed bool) error {
	if len(ip) == 0 {
		return nil
	}
	if isAllowed {
		f.allowList = append(f.allowList, New(ip))
	} else {
		f.blockList = append(f.blockList, New(ip))
	}
	return nil
}

func (f *Filter) IsAllowed(ip string) bool {
	trial := net.ParseIP(ip)
	if trial.To4() == nil {
		log.Warnf("%v is not an IPv4 address", ip)
		return false
	}
	if len(f.allowList) > 0 {
		for _, allow := range f.allowList {
			if allow.Contains(trial) {
				return true
			}
		}
		return false
	}
	for _, block := range f.blockList {
		if block.Contains(trial) {
			return false
		}
	}
	return true
}
