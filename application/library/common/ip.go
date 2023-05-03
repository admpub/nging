package common

import "net"

// GetLocalIP 获取本机网卡IP
func GetLocalIP(ver ...int) (string, error) {
	var v int
	if len(ver) > 0 {
		v = ver[0]
	}
	if v == 6 {
		return GetLocalIPv6()
	}

	return GetLocalIPv4()
}

// GetLocalIPv4 获取本机网卡IPv4
func GetLocalIPv4() (ipv4 string, err error) {
	var (
		addrs   []net.Addr
		addr    net.Addr
		ipNet   *net.IPNet // IP地址
		isIPNet bool
	)
	ipv4 = `127.0.0.1`
	// 获取所有网卡
	if addrs, err = net.InterfaceAddrs(); err != nil {
		return
	}
	// 取第一个非lo的网卡IP
	for _, addr = range addrs {
		// 这个网络地址是IP地址: ipv4, ipv6
		if ipNet, isIPNet = addr.(*net.IPNet); isIPNet && !ipNet.IP.IsLoopback() {
			// 跳过IPV6
			if ipNet.IP.To4() != nil {
				ipv4 = ipNet.IP.String() // 192.168.1.1
				return
			}
		}
	}
	return
}

// GetLocalIPv6 获取本机网卡IPv6
func GetLocalIPv6() (ipv6 string, err error) {
	var (
		addrs   []net.Addr
		addr    net.Addr
		ipNet   *net.IPNet // IP地址
		isIPNet bool
	)
	ipv6 = `::1`
	// 获取所有网卡
	if addrs, err = net.InterfaceAddrs(); err != nil {
		return
	}
	// 取第一个非lo的网卡IP
	for _, addr = range addrs {
		// 这个网络地址是IP地址: ipv4, ipv6
		if ipNet, isIPNet = addr.(*net.IPNet); isIPNet && !ipNet.IP.IsLoopback() {
			// IPV6
			if ipNet.IP.To4() == nil {
				ipv6 = ipNet.IP.String()
				return
			}
		}
	}
	return
}

// GetLocalIPs 获取本机网卡IP
func GetLocalIPs() (ipv4 []string, ipv6 []string, err error) {
	var (
		addrs   []net.Addr
		addr    net.Addr
		ipNet   *net.IPNet // IP地址
		isIPNet bool
	)
	// 获取所有网卡
	if addrs, err = net.InterfaceAddrs(); err != nil {
		return
	}
	for _, addr = range addrs {
		// 这个网络地址是IP地址: ipv4, ipv6
		if ipNet, isIPNet = addr.(*net.IPNet); isIPNet && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				ipv4 = append(ipv4, ipNet.IP.String()) // 192.168.1.1
			} else {
				ipv6 = append(ipv6, ipNet.IP.String())
			}
		}
	}
	return
}
