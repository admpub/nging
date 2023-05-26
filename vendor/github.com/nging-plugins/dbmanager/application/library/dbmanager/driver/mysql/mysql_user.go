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
package mysql

import (
	"database/sql"
	"errors"
	"strconv"
	"strings"

	"github.com/webx-top/com"
	"github.com/webx-top/echo"
)

func (m *mySQL) dropUser(user string, host string) error {
	if len(host) > 0 {
		user = quoteVal(user) + `@` + quoteVal(host)
	} else {
		user = quoteVal(user) + `@''`
	}
	r := &Result{}
	r.SQL = "DROP USER " + user
	r.Exec(m.newParam())
	m.AddResults(r)
	if r.err != nil {
		return r.err
	}
	r2 := &Result{}
	r2.SQL = "FLUSH PRIVILEGES"
	r2.Exec(m.newParam())
	m.AddResults(r2)
	return r2.err
}

func (m *mySQL) modifyPassword(user string, host string, password string) error {
	v8plus := m.isV8Plus()
	if len(password) == 0 {
		return errors.New(m.T(`密码不能为空`))
	}
	userAndHost := quoteVal(user) + `@` + quoteVal(host)
	r := &Result{}
	if v8plus {
		r.SQL = `ALTER USER ` + userAndHost + ` IDENTIFIED WITH MYSQL_NATIVE_PASSWORD BY ` + quoteVal(password)
	} else {
		r.SQL = `SET PASSWORD FOR ` + userAndHost + `=PASSWORD(` + quoteVal(password) + `)`
	}
	r.Exec(m.newParam())
	m.AddResults(r)
	return r.err
}

func (m *mySQL) addUser(user string, host string, password string) error {
	v8plus := m.isV8Plus()
	r := &Result{}
	if len(password) == 0 {
		return errors.New(m.T(`密码不能为空。请注意：修改用户名的时候，必须设置密码`))
	}
	userAndHost := quoteVal(user) + `@` + quoteVal(host)
	if v8plus {
		r.SQL = `CREATE USER ` + userAndHost + ` IDENTIFIED WITH MYSQL_NATIVE_PASSWORD BY ` + quoteVal(password)
	} else {
		r.SQL = `GRANT USAGE ON *.* TO`
		if com.VersionCompare(m.getVersion(), `5`) >= 0 {
			r.SQL = `CREATE USER`
		}
		r.SQL += ` ` + userAndHost + ` IDENTIFIED BY ` + quoteVal(password)
	}
	r.Exec(m.newParam())
	m.AddResults(r)
	return r.err
}

func (m *mySQL) isV8Plus() bool {
	if !m._isV8Plus.Valid {
		if !strings.Contains(m.getVersion(), `MariaDB`) {
			m._isV8Plus.Bool = com.VersionCompare(m.getVersion(), `8.0.11`) >= 0
		} else {
			m._isV8Plus.Bool = false
			// m._isV8Plus.Bool = com.VersionCompare(m.getVersion(), `10.6`) >= 0 // maybe
		}
		m._isV8Plus.Valid = true
	}
	return m._isV8Plus.Bool
}

func (m *mySQL) editUser(oldUser string, oldHost string, newUser string, newHost string, newPasswd string, modifyPassword bool) error {
	oldUserAndHost := quoteVal(oldUser) + `@` + quoteVal(oldHost)
	if len(newUser) == 0 {
		return errors.New(m.T(`用户名不能为空`))
	}

	_, grants, _, err := m.getUserGrants(oldHost, oldUser)
	if err != nil {
		return err
	}
	newUserAndHost := quoteVal(newUser) + `@` + quoteVal(newHost)
	var created bool
	onerror := func(err error) error {
		return err
	}
	if oldUserAndHost != newUserAndHost { // 新建账号
		created = true
		if err = m.addUser(newUser, newHost, newPasswd); err != nil {
			return err
		}
		onerror = func(err error) error {
			r2 := &Result{}
			r2.SQL = "DROP USER " + newUserAndHost
			r2.Exec(m.newParam())
			if r2.err != nil {
				m.Echo().Logger().Error(r2.err)
			}
			m.AddResults(r2)
			return err
		}
	} else if modifyPassword {
		if err = m.modifyPassword(newUser, newHost, newPasswd); err != nil {
			return err
		}
	}

	// 更改权限
	scopes := m.FormValues(`scopes[]`)       // 作用范围
	databases := m.FormValues(`databases[]`) // 数据库名
	tables := m.FormValues(`tables[]`)       // 数据表名
	columns := m.FormValues(`columns[]`)     // 列名
	scopeMaxIndex := len(scopes) - 1
	databaseMaxIndex := len(databases) - 1
	tableMaxIndex := len(tables) - 1
	columnMaxIndex := len(columns) - 1
	objects := m.FormValues(`objects[]`) // 赋予权限的对象
	newGrants := map[string]*Grant{}

	mapx := echo.NewMapx(m.Forms())
	mapx = mapx.Get(`grants`)
	logger := m.Echo().Logger()
	//objects: objects[0|1|...]=`*.*|db.*|db.table|db.table.col1,col2`
	for k, v := range objects {
		if k > scopeMaxIndex {
			logger.Debugf(`k > scopeMaxIndex: %v > %v`, k, scopeMaxIndex)
			continue
		}
		if k > databaseMaxIndex {
			logger.Debugf(`k > databaseMaxIndex: %v > %v`, k, databaseMaxIndex)
			continue
		}
		if k > tableMaxIndex {
			logger.Debugf(`k > tableMaxIndex: %v > %v`, k, tableMaxIndex)
			continue
		}
		if k > columnMaxIndex {
			logger.Debugf(`k > columnMaxIndex: %v > %v`, k, columnMaxIndex)
			continue
		}
		if len(scopes[k]) == 0 {
			logger.Debugf(`scopes[%v] is not set`, k)
			continue
		}
		gr := &Grant{
			Scope:    scopes[k],
			Value:    v,
			Database: databases[k],
			Table:    tables[k],
			Columns:  columns[k],
			Settings: map[string]string{},
		}
		v = gr.String()
		if oldGr, ok := newGrants[v]; !ok {
			newGrants[v] = gr
		} else {
			for k, v := range oldGr.Settings {
				gr.Settings[k] = v
			}
		}
		if mapx == nil {
			newGrants[v] = gr
			continue
		}
		mp := mapx.Get(strconv.Itoa(k))
		if mp != nil {
			for group, settings := range mp.Map {
				if settings.Map == nil || !gr.IsValid(group, settings.Map) {
					continue
				}
				for name, m := range settings.Map {
					gr.Settings[name] = m.Value()
				}
			}
		}
	}
	//panic(echo.Dump(echo.H{`old`: grants, `new`: newGrants}, false))
	hasURLGrantValue := len(m.Form(`grant`)) > 0
	operations := []*Grant{}
	//newGrants: newGrants[*.*|db.*|db.table|db.table(col1,col2)][DROP|...]=`0|1`
	for object, grant := range newGrants {
		onAndCol := reGrantColumn.FindStringSubmatch(object)
		//fmt.Printf("object: %v matched: %#v\n", object, onAndCol)
		if len(onAndCol) < 3 {
			continue
		}
		grant.Operation = &Operation{
			Grant:   []string{},
			Revoke:  []string{},
			Columns: onAndCol[2],
			On:      onAndCol[1],
			User:    newUserAndHost,
			Scope:   grant.Scope,
		}
		if hasURLGrantValue { // 下拉菜单模式
			for key, val := range grant.Settings {
				if val != `1` && key != `ALL PRIVILEGES` {
					grant.Revoke = append(grant.Revoke, key) // 清理取消的权限
				}
			}
		} else { // 勾选模式
			if !created { // 编辑模式
				if vals, ok := grants[object]; ok {
					for key := range vals {
						if _, ok := grant.Settings[key]; !ok { // 在新提交的数据中没有勾选旧权限项时，取消旧权限项
							grant.Revoke = append(grant.Revoke, key)
						}
					}
					for key := range grant.Settings {
						if _, ok := vals[key]; !ok { // 在新提交的数据中勾选了旧权限没有的项时，增加新权限项
							grant.Grant = append(grant.Grant, key)
						}
					}
				} else { // 没有旧权限时，作为新权限添加
					for key := range grant.Settings {
						grant.Grant = append(grant.Grant, key)
					}
				}
			} else { // 新建账号模式下，添加新权限
				for key := range grant.Settings {
					grant.Grant = append(grant.Grant, key)
				}
			}
		}
		operations = append(operations, grant)
		if _, ok := grants[object]; ok {
			delete(grants, object) // 清理掉本次提交的项，剩下的就是需要取消的权限
		}
	}
	//panic(echo.Dump(echo.H{`new`: operations, `revoke`: grants}, false))
	if len(oldUser) > 0 && (!hasURLGrantValue && !created) {
		for object, revoke := range grants { // 删掉没有勾选的权限
			onAndCol := reGrantColumn.FindStringSubmatch(object)
			if len(onAndCol) < 3 {
				continue
			}
			op := &Operation{
				Grant:   []string{},
				Revoke:  []string{},
				Columns: onAndCol[2],
				On:      onAndCol[1],
				User:    newUserAndHost,
			}
			for k := range revoke {
				op.Revoke = append(op.Revoke, k)
			}
			if err := op.Apply(m); err != nil {
				return err
			}
		}
	}
	for _, op := range operations { // 执行操作的权限
		if err := op.Apply(m); err != nil {
			return onerror(err)
		}
	}
	if len(oldUser) > 0 { // 如果是在旧账号的基础上创建新账号，则删除旧账号
		if created {
			r := &Result{}
			r.SQL = "DROP USER " + oldUserAndHost
			r.Exec(m.newParam())
			m.AddResults(r)
			if r.err != nil {
				return onerror(err)
			}
		}
	}
	r2 := &Result{}
	r2.SQL = "FLUSH PRIVILEGES"
	r2.Exec(m.newParam())
	m.AddResults(r2)
	return nil
}

func (m *mySQL) getUserGrants(host, user string) (string, map[string]map[string]bool, []string, error) {
	r := map[string]map[string]bool{}
	var (
		sortNumber []string
		oldPass    string
		err        error
	)
	if len(host) > 0 {
		sqlStr := "SHOW GRANTS FOR " + quoteVal(user) + "@" + quoteVal(host)
		rows, err := m.newParam().SetCollection(sqlStr).Query()
		if err != nil {
			return oldPass, r, sortNumber, err
		}
		defer rows.Close()
		for rows.Next() {
			var v sql.NullString
			err = rows.Scan(&v)
			if err != nil {
				break
			}
			matchOn := reGrantOn.FindStringSubmatch(v.String)
			if len(matchOn) > 0 {
				matchBrackets := reGrantBrackets.FindAllStringSubmatch(matchOn[1], -1)
				if len(matchBrackets) > 0 {
					for _, val := range matchBrackets {
						if val[1] != `USAGE` {
							k := matchOn[2] + val[2]
							if _, ok := r[k]; !ok {
								r[k] = map[string]bool{}
								sortNumber = append(sortNumber, k)
							}
							if val[1] == `PROXY` {
								r[k]["ALL PRIVILEGES"] = true
							}
							r[k][val[1]] = true
						}
						if reGrantOption.MatchString(v.String) {
							k := matchOn[2] + val[2]
							if _, ok := r[k]; !ok {
								r[k] = map[string]bool{}
								sortNumber = append(sortNumber, k)
							}
							r[k]["GRANT OPTION"] = true
						}
					}
				}
			}
			matchIdent := reGrantIdent.FindStringSubmatch(v.String)
			if len(matchIdent) > 0 {
				oldPass = matchIdent[1]
			}
		}
	} else {
		sqlStr := "SELECT SUBSTRING_INDEX(CURRENT_USER, '@', -1)"
		row, err := m.newParam().SetCollection(sqlStr).QueryRow()
		if err != nil {
			return oldPass, r, sortNumber, err
		}
		var v sql.NullString
		err = row.Scan(&v)
		if err != nil {
			return oldPass, r, sortNumber, err
		}
		m.Request().Form().Set(`host`, v.String)
	}
	var key string
	if len(m.dbName) == 0 || len(r) > 0 {
	} else {
		key = com.AddCSlashes(m.dbName, '%', '_', '\\') + ".*"
	}
	r[key] = map[string]bool{}
	sortNumber = append(sortNumber, key)
	return oldPass, r, sortNumber, err
}
