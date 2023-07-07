package packer

import (
	"fmt"

	"github.com/admpub/osinfo"
)

func DetectManager() (Manager, error) {
	osversion := osinfo.GetVersion()
	//fmt.Println(osversion.String())
	opsystem := osversion.Runtime
	mgrs, ok := managers[opsystem]
	if !ok {
		return empty, fmt.Errorf("%s is %w", opsystem, ErrUnsupported)
	}
	pkgMgr := osversion.Linux.PkgManager
	if len(pkgMgr) > 0 {
		for _, managers := range mgrs {
			for _, mgr := range managers {
				if mgr.Name == pkgMgr {
					if Check(mgr.Name) {
						return mgr, nil
					}
				}
			}
		}
	}
	distro := osversion.Linux.Distro
	list, ok := mgrs[distro]
	if !ok {
		if len(distro) == 0 {
			return empty, fmt.Errorf("%s is %w", opsystem, ErrUnsupported)
		}
		list, ok = mgrs[""]
		if !ok {
			return empty, fmt.Errorf("%s %s is %w", opsystem, distro, ErrUnsupported)
		}
	}
	for _, mgr := range list {
		if Check(mgr.Name) {
			return mgr, nil
		}
	}
	return empty, ErrNotFound
}
