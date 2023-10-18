package set

import (
	"fmt"
	"net"
	"net/netip"
	"strconv"
	"strings"
	"time"
)

// SetData is a struct that is used to create elements of a given set based on the key type of the set
type SetData struct {
	Port              uint16
	PortRangeStart    uint16
	PortRangeEnd      uint16
	Address           netip.Addr
	AddressRangeStart netip.Addr
	AddressRangeEnd   netip.Addr
	Prefix            netip.Prefix
	Timeout           time.Duration
}

// Convert a string address to the SetData type
func AddressStringToSetData(addressString string, timeout ...time.Duration) (SetData, error) {
	address, err := netip.ParseAddr(addressString)
	if err != nil {
		return SetData{}, err
	}

	var t time.Duration
	if len(timeout) > 0 {
		t = timeout[0]
	}
	return SetData{Address: address, Timeout: t}, nil
}

// Convert a string prefix/CIDR to the SetData type
func PrefixStringToSetData(prefixString string, timeout ...time.Duration) (SetData, error) {
	prefix, err := netip.ParsePrefix(prefixString)
	if err != nil {
		return SetData{}, err
	}

	var t time.Duration
	if len(timeout) > 0 {
		t = timeout[0]
	}
	return SetData{Prefix: prefix, Timeout: t}, nil
}

// Convert a string address range to the SetData type
func AddressRangeStringToSetData(startString string, endString string, timeout ...time.Duration) (SetData, error) {
	start, err := netip.ParseAddr(startString)
	if err != nil {
		return SetData{}, err
	}

	end, err := netip.ParseAddr(endString)
	if err != nil {
		return SetData{}, err
	}

	var t time.Duration
	if len(timeout) > 0 {
		t = timeout[0]
	}
	return SetData{
		AddressRangeStart: start,
		AddressRangeEnd:   end,
		Timeout:           t,
	}, nil
}

// Convert a list of string addresses to the SetData type
func AddressStringsToSetData(addressStrings []string, timeout ...time.Duration) ([]SetData, error) {
	data := []SetData{}

	for _, addressString := range addressStrings {
		if strings.Contains(addressString, "/") {
			// if it includes / we assume prefix i.e. 4.4.4.4/32
			prefix, err := PrefixStringToSetData(addressString, timeout...)
			if err != nil {
				return data, err
			}
			data = append(data, prefix)
			continue
		}
		if strings.Contains(addressString, "-") {
			// if it includes - we assume a range i.e. 10.10.10.10-10.10.10.15
			split := strings.Split(addressString, "-")
			addressRange, err := AddressRangeStringToSetData(split[0], split[1], timeout...)
			if err != nil {
				return data, err
			}
			data = append(data, addressRange)
			continue
		}
		// if we get here assume its just a normal IP i.e. 1.1.1.1
		address, err := AddressStringToSetData(addressString, timeout...)
		if err != nil {
			return data, err
		}
		data = append(data, address)

	}

	return data, nil
}

// Convert a string port to the SetData type
func PortStringToSetData(portString string, timeout ...time.Duration) (SetData, error) {
	port, err := strconv.ParseUint(portString, 10, 16)
	if err != nil {
		return SetData{}, err
	}

	var t time.Duration
	if len(timeout) > 0 {
		t = timeout[0]
	}
	return SetData{Port: uint16(port), Timeout: t}, nil
}

// Convert a string port range to the SetData type
func PortRangeStringToSetData(startString string, endString string, timeout ...time.Duration) (SetData, error) {
	start, err := strconv.ParseUint(startString, 10, 16)
	if err != nil {
		return SetData{}, err
	}

	end, err := strconv.ParseUint(endString, 10, 16)
	if err != nil {
		return SetData{}, err
	}

	var t time.Duration
	if len(timeout) > 0 {
		t = timeout[0]
	}
	return SetData{
		PortRangeStart: uint16(start),
		PortRangeEnd:   uint16(end),
		Timeout:        t,
	}, nil
}

// Convert a list string ports to the SetData type
func PortStringsToSetData(portStrings []string, timeout ...time.Duration) ([]SetData, error) {
	data := []SetData{}

	for _, portString := range portStrings {
		if strings.Contains(portString, "-") {
			// if it includes - we assume a range i.e. 10000-30000
			split := strings.Split(portString, "-")
			portRange, err := PortRangeStringToSetData(split[0], split[1], timeout...)
			if err != nil {
				return data, err
			}
			data = append(data, portRange)
		} else {
			// assume its just a normal port i.e. 80
			port, err := PortStringToSetData(portString, timeout...)
			if err != nil {
				return data, err
			}
			data = append(data, port)
		}
	}

	return data, nil
}

// Convert net.IPNet to the SetData type
func NetIPNetToSetData(net *net.IPNet, timeout ...time.Duration) (SetData, error) {
	ones, _ := net.Mask.Size()
	ip, ok := netip.AddrFromSlice(net.IP)

	if !ok {
		return SetData{}, fmt.Errorf("could not parse %v", net.String())
	}

	var t time.Duration
	if len(timeout) > 0 {
		t = timeout[0]
	}
	return SetData{Prefix: netip.PrefixFrom(ip, ones), Timeout: t}, nil
}

// Convert a list of net.IPNet to the SetData type
func NetIPNetsToSetData(nets []*net.IPNet, timeout ...time.Duration) ([]SetData, error) {
	data := []SetData{}

	for _, net := range nets {
		prefix, err := NetIPNetToSetData(net, timeout...)
		if err != nil {
			return data, err
		}
		data = append(data, prefix)
	}

	return data, nil
}

// Convert net.IP to the SetData type
func NetIPToSetData(ip net.IP, timeout ...time.Duration) (SetData, error) {
	netip, ok := netip.AddrFromSlice(ip)
	if !ok {
		return SetData{}, fmt.Errorf("could not parse ip: %v", ip)
	}

	var t time.Duration
	if len(timeout) > 0 {
		t = timeout[0]
	}
	return SetData{Address: netip, Timeout: t}, nil
}

// Convert a list of net.IP to the SetData type
func NetIPsToSetData(ips []net.IP, timeout ...time.Duration) ([]SetData, error) {
	data := []SetData{}

	for _, ip := range ips {
		netip, err := NetIPToSetData(ip, timeout...)
		if err != nil {
			return data, err
		}
		data = append(data, netip)
	}

	return data, nil
}

// Convert a list of netip.Addr to SetData type
func NetipAddrsToSetData(ips []netip.Addr, timeout ...time.Duration) ([]SetData, error) {
	data := []SetData{}

	for _, ip := range ips {
		netip, err := NetipAddrToSetData(ip, timeout...)
		if err != nil {
			return data, err
		}
		data = append(data, netip)
	}

	return data, nil
}

// Convert netip.Addr to SetData type
func NetipAddrToSetData(ip netip.Addr, timeout ...time.Duration) (SetData, error) {
	var t time.Duration
	if len(timeout) > 0 {
		t = timeout[0]
	}
	return SetData{Address: ip, Timeout: t}, nil
}

// Convert a list of netip.Prefix to SetData type
func NetipPrefixesToSetData(prefixes []netip.Prefix, timeout ...time.Duration) ([]SetData, error) {
	data := []SetData{}

	for _, prefix := range prefixes {
		netip, err := NetipPrefixToSetData(prefix, timeout...)
		if err != nil {
			return data, err
		}
		data = append(data, netip)
	}

	return data, nil
}

// Convert netip.Prefix to SetData type
func NetipPrefixToSetData(prefix netip.Prefix, timeout ...time.Duration) (SetData, error) {
	var t time.Duration
	if len(timeout) > 0 {
		t = timeout[0]
	}
	return SetData{Prefix: prefix, Timeout: t}, nil
}

// Convert a list of netip.AddrPort to SetData type, returns a list of addresses and a list of ports
func NetipAddrPortsToSetData(addrports []netip.AddrPort, timeout ...time.Duration) ([]SetData, []SetData, error) {
	addrs := []SetData{}
	ports := []SetData{}

	for _, addrport := range addrports {
		addr, port, err := NetipAddrPortToSetData(addrport, timeout...)
		if err != nil {
			return addrs, ports, err
		}
		addrs = append(addrs, addr)
		ports = append(ports, port)
	}

	return addrs, ports, nil
}

// Convert netip.AddrPort to SetData type, returns a address and a port
func NetipAddrPortToSetData(addrport netip.AddrPort, timeout ...time.Duration) (SetData, SetData, error) {
	var t time.Duration
	if len(timeout) > 0 {
		t = timeout[0]
	}
	return SetData{Address: addrport.Addr()}, SetData{Port: uint16(addrport.Port()), Timeout: t}, nil
}
