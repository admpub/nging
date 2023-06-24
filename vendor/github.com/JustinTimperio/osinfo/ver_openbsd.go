package osinfo

import (
	"os/exec"
	"runtime"
)

// GetVersion OpenBSD returns version info
// fetching info for this os is fairly simple
// version information is all fetched via `uname`
// Returns:
//		- r.Runtime
//		- r.Arch
//		- r.Name
//		- r.Version
//		- r.BSD.Kernel
//		- r.BSD.PkgManager
func GetVersion() *Release {

	info := &Release{
		Runtime: runtime.GOOS,
		Arch:    runtime.GOARCH,
		Name:    "unknown",
		Version: "unknown",

		BSD: bsdRelease{
			Kernel:     "unknown",
			PkgManager: "pkg_add",
		},
	}

	fullName, err := exec.Command("uname", "-sr").Output()
	if err == nil {
		info.Name = cleanString(string(fullName))
	}

	version, err := exec.Command("uname", "-r").Output()
	if err == nil {
		info.Version = cleanString(string(version))
	}

	kernel, _ := exec.Command("uname", "-v").Output()
	if err == nil {
		info.BSD.Kernel = cleanString(string(kernel))
	}

	return info
}
