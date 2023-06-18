package nftablesutils

import (
	"fmt"
	"net/netip"
)

// Validates start and end port numbers
func ValidatePortRange(start uint16, end uint16) error {
	if err := ValidatePort(start); err != nil {
		return err
	}

	if err := ValidatePort(end); err != nil {
		return err
	}

	if end < start {
		return fmt.Errorf("starting port (%v) is higher than ending port (%v)", start, end)
	}

	return nil
}

// Validates a port number
func ValidatePort(port uint16) error {
	if port < 1 {
		return fmt.Errorf("port (%v) less than 1", port)
	}

	if port > 65535 {
		return fmt.Errorf("port (%v) greater than 65535", port)
	}

	return nil
}

// Validates an IP address range
func ValidateAddressRange(start netip.Addr, end netip.Addr) error {
	if err := ValidateAddress(start); err != nil {
		return err
	}

	if err := ValidateAddress(end); err != nil {
		return err
	}

	if end.Less(start) {
		return fmt.Errorf("start address (%v) is after end address (%v)", start.String(), end.String())
	}

	return nil
}

// Validates an IP address
func ValidateAddress(ip netip.Addr) error {
	if !ip.IsValid() {
		return fmt.Errorf("address is zero")
	}

	if ip.IsUnspecified() {
		return fmt.Errorf("address is unspecified %v", ip.String())
	}

	return nil
}

// Validates a Prefix/CIDR
func ValidatePrefix(prefix netip.Prefix) error {
	if !prefix.IsValid() {
		return fmt.Errorf("prefix is invalid %v", prefix.String())
	}

	return nil
}
