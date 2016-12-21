package mysql

import (
	"fmt"
	"strings"

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
	if len(settings.Database) == 0 {
		settings.Database = m.Form(`db`)
	}
	m.dbName = settings.Database
	db, err := mysql.Open(settings)
	if err != nil {
		return err
	}
	cluster := factory.NewCluster().AddW(db)
	m.db.SetCluster(0, cluster)
	m.Set(`dbName`, settings.Database)
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
	r, e := m.processList()
	m.Set(`processList`, r)
	return m.Render(`db/mysql/proccess_list`, e)
}

func (m *mySQL) returnTo() error {
	returnTo := m.Form(`return_to`)
	if len(returnTo) == 0 {
		returnTo = m.Request().URI()
	}
	return m.Redirect(returnTo)
}

func (m *mySQL) Privileges() error {
	var ret interface{}
	var err error
	act := m.Form(`act`)
	if len(act) > 0 {
		switch act {
		case `drop`:
		case `edit`:
			if m.IsPost() {
				isHashed := len(m.Form(`hashed`)) > 0
				user := m.Form(`oldUser`)
				host := m.Form(`host`)
				newUser := m.Form(`user`)
				oldPasswd := m.Form(`oldPass`)
				newPasswd := m.Form(`pass`)
				r := m.editUser(user, host, newUser, oldPasswd, newPasswd, isHashed)
				if r.Error == nil {
					m.Session().AddFlash(common.Ok(m.T(`操作成功`)))
					return m.returnTo()
				}
				m.Session().AddFlash(r.Error.Error())
			}
			privs, err := m.showPrivileges()
			if err == nil {
				privs.Parse()
			}
			m.Set(`list`, privs.privileges)
			user := m.Form(`user`)
			host := m.Form(`host`)
			var (
				oldPass string
				grants  map[string]map[string]bool
				sorts   []string
				oldUser string
			)
			if len(host) > 0 {
				oldPass, grants, sorts, err = m.getUserGrants(host, user)
				if _, ok := grants["*.*"]; ok {
					m.Set(`serverAdminObject`, "*.*")
				} else {
					m.Set(`serverAdminObject`, ".*")
				}
				if err == nil {
					oldUser = user
				}
			}
			m.Set(`sorts`, sorts)
			m.Set(`grants`, grants)
			m.Set(`hashed`, true)
			m.Set(`oldPass`, oldPass)
			m.Set(`oldUser`, oldUser)
			m.Request().Form().Set(`pass`, oldPass)
			m.SetFunc(`getGrantByPrivilege`, func(grant map[string]bool, index int, privilege string) bool {
				priv := strings.ToUpper(privilege)
				value := m.Form(fmt.Sprintf(`grants[%v][%v]`, index, priv))
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
			m.SetFunc(`getScope`, func(object string) *Grant {
				g := &Grant{Value: object}
				if object == `*.*` {
					g.Scope = `all`
					return g
				}
				strs := strings.SplitN(object, `.`, 2)
				for i, v := range strs {
					v = strings.Trim(v, "`")
					switch i {
					case 0:
						g.Database = v
					case 1:
						if v == `*` {
							g.Scope = `database`
						} else if strings.HasSuffix(v, `)`) {
							vs := strings.SplitN(v, `(`, 2)
							switch len(vs) {
							case 2:
								g.Table = strings.TrimSpace(vs[0])
								g.Table = strings.TrimSuffix(g.Table, "`")
								g.Columns = strings.TrimSuffix(vs[1], `)`)
								g.Scope = `column`
							}
						} else {
							g.Table = strings.TrimSpace(v)
							g.Table = strings.TrimSuffix(g.Table, "`")
							g.Scope = `table`
						}
					}
				}
				return g
			})
			m.SetFunc(`fieldName`, func(index int, privilege string) string {
				return fmt.Sprintf(`grants[%v][%v]`, index, strings.ToUpper(privilege))
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
		if res.Error != nil {
			data.SetError(res.Error)
		} else {
			data.SetData(res)
		}
	}
	return m.JSON(data)
}
func (m *mySQL) ModifyDb() error {
	return nil
}
func (m *mySQL) ListDb() error {
	switch m.Form(`json`) {
	case `drop`:
		data := m.NewData()
		dbs := m.FormValues(`db[]`)
		rs := []*Result{}
		code := 1
		for _, db := range dbs {
			r := m.dropDatabase(db)
			rs = append(rs, r)
			if r.Error != nil {
				data.SetError(r.Error)
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
			tableStatus, err = m.getTableStatus(dbName, ``, true)
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
	//return m.String(`OK`)
	return m.Render(`db/mysql/list_db`, err)
}
func (m *mySQL) CreateTable() error {
	return nil
}
func (m *mySQL) ModifyTable() error {
	return nil
}
func (m *mySQL) ListTable() error {
	var err error
	if len(m.dbName) > 0 {
		if _, ok := m.Get(`tableList`).([]string); !ok {
			tableList, err := m.getTables()
			if err != nil {
				return err
			}
			m.Set(`tableList`, tableList)
		}
	}
	return m.Render(`db/mysql/list_table`, err)
}
func (m *mySQL) ViewTable() error {
	return nil
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
