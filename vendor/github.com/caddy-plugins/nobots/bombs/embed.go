package bombs

import "embed"

//go:embed 1G.gzip
//go:embed 10G.gzip
//go:embed 1T.gzip
var Bombs embed.FS

var BombFileNameList = []string{`1G`, `10G`, `1T`}

func Exists(filename string) bool {
	for _, v := range BombFileNameList {
		if v == filename {
			return true
		}
	}
	return false
}
