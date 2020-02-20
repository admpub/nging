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
	"bytes"
	"database/sql"
	"fmt"
	"io"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/webx-top/com"
	"github.com/webx-top/db/lib/factory"
	"github.com/webx-top/db/mysql"
	"github.com/webx-top/echo"
	"github.com/webx-top/pagination"

	"github.com/admpub/errors"
	"github.com/admpub/nging/application/library/common"
	"github.com/admpub/nging/application/library/dbmanager/driver"
)

func init() {
	driver.Register(`mysql`, &mySQL{
		TriggerOptions: []*TriggerOption{
			&TriggerOption{
				Type:    `Timing`,
				Options: []string{"BEFORE", "AFTER"},
			},
			&TriggerOption{
				Type:    `Event`,
				Options: []string{"INSERT", "UPDATE", "DELETE"},
			},
			&TriggerOption{
				Type:    `Type`,
				Options: []string{"FOR EACH ROW"},
			},
		},
		supportSQL: true,
	})
}

type mySQL struct {
	*driver.BaseDriver
	db             *factory.Factory
	dbName         string
	version        string
	TriggerOptions TriggerOptions
	supportSQL     bool
}

func (m *mySQL) Name() string {
	return `MySQL`
}

func (m *mySQL) Init(ctx echo.Context, auth *driver.DbAuth) {
	m.BaseDriver = driver.NewBaseDriver()
	m.BaseDriver.Init(ctx, auth)
	m.Set(`supportSQL`, m.supportSQL)
}

func (m *mySQL) IsSupported(operation string) bool {
	return true
}

func (m *mySQL) Login() error {
	factoryDB := factory.New()
	settings := mysql.ConnectionURL{
		User:     m.DbAuth.Username,
		Password: m.DbAuth.Password,
		Host:     m.DbAuth.Host,
		Database: m.DbAuth.Db,
	}
	var dbNameIsEmpty bool
	if len(settings.Database) == 0 {
		dbNameIsEmpty = true
		settings.Database = m.Form(`db`)
	}
	m.dbName = settings.Database
	db, err := mysql.Open(settings)
	if err != nil {
		if dbNameIsEmpty {
			m.fail(err.Error())
			return m.returnTo(`/db`)
		}
		settings.Password = strings.Repeat(`*`, len(settings.Password))
		return errors.Wrap(err, m.T(`连接数据库出错`)+`: `+echo.Dump(settings, false))
	}
	err = db.Ping()
	if err != nil {
		return errors.Wrap(err, m.T(`连接数据库出错`))
	}
	cluster := factory.NewCluster().AddMaster(db)
	factoryDB.SetCluster(0, cluster)
	m.db = factoryDB
	m.Set(`dbName`, m.dbName)
	m.Set(`table`, m.Form(`table`))
	if len(settings.Database) > 0 {
		m.Set(`dbList`, []string{settings.Database})
	}
	return m.baseInfo()
}

func (m *mySQL) Logout() error {
	if m.db != nil {
		m.db.CloseAll()
		m.db = nil
	}
	return nil
}
func (m *mySQL) ProcessList() error {
	var ret interface{}
	if m.IsPost() {
		pids := m.FormValues(`pid[]`)
		for _, pid := range pids {
			i, e := strconv.ParseInt(pid, 10, 64)
			if e == nil {
				e = m.killProcess(i)
			}
		}
	}
	r, e := m.processList()
	ret = common.Err(m.Context, e)
	m.Set(`processList`, r)
	return m.Render(`db/mysql/process_list`, ret)
}

func (m *mySQL) returnTo(rets ...string) error {
	m.EnableFlashSession()
	return m.ReturnTo(rets...)
}

func (m *mySQL) Privileges() error {
	var ret interface{}
	var err error
	act := m.Form(`act`)
	if len(act) > 0 {
		switch act {
		case `drop`:
			host := m.Form(`host`)
			user := m.Form(`user`)
			if len(user) < 1 {
				m.Session().AddFlash(errors.New(m.T(`用户名不正确`)))
				return m.returnTo()
			}
			if user == `root` {
				m.Session().AddFlash(errors.New(m.T(`root 用户不可删除`)))
				return m.returnTo()
			}
			err = m.dropUser(user, host)
			if err != nil {
				m.fail(err.Error())
			}
			return m.returnTo(m.GenURL(`privileges`))
		case `edit`:
			if m.IsPost() {
				isHashed := len(m.Form(`hashed`)) > 0
				user := m.Form(`oldUser`)
				host := m.Query(`host`)
				newHost := m.Form(`host`)
				newUser := m.Form(`user`)
				oldPasswd := m.Form(`oldPass`)
				newPasswd := m.Form(`pass`)
				err = m.editUser(user, host, newUser, newHost, oldPasswd, newPasswd, isHashed)
				if err == nil {
					m.ok(m.T(`操作成功`))
					return m.returnTo(m.GenURL(`privileges`) + `&act=edit&user=` + url.QueryEscape(newUser) + `&host=` + url.QueryEscape(newHost))
				}
				m.fail(err.Error())
			}
			privs, err := m.showPrivileges()
			if err == nil {
				privs.Parse()
			}
			m.Set(`list`, privs.privileges)
			m.Set(`groups`, []*KV{
				&KV{`_Global_`, ``},
				&KV{`Server_Admin`, m.T(`服务器`)},
				&KV{`Databases`, m.T(`数据库`)},
				&KV{`Tables`, m.T(`表`)},
				&KV{`Columns`, m.T(`列`)},
				&KV{`Procedures`, m.T(`子程序`)},
			})
			user := m.Form(`user`)
			host := m.Form(`host`)
			var oldUser string
			oldPass, grants, sorts, err := m.getUserGrants(host, user)
			if _, ok := grants["*.*"]; ok {
				m.Set(`hasGlobalScope`, true)
			} else {
				m.Set(`hasGlobalScope`, false)
			}
			if err == nil {
				oldUser = user
			}

			m.Set(`sorts`, sorts)
			m.Set(`grants`, grants)
			if oldPass != `` {
				m.Set(`hashed`, true)
			} else {
				m.Set(`hashed`, false)
			}
			m.Set(`oldPass`, oldPass)
			m.Set(`oldUser`, oldUser)
			m.Request().Form().Set(`pass`, oldPass)
			m.SetFunc(`getGrantByPrivilege`, func(grant map[string]bool, index int, group string, privilege string) bool {
				priv := strings.ToUpper(privilege)
				value := m.Form(fmt.Sprintf(`grants[%v][%v][%v]`, index, group, priv))
				if len(value) > 0 && value == `1` {
					return true
				}
				return grant[priv]
			})
			m.SetFunc(`getGrantsByKey`, func(key string) map[string]bool {
				if vs, ok := grants[key]; ok {
					return vs
				}
				return map[string]bool{}
			})
			m.SetFunc(`getScope`, m.getScopeGrant)
			m.SetFunc(`fieldName`, func(index int, group string, privilege string) string {
				return fmt.Sprintf(`grants[%v][%v][%v]`, index, group, strings.ToUpper(privilege))
			})
			ret = common.Err(m.Context, err)
			return m.Render(`db/mysql/privilege_edit`, ret)
		}
	}
	ret = common.Err(m.Context, err)
	isSysUser, list, err := m.listPrivileges()
	m.Set(`isSysUser`, isSysUser)
	m.Set(`list`, list)
	return m.Render(`db/mysql/privileges`, ret)
}
func (m *mySQL) Info() error {
	var r []map[string]string
	var e error
	switch m.Form(`type`) {
	case `variables`:
		r, e = m.showVariables()
	default:
		r, e = m.showStatus()
	}
	m.Set(`list`, r)
	return m.Render(`db/mysql/info`, e)
}
func (m *mySQL) CreateDb() error {
	dbName := m.Form(`name`)
	collate := m.Form(`collation`)
	data := m.Data()
	if len(dbName) < 1 {
		data.SetZone(`name`).SetInfo(m.T(`数据库名称不能为空`), 0)
	} else {
		res := m.createDatabase(dbName, collate)
		if res.err != nil {
			data.SetError(res.err)
		} else {
			data.SetData(res)
		}
	}
	return m.JSON(data)
}
func (m *mySQL) ModifyDb() error {
	opType := m.Form(`json`)
	if len(opType) > 0 {
		switch opType {
		case `collations`:
			return m.listDbAjax(opType)
		}
		return nil
	}
	if len(m.dbName) < 1 {
		m.fail(m.T(`请先选择一个数据库`))
		return m.returnTo(m.GenURL(`listDb`))
	}
	var err error
	if m.IsPost() {
		name := m.Form(`name`)
		collation := m.Form(`collation`)
		if name != m.dbName {
			results := m.renameDatabase(name, collation)
			for _, r := range results {
				m.AddResults(r)
			}
		} else {
			m.AddResults(m.alterDatabase(name, collation))
		}
		return m.returnTo(m.GenURL(`listDb`))
	}
	form := m.Request().Form()
	form.Set(`name`, m.dbName)
	collation, err := m.getCollation(m.dbName, nil)
	form.Set(`collation`, collation)

	return m.Render(`db/mysql/modify_db`, err)
}
func (m *mySQL) listDbAjax(opType string) error {
	switch opType {
	case `drop`:
		data := m.Data()
		dbs := m.FormValues(`db[]`)
		rs := []*Result{}
		code := 1
		for _, db := range dbs {
			r := m.dropDatabase(db)
			rs = append(rs, r)
			if r.err != nil {
				data.SetError(r.err)
				code = 0
				break
			}
		}
		data.SetData(rs, code)
		return m.JSON(data)
	case `create`:
		return m.CreateDb()
	case `collations`:
		data := m.Data()
		collations, err := m.getCollations()
		if err != nil {
			data.SetError(err)
		} else {
			data.SetData(collations.Collations)
		}
		return m.JSON(data)
	}
	return nil
}
func (m *mySQL) ListDb() error {
	opType := m.Form(`json`)
	if len(opType) > 0 {
		return m.listDbAjax(opType)
	}
	var err error
	dbList, ok := m.Get(`dbList`).([]string)
	if !ok {
		dbList, err = m.getDatabases()
		if err != nil {
			return err
		}
		m.Set(`dbList`, dbList)
	}
	colls := make([]string, len(dbList))
	sizes := make([]int64, len(dbList))
	tables := make([]int, len(dbList))
	collations, err := m.getCollations()
	if err != nil {
		return err
	}
	for index, dbName := range dbList {
		colls[index], err = m.getCollation(dbName, collations)
		if err == nil {
			var tableStatus map[string]*TableStatus
			tableStatus, _, err = m.getTableStatus(dbName, ``, true)
			if err == nil {
				tables[index] = len(tableStatus)
				for _, tableStat := range tableStatus {
					sizes[index] += tableStat.Size()
				}
			}
		}
		if err != nil {
			return err
		}
	}
	m.Set(`dbColls`, colls)
	m.Set(`dbSizes`, sizes)
	m.Set(`dbTables`, tables)
	return m.Render(`db/mysql/list_db`, m.checkErr(err))
}
func (m *mySQL) CreateTable() error {
	opType := m.Form(`json`)
	if len(opType) > 0 {
		switch opType {
		case `collations`:
			return m.listDbAjax(opType)
		}
		return nil
	}

	referencablePrimary, _, err := m.referencablePrimary(``)
	foreignKeys := map[string]string{}
	for tblName, field := range referencablePrimary {
		foreignKeys[strings.Replace(tblName, "`", "``", -1)+"`"+strings.Replace(field.Field, "`", "``", -1)] = tblName
	}
	partitions := map[string]string{}
	for _, p := range PartitionTypes {
		partitions[p] = p
	}
	postFields := []*Field{}
	if m.IsPost() {
		table := m.Form(`name`)
		engine := m.Form(`engine`)
		collation := m.Form(`collation`)
		autoIncrementStartValue := m.Form(`ai_start_val`)
		autoIncrementStart := sql.NullInt64{Valid: len(autoIncrementStartValue) > 0}
		if autoIncrementStart.Valid {
			autoIncrementStart.Int64, _ = strconv.ParseInt(autoIncrementStartValue, 10, 64)
		}
		comment := m.Form(`comment`)
		aiIndex := m.Formx(`auto_increment`)
		aiIndexInt := aiIndex.Int()
		aiIndexStr := aiIndex.String()
		mapx := echo.NewMapx(m.Forms())
		f := mapx.Get(`fields`)
		allFields := []*fieldItem{}
		after := " FIRST"
		foreign := map[string]string{}
		if err == nil && f != nil {
			size := len(f.Map)
			for i := 0; i < size; i++ {
				ii := strconv.Itoa(i)
				fieldName := f.Value(ii, `field`)
				if len(fieldName) == 0 {
					break
				}
				field := &Field{}
				field.Field = fieldName
				field.Type = f.Value(ii, `type`)
				field.Length = f.Value(ii, `length`)
				field.Unsigned = f.Value(ii, `unsigned`)
				field.Collation = f.Value(ii, `collation`)
				field.On_delete = f.Value(ii, `on_delete`)
				field.On_update = f.Value(ii, `on_update`)
				field.Null, _ = strconv.ParseBool(f.Value(ii, `null`))
				field.Comment = f.Value(ii, `comment`)
				field.Default = sql.NullString{
					String: f.Value(ii, `default`),
					Valid:  f.Value(ii, `has_default`) == `1`,
				}
				field.AutoIncrement = sql.NullString{
					Valid: len(aiIndexStr) > 0 && aiIndexInt == i,
				}
				if field.AutoIncrement.Valid {
					field.AutoIncrement.String = autoIncrementStartValue
				}
				var typeField *Field
				if foreignKey, ok := foreignKeys[field.Type]; ok {
					typeField, _ = referencablePrimary[foreignKey]
					foreignK, err := m.formatForeignKey(&ForeignKeyParam{
						Table:  foreignKey,
						Source: []string{field.Field},
						Target: []string{field.On_delete},
					})
					if err != nil {
						return err
					}
					foreign[quoteCol(field.Field)] = ` ` + foreignK
				}
				if typeField == nil {
					typeField = field
				}
				item := &fieldItem{
					Original:     ``,
					ProcessField: []string{},
					After:        after,
				}
				item.ProcessField, err = m.processField(``, field, typeField, aiIndexStr)
				if err != nil {
					return err
				}
				allFields = append(allFields, item)
				after = " AFTER " + quoteCol(field.Field)
				postFields = append(postFields, field)
			}
		}
		partitioning := m.tablePartitioning(partitions, nil)
		err = m.alterTable(``, table, allFields, foreign,
			sql.NullString{String: comment, Valid: len(comment) > 0},
			engine, collation,
			autoIncrementStart,
			partitioning)
		if err == nil {
			return m.returnTo()
		}
	}
	engines, err := m.getEngines()
	m.Set(`engines`, engines)
	m.Set(`typeGroups`, typeGroups)
	m.Set(`foreignKeys`, foreignKeys)
	m.Set(`onActions`, strings.Split(OnActions, `|`))
	m.Set(`unsignedTags`, UnsignedTags)
	if m.Form(`engine`) == `` {
		m.Request().Form().Set(`engine`, `InnoDB`)
	}
	if len(postFields) == 0 {
		postFields = append(postFields, &Field{})
	}
	m.Set(`postFields`, postFields)
	m.SetFunc(`isString`, reFieldTypeText.MatchString)
	m.SetFunc(`isNumeric`, reFieldTypeNumber.MatchString)
	supportPartitioning := m.support(`partitioning`)
	if supportPartitioning {
		partition := &Partition{
			Names:  []string{``},
			Values: []string{``},
		}
		m.Set(`partition`, partition)
	}
	m.Set(`supportPartitioning`, supportPartitioning)
	m.Set(`partitionTypes`, PartitionTypes)
	return m.Render(`db/mysql/create_table`, err)
}
func (m *mySQL) ModifyTable() error {
	opType := m.Form(`json`)
	if len(opType) > 0 {
		switch opType {
		case `collations`:
			return m.listDbAjax(opType)
		}
		return nil
	}

	oldTable := m.Form(`table`)
	if len(oldTable) < 1 {
		m.fail(m.T(`table参数不能为空`))
		return m.returnTo(`listDb`)
	}

	referencablePrimary, _, err := m.referencablePrimary(``)
	foreignKeys := map[string]string{}
	for tblName, field := range referencablePrimary {
		foreignKeys[strings.Replace(tblName, "`", "``", -1)+"`"+strings.Replace(field.Field, "`", "``", -1)] = tblName
	}
	postFields := []*Field{}
	var origFields map[string]*Field
	var sortFields []string
	var tableStatus *TableStatus
	val, sort, err := m.tableFields(oldTable)
	if err != nil {
		return err
	}
	origFields = val
	sortFields = sort
	stt, _, err := m.getTableStatus(m.dbName, oldTable, false)
	if err != nil {
		return err
	}
	if ts, ok := stt[oldTable]; ok {
		tableStatus = ts
	}
	partitions := map[string]string{}
	for _, p := range PartitionTypes {
		partitions[p] = p
	}
	if m.IsPost() {
		table := m.Form(`name`)
		engine := m.Form(`engine`)
		collation := m.Form(`collation`)
		autoIncrementStartValue := m.Form(`ai_start_val`)
		autoIncrementStart := sql.NullInt64{Valid: len(autoIncrementStartValue) > 0}
		if autoIncrementStart.Valid {
			autoIncrementStart.Int64, _ = strconv.ParseInt(autoIncrementStartValue, 10, 64)
		}
		comment := m.Form(`comment`)
		aiIndex := m.Formx(`auto_increment`)
		aiIndexInt := aiIndex.Int()
		aiIndexStr := aiIndex.String()
		mapx := echo.NewMapx(m.Forms())
		f := mapx.Get(`fields`)
		var origField *Field
		origFieldsNum := len(sortFields)
		if origFieldsNum > 0 {
			fieldName := sortFields[0]
			origField = origFields[fieldName]
		}
		var useAllFields bool
		fields := []*fieldItem{}
		allFields := []*fieldItem{}
		after := " FIRST"
		foreign := map[string]string{}
		driverName := strings.ToLower(m.DbAuth.Driver)
		j := 1
		if err == nil && f != nil {
			size := len(f.Map)
			for i := 0; i < size; i++ {
				ii := strconv.Itoa(i)
				fieldName, posted := f.ValueOk(ii, `field`)
				orig, exists := f.ValueOk(ii, `orig`)
				if !posted && !exists {
					break
				}
				if len(fieldName) < 1 {
					if len(orig) > 0 {
						useAllFields = true
						item := &fieldItem{
							Original:     orig,
							ProcessField: []string{},
						}
						fields = append(fields, item)
					}
				} else {
					field := &Field{}
					field.Field = fieldName
					field.Type = f.Value(ii, `type`)
					field.Length = f.Value(ii, `length`)
					field.Unsigned = f.Value(ii, `unsigned`)
					field.Collation = f.Value(ii, `collation`)
					field.On_delete = f.Value(ii, `on_delete`)
					field.On_update = f.Value(ii, `on_update`)
					field.Null, _ = strconv.ParseBool(f.Value(ii, `null`))
					field.Comment = f.Value(ii, `comment`)
					field.Default = sql.NullString{
						String: f.Value(ii, `default`),
						Valid:  f.Value(ii, `has_default`) == `1`,
					}
					field.AutoIncrement = sql.NullString{
						Valid: len(aiIndexStr) > 0 && aiIndexInt == i,
					}
					if field.AutoIncrement.Valid {
						field.AutoIncrement.String = autoIncrementStartValue
					}
					var typeField *Field
					if foreignKey, ok := foreignKeys[field.Type]; ok {
						typeField, _ = referencablePrimary[foreignKey]
						foreignK, err := m.formatForeignKey(&ForeignKeyParam{
							Table:    foreignKey,
							Source:   []string{field.Field},
							Target:   []string{typeField.Field},
							OnDelete: field.On_delete,
						})
						if err != nil {
							return err
						}
						if driverName == `sqlite` || len(oldTable) == 0 {
							foreign[quoteCol(field.Field)] = ` ` + foreignK
						} else {
							foreign[quoteCol(field.Field)] = `ADD` + foreignK
						}
					}
					if typeField == nil {
						typeField = field
					}
					field.Original = f.Value(ii, `orig`)
					item := &fieldItem{
						Original:     field.Original,
						ProcessField: []string{},
						After:        after,
					}
					item.ProcessField, err = m.processField(oldTable, field, typeField, aiIndexStr)
					if err != nil {
						return err
					}
					allFields = append(allFields, item)
					processField, err := m.processField(oldTable, origField, origField, aiIndexStr)
					if err != nil {
						return err
					}
					//fmt.Printf(`%#v`+"\n", item.ProcessField)
					//fmt.Printf(`%#v`+"\n", processField)
					isChanged := fmt.Sprintf(`%#v`, item.ProcessField) != fmt.Sprintf(`%#v`, processField)
					if isChanged {
						fields = append(fields, item)
						if len(field.Original) > 0 || len(after) > 0 {
							useAllFields = true
						}
					}
					after = " AFTER " + quoteCol(field.Field)
					postFields = append(postFields, field)
				}
				if len(orig) > 0 {
					if origFieldsNum > j {
						origField = origFields[sortFields[j]]
						j++
					} else {
						after = ``
					}
				}
			}
		}
		partitioning := m.tablePartitioning(partitions, tableStatus)
		if tableStatus != nil {
			if comment == tableStatus.Comment.String {
				comment = ``
			}
			if engine == tableStatus.Engine.String {
				engine = ``
			}
			if collation == tableStatus.Collation.String {
				collation = ``
			}
		}
		if driverName == `sqlite` && (useAllFields || len(foreign) > 0) {
			err = m.alterTable(oldTable, table, allFields, foreign,
				sql.NullString{String: comment, Valid: len(comment) > 0},
				engine, collation,
				autoIncrementStart,
				partitioning)
		} else {
			err = m.alterTable(oldTable, table, fields, foreign,
				sql.NullString{String: comment, Valid: len(comment) > 0},
				engine, collation,
				autoIncrementStart,
				partitioning)
		}
		if err == nil {
			return m.returnTo()
		}
	} else {
		postFields = make([]*Field, len(sortFields))
		for k, v := range sortFields {
			postFields[k] = origFields[v]
		}
	}
	engines, err := m.getEngines()
	m.Set(`engines`, engines)
	m.Set(`typeGroups`, typeGroups)
	m.Set(`typeGroups`, typeGroups)
	m.Set(`foreignKeys`, foreignKeys)
	m.Set(`onActions`, strings.Split(OnActions, `|`))
	m.Set(`unsignedTags`, UnsignedTags)
	if tableStatus != nil {
		form := m.Request().Form()
		form.Set(`engine`, tableStatus.Engine.String)
		form.Set(`name`, tableStatus.Name.String)
		form.Set(`collation`, tableStatus.Collation.String)
		form.Set(`comment`, tableStatus.Comment.String)
	}
	if len(postFields) == 0 {
		postFields = append(postFields, &Field{})
	}
	m.Set(`postFields`, postFields)
	m.SetFunc(`isString`, reFieldTypeText.MatchString)
	m.SetFunc(`isNumeric`, reFieldTypeNumber.MatchString)
	supportPartitioning := m.support(`partitioning`)
	if supportPartitioning {
		partition, err := m.tablePartitions(oldTable)
		if err != nil {
			supportPartitioning = false
		}
		partition.Names = append(partition.Names, ``)
		partition.Values = append(partition.Values, ``)
		m.Set(`partition`, partition)
	}
	m.Set(`supportPartitioning`, supportPartitioning)
	m.Set(`partitionTypes`, PartitionTypes)
	return m.Render(`db/mysql/create_table`, err)
}
func (m *mySQL) listTableAjax(opType string) error {
	switch opType {
	case `analyze`, `optimize`, `check`, `repair`:
		tables := m.FormValues(`table[]`)
		views := m.FormValues(`view[]`)
		data := m.Data()
		err := m.optimizeTables(append(tables, views...), opType)
		if err != nil {
			data.SetError(err)
		} else {
			data.SetData(m.SavedResults())
		}
		return m.JSON(data)
	case `truncate`:
		tables := m.FormValues(`table[]`)
		//views := m.FormValues(`view[]`)
		data := m.Data()
		var err error
		if len(tables) > 0 {
			err = m.truncateTables(tables)
		}
		if err != nil {
			data.SetError(err)
		} else {
			data.SetData(m.SavedResults())
		}
		return m.JSON(data)
	case `drop`:
		tables := m.FormValues(`table[]`)
		views := m.FormValues(`view[]`)
		data := m.Data()
		var err error
		if len(tables) > 0 {
			err = m.dropTables(tables, false)
		}
		if len(views) > 0 {
			err = m.dropTables(views, true)
		}
		if err != nil {
			data.SetError(err)
		} else {
			data.SetData(m.SavedResults())
		}
		return m.JSON(data)
	case `copy`:
		destDb := m.Form(`dbName`)
		tables := m.FormValues(`table[]`)
		views := m.FormValues(`view[]`)
		data := m.Data()
		var err error
		if len(tables) > 0 {
			err = m.copyTables(tables, destDb, false)
		}
		if len(views) > 0 {
			err = m.copyTables(views, destDb, true)
		}
		if err != nil {
			data.SetError(err)
		} else {
			data.SetData(m.SavedResults())
		}
		return m.JSON(data)
	case `move`:
		destDb := m.Form(`dbName`)
		tables := m.FormValues(`table[]`)
		views := m.FormValues(`view[]`)
		data := m.Data()
		err := m.moveTables(append(tables, views...), destDb)
		if err != nil {
			data.SetError(err)
		} else {
			data.SetData(m.SavedResults())
		}
		return m.JSON(data)
	case `dbs`:
		data := m.Data()
		dbList, err := m.getDatabases()
		if err != nil {
			data.SetError(err)
		} else {
			data.SetData(dbList)
		}
		return m.JSON(data)
	}
	return nil
}
func (m *mySQL) ListTable() error {
	opType := m.Form(`json`)
	if len(opType) > 0 {
		return m.listTableAjax(opType)
	}
	var err error
	if len(m.dbName) > 0 {
		tableList, ok := m.Get(`tableList`).([]string)
		if !ok {
			tableList, err = m.getTables()
			if err != nil {
				m.fail(err.Error())
				return m.returnTo(m.GenURL(`listDb`))
			}
			m.Set(`tableList`, tableList)
		}
		var tableStatus map[string]*TableStatus
		tableStatus, _, err = m.getTableStatus(m.dbName, ``, true)
		if err != nil {
			m.fail(err.Error())
			return m.returnTo(m.GenURL(`listDb`))
		}
		m.Set(`tableStatus`, tableStatus)
	}
	return m.Render(`db/mysql/list_table`, err)
}
func (m *mySQL) ViewTable() error {
	var err error
	oldTable := m.Form(`table`)
	if m.Formx(`ddl`).Bool() {
		data := m.Data()
		ddl, err := m.tableDDL(oldTable)
		if err != nil {
			return m.JSON(data.SetError(err))
		}
		return m.JSON(data.SetData(echo.H{`ddl`: ddl}))
	}
	foreignKeys, sortForeignKeys, err := m.tableForeignKeys(oldTable)
	if err != nil {
		return err
	}
	var (
		origFields   map[string]*Field
		sortFields   []string
		origIndexes  map[string]*Indexes
		sortIndexes  []string
		origTriggers map[string]*Trigger
		sortTriggers []string
		tableStatus  *TableStatus
	)
	if len(oldTable) > 0 {
		val, sort, err := m.tableFields(oldTable)
		if err != nil {
			return err
		}
		origFields = val
		sortFields = sort
		stt, _, err := m.getTableStatus(m.dbName, oldTable, false)
		if err != nil {
			return err
		}
		if ts, ok := stt[oldTable]; ok {
			tableStatus = ts
		}
		val2, sort2, err := m.tableIndexes(oldTable)
		if err != nil {
			return err
		}
		origIndexes = val2
		sortIndexes = sort2
	} else {
		origFields = map[string]*Field{}
		sortFields = []string{}
		origIndexes = map[string]*Indexes{}
		sortIndexes = []string{}
	}
	if tableStatus == nil {
		tableStatus = &TableStatus{}
	}
	postFields := make([]*Field, len(sortFields))
	for k, v := range sortFields {
		postFields[k] = origFields[v]
	}
	indexes := make([]*Indexes, len(sortIndexes))
	for k, v := range sortIndexes {
		indexes[k] = origIndexes[v]
	}
	forkeys := make([]*ForeignKeyParam, len(sortForeignKeys))
	for k, v := range sortForeignKeys {
		forkeys[k] = foreignKeys[v]
	}
	m.Set(`tableStatus`, tableStatus)
	m.Set(`postFields`, postFields)
	m.Set(`indexes`, indexes)
	m.Set(`version`, m.getVersion())
	m.Set(`foreignKeys`, forkeys)
	triggerName := `trigger`
	if tableStatus.IsView() {
		triggerName = `view_trigger`
	}
	supported := m.support(triggerName)
	m.Set(`supportTrigger`, supported)
	if supported {
		origTriggers, sortTriggers, err = m.tableTriggers(oldTable)
		if err != nil {
			return err
		}
		triggers := make([]*Trigger, len(sortTriggers))
		for k, v := range sortTriggers {
			triggers[k] = origTriggers[v]
		}
		m.Set(`triggers`, triggers)
	}
	return m.Render(`db/mysql/view_table`, m.checkErr(err))
}
func (m *mySQL) ListData() error {
	var err error
	table := m.Form(`table`)
	limit := m.Formx(`limit`).Int()
	page := m.Formx(`page`).Int()
	totalRows := m.Formx(`rows`).Int()
	textLength := m.Formx(`text_length`).Int()
	if limit < 1 {
		limit = 50
		m.Request().Form().Set(`limit`, strconv.Itoa(limit))
	}
	if page < 1 {
		page = 1
	}
	if textLength < 1 {
		textLength = 100
		m.Request().Form().Set(`text_length`, strconv.Itoa(textLength))
	}
	selectFuncs := m.FormValues(`columns[fun][]`)
	if len(selectFuncs) == 0 {
		m.Request().Form().Set(`columns[fun][]`, ``)
	}
	selectCols := m.FormValues(`columns[col][]`)

	whereCols := m.FormValues(`where[col][]`)
	whereOperators := m.FormValues(`where[op][]`)
	whereVals := m.FormValues(`where[val][]`)
	if len(whereCols) == 0 {
		m.Request().Form().Set(`where[col][]`, ``)
	}
	if len(whereVals) == 0 {
		m.Request().Form().Set(`where[val][]`, ``)
	}

	opNum := len(whereOperators)
	valNum := len(whereVals)
	orderFields := m.FormValues(`order[]`)

	if len(orderFields) == 0 {
		m.Request().Form().Set(`order[]`, ``)
	}

	descs := m.FormValues(`desc[]`)
	if sort := m.Form(`sort`); len(sort) > 0 {
		if sort[0] == '-' {
			orderFields = []string{sort[1:]}
			descs = []string{`1`}
		} else {
			orderFields = []string{sort}
			descs = []string{`0`}
		}
	}
	var wheres []string
	fields, sortFields, err := m.tableFields(table)
	if err != nil {
		return err
	}
	for index, colName := range whereCols {
		if index >= opNum || index >= valNum {
			break
		}
		invalidOperator := true
		for _, op := range operators {
			if op == whereOperators[index] {
				invalidOperator = false
				break
			}
		}
		if invalidOperator {
			continue
		}
		field, ok := fields[colName]
		if !ok {
			continue
		}
		op := whereOperators[index]
		val := whereVals[index]
		cond := ` ` + op
		switch op {
		case `SQL`:
			cond = ` ` + val
		case `LIKE %%`:
			cond = ` LIKE ` + processInput(field, `%`+val+`%`, ``)
		case `ILIKE %%`:
			cond = ` ILIKE ` + processInput(field, `%`+val+`%`, ``)
		default:
			if strings.HasSuffix(op, `IN`) {
				in, er := m.processLength(val)
				if er != nil {
					return er
				}
				if len(in) > 0 {
					cond += ` ` + in
				} else {
					cond += ` (NULL)`
				}
			} else if !strings.HasSuffix(op, `NULL`) {
				cond += ` ` + processInput(field, val, ``)
			}
		}

		if len(colName) == 0 {
			cols := []string{}
			charset := getCharset(m.getVersion())
			for _, fieldName := range sortFields {
				field := fields[fieldName]
				isText := reFieldTypeText.MatchString(field.Type)
				if (reOnlyNumber.MatchString(val) || !reFieldTypeNumber.MatchString(field.Type)) &&
					(!reChineseAndPunctuation.MatchString(val) || isText) {
					name := quoteCol(fieldName)
					col := name
					if m.supportSQL && isText && !strings.HasPrefix(field.Collation, `utf8_`) {
						col = "CONVERT(" + name + " USING " + charset + ")"
					}
					cols = append(cols, col)
				}
			}
			if len(cols) > 0 {
				wheres = append(wheres, `(`+strings.Join(cols, cond+` OR `)+cond+`)`)
			} else {
				wheres = append(wheres, `0`)
			}
		} else {
			wheres = append(wheres, quoteCol(colName)+cond)
		}
	}
	if m.IsPost() {
		save := m.Form(`save`)
		inputName := `check[]`
		multiSelection := true
		if save == `set` {
			inputName = `pk`
			multiSelection = false
		}
		condition, err := m.genCheckedCond(fields, wheres, multiSelection, inputName)
		if err != nil {
			return err
		}
		if len(condition) > 0 {
			switch save {
			case `delete`:
				condition = ` WHERE ` + condition
				err = m.delete(table, condition, 0)
			case `export`:
				return m.exportData(fields, table, selectFuncs, selectCols, []string{condition}, orderFields, descs, 1, limit, totalRows, 0)
			case `set`:
				condition = ` WHERE ` + condition
				key := m.Form(`name`)
				values := m.FormValues(`value[]`)
				value := m.Form(`value`)
				if len(values) > 0 {
					value = strings.Join(values, ",")
				}
				err = m.set(table, condition, key, value, 1)
				data := m.Data()
				if err != nil {
					data.SetError(err)
				} else {
					data.SetInfo(`修改成功`)
				}
				return m.JSON(data)
			}
			if err == nil {
				return m.returnTo()
			}
		}
	}
	var (
		columns []string
		values  []map[string]*sql.NullString
	)
	columns, values, totalRows, err = m.listData(nil, table, selectFuncs, selectCols, wheres, orderFields, descs, page, limit, totalRows, textLength)
	m.Set(`sortFields`, sortFields)
	m.Set(`fields`, fields)
	m.Set(`columns`, columns)
	m.Set(`values`, values)
	m.Set(`functions`, functions)
	m.Set(`grouping`, grouping)
	m.Set(`operators`, operators)
	m.Set(`total`, totalRows)
	m.SetFunc(`isBlobData`, func(colName string) bool {
		f, y := fields[colName]
		if !y {
			return false
		}
		return reFieldTypeBlob.MatchString(f.Type)
	})
	indexes, _, err := m.tableIndexes(table)
	m.SetFunc(`uniqueIdf`, func(row map[string]*sql.NullString) string {
		idf := ``
		uniqueArr := uniqueArray(row, indexes)
		if len(uniqueArr) == 0 {
			uniqueArr = map[string]*sql.NullString{}
			for key, val := range row {
				if !reSQLFunction.MatchString(key) {
					uniqueArr[key] = val
				}
			}
		}
		for key, val := range uniqueArr {
			field, y := fields[key]
			if !y {
				fmt.Printf(`not exists: %v in %#v`+"\n", key, fields)
				return idf
			}
			if (m.supportSQL || m.DbAuth.Driver == "pgsql") && len(val.String) > 64 {
				//! columns looking like functions
				if strings.Index(key, `(`) <= 0 {
					key = quoteCol(key)
				}
				if m.supportSQL && strings.HasPrefix(field.Collation, `utf8_`) {
					key = "MD5(" + key + ")"
				} else {
					key = "MD5(CONVERT(" + key + " USING " + getCharset(m.getVersion()) + "))"
				}
				val.String = com.Md5(val.String)
			}
			if val.Valid {
				idf += "&" + url.QueryEscape("where["+bracketEscape(key, false)+"]") + "=" + url.QueryEscape(val.String)
			} else {
				idf += "&null%5B%5D=" + url.QueryEscape(key)
			}
		}
		return idf
	})
	q := m.Request().URL().Query()
	q.Del(`page`)
	q.Del(`rows`)
	q.Del(`_pjax`)
	m.Set(`pagination`, pagination.New(m.Context).SetURL(`/db?`+q.Encode()+`&page={page}&rows={rows}`).SetPage(page).SetRows(totalRows))
	return m.Render(`db/mysql/list_data`, m.checkErr(err))
}
func (m *mySQL) genCheckedCond(fields map[string]*Field, wheres []string, multiSelection bool, inputNames ...string) (condition string, err error) {
	var conds []string
	var inputName string
	for _, inName := range inputNames {
		if len(inName) > 0 {
			inputName = inName
			break
		}
	}
	if len(inputName) == 0 {
		inputName = `check[]`
	}
	if multiSelection {
		conds = m.FormValues(inputName)
	} else {
		conds = []string{m.Form(inputName)}
	}
	datas := []string{}
	for _, cond := range conds {
		cond = strings.TrimLeft(cond, `&`)
		cond, err = url.QueryUnescape(cond)
		if err != nil {
			return
		}
		values, err := url.ParseQuery(cond)
		if err != nil {
			return ``, err
		}
		mpx := echo.NewMapx(values)
		where := mpx.Get(`where`)
		null := mpx.Get(`null`)
		if where == nil && null == nil {
			continue
		}
		cond = m.whereByMapx(where, null, fields)
		if len(cond) < 1 {
			continue
		}
		datas = append(datas, cond)
	}
	if len(datas) > 0 {
		condition = `(` + strings.Join(datas, `) OR (`) + `)`
		if len(wheres) > 0 {
			condition = `(` + strings.Join(wheres, ` AND `) + `) AND (` + condition + `)`
		}
	}
	return
}
func (m *mySQL) CreateData() error {
	var err error
	saveType := m.Form(`save`)
	clone := m.Formx(`clone`).Bool()
	table := m.Form(`table`)
	fields, sortFields, err := m.tableFields(table)
	if err != nil {
		return err
	}
	var condition string
	var where *echo.Mapx
	condition, err = m.genCheckedCond(fields, nil, false)
	if err != nil {
		return err
	}
	if len(condition) == 0 {
		mapx := echo.NewMapx(m.Forms())
		where = mapx.Get(`where`)
		null := mapx.Get(`null`)
		if where != nil || null != nil {
			condition = m.whereByMapx(where, null, fields)
		}
	}
	var columns []string
	values := map[string]*sql.NullString{}
	sqlStr := `SELECT * FROM ` + quoteCol(table)
	var whereStr string
	if len(condition) > 0 {
		whereStr = ` WHERE ` + condition
	}
	if m.IsPost() && (saveType == `save` || saveType == `saveAndContinue` || saveType == `delete`) {
		indexes, _, err := m.tableIndexes(table)
		if err != nil {
			return err
		}
		wheres := map[string]*sql.NullString{}
		if where != nil {
			for k, v := range where.Map {
				val := &sql.NullString{}
				val.String, val.Valid = v.ValueOk()
				wheres[k] = val
			}
		}
		uniqueArr := uniqueArray(wheres, indexes)
		var limit int
		if len(uniqueArr) > 0 {
			limit = 1
		}
		if saveType == `delete` {
			err = m.delete(table, whereStr, limit)
		} else {
			set := map[string]string{}
			for _, col := range sortFields {
				field, ok := fields[col]
				if !ok {
					continue
				}
				v, y := m.processInputFieldValue(field)
				if !y {
					continue
				}
				set[col] = v
			}
			if len(whereStr) > 0 && !clone {
				err = m.update(table, set, whereStr, limit)
			} else {
				err = m.insert(table, set)
			}
		}
		if err == nil && (saveType == `save` || saveType == `delete`) {
			return m.returnTo(m.GenURL(`listData`, m.dbName, table))
		}
	}
	if len(whereStr) > 0 {
		rows, err := m.newParam().SetCollection(sqlStr + whereStr).Query()
		if err != nil {
			return err
		}
		columns, err = rows.Columns()
		size := len(columns)
		for rows.Next() {
			recv := make([]interface{}, size)
			for i := 0; i < size; i++ {
				recv[i] = &sql.NullString{}
			}
			err = rows.Scan(recv...)
			if err != nil {
				continue
			}
			for k, colName := range columns {
				values[colName] = recv[k].(*sql.NullString)
			}
			break
		}
	} else {
		columns = sortFields
		for _, v := range sortFields {
			values[v] = &sql.NullString{}
		}
	}
	m.Set(`columns`, columns)
	m.Set(`values`, values)
	m.Set(`fields`, fields)
	if clone {
		m.Set(`saveType`, `copy`)
	} else {
		m.Set(`saveType`, saveType)
	}
	m.SetFunc(`isNumber`, func(typ string) bool {
		return reFieldTypeNumber.MatchString(typ)
	})
	m.SetFunc(`mumberStep`, func(field *Field) string {
		if field.Precision < 1 {
			return `1`
		}
		return fmt.Sprintf(`0.%0*d`, field.Precision, 1)
	})
	m.SetFunc(`isBlob`, func(typ string) bool {
		return reFieldTypeBlob.MatchString(typ)
	})
	m.SetFunc(`isText`, func(typ string) bool {
		return reFieldTextValue.MatchString(typ)
	})
	m.SetFunc(`enumValues`, func(field *Field) []*Enum {
		return enumValues(field)
	})
	m.SetFunc(`functions`, m.editFunctions)
	m.SetFunc(`isSelectedFunc`, func(function string, value *sql.NullString) bool {
		if len(function) == 0 {
			if len(value.String) > 0 {
				return true
			}
			if value.Valid {
				return true
			}
		}
		return false
	})
	return m.Render(`db/mysql/edit_data`, m.checkErr(err))
}
func (m *mySQL) Indexes() error {
	return m.modifyIndexes()
}
func (m *mySQL) modifyIndexes() error {
	table := m.Form(`table`)
	indexTypes := []string{"PRIMARY", "UNIQUE", "INDEX"}
	rule := `(?i)MyISAM|M?aria`
	if com.VersionCompare(m.getVersion(), `5.6`) >= 0 {
		rule += `|InnoDB`
	}
	re, err := regexp.Compile(rule)
	if err != nil {
		return m.String(err.Error())
	}
	status, _, err := m.getTableStatus(m.dbName, table, true)
	if err != nil {
		return m.String(err.Error())
	}
	tableStatus, ok := status[table]
	if ok && re.MatchString(tableStatus.Engine.String) {
		indexTypes = append(indexTypes, "FULLTEXT")
	}
	indexes, sorts, err := m.tableIndexes(table)
	if err != nil {
		return m.String(err.Error())
	}
	if m.IsPost() {
		mapx := echo.NewMapx(m.Forms())
		mapx = mapx.Get(`indexes`)
		alter := []*indexItems{}
		if mapx != nil {
			size := len(mapx.Map)
			for i := 0; i < size; i++ {
				ii := strconv.Itoa(i)
				item := &indexItems{
					Indexes: &Indexes{
						Name:    mapx.Value(ii, `name`),
						Type:    mapx.Value(ii, `type`),
						Columns: mapx.Values(ii, `columns`),
						Lengths: mapx.Values(ii, `lengths`),
						Descs:   mapx.Values(ii, `descs`),
					},
					Set: []string{},
				}
				var typeOk bool
				for _, indexType := range indexTypes {
					if item.Type == indexType {
						typeOk = true
						break
					}
				}
				if !typeOk {
					continue
				}
				lenSize := len(item.Lengths)
				descSize := len(item.Descs)
				columns := []string{}
				lengths := []string{}
				descs := []string{}
				for key, col := range item.Columns {
					if len(col) == 0 {
						continue
					}
					var length, desc string
					if key < lenSize {
						length = item.Lengths[key]
					}
					if key < descSize {
						desc = item.Descs[key]
					}
					set := quoteCol(col)
					if len(length) > 0 {
						set += `(` + length + `)`
					}
					if len(desc) > 0 {
						set += ` DESC`
					}
					item.Set = append(item.Set, set)
					columns = append(columns, col)
					lengths = append(lengths, length)
					descs = append(descs, desc)
				}
				if len(columns) < 1 {
					continue
				}
				if existing, ok := indexes[item.Name]; ok {
					/*
						fmt.Println(item.Type, `==`, existing.Type)
						fmt.Printf(`columns：%#v`+" == %#v\n", columns, existing.Columns)
						fmt.Printf(`lengths：%#v`+" == %#v\n", lengths, existing.Lengths)
						fmt.Printf(`descs：%#v`+" == %#v\n", descs, existing.Descs)
					// */
					if item.Type == existing.Type && fmt.Sprintf(`%#v`, columns) == fmt.Sprintf(`%#v`, existing.Columns) &&
						fmt.Sprintf(`%#v`, lengths) == fmt.Sprintf(`%#v`, existing.Lengths) &&
						fmt.Sprintf(`%#v`, descs) == fmt.Sprintf(`%#v`, existing.Descs) {
						delete(indexes, item.Name)
						continue
					}
				}
				alter = append(alter, item)
			}
		}
		for name, existing := range indexes {
			alter = append(alter, &indexItems{
				Indexes: &Indexes{
					Name: name,
					Type: existing.Type,
				},
				Set:       []string{},
				Operation: `DROP`,
			})
		}
		if len(alter) > 0 {
			err = m.alterIndexes(table, alter)
		}
		if err != nil {
			m.fail(err.Error())
		}
		return m.returnTo(m.GenURL(`viewTable`, m.dbName, table))
	}
	indexesSlice := make([]*Indexes, len(sorts))
	for k, name := range sorts {
		indexesSlice[k] = indexes[name]
		indexesSlice[k].Columns = append(indexesSlice[k].Columns, "")
		indexesSlice[k].Lengths = append(indexesSlice[k].Lengths, "")
		indexesSlice[k].Descs = append(indexesSlice[k].Descs, "")
	}
	indexesSlice = append(indexesSlice, &Indexes{
		Columns: []string{""},
		Lengths: []string{""},
		Descs:   []string{""},
	})
	fields, sortFields, err := m.tableFields(table)
	if err != nil {
		return m.String(err.Error())
	}
	fieldsSlice := make([]*Field, len(sortFields))
	for k, name := range sortFields {
		fieldsSlice[k] = fields[name]
	}
	m.Set(`indexes`, indexesSlice)
	m.Set(`indexTypes`, indexTypes)
	m.Set(`fields`, fieldsSlice)
	return m.Render(`db/mysql/modify_index`, m.checkErr(err))
}
func (m *mySQL) Foreign() error {
	return m.modifyForeignKeys()
}
func (m *mySQL) modifyForeignKeys() error {
	table := m.Form(`table`)
	name := m.Form(`name`)
	foreignTable := m.Form(`foreign_table`)
	if len(foreignTable) == 0 {
		foreignTable = table
	}
	_, sortFields, err := m.tableFields(table)
	if err != nil {
		return m.String(err.Error())
	}
	status, sortStatus, err := m.getTableStatus(m.dbName, ``, true)
	if err != nil {
		return m.String(err.Error())
	}
	var referencable []string
	for _, tableName := range sortStatus {
		tableStatus := status[tableName]
		if tableStatus.FKSupport(m.getVersion()) {
			referencable = append(referencable, tableName)
		}
	}
	var foreignKey *ForeignKeyParam
	if len(name) > 0 {
		fkeys, _, err := m.tableForeignKeys(table)
		if err != nil {
			return m.String(err.Error())
		}
		var ok bool
		foreignKey, ok = fkeys[name]
		if !ok {
			return m.String(m.T(`外键不存在`))
		}
	} else {
		foreignKey = &ForeignKeyParam{
			Table:  foreignTable,
			Source: []string{},
			Target: []string{},
		}
	}
	drop := m.Form(`drop`)
	isDrop := len(drop) > 0
	if isDrop || m.IsPost() {
		targets := m.FormValues(`target[]`)
		endIndex := len(targets) - 1
		foreignKey.Source = []string{}
		foreignKey.Target = []string{}
		foreignKey.OnDelete = m.Form(`on_delete`)
		foreignKey.OnUpdate = m.Form(`on_update`)
		for i, source := range m.FormValues(`source[]`) {
			if len(source) == 0 {
				continue
			}
			if i > endIndex || len(targets[i]) == 0 {
				continue
			}
			foreignKey.Source = append(foreignKey.Source, source)
			foreignKey.Target = append(foreignKey.Target, targets[i])
		}
		if len(name) > 0 && len(foreignKey.Source) == 0 {
			isDrop = true
		}
		err = m.alterForeignKeys(table, foreignKey, isDrop)
		if err != nil {
			m.fail(err.Error())
		}
		return m.returnTo(m.GenURL(`viewTable`, m.dbName, table))
	}
	foreignKey.Source = append(foreignKey.Source, "")
	foreignKey.Target = append(foreignKey.Target, "")
	var target []string
	if foreignKey.Table == table {
		target = sortFields
	} else {
		_, target, err = m.tableFields(foreignKey.Table)
		if err != nil {
			return m.String(err.Error())
		}
	}
	m.Set(`source`, sortFields)         //源(当前表中的字段)
	m.Set(`target`, target)             //目标(外部表中的字段)
	m.Set(`referencable`, referencable) //可以使用的目标表
	m.Set(`onActions`, strings.Split(OnActions, `|`))
	m.Set(`foreign`, foreignKey)
	return m.Render(`db/mysql/modify_foreign`, m.checkErr(err))
}
func (m *mySQL) Trigger() error {
	return m.modifyTrigger()
}
func (m *mySQL) modifyTrigger() error {
	var err error
	table := m.Form(`table`)
	name := m.Form(`name`)
	var trigger *Trigger
	if len(name) > 0 {
		trigger, err = m.tableTrigger(name)
		if err != nil {
			return err
		}
	}
	if trigger == nil {
		trigger = &Trigger{}
	}
	if m.IsPost() {
		if len(name) > 0 {
			err = m.dropTrigger(table, name)
			if len(m.Form(`drop`)) > 0 {
				return m.returnTo(m.GenURL(`viewTable`, m.dbName, table))
			}
		}
		trigger.Timing.String = m.Form(`timing`)
		trigger.Event.String = m.Form(`event`)
		trigger.Type = m.Form(`type`)
		trigger.Of = m.Form(`of`)
		trigger.Trigger.String = m.Form(`trigger`)
		trigger.Statement.String = m.Form(`statement`)
		err = m.createTrigger(table, trigger)
		return m.returnTo(m.GenURL(`viewTable`, m.dbName, table))
	}
	if len(trigger.Trigger.String) == 0 {
		trigger.Trigger.String = table + `_bi`
	}
	m.Set(`trigger`, trigger)
	m.Set(`triggerOptions`, m.TriggerOptions)
	return m.Render(`db/mysql/modify_trigger`, m.checkErr(err))
}
func (m *mySQL) RunCommand() error {
	var err error
	selects := []*SelectData{}
	if m.IsPost() {
		query := m.Form(`query`)
		query = strings.TrimSpace(query)
		errorStops := m.Formx(`error_stops`).Bool()
		onlyErrors := m.Formx(`only_errors`).Bool()
		limit := m.Formx(`limit`).Int()
		if limit <= 0 {
			limit = 50
		}
		var reader *bytes.Reader
		reader = bytes.NewReader([]byte(query))
		space := "(?:\\s|/\\*[\\s\\S]*?\\*/|(?:#|-- )[^\\n]*\\n?|--\\r?\\n)"
		delimiter := ";"
		parse := `['"`
		empty := true
		switch m.DbAuth.Driver {
		case `sqlite`:
			parse += "`["
		case `mssql`:
			parse += "["
		default:
			if strings.Contains(m.DbAuth.Driver, `sql`) {
				parse += "`#"
			}
		}
		parse += "]|/\\*|-- |$"
		switch m.DbAuth.Driver {
		case `sqlite`:
			parse += "|\\$[^$]*\\$"
		}
		buf := make([]byte, 1e6)
		query = ``
		offset := 0
		for {
			n, e := reader.Read(buf)
			if e != nil {
				if e == io.EOF {
					break
				}
				m.Logger().Error(err)
			}
			q := string(buf[0:n])
			if offset == 0 {
				if match := regexp.MustCompile("(?i)^" + space + "*DELIMITER\\s+(\\S+)").FindStringSubmatch(q); len(match) > 1 {
					delimiter = match[1]
					q = q[len(match[0]):]
					query += q
					offset += n
					continue
				}
			}
			query += q
			offset += n

			/*/ 跳过注释和空白
			match := regexp.MustCompile("(" + regexp.QuoteMeta(delimiter) + "\\s*|" + parse + ")").FindStringSubmatch(query)
			com.Dump(match)
			if len(match) > 1 {
				found := match[1]
				if strings.TrimRight(query, " \t\n\r") != delimiter {
					rule := `(?s)`
					switch found {
					case `/*`:
						rule += "\\*\/"
					case `[`:
						rule += `]`
					default:
						match := regexp.MustCompile("^-- |^#").FindStringSubmatch(found)
						if len(match) > 1 {
							rule += "\n"
						} else {
							rule += regexp.QuoteMeta(found) + "|\\\\."
						}
					}
					pos := strings.Index(query, found)
					query = query[:pos]
					rule += `|$`
					match := regexp.MustCompile(rule).FindStringSubmatch(query)
					for len(match) > 0 {
						n, e := reader.Read(buf)
						if e != nil {
							if e == io.EOF {
								break
							}
							m.Logger().Error(err)
						}
						q := string(buf[0:n])
						if len(match) > 1 && len(match[1]) > 0 && match[1][0] != '\\' {
							break
						}
						match = regexp.MustCompile(rule).FindStringSubmatch(q)
					}
				}
			}
			// */

			empty = false
			if m.DbAuth.Driver == `sqlite` && regexp.MustCompile(`(?i)^`+space+`*ATTACH\b`).MatchString(query) {
				if errorStops {
					err = errors.New(m.T(`ATTACH queries are not supported.`))
					break
				}
			}

			if regexp.MustCompile(`(?i)^` + space + `*USE\b`).MatchString(query) {
				_, err = m.newParam().DB().Exec(query)
				if err != nil {
					m.Logger().Error(err, query)
					if onlyErrors {
						return err
					}
				}
				continue
			}

			if regexp.MustCompile(`(?i)^` + space + `*(CREATE|DROP|ALTER)` + space + `+(DATABASE|SCHEMA)\b`).MatchString(query) {
				_, err = m.newParam().DB().Exec(query)
				if err != nil {
					m.Logger().Error(err, query)
					if onlyErrors {
						return err
					}
				}
				continue
			}

			if !regexp.MustCompile(`(?i)^(` + space + `|\()*(SELECT|SHOW|EXPLAIN)\b`).MatchString(query) {
				var sqlStr string
				execute := func(line string) (rErr error) {
					if strings.HasPrefix(line, `--`) {
						return nil
					}
					if strings.HasPrefix(line, `/*`) && strings.HasSuffix(line, `*/;`) {
						return nil
					}
					line = strings.TrimSpace(line)
					sqlStr += line
					if strings.HasSuffix(line, `;`) && len(sqlStr) > 0 {
						defer func() {
							sqlStr = ``
						}()
						r := &Result{
							SQL: sqlStr,
						}
						r.Exec(m.newParam())
						m.AddResults(r)
						return r.Error()
					}
					return nil
				}
				for _, line := range strings.Split(query, "\n") {
					line = strings.TrimSpace(line)
					if len(line) == 0 {
						continue
					}
					err = execute(line)
					if err != nil {
						m.Logger().Error(err, line)
						if onlyErrors {
							return err
						}
					}
				}
				continue
			}
			r := &Result{
				SQL: query,
			}
			dt := &DataTable{}
			r.Query(m.newParam(), func(rows *sql.Rows) error {
				dt.Columns, dt.Values, err = m.selectTable(rows, limit)
				return err
			})
			if r.err != nil {
				m.Logger().Error(r.err, query)
				if onlyErrors {
					return err
				}
				continue
			}
			selectData := &SelectData{Result: r, Data: dt}
			if regexp.MustCompile(`(?i)^(` + space + `|\()*SELECT\b`).MatchString(query) {
				rows, err := m.newParam().DB().Query(`EXPLAIN ` + query)
				if err != nil {
					m.Logger().Error(err, `EXPLAIN `+query)
					if onlyErrors {
						return err
					}
					continue
				}
				dt := &DataTable{}
				dt.Columns, dt.Values, err = m.selectTable(rows, limit)
				selectData.Explain = dt
			}
			selects = append(selects, selectData)
			/*
				com.Dump(columns)
				com.Dump(values)
			// */
		}
		_ = delimiter
		_ = empty
	}
	m.Set(`selects`, selects)
	return m.Render(`db/mysql/sql`, m.checkErr(err))
}
