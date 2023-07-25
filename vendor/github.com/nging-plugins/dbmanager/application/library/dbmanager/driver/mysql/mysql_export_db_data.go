package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/admpub/errors"
	"github.com/webx-top/com"
	"github.com/webx-top/db"

	"github.com/admpub/nging/v5/application/library/config"
	"github.com/admpub/nging/v5/application/library/notice"
	"github.com/nging-plugins/dbmanager/application/library/dbmanager/driver"
)

const maxInsertBytes = 1024 * 1024

// exportDBData 导出表数据s
func (m *mySQL) exportDBData(ctx context.Context, noticer notice.Noticer,
	cfg *driver.DbAuth, tables []string, dataWriter interface{}, mysqlVersion string) error {
	if !com.InSlice(cfg.Charset, Charsets) {
		return fmt.Errorf(`字符集charset值无效: %v`, cfg.Charset)
	}
	var (
		w   io.Writer
		err error
	)
	switch v := dataWriter.(type) {
	case io.Writer:
		w = v
	case string:
		dir := filepath.Dir(v)
		err = com.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return fmt.Errorf(`failed to backup: %v`, err)
		}
		w, err = os.Create(v)
		if err != nil {
			return fmt.Errorf(`failed to backup: %v`, err)
		}
	default:
		return errors.Wrapf(db.ErrUnsupported, `SQL Writer Error: %T`, v)
	}
	var (
		selectFuncs []string
		selectCols  []string
		wheres      []string
		orderFields []string
		descs       []string
		page        int = 1
		limit       int = -1 // 不限数量
		exportStyle string
	)
	_, err = w.Write([]byte(`-- Nging DBManager MySQL dump ` + config.Version.Number + `, for ` + runtime.GOOS + ` (` + runtime.GOARCH + `)
--
-- Host: ` + cfg.Host + `    Database: ` + cfg.Db + `
-- ------------------------------------------------------
-- Server version	` + mysqlVersion + `

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES ` + cfg.Charset + ` */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

`))
	if err != nil {
		return err
	}
	dbfactory, err := connect(cfg)
	if err != nil {
		return err
	}
	defer dbfactory.CloseAll()
	for _, table := range tables {
		fields, _, err := m.tableFields(table, dbfactory)
		if err != nil {
			return err
		}
		var (
			insert     string
			suffix     string
			hasValues  bool
			totalBytes int
		)
		_, _, _, err = m.listData(dbfactory, func(cols []string, row map[string]*sql.NullString) error {
			if len(insert) == 0 {
				keys := make([]string, len(cols))
				vals := make([]string, len(cols))
				for idx, key := range cols {
					key = quoteCol(key)
					keys[idx] = key
					vals[idx] = key + " = VALUES(" + key + ")"
				}
				if exportStyle == `INSERT+UPDATE` {
					suffix = " ON DUPLICATE KEY UPDATE " + strings.Join(vals, ", ")
				} else {
					suffix = ""
				}
				suffix += ";\n"
				insert = "INSERT INTO " + quoteCol(table) + " (" + strings.Join(keys, `, `) + ") VALUES "
			}
			var values, sep string
			for _, col := range cols {
				val := row[col]
				if !val.Valid {
					values += sep + `NULL`
				} else {
					field, ok := fields[col]
					var v string
					if ok && reFieldTypeNumber.MatchString(field.Type) && len(val.String) > 0 && !strings.HasPrefix(field.Full_type, `[`) {
						v = val.String
					} else {
						v = field.Format(val.String)
						v = quoteVal(v)
						v = com.AddRSlashes(v)
					}
					values += sep + unconvertField(field, v)
				}
				sep = `, `
			}
			s := "(" + values + ")"
			if !hasValues {
				totalBytes += len(s)
				s = insert + s
				hasValues = true
			} else if totalBytes > maxInsertBytes {
				totalBytes = len(s)
				s = suffix + insert + s
			} else {
				s = "," + s
				totalBytes += len(s)
			}
			_, err = w.Write(com.Str2bytes(s))
			return err
		}, table, selectFuncs, selectCols, wheres, orderFields, descs, page, limit, 0, false)
		if err == nil && hasValues {
			_, err = w.Write(com.Str2bytes(suffix))
		}
		if err != nil {
			return err
		}
	}

	_, err = w.Write([]byte(`/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on ` + time.Now().Format("2006-01-02 15:04:05") + `

`))
	if c, y := w.(io.Closer); y {
		c.Close()
	}
	return err
}
