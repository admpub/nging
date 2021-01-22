/*
   Nging is a toolbox for webmasters
   Copyright (C) 2018-present  Wenhui Shen <swh@admpub.com>

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU Affero General Public License as published
   by the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Affero General Public License for more details.

   You should have received a copy of the GNU Affero General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package sqlite

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/admpub/nging/application/library/common"
)

var sqlStr = `CREATE TABLE ` + "`" + `nging_login_log` + "`" + ` (
	` + "`" + `owner_type` + "`" + ` enum('customer','user') COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'user' COMMENT '用户类型(user-后台用户;customer-前台客户)',
	` + "`" + `owner_id` + "`" + ` bigint unsigned NOT NULL DEFAULT '0' COMMENT '用户ID',
	` + "`" + `username` + "`" + ` varchar(60) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '登录名',
	` + "`" + `errpwd` + "`" + ` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '错误密码',
	` + "`" + `ip_address` + "`" + ` varchar(46) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT 'ip地址',
	` + "`" + `ip_location` + "`" + ` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT 'ip定位',
	` + "`" + `user_agent` + "`" + ` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '浏览器代理',
	` + "`" + `success` + "`" + ` enum('Y','N') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'N' COMMENT '是否登录成功',
	` + "`" + `failmsg` + "`" + ` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '失败信息',
	` + "`" + `day` + "`" + ` int unsigned NOT NULL DEFAULT '0' COMMENT '日期(Ymd)',
	` + "`" + `created` + "`" + ` int unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
	KEY ` + "`" + `ip_address` + "`" + ` (` + "`" + `ip_address` + "`" + `,` + "`" + `day` + "`" + `),
	KEY ` + "`" + `owner_type` + "`" + ` (` + "`" + `owner_type` + "`" + `,` + "`" + `owner_id` + "`" + `),
	KEY ` + "`" + `created` + "`" + ` (` + "`" + `created` + "`" + `)
  ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='登录日志';`

func parseTestSQL(sql string) error {
	sqls, err := covertCreateTableSQL(sql)
	if err != nil {
		return err
	}
	for _, sql := range sqls {
		fmt.Println(sql)
	}
	return nil
}

func TestMySQLToSQLite(t *testing.T) {
	fmt.Println(`============= multi-line:`)
	err := common.ParseSQL(sqlStr, false, parseTestSQL)
	if err != nil {
		panic(err)
	}
	fmt.Println(`============= single-line:`)
	err = parseTestSQL(sqlStr)
	if err != nil {
		panic(err)
	}
}

func TestMySQLToSQLiteFile(t *testing.T) {
	//return
	var err error
	err = ConvertMySQLFile(filepath.Join(os.Getenv("GOPATH"), `src/github.com/admpub/nging/config/install.sql`), `/Users/hank/Downloads/nging.sqlite.sql`)
	if err != nil {
		panic(err)
	}
}
