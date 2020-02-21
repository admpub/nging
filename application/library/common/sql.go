package common

import (
	"strings"

	"github.com/webx-top/com"
)

func ParseSQL(sqlFile string, isFile bool, installer func(string) error) (err error) {
	var sqlStr string
	installFunction := func(line string) (rErr error) {
		if strings.HasPrefix(line, `--`) {
			return nil
		}
		if strings.HasPrefix(line, `/*`) && strings.HasSuffix(line, `*/;`) {
			return nil
		}
		line = strings.TrimSpace(line)
		sqlStr += line
		if strings.HasSuffix(line, `;`) && len(sqlStr) > 0 {
			defer func() {
				sqlStr = ``
			}()
			return installer(sqlStr)
		}
		return nil
	}
	if isFile {
		return com.SeekFileLines(sqlFile, installFunction)
	}
	sqlContent := sqlFile
	for _, line := range strings.Split(sqlContent, "\n") {
		err = installFunction(line)
		if err != nil {
			return err
		}
	}
	return err
}
