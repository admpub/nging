package mysql

import (
	"context"
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
	"github.com/nging-plugins/dbmanager/application/library/dbmanager/driver/mysql/utils"
)

// exportDBStruct 导出表结构
func (m *mySQL) exportDBStruct(ctx context.Context, noticer notice.Noticer,
	cfg *driver.DbAuth, tables []string, structWriter interface{}, mysqlVersion string, resetAutoIncrements ...bool) error {
	if !com.InSlice(cfg.Charset, Charsets) {
		return fmt.Errorf(`字符集charset值无效: %v`, cfg.Charset)
	}
	var (
		w                  io.Writer
		err                error
		onFinish           func(string) string
		resetAutoIncrement bool
	)
	if len(resetAutoIncrements) > 0 {
		resetAutoIncrement = resetAutoIncrements[0]
	}
	switch v := structWriter.(type) {
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
		onFinish = func(ddl string) string {
			if resetAutoIncrement {
				return utils.RemoveAutoIncrementValue(v)
			}
			return ddl
		}
	default:
		return errors.Wrapf(db.ErrUnsupported, `SQL Writer Error: %T`, v)
	}
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
		ddl, err := m.tableDDL(table, dbfactory)
		if err != nil {
			return err
		}
		if onFinish != nil {
			ddl = onFinish(ddl)
		}
		if !strings.HasSuffix(ddl, `;`) {
			ddl += `;`
		}
		_, err = w.Write([]byte("DROP TABLE IF EXISTS " + quoteCol(table) + ";\n"))
		if err != nil {
			return err
		}
		_, err = w.Write([]byte(`/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = ` + cfg.Charset + ` */;
`))
		_, err = w.Write([]byte(ddl + "\n/*!40101 SET character_set_client = @saved_cs_client */;\n\n"))
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
