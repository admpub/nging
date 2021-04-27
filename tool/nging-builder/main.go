package main

import "os"

var targetNames = map[string]string{
	`linux_386`:     `linux/386`,
	`linux_amd64`:   `linux/amd64`,
	`linux_arm5`:    `linux/arm-5`,
	`linux_arm6`:    `linux/arm-6`,
	`linux_arm7`:    `linux/arm-7`,
	`linux_arm64`:   `linux/arm64`,
	`darwin_amd64`:  `darwin/amd64`,
	`windows_386`:   `windows/386`,
	`windows_amd64`: `windows/amd64`,
}

func main() {
	var targets []string
	var minify bool
	switch len(os.Args) {
	case 3:
		minify = os.Args[2] == `m` || os.Args[2] == `min`
		if t, y := targetNames[os.Args[1]]; y {
			targets = append(targets, t)
		}
	case 2:
		if os.Args[1] == `m` || os.Args[1] == `min` {
			minify = true
		} else {
			if t, y := targetNames[os.Args[1]]; y {
				targets = append(targets, t)
			}
		}
	case 1:
		for _, t := range targetNames {
			targets = append(targets, t)
		}
	}
	_ = minify
	_ = targets
}
