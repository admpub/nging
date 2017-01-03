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

import (
	"errors"
	"fmt"
	"strings"

	"strconv"

	"database/sql"

	"github.com/admpub/nging/application/library/common"
	"github.com/admpub/nging/application/library/dbmanager/driver"
	"github.com/webx-top/db/lib/factory"
	"github.com/webx-top/db/mysql"
	"github.com/webx-top/echo"
)

func init() {
	driver.Register(`MySQL`, &mySQL{})
}

type mySQL struct {
	*driver.BaseDriver
	db      *factory.Factory
	dbName  string
	version string
}

func (m *mySQL) Init(ctx echo.Context, auth *driver.DbAuth) {
	m.BaseDriver = driver.NewBaseDriver()
	m.BaseDriver.Init(ctx, auth)
}

func (m *mySQL) IsSupported(operation string) bool {
	return true
}

func (m *mySQL) Login() error {
	m.db = factory.New()
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
		return err
	}
	cluster := factory.NewCluster().AddW(db)
	m.db.SetCluster(0, cluster)
	m.Set(`dbName`, m.dbName)
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
	return m.Render(`db/mysql/proccess_list`, ret)
}

func (m *mySQL) returnTo(rets ...string) error {
	returnTo := m.Form(`return_to`)
	if len(returnTo) == 0 {
		if len(rets) > 0 {
			returnTo = rets[0]
		} else {
			returnTo = m.Request().Referer()
		}
	}
	m.SaveResults()
	return m.Redirect(returnTo)
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
			r := m.dropUser(user, host)
			m.AddResults(r)
			return m.returnTo()
		case `edit`:
			if m.IsPost() {
				isHashed := len(m.Form(`hashed`)) > 0
				user := m.Form(`oldUser`)
				host := m.Form(`host`)
				newUser := m.Form(`user`)
				oldPasswd := m.Form(`oldPass`)
				newPasswd := m.Form(`pass`)
				err = m.editUser(user, host, newUser, oldPasswd, newPasswd, isHashed)
				if err == nil {
					m.ok(m.T(`操作成功`))
					return m.returnTo(m.GenURL(`privileges`) + `&act=edit&user=` + newUser + `&host=` + host)
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
	data := m.NewData()
	if len(dbName) < 1 {
		data.SetZone(`name`).SetInfo(m.T(`数据库名称不能为空`))
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
		data := m.NewData()
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
		data := m.NewData()
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
		mapx := NewMapx(m.Forms())
		f := mapx.Get(`fields`)
		allFields := []*fieldItem{}
		after := " FIRST"
		foreign := map[string]string{}
		if err == nil && f != nil {
			for i := 0; ; i++ {
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
					String: f.Value(ii, `has_default`),
					Valid:  f.Value(ii, `default`) == `1`,
				}
				field.AutoIncrement = sql.NullString{
					Valid: aiIndexInt == i,
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

	referencablePrimary, _, err := m.referencablePrimary(``)
	foreignKeys := map[string]string{}
	for tblName, field := range referencablePrimary {
		foreignKeys[strings.Replace(tblName, "`", "``", -1)+"`"+strings.Replace(field.Field, "`", "``", -1)] = tblName
	}
	postFields := []*Field{}
	oldTable := m.Form(`table`)
	var origFields map[string]*Field
	var sortFields []string
	var tableStatus *TableStatus
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
	} else {
		origFields = map[string]*Field{}
		sortFields = []string{}
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
		mapx := NewMapx(m.Forms())
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
			for i := 0; ; i++ {
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
						Valid: aiIndexInt == i,
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
		data := m.NewData()
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
		data := m.NewData()
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
		data := m.NewData()
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
		data := m.NewData()
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
		data := m.NewData()
		err := m.moveTables(append(tables, views...), destDb)
		if err != nil {
			data.SetError(err)
		} else {
			data.SetData(m.SavedResults())
		}
		return m.JSON(data)
	case `dbs`:
		data := m.NewData()
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
	return m.Render(`db/mysql/view_table`, err)
}
func (m *mySQL) ListData() error {
	return nil
}
func (m *mySQL) CreateData() error {
	return nil
}
func (m *mySQL) Indexes() error {
	return nil
}
func (m *mySQL) Foreign() error {
	return nil
}
func (m *mySQL) Trigger() error {
	return nil
}
func (m *mySQL) RunCommand() error {
	return nil
}
func (m *mySQL) Import() error {
	return nil
}
func (m *mySQL) Export() error {
	return nil
}
