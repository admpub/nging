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

func (m *mySQL) listPrivileges() (bool, []map[string]string, error) {
	sqlStr := "SELECT User, Host FROM mysql."
	if len(m.dbName) == 0 {
		sqlStr += `user`
	} else {
		sqlStr += "db WHERE Db LIKE " + quoteVal(m.dbName)
	}
	sqlStr += " ORDER BY Host, User"
	res, err := m.kvVal(sqlStr)
	sysUser := true
	if err != nil || res == nil || len(res) == 0 {
		sysUser = false
		sqlStr = `SELECT SUBSTRING_INDEX(CURRENT_USER, '@', 1) AS User, SUBSTRING_INDEX(CURRENT_USER, '@', -1) AS Host`
		res, err = m.kvVal(sqlStr)
	}
	return sysUser, res, err
}

func (m *mySQL) showPrivileges() (*Privileges, error) {
	r := NewPrivileges()
	sqlStr := "SHOW PRIVILEGES"
	rows, err := m.newParam().SetCollection(sqlStr).Query()
	if err != nil {
		return r, err
	}
	defer rows.Close()
	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	n := len(cols)
	for rows.Next() {
		v := &Privilege{}
		err = safeScan(rows, n, &v.Privilege, &v.Context, &v.Comment)
		if err != nil {
			break
		}
		r.Privileges = append(r.Privileges, v)
	}
	return r, err
}
