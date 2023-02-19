package sqlite

import (
	"bytes"
	"strings"

	"github.com/webx-top/com"
)

func isSameSchemaItem(src, dest string) bool {
	equal := src == dest
	return equal
}

func parseIndex(line string) (string, bool) {
	indexMatches := indexReg.FindStringSubmatch(line)
	var matched bool
	if len(indexMatches) > 0 {
		if len(indexMatches[2]) == 0 {
			line = indexMatches[1] + `IF NOT EXISTS ` + strings.TrimPrefix(line, indexMatches[1]) // 强制加“IF NOT EXISTS”
		}
		matched = true
	}
	return line, matched
}

func ParseIndexDDL(schemas string) []string {
	var indexes []string
	com.SeekLines(bytes.NewBufferString(schemas), func(line string) error {
		if len(line) == 0 {
			return nil
		}
		var matched bool
		line, matched = parseIndex(line)
		if !matched {
			return nil
		}
		line = strings.TrimSuffix(line, `;`)
		indexes = append(indexes, line)
		return nil
	})
	return indexes
}
