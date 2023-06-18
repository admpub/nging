package nftablesutils

import (
	"fmt"
	"net"

	"github.com/admpub/log"
	"github.com/vishvananda/netlink"
)

type Logger interface {
	Debugf(format string, a ...interface{})
}

// Create network link for interface.
func CreateIface(
	log Logger,
	iface, linkType string,
	ip net.IP, ipNet *net.IPNet,
) error {
	log.Debugf("%q creating…", iface)

	_, err := net.InterfaceByName(iface)
	if err == nil {
		log.Debugf("%q already exists", iface)
		// we should remove it first
		err = RemoveIface(log, iface)
		if err != nil {
			return err
		}
	}

	la := netlink.NewLinkAttrs()
	la.Name = iface
	link := &netlink.GenericLink{LinkAttrs: la, LinkType: linkType}
	err = netlink.LinkAdd(link)
	if err != nil {
		return fmt.Errorf("%q can't add link: %s", iface, err)
	}
	log.Debugf("%q link added", iface)

	addr := &netlink.Addr{IPNet: &net.IPNet{IP: ip, Mask: ipNet.Mask}}
	err = netlink.AddrAdd(link, addr)
	if err != nil {
		return fmt.Errorf("%q can't add addr: %v", iface, err)
	}
	log.Debugf("%q ip %q, net %q was set", iface, ip, ipNet)

	err = netlink.LinkSetUp(link)
	if err != nil {
		return fmt.Errorf("%s can't link set up: %s", iface, err)
	}
	log.Debugf("%q link is up", iface)

	return nil
}

// Remove network link for interface.
func RemoveIface(log Logger, iface string) error {
	log.Debugf("%q removing…", iface)

	link, err := netlink.LinkByName(iface)
	if err != nil {
		return fmt.Errorf("%q can't find: %v", iface, err)
	}

	err = netlink.LinkSetDown(link)
	if err != nil {
		return fmt.Errorf("%s can't link set down: %s", iface, err)
	}
	log.Debugf("%q link is down", iface)

	err = netlink.LinkDel(link)
	if err != nil {
		return fmt.Errorf("%q can't del link: %s", iface, err)
	}
	log.Debugf("%q link removed", iface)

	return nil
}

// NetInterface 本机网络
type NetInterface struct {
	Name    string
	Address []string
}

var ipv6Unicast *net.IPNet

func init() {
	var err error
	// https://en.wikipedia.org/wiki/IPv6_address#General_allocation
	_, ipv6Unicast, err = net.ParseCIDR("2000::/3")
	if err != nil {
		panic(err)
	}
}

// GetNetInterface 获得网卡地址 (返回ipv4, ipv6地址)
func GetNetInterface(interfaceName string) (ipv4NetInterfaces []NetInterface, ipv6NetInterfaces []NetInterface, err error) {
	var allNetInterfaces []net.Interface
	if len(interfaceName) > 0 {
		var ifaces *net.Interface
		ifaces, err = net.InterfaceByName(interfaceName)
		if err == nil {
			allNetInterfaces = append(allNetInterfaces, *ifaces)
		}
	} else {
		allNetInterfaces, err = net.Interfaces()
	}
	if err != nil {
		log.Error("net.Interfaces failed, err: ", err.Error())
		return ipv4NetInterfaces, ipv6NetInterfaces, err
	}

	for _, netInter := range allNetInterfaces {
		if (netInter.Flags & net.FlagUp) == 0 {
			continue
		}
		addrs, _ := netInter.Addrs()
		var ipv4 []string
		var ipv6 []string

		for _, address := range addrs {
			ipnet, ok := address.(*net.IPNet)
			if !ok || !ipnet.IP.IsGlobalUnicast() {
				continue
			}
			// 需匹配全局单播地址
			ones, bits := ipnet.Mask.Size()
			switch bits / 8 {
			case net.IPv6len:
				if ones < bits && ipv6Unicast.Contains(ipnet.IP) {
					ipv6 = append(ipv6, ipnet.IP.String())
				}
			case net.IPv4len:
				ipv4 = append(ipv4, ipnet.IP.String())
			}
		}

		if len(ipv4) > 0 {
			ipv4NetInterfaces = append(
				ipv4NetInterfaces,
				NetInterface{
					Name:    netInter.Name,
					Address: ipv4,
				},
			)
		}

		if len(ipv6) > 0 {
			ipv6NetInterfaces = append(
				ipv6NetInterfaces,
				NetInterface{
					Name:    netInter.Name,
					Address: ipv6,
				},
			)
		}
	}

	return ipv4NetInterfaces, ipv6NetInterfaces, nil
}
