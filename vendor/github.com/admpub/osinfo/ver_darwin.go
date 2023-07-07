package osinfo

import (
	"os/exec"
	"runtime"
	"strings"
)

// GetVersion Darwin returns version info
// fetching info for this os is fairly simple
// version information is all fetched via `sw_vers`
// Returns:
//   - r.Runtime
//   - r.Arch
//   - r.Name
//   - r.Version
//   - r.Mac.VersionName
func GetVersion() Release {
	info := Release{
		Runtime: runtime.GOOS,
		Arch:    runtime.GOARCH,
		Name:    "Mac OS X",
		Version: "unknown",

		MacOs: macOsRelease{
			VersionName: "unknown",
		},
	}

	version, err := exec.Command("sw_vers").Output()
	if err == nil {
		str := strings.Split(string(version), "\n")
		for _, s := range str {
			if strings.HasPrefix(s, "ProductVersion:\t") {
				info.Version = strings.TrimPrefix(s, "ProductVersion:\t")
			}
		}
	}

	var name string
	idx := strings.LastIndex(info.Version, ".")
	ver := info.Version[0:idx]
	switch ver {
	case "10.6":
		name = "MacOS - Snow Leopard"
	case "10.7":
		name = "MacOS - Lion"
	case "10.8":
		name = "MacOS - Mountain Lion"
	case "10.9":
		name = "MacOS - Mavericks"
	case "10.10":
		name = "MacOS - Yosemite"
	case "10.11":
		name = "MacOS - El Capitan"
	case "10.12":
		name = "MacOS - Sierra"
	case "10.13":
		name = "MacOS - High Sierra"
	case "10.14":
		name = "MacOS - Mojave"
	case "10.15":
		name = "MacOS - Catalina"
	default:
		switch strings.SplitN(ver, `.`, 2)[0] {
		case "11":
			name = "MacOS - Big Sur"
		case "12":
			name = "MacOS - Monterey"
		case "13":
			name = "MacOS - Ventura"
		}
	}
	info.MacOs = macOsRelease{
		VersionName: name,
	}
	return info
}
