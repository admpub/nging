/*

   Copyright 2016 Wenhui Shen <www.webx.top>

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.

*/
package mysql

func (m *mySQL) listPrivileges() (bool, []map[string]string, error) {
	sqlStr := "SELECT User, Host FROM mysql."
	if len(m.dbName) == 0 {
		sqlStr += `user`
	} else {
		sqlStr += "db WHERE " + quoteCol(m.dbName) + " LIKE Db"
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
	for rows.Next() {
		v := &Privilege{}
		err = rows.Scan(&v.Privilege, &v.Context, &v.Comment)
		if err != nil {
			break
		}
		r.Privileges = append(r.Privileges, v)
	}
	return r, err
}
