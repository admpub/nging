package iptables

import (
	"io"
	"net"
	"strings"
)

func (ipt *IPTables) ExecuteList(args []string) ([]string, error) {
	return ipt.executeList(args)
}

// Run runs an iptables command with the given arguments, ignoring
// any stdout output
func (ipt *IPTables) Run(args ...string) error {
	return ipt.run(args...)
}

// runWithOutput runs an iptables command with the given arguments,
// writing any stdout output to the given writer
func (ipt *IPTables) RunWithOutput(args []string, stdout io.Writer) error {
	return ipt.runWithOutput(args, stdout)
}

// StatsWithLineNumber lists rules including the byte and packet counts
func (ipt *IPTables) StatsWithLineNumber(table, chain string) ([]map[string]string, error) {
	args := []string{"-t", table, "-L", chain, "-n", "-v", "-x", "--line-number"}
	lines, err := ipt.executeList(args)
	if err != nil {
		return nil, err
	}

	appendSubnet := func(addr string) string {
		if strings.IndexByte(addr, byte('/')) < 0 {
			if strings.IndexByte(addr, '.') < 0 {
				return addr + "/128"
			}
			return addr + "/32"
		}
		return addr
	}

	ipv6 := ipt.proto == ProtocolIPv6

	rows := []map[string]string{}
	var keys []string
	for i, line := range lines {
		// Skip over chain name and field header
		if i < 2 {
			if i == 1 {
				line = strings.TrimSpace(line)
				keys = strings.Fields(line)
			}
			continue
		}

		// Fields:
		// 0=num 1=pkts 2=bytes 3=target 4=prot 5=opt 6=in 7=out 8=source 9=destination
		line = strings.TrimSpace(line)
		fields := strings.Fields(line)
		// The ip6tables verbose output cannot be naively split due to the default "opt"
		// field containing 2 single spaces.
		if ipv6 {
			// Check if field 7 is "opt" or "source" address
			dest := fields[7]
			ip, _, _ := net.ParseCIDR(dest)
			if ip == nil {
				ip = net.ParseIP(dest)
			}

			// If we detected a CIDR or IP, the "opt" field is empty.. insert it.
			if ip != nil {
				f := []string{}
				f = append(f, fields[:5]...)
				f = append(f, "  ") // Empty "opt" field for ip6tables
				f = append(f, fields[5:]...)
				fields = f
			}
		}

		// Adjust "source" and "destination" to include netmask, to match regular
		// List output
		fields[8] = appendSubnet(fields[8])
		fields[9] = appendSubnet(fields[9])

		// Combine "options" fields 10... into a single space-delimited field.
		options := fields[10:]
		fields = fields[:10]
		fields = append(fields, strings.Join(options, " "))

		fieldsSize := len(fields)
		row := map[string]string{}
		for index, key := range keys {
			if index < fieldsSize {
				row[key] = fields[index]
			}
		}
		if fieldsSize > len(keys) {
			row[`options`] = strings.Join(fields[len(keys):], " ")
		}

		rows = append(rows, row)
	}
	return rows, nil
}

// GetIptablesCommand returns the correct command for the given protocol, either "iptables" or "ip6tables".
func GetIptablesCommand(proto Protocol) string {
	return getIptablesCommand(proto)
}
