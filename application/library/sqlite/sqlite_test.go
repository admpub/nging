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
	"testing"

	"github.com/admpub/nging/application/library/common"
)

var sqlStr = `CREATE TABLE ` + "`" + `forever_process` + "`" + ` (
	` + "`" + `id` + "`" + ` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
	` + "`" + `name` + "`" + ` varchar(60) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '名称',
	` + "`" + `command` + "`" + ` varchar(300) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '命令',
	` + "`" + `workdir` + "`" + ` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '工作目录',
	` + "`" + `env` + "`" + ` text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '环境变量',
	` + "`" + `args` + "`" + ` text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '命令参数',
	` + "`" + `pidfile` + "`" + ` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT 'PID记录文件',
	` + "`" + `logfile` + "`" + ` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '日志记录文件',
	` + "`" + `errfile` + "`" + ` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '错误记录文件',
	` + "`" + `respawn` + "`" + ` int(11) unsigned NOT NULL DEFAULT '1' COMMENT '重试次数(进程被外部程序结束后自动启动)',
	` + "`" + `delay` + "`" + ` varchar(30) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '延迟启动(例如1ms/1s/1m/1h)',
	` + "`" + `ping` + "`" + ` varchar(30) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '心跳时间(例如1ms/1s/1m/1h)',
	` + "`" + `pid` + "`" + ` int(8) NOT NULL DEFAULT '0' COMMENT 'PID',
	` + "`" + `status` + "`" + ` enum('started','running','stopped','restarted','exited','killed','idle') CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT 'idle' COMMENT '进程运行状态',
	` + "`" + `debug` + "`" + ` enum('Y','N') CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT 'N' COMMENT 'DEBUG',
	` + "`" + `disabled` + "`" + ` enum('Y','N') CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT 'N' COMMENT '是否禁用',
	` + "`" + `created` + "`" + ` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
	` + "`" + `updated` + "`" + ` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '修改时间',
	` + "`" + `error` + "`" + ` varchar(300) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '错误信息',
	` + "`" + `lastrun` + "`" + ` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '上次运行时间',
	` + "`" + `description` + "`" + ` varchar(500) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '说明',
	PRIMARY KEY (` + "`" + `id` + "`" + `),
	UNIQUE KEY ` + "`" + `name` + "`" + ` (` + "`" + `name` + "`" + `)
  ) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='持久进程';`

func TestMySQLToSQLite(t *testing.T) {
	err := common.ParseSQL(sqlStr, false, func(sql string) error {
		sqls, err := covertCreateTableSQL(sql)
		if err != nil {
			return err
		}
		for _, sql := range sqls {
			fmt.Println(sql)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
}
