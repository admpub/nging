package ip2region

import (
	"strings"
)

type IpInfo struct {
	Country  string
	Region   string
	Province string
	City     string
	ISP      string
}

func (ip IpInfo) String() string {
	return ip.Country + "|" + ip.Region + "|" + ip.Province + "|" + ip.City + "|" + ip.ISP
}

func getIpInfo(line string) IpInfo {
	lineSlice := strings.Split(line, "|")
	ipInfo := IpInfo{}
	length := len(lineSlice)
	if length < 5 {
		for i := 0; i <= 5-length; i++ {
			lineSlice = append(lineSlice, "")
		}
	}

	ipInfo.Country = lineSlice[0]
	ipInfo.Region = lineSlice[1]
	ipInfo.Province = lineSlice[2]
	ipInfo.City = lineSlice[3]
	ipInfo.ISP = lineSlice[4]
	return ipInfo
}
