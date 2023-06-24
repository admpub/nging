package osinfo

import (
	"os/exec"
	"regexp"
	"runtime"
	"strings"
)

// GetVersion Linux returns version info
// fetching os info for linux distros is more complicated process than it should be
// version information is collected via `uname` and `/etc/os-release`
// Returns:
//		- r.Runtime
//		- r.Arch
//		- r.Name
//		- r.Version
//		- r.Linux.Kernel
//		- r.Linux.Distro
//		- r.Linux.PkgMng
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

	var (
		nameField  = regexp.MustCompile(`NAME=(.*?)\n|\nNAME=(.*?)\n`)
		pnameField = regexp.MustCompile(`PRETTY_NAME=(.*?)\n|\nPRETTY_NAME=(.*?)\n`)
		verField   = regexp.MustCompile(`VERSION_ID=(.*?)\n|\nVERSION_ID=(.*?)\n`)
	)

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

	var (
		suse   = regexp.MustCompile(`SLES|openSUSE`)
		debian = regexp.MustCompile(`Debian|Ubuntu|Kali|Parrot|Mint`)
		rhl    = regexp.MustCompile(`Red Hat|CentOS|Fedora|Oracle`)
		arch   = regexp.MustCompile(`Arch|Manjaro`)
		alpine = regexp.MustCompile(`Alpine`)
	)

	var name string
	if namef := strings.Split(namef, "="); len(namef) >= 2 {
		name = namef[1]
	} else {
		return info
	}

	switch {
	case suse.MatchString(name):
		info.Linux.Distro = "opensuse"
		info.Linux.PkgManager = "zypper"

	case debian.MatchString(name):
		info.Linux.Distro = "debian"
		info.Linux.PkgManager = "apt"

	case rhl.MatchString(name):
		info.Linux.Distro = "fedora"
		info.Linux.PkgManager = "yum"

	case arch.MatchString(name):
		info.Version = "rolling"
		info.Linux.Distro = "arch"
		info.Linux.PkgManager = "pacman"

	case alpine.MatchString(name):
		info.Linux.Distro = "alpine"
		info.Linux.PkgManager = "apk"
	}

	return info
}
