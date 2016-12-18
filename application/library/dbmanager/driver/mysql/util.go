package mysql

import (
	"database/sql"
	"time"

	"regexp"

	"strings"

	"github.com/webx-top/com"
	"github.com/webx-top/db/lib/factory"
)

type Result struct {
	SQL          string
	RowsAffected int64
	timeStart    time.Time
	timeEnd      time.Time
	Started      string
	Elapsed      string
}

func (r *Result) elapsed() time.Duration {
	return r.timeEnd.Sub(r.timeStart)
}

func (r *Result) Exec(p *factory.Param) (*Result, error) {
	r.timeStart = time.Now()
	defer func() {
		r.timeEnd = time.Now()
		r.Started = r.timeStart.Format(`2006-01-02 15:04:05`)
		r.Elapsed = r.elapsed().String()
	}()
	result, err := p.SetCollection(r.SQL).Exec()
	if err != nil {
		return r, err
	}
	r.RowsAffected, err = result.RowsAffected()
	return r, err
}

func (m *mySQL) createDatabase(dbName, collate string) (*Result, error) {
	r := &Result{}
	r.SQL = "CREATE DATABASE `" + strings.Replace(dbName, "`", "", -1) + "`"
	if len(collate) > 0 {
		r.SQL += " COLLATE '" + com.AddSlashes(collate) + "'"
	}
	_, err := r.Exec(m.newParam())
	return r, err
}

func (m *mySQL) setLastSQL(sqlStr string) {
	m.Session().AddFlash(sqlStr, `lastSQL`)
}

func (m *mySQL) lastSQL() interface{} {
	return m.Flash(`lastSQL`)
}

// 获取数据库列表
func (m *mySQL) getDatabases() ([]string, error) {
	sqlStr := `SELECT SCHEMA_NAME FROM information_schema.SCHEMATA`
	if com.VersionCompare(m.getVersion(), `5`) < 0 {
		sqlStr = `SHOW DATABASES`
	}
	rows, err := m.newParam().SetCollection(sqlStr).Query()
	if err != nil {
		return nil, err
	}
	ret := []string{}
	for rows.Next() {
		var v sql.NullString
		err := rows.Scan(&v)
		if err != nil {
			return nil, err
		}
		ret = append(ret, v.String)
	}
	return ret, nil
}

// 获取数据表列表
func (m *mySQL) getTables() ([]string, error) {
	sqlStr := `SHOW TABLES`
	rows, err := m.newParam().SetCollection(sqlStr).Query()
	if err != nil {
		return nil, err
	}
	ret := []string{}
	for rows.Next() {
		var v sql.NullString
		err := rows.Scan(&v)
		if err != nil {
			return nil, err
		}
		ret = append(ret, v.String)
	}
	return ret, nil
}

type TableStatus struct {
	Name            sql.NullString
	Engine          sql.NullString
	Version         sql.NullString
	Row_format      sql.NullString
	Rows            sql.NullInt64
	Avg_row_length  sql.NullInt64
	Data_length     sql.NullInt64
	Max_data_length sql.NullInt64
	Index_length    sql.NullInt64
	Data_free       sql.NullInt64
	Auto_increment  sql.NullInt64
	Create_time     sql.NullString
	Update_time     sql.NullString
	Check_time      sql.NullString
	Collation       sql.NullString
	Checksum        sql.NullString
	Create_options  sql.NullString
	Comment         sql.NullString
}

func (t *TableStatus) IsView() bool {
	return t.Engine.Valid == false
}

func (t *TableStatus) FKSupport(currentVersion string) bool {
	switch t.Engine.String {
	case `InnoDB`, `IBMDB2I`, `NDB`:
		if com.VersionCompare(currentVersion, `5.6`) >= 0 {
			return true
		}
	}
	return false
}

func (t *TableStatus) Size() int64 {
	return t.Data_length.Int64 + t.Index_length.Int64
}

type Collation struct {
	Collation sql.NullString
	Charset   sql.NullString `json:"-"`
	Id        sql.NullInt64  `json:"-"`
	Default   sql.NullString `json:"-"`
	Compiled  sql.NullString `json:"-"`
	Sortlen   sql.NullInt64  `json:"-"`
}

type Collations struct {
	Collations map[string][]Collation
	Defaults   map[string]int
}

func NewCollations() *Collations {
	return &Collations{
		Collations: make(map[string][]Collation),
		Defaults:   make(map[string]int),
	}
}

// 获取支持的字符集
func (m *mySQL) getCollations() (*Collations, error) {
	sqlStr := `SHOW COLLATION`
	rows, err := m.newParam().SetCollection(sqlStr).Query()
	if err != nil {
		return nil, err
	}
	ret := NewCollations()
	for rows.Next() {
		var v Collation
		err := rows.Scan(&v.Collation, &v.Charset, &v.Id, &v.Default, &v.Compiled, &v.Sortlen)
		if err != nil {
			return nil, err
		}
		coll, ok := ret.Collations[v.Charset.String]
		if !ok {
			coll = []Collation{}
		}
		if v.Default.Valid && len(v.Default.String) > 0 {
			ret.Defaults[v.Charset.String] = len(coll)
		}
		coll = append(coll, v)
		ret.Collations[v.Charset.String] = coll
	}
	return ret, nil
}

var (
	reCollate       = regexp.MustCompile(` COLLATE ([^ ]+)`)
	reCharacter     = regexp.MustCompile(` CHARACTER SET ([^ ]+)`)
	reInnoDBComment = regexp.MustCompile(`(?:(.+); )?InnoDB free: .*`)
)

func (m *mySQL) getCollation(dbName string, collations *Collations) (string, error) {
	var err error
	if collations == nil {
		collations, err = m.getCollations()
		if err != nil {
			return ``, err
		}
	}
	sqlStr := "SHOW CREATE DATABASE `" + strings.Replace(dbName, "`", "", -1) + "`"
	row, err := m.newParam().SetCollection(sqlStr).QueryRow()
	if err != nil {
		return ``, err
	}
	var database sql.NullString
	var createDb sql.NullString
	err = row.Scan(&database, &createDb)
	if err != nil {
		return ``, err
	}
	matches := reCollate.FindStringSubmatch(createDb.String)
	if len(matches) > 1 {
		return matches[1], nil
	}
	matches = reCharacter.FindStringSubmatch(createDb.String)
	if len(matches) > 1 {
		if idx, ok := collations.Defaults[matches[1]]; ok {
			return collations.Collations[matches[1]][idx].Collation.String, nil
		}
	}

	return ``, nil
}

func (m *mySQL) getTableStatus(dbName string, tableName string, fast bool) (map[string]*TableStatus, error) {
	sqlStr := `SHOW TABLE STATUS`
	if len(dbName) > 0 {
		sqlStr += " FROM `" + strings.Replace(dbName, "`", "", -1) + "`"
	}
	if len(tableName) > 0 {
		tableName = com.AddSlashes(tableName, '_', '%')
		tableName = `'` + tableName + `'`
		sqlStr += ` LIKE ` + tableName
	}
	rows, err := m.newParam().SetCollection(sqlStr).Query()
	if err != nil {
		return nil, err
	}
	ret := map[string]*TableStatus{}
	for rows.Next() {
		v := &TableStatus{}
		err := rows.Scan(&v.Name, &v.Engine, &v.Version, &v.Row_format, &v.Rows, &v.Avg_row_length, &v.Data_length, &v.Max_data_length, &v.Index_length, &v.Data_free, &v.Auto_increment, &v.Create_time, &v.Update_time, &v.Check_time, &v.Collation, &v.Checksum, &v.Create_options, &v.Comment)
		if err != nil {
			return nil, err
		}
		if v.Engine.String == `InnoDB` {
			v.Comment.String = reInnoDBComment.ReplaceAllString(v.Comment.String, `$1`)
		}
		ret[v.Name.String] = v
		if len(tableName) > 0 {
			return ret, nil
		}
	}
	return ret, nil
}

func (m *mySQL) newParam() *factory.Param {
	return factory.NewParam(m.db)
}

func (m *mySQL) getVersion() string {
	if len(m.version) > 0 {
		return m.version
	}
	row, err := m.newParam().SetCollection(`SELECT version()`).QueryRow()
	if err != nil {
		return err.Error()
	}
	var v sql.NullString
	err = row.Scan(&v)
	if err != nil {
		return err.Error()
	}
	m.version = v.String
	return v.String
}

func (m *mySQL) baseInfo() error {
	dbList, err := m.getDatabases()
	if err != nil {
		return err
	}
	m.Set(`dbList`, dbList)
	if len(m.DbAuth.Db) > 0 {
		tableList, err := m.getTables()
		if err != nil {
			return err
		}
		m.Set(`tableList`, tableList)
	}

	m.Set(`dbVersion`, m.getVersion())
	return nil
}
