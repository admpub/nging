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
	exec := func(sqlFile string, callback func(strLen int)) func(string) error {
		return common.SQLLineParser(func(sqlStr string) error {
			_, err := m.exec(sqlStr, dbfactory)
			if err != nil {
				noticer.Failure(`[FAILURE] ` + err.Error() + `: ` + com.HTMLEncode(sqlStr) + `: ` + filepath.Base(sqlFile))
			} else {
				noticer.Success(`[SUCCESS] ` + filepath.Base(sqlFile))
			}
			callback(len(sqlStr))
			return err
		})
	}
	for _, sqlFile := range sqlFiles {
		err = execSQLFileWithProgress(noticer, sqlFile, exec)
		if err != nil {
			return
		}
	}
	return err
}

func execSQLFileWithProgress(
	noticer *notice.NoticeAndProgress,
	sqlFile string,
	exec func(sqlFile string, callback func(strLen int)) func(string) error,
) error {
	fileSize, _ := com.FileSize(sqlFile)
	var totalPercent int64
	var perByteVal float64
	if fileSize > 0 {
		perByteVal = float64(100) / float64(fileSize)
	}
	var callback = func(strLen int) {
		if perByteVal > 0 {
			percent := int64(perByteVal * float64(strLen))
			totalPercent += percent
			noticer.Done(percent)
		}
	}
	err := com.SeekFileLines(sqlFile, exec(sqlFile, callback))
	if totalPercent < 100 {
		noticer.Done(100 - totalPercent)
	}
	return err
}
