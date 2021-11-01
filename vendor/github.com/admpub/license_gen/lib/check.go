package lib

import (
	"strings"

	"github.com/webx-top/com"
)

func CheckMAC(hashedMAC string) error {
	addrs, err := MACAddresses(false)
	if err != nil {
		return err
	}
	for _, addr := range addrs {
		if hashedMAC == Hash(addr) {
			return nil
		}
	}
	return InvalidMachineID
}

func CheckVersion(version string, versionRule string) bool {
	if len(versionRule) == 0 {
		return true
	}
	if len(versionRule) < 2 {
		return versionRule == version
	}
	switch versionRule[0] {
	case '>':
		if len(versionRule) > 2 {
			if versionRule[1] == '=' {
				return com.VersionComparex(version, versionRule[2:], `>=`)
			}
		}
		return com.VersionComparex(version, versionRule[1:], `>`)
	case '<':
		if len(versionRule) > 2 {
			if versionRule[1] == '=' {
				return com.VersionComparex(version, versionRule[2:], `<=`)
			}
		}
		return com.VersionComparex(version, versionRule[1:], `<`)
	case '!':
		if len(versionRule) > 2 {
			if versionRule[1] == '=' {
				return versionRule[2:] != version
			}
		}
		return versionRule[1:] != version
	case '=':
		return versionRule[1:] == version
	default:
		return versionRule == version
	}
}

func CheckDomain(fullDomain string, rootDomain string) bool {
	rootDomain = strings.Trim(rootDomain, `.`)
	rootParts := strings.Split(rootDomain, `.`)
	fullParts := strings.Split(fullDomain, `.`)
	l := len(fullParts) - len(rootParts)
	if l < 0 {
		return false
	}
	for i, j := 0, len(rootParts); i < j; i++ {
		if rootParts[i] != fullParts[i+l] {
			return false
		}
	}
	return true
}
