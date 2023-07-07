package osinfo

import (
	"os/exec"
	"regexp"
	"runtime"
	"strings"
)

var (
	nameField  = regexp.MustCompile(`NAME=(.*?)\n|\nNAME=(.*?)\n`)
	pnameField = regexp.MustCompile(`PRETTY_NAME=(.*?)\n|\nPRETTY_NAME=(.*?)\n`)
	verField   = regexp.MustCompile(`VERSION_ID=(.*?)\n|\nVERSION_ID=(.*?)\n`)
)

type LinuxMatcher struct {
	MatchRegexp *regexp.Regexp
	Distro      string
	PkgManager  string
	Version     string
}

var linuxMatchers = []LinuxMatcher{suse, debian, rhl, freebsd, arch, alpine, openwrt}

func RegisterLinuxMatcher(m LinuxMatcher) {
	linuxMatchers = append(linuxMatchers, m)
}

var (
	suse = LinuxMatcher{
		MatchRegexp: regexp.MustCompile(`SLES|openSUSE`),
		Distro:      `opensuse`,
		PkgManager:  `zypper`,
	}
	debian = LinuxMatcher{
		MatchRegexp: regexp.MustCompile(`Debian|Ubuntu|Kali|Parrot|Mint`),
		Distro:      `debian`,
		PkgManager:  `apt`,
	}
	rhl = LinuxMatcher{
		MatchRegexp: regexp.MustCompile(`Red Hat|CentOS|Fedora|Oracle`),
		Distro:      `fedora`,
		PkgManager:  `yum`,
	}
	freebsd = LinuxMatcher{
		MatchRegexp: regexp.MustCompile(`FreeBSD`),
		Distro:      `freebsd`,
		PkgManager:  `pkg`,
	}
	arch = LinuxMatcher{
		MatchRegexp: regexp.MustCompile(`Arch|Manjaro`),
		Distro:      `arch`,
		PkgManager:  `pacman`,
		Version:     `rolling`,
	}
	alpine = LinuxMatcher{
		MatchRegexp: regexp.MustCompile(`Alpine`),
		Distro:      `alpine`,
		PkgManager:  `apk`,
	}
	openwrt = LinuxMatcher{
		MatchRegexp: regexp.MustCompile(`OpenWrt`),
		Distro:      `openwrt`,
		PkgManager:  `opkg`,
	}
)

// GetVersion Linux returns version info
// fetching os info for linux distros is more complicated process than it should be
// version information is collected via `uname` and `/etc/os-release`
// Returns:
//   - r.Runtime
//   - r.Arch
//   - r.Name
//   - r.Version
//   - r.Linux.Kernel
//   - r.Linux.Distro
//   - r.Linux.PkgMng
func GetVersion() Release {
	info := Release{
		Runtime: runtime.GOOS,
		Arch:    runtime.GOARCH,
		Name:    "unknown",
		Version: "unknown",

		Linux: linuxRelease{
			Distro:     "unknown",
			Kernel:     "unknown",
			PkgManager: "unknown",
		},
	}

	out, err := exec.Command("uname", "-r").Output()
	if err == nil {
		info.Linux.Kernel = cleanString(string(out))
	}

	if !pathExists("/etc/os-release") {
		return info
	}

	f := readFile("/etc/os-release")
	var (
		namef  = cleanString(nameField.FindString(f))
		pnamef = cleanString(pnameField.FindString(f))
		verf   = cleanString(verField.FindString(f))
	)

	if pnamef := strings.Split(pnamef, "="); len(pnamef) >= 2 {
		info.Name = pnamef[1]
	}

	if verf := strings.Split(verf, "="); len(verf) >= 2 {
		info.Version = verf[1]
	}

	var name string
	if namef := strings.Split(namef, "="); len(namef) >= 2 {
		name = namef[1]
	} else {
		return info
	}

	for _, m := range linuxMatchers {
		if m.MatchRegexp.MatchString(name) {
			info.Linux.Distro = m.Distro
			info.Linux.PkgManager = m.PkgManager
			if len(m.Version) > 0 {
				info.Version = m.Version
			}
			return info
		}
	}
	return info
}
