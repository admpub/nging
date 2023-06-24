package packer

var (
	apk     = Manager{"apk", "add", "update", "del"}
	apt     = Manager{"apt", "-y install", "update", "remove"}
	brew    = Manager{"brew", "install", "update", "remove"}
	dnf     = Manager{"dnf", "install", "upgrade", "erase"}
	flatpak = Manager{"flatpak", "install", "update", "uninstall"}
	snap    = Manager{"snap", "install", "upgrade", "remove"}
	pacman  = Manager{"pacman", "--noconfirm -S", "--noconfirm -Syuu", "--noconfirm -Rscn"}
	paru    = Manager{"paru", "-S", "-Syuu", "-R"}
	yay     = Manager{"yay", "-S", "-Syuu", "-R"}
	zypper  = Manager{"zypper", "-n install", "update", "-n remove"}
	choco   = Manager{"choco", "install -y", "update", "uninstall"} // adminstrator
	scoop   = Manager{"scoop", "install", "update", "uninstall"}
	winget  = Manager{"winget", "install", "source update", "uninstall"}
)

var managers = map[string]map[string][]Manager{
	"windows": {
		"": []Manager{winget, scoop, choco},
	},
	"darwin": {
		"": []Manager{brew},
	},
	"linux": {
		"arch":     []Manager{pacman, yay, paru},
		"alpine":   []Manager{apk},
		"fedora":   []Manager{dnf},
		"opensuse": []Manager{zypper},
		"debian":   []Manager{apt, snap},
		"":         []Manager{snap, flatpak},
	},
}

func Register(system string, distro string, mgr Manager) {
	if _, ok := managers[system]; !ok {
		managers[system] = map[string][]Manager{}
	}
	if _, ok := managers[system][distro]; !ok {
		managers[system][distro] = []Manager{mgr}
		return
	}
	for i, m := range managers[system][distro] {
		if m.Name == mgr.Name {
			managers[system][distro][i] = mgr
			return
		}
	}
	managers[system][distro] = append(managers[system][distro], mgr)
}
