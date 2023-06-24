package nftablesutils

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// IPAddr returns default gw iface name, gw ip address
// and wan ip address.
func IPAddr() (string, net.IP, net.IP, error) {
	// See http://man7.org/linux/man-pages/man8/route.8.html
	const file = "/proc/net/route"
	f, err := os.Open(file)
	if err != nil {
		return ``, nil, nil, fmt.Errorf("can't access %s: %w", file, err)
	}
	defer f.Close()

	bytes, err := ioutil.ReadAll(f)
	if err != nil {
		return ``, nil, nil, fmt.Errorf("can't read %s: %w", file, err)
	}

	wanIface, gatewayIP, err := parseLinuxProcNetRoute(bytes)
	if err != nil {
		return wanIface, gatewayIP, nil, err
	}

	iface, err := net.InterfaceByName(wanIface)
	if err != nil {
		return wanIface, gatewayIP, nil, fmt.Errorf("can't get interface by name: %w", err)
	}

	addrs, err := iface.Addrs()
	if err != nil {
		return wanIface, gatewayIP, nil, fmt.Errorf("can't get iface addrs: %w", err)
	}

	wanIP := net.IP{}
	found := false
	for _, addr := range addrs {
		ipnet, ok := addr.(*net.IPNet)
		if !ok || ipnet.IP.DefaultMask() == nil {
			continue
		}
		wanIP = ipnet.IP.To4()
		found = true
		break
	}

	if !found {
		return wanIface, gatewayIP, nil, errors.New("can't found public ip")
	}

	return wanIface, gatewayIP, wanIP, nil
}

func parseLinuxProcNetRoute(f []byte) (string, net.IP, error) {
	/* /proc/net/route file:
	   Iface   Destination Gateway     Flags   RefCnt  Use Metric  Mask
	   eno1    00000000    C900A8C0    0003    0   0   100 00000000    0   00
	   eno1    0000A8C0    00000000    0001    0   0   100 00FFFFFF    0   00
	*/
	const (
		sep   = "\t" // field separator
		field = 2    // field containing hex gateway address
	)
	scanner := bufio.NewScanner(bytes.NewReader(f))
	if scanner.Scan() {
		// Skip header line
		if !scanner.Scan() {
			return "", nil, errors.New("invalid linux route file")
		}

		// get field containing gateway address
		tokens := strings.Split(scanner.Text(), sep)
		if len(tokens) <= field {
			return "", nil, errors.New("invalid linux route file")
		}
		gatewayHex := "0x" + tokens[field]
		wanIface := tokens[0]

		// cast hex address to uint32
		d, _ := strconv.ParseInt(gatewayHex, 0, 64)
		d32 := uint32(d)

		// make net.IP address from uint32
		ipd32 := make(net.IP, 4)
		binary.LittleEndian.PutUint32(ipd32, d32)

		// format net.IP to dotted ipV4 string
		return wanIface, net.IP(ipd32), nil
	}
	return "", nil, errors.New("failed to parse linux route file")
}

func IPv6Addr() (string, net.IP, net.IP, error) {
	cmd := exec.Command(`ip`, `-6`, `route`)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return ``, nil, nil, err
	}
	/* output:
	fe80::/64 dev eth0 proto kernel metric 256 pref medium
	multicast ff00::/8 dev eth0 proto kernel metric 256 pref medium
	*/
	var wanIface string
	var gatewayIP net.IP
	for _, lineBytes := range bytes.Split(output, []byte(`\n`)) {
		parts := bytes.Split(lineBytes, []byte(` `))
		gatewayIP, _, err = net.ParseCIDR(string(parts[0]))
		if err == nil {
			wanIface = string(parts[2])
			break
		}
	}

	iface, err := net.InterfaceByName(wanIface)
	if err != nil {
		return wanIface, gatewayIP, nil, fmt.Errorf("can't get interface by name: %w", err)
	}

	addrs, err := iface.Addrs()
	if err != nil {
		return wanIface, gatewayIP, nil, fmt.Errorf("can't get iface addrs: %w", err)
	}

	wanIP := net.IP{}
	found := false
	for _, addr := range addrs {
		ipnet, ok := addr.(*net.IPNet)
		//fmt.Printf("================>%s \n", ipnet.IP)
		if !ok || ipnet.IP.IsGlobalUnicast() || ipnet.IP.To4() != nil {
			continue
		}
		wanIP = ipnet.IP
		found = true
		break
	}

	if !found {
		return wanIface, gatewayIP, nil, errors.New("can't found public ip")
	}

	return wanIface, gatewayIP, wanIP, nil
}
