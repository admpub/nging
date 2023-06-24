package mysql

import (
	"context"
	"path/filepath"
	"strings"

	"github.com/admpub/errors"
	"github.com/admpub/nging/v5/application/library/notice"
	"github.com/nging-plugins/dbmanager/application/library/dbmanager/driver"
	"github.com/nging-plugins/dbmanager/application/library/dbmanager/driver/mysql/utils"
	"github.com/webx-top/com"
)

func (m *mySQL) importDB(c context.Context, noticer *notice.NoticeAndProgress,
	cfg *driver.DbAuth, cacheDir string, files []string) (err error) {
	if !com.InSlice(cfg.Charset, Charsets) {
		return errors.New(`字符集charset值无效`)
	}
	names := make([]string, len(files))
	for i, file := range files {
		names[i] = filepath.Base(file)
	}
	noticer.Success(`开始导入: ` + strings.Join(names, ", "))
	var ifi *utils.ImportFile
	ifi, err = utils.ParseImportFile(cacheDir, files)
	if err != nil {
		return
	}
	dbfactory, err := connect(cfg)
	if err != nil {
		return err
	}
	defer dbfactory.CloseAll()
	for _, sqlStr := range []string{
		`SET NAMES ` + cfg.Charset + `;`,
		`SET FOREIGN_KEY_CHECKS=0;`,
		`SET UNIQUE_CHECKS=0;`,
	} {
		_, err := m.exec(sqlStr, dbfactory)
		if err != nil {
			noticer.Failure(`[FAILURE] ` + err.Error() + `: ` + sqlStr)
		} else {
			noticer.Success(`[SUCCESS] ` + sqlStr)
		}
	}
	noticer.Add(int64(ifi.AllSqlFileNum()))
	if len(ifi.StructFiles) > 0 {
		err = m.importDBStruct(c, noticer, dbfactory, ifi.StructFiles)
	}
	if err == nil && len(ifi.DataFiles) > 0 {
		err = m.importDBData(c, noticer, dbfactory, ifi.StructFiles)
	}
	for _, sqlStr := range []string{
		`SET FOREIGN_KEY_CHECKS=1;`,
		`SET UNIQUE_CHECKS=1;`,
	} {
		_, err := m.exec(sqlStr, dbfactory)
		if err != nil {
			noticer.Failure(`[FAILURE] ` + err.Error() + `: ` + sqlStr)
		} else {
			noticer.Success(`[SUCCESS] ` + sqlStr)
		}
	}
	return
}
