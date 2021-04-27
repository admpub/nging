package main

import (
	"os"
	"regexp"

	"github.com/webx-top/com"
)

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

var armRegexp = regexp.MustCompile(`/arm`)

func main() {
	var targets []string
	var armTargets []string
	addTarget := func(target string, notNames ...bool) {
		if len(notNames) == 0 || !notNames[0] {
			target = getTarget(target)
			if len(target) == 0 {
				return
			}
		}
		if armRegexp.MatchString(target) {
			armTargets = append(armTargets, target)
		} else {
			targets = append(targets, target)
		}
	}
	var minify bool
	switch len(os.Args) {
	case 3:
		minify = isMinified(os.Args[2])
		addTarget(os.Args[1])
	case 2:
		if isMinified(os.Args[1]) {
			minify = true
		} else {
			addTarget(os.Args[1])
		}
	case 1:
		for _, t := range targetNames {
			addTarget(t, true)
		}
	default:
		com.ExitOnFailure(`参数错误`)
	}
	_ = minify
	_ = targets
	_ = armTargets
}

func getTarget(target string) string {
	if t, y := targetNames[target]; y {
		return t
	}
	for _, t := range targetNames {
		if t == target {
			return t
		}
	}
	return ``
}

func isMinified(arg string) bool {
	return arg == `m` || arg == `min`
}
