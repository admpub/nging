package mysql

import (
	"context"
	"path/filepath"

	"github.com/admpub/nging/v5/application/library/common"
	"github.com/admpub/nging/v5/application/library/notice"
	"github.com/nging-plugins/dbmanager/application/library/dbmanager/driver"
	"github.com/webx-top/com"
)

// importDBData 导出表结构
func (m *mySQL) importDBData(ctx context.Context, noticer *notice.NoticeAndProgress,
	cfg *driver.DbAuth, sqlFiles []string) (err error) {
	exec := func(sqlFile string) func(string) error {
		return common.SQLLineParser(func(sqlStr string) error {
			_, err := m.exec(sqlStr)
			if err != nil {
				noticer.Failure(`[FAILURE] ` + err.Error() + `: ` + sqlStr + `: ` + filepath.Base(sqlFile))
			} else {
				noticer.Success(`[SUCCESS] ` + filepath.Base(sqlFile))
			}
			return err
		})
	}
	for _, sqlFile := range sqlFiles {
		err = com.SeekFileLines(sqlFile, exec(sqlFile))
		noticer.Done(1)
		if err != nil {
			return
		}
	}
	return err
}
