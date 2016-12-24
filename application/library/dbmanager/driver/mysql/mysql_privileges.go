package mysql

import "strings"

func (m *mySQL) listPrivileges() (bool, []map[string]string, error) {
	sqlStr := "SELECT User, Host FROM mysql."
	if len(m.dbName) == 0 {
		sqlStr += `user`
	} else {
		sqlStr += "db WHERE `" + strings.Replace(m.dbName, "`", "", -1) + "` LIKE Db"
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
