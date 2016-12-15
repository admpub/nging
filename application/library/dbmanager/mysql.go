package dbmanager

import (
	"database/sql"
	"fmt"

	"github.com/webx-top/db/lib/factory"
	"github.com/webx-top/db/mysql"
	"github.com/webx-top/echo"
)

func init() {
	Register(`MySQL`, &mySQL{})
}

type mySQL struct {
	*BaseDriver
	db *factory.Factory
}

func (m *mySQL) Init(ctx echo.Context) {
	m.BaseDriver = NewBaseDriver()
	m.BaseDriver.Init(ctx)
}

func (m *mySQL) IsSupported(operation string) bool {
	return true
}
func (m *mySQL) Login() error {
	m.db = factory.New()
	settings := mysql.ConnectionURL{
		User:     m.Form(`username`),
		Password: m.Form(`password`),
		Host:     m.Form(`host`),
		Database: m.Form(`db`),
	}
	if len(settings.User) == 0 {
		settings.User = `root`
	}
	if len(settings.Host) == 0 {
		settings.Host = `127.0.0.1:3306`
	}
	m.Echo().Logger().Debugf("db settings: %#v", settings)
	db, err := mysql.Open(settings)
	if err != nil {
		return err
	}
	cluster := factory.NewCluster().AddW(db)
	m.db.SetCluster(0, cluster)
	rows, err := m.db.Query(factory.NewParam(m.db).SetCollection(`SHOW DATABASES`))
	if err != nil {
		return err
	}
	for rows.Next() {
		var v sql.NullString
		err := rows.Scan(&v)
		if err != nil {
			return err
		}
		fmt.Println(`database: `, v.String)
	}
	return nil
}
func (m *mySQL) Logout() error {
	if m.db != nil {
		m.db.CloseAll()
		m.db = nil
	}
	return nil
}
func (m *mySQL) ProcessList() error {
	return nil
}
func (m *mySQL) Privileges() error {
	return nil
}
func (m *mySQL) Info() error {
	return nil
}
func (m *mySQL) CreateDb() error {
	return nil
}
func (m *mySQL) ModifyDb() error {
	return nil
}
func (m *mySQL) ListDb() error {
	return nil
}
func (m *mySQL) CreateTable() error {
	return nil
}
func (m *mySQL) ModifyTable() error {
	return nil
}
func (m *mySQL) ListTable() error {
	return nil
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
