package mysql

import (
	"context"
	"path/filepath"

	"github.com/admpub/nging/v5/application/library/common"
	"github.com/admpub/nging/v5/application/library/notice"
	"github.com/webx-top/com"
	"github.com/webx-top/db/lib/factory"
)

func (m *mySQL) exec(sqlStr string, dbfactory ...*factory.Factory) (int64, error) {
	result, err := m.newParam(dbfactory...).SetCollection(sqlStr).Exec()
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// importDBStruct 导出表结构
func (m *mySQL) importDBStruct(ctx context.Context, noticer *notice.NoticeAndProgress,
	dbfactory *factory.Factory, sqlFiles []string) (err error) {
	exec := func(sqlFile string) func(string) error {
		return common.SQLLineParser(func(sqlStr string) error {
			_, err := m.exec(sqlStr, dbfactory)
			if err != nil {
				noticer.Failure(`[FAILURE] ` + err.Error() + `: ` + com.HTMLEncode(sqlStr) + `: ` + filepath.Base(sqlFile))
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
