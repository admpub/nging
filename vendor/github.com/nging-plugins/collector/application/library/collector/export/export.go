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

package export

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/admpub/gopiper"
	"github.com/admpub/log"
	"github.com/admpub/marmot/miner"
	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/sqlbuilder"
	"github.com/webx-top/db/mysql"
	"github.com/webx-top/echo"

	"github.com/nging-plugins/collector/application/dbschema"
	"github.com/nging-plugins/collector/application/library/collector/exec"
	"github.com/nging-plugins/collector/application/library/collector/sender"
	dbmgrdbschema "github.com/nging-plugins/dbmanager/application/dbschema"
)

var emptyPipeItem = &gopiper.PipeItem{}

// Export 导出到数据库
// result - 结果
// data - 采集结果
func (m *Mappings) Export(result *exec.Recv, data echo.Store, config *dbschema.NgingCollectorExport, noticeSender sender.Notice) error {
	switch config.DestType {
	case `API`:
		return m.Export2API(result, data, config, noticeSender)
	case `DSN`, `dbAccountID`:
		return m.Export2DB(result, data, config, noticeSender)
	default:
		return errors.New(`Unsupported DestType: ` + config.DestType)
	}
}

// Export2DB 导出到数据库
// result - 结果
// data - 本次采集结果
func (m *Mappings) Export2DB(result *exec.Recv, data echo.Store, config *dbschema.NgingCollectorExport, noticeSender sender.Notice) error {
	var (
		dbs sqlbuilder.Database
		err error
	)
	if config.DestType == `dbAccountID` {
		accountM := dbmgrdbschema.NewNgingDbAccount(config.Context())
		err = accountM.Get(nil, db.Cond{`id`: config.Dest})
		if err != nil {
			return err
		}
		settings := mysql.ConnectionURL{
			User:     accountM.User,
			Password: accountM.Password,
			Host:     accountM.Host,
			Database: accountM.Name,
			Options:  map[string]string{},
		}
		if len(accountM.Options) > 0 {
			err = com.JSONDecode([]byte(accountM.Options), &accountM.Options)
			if err != nil {
				return err
			}
		}
		dbs, err = DBConn(settings)
	} else {
		dbs, err = DBConn(config.Dest)
	}
	if err != nil {
		return err
	}
	savedTable := make(echo.Store)
	for _, table := range m.TableNames {
		keys := m.TableKeys[table]
		row := echo.Store{}
		for _, index := range keys {
			mapping := m.Slice[index]
			newData, ok, err := m.get(result, data, mapping, savedTable)
			if err != nil {
				return err
			}
			if !ok {
				continue
			}
			row.Set(mapping.ToField, newData)
		}
		err := dbs.Collection(table).InsertReturning(&row)
		if err != nil {
			return err
		}
		logM := &dbschema.NgingCollectorExportLog{
			PageId:   config.PageId,
			ExportId: config.Id,
			Result:   ``,
			Status:   `idle`,
		}
		if err != nil {
			logM.Result = err.Error()
			if sendErr := noticeSender(`[`+config.Name+`]导入数据库失败: `+logM.Result, 0); sendErr != nil {
				return sendErr
			}
			logM.Status = `failure`
		} else {
			logM.Result = ``
			if sendErr := noticeSender(`[`+config.Name+`]导入数据库成功`, 1); sendErr != nil {
				return sendErr
			}
			logM.Status = `success`
		}
		_, err = logM.Insert()
		if err != nil {
			log.Error(err)
		}
		savedTable.Set(table, row)
	}
	config.UpdateField(nil, `exported`, int(time.Now().Unix()), db.Cond{`id`: config.Id})
	return nil
}

func postAPI(data echo.Store, config *dbschema.NgingCollectorExport, logM *dbschema.NgingCollectorExportLog, noticeSender sender.Notice) error {
	body, err := com.JSONEncode(data)
	if err != nil {
		return err
	}
	worker := miner.NewAPI()
	worker.SetURL(config.Dest)
	worker.SetBinary(body)
	//worker.SetForm(url.Values{`data`: []string{string(body)}})
	body, err = worker.PostJSON()
	if err != nil {
		logM.Result = err.Error()
		if sendErr := noticeSender(`[`+config.Name+`]导入API失败: `+"\n"+logM.Result, 0); sendErr != nil {
			return sendErr
		}
		logM.Status = `failure`
	} else {
		logM.Result = string(body)
		if worker.StatusCode != http.StatusOK {
			logM.Status = `failure`
			if sendErr := noticeSender(`[`+config.Name+`]导入API失败: `+logM.Result, 0); sendErr != nil {
				return sendErr
			}
		} else {
			logM.Status = `success`
			if sendErr := noticeSender(`[`+config.Name+`]导入API成功: `+logM.Result, 1); sendErr != nil {
				return sendErr
			}
		}
	}
	return nil
}

// Export2API 导出到WebAPI
// result - 结果
// data - 本次采集结果
func (m *Mappings) Export2API(result *exec.Recv, data echo.Store, config *dbschema.NgingCollectorExport, noticeSender sender.Notice) error {
	row := make(echo.Store)
	for _, mapping := range m.Slice {
		newData, ok, err := m.get(result, data, mapping)
		if err != nil {
			return err
		}
		if !ok {
			continue
		}
		row.Set(mapping.ToField, newData)
	}
	logM := &dbschema.NgingCollectorExportLog{
		PageId:   config.PageId,
		ExportId: config.Id,
		Result:   ``,
		Status:   `idle`,
	}
	err := postAPI(row, config, logM, noticeSender)
	if err != nil {
		return err
	}
	//panic(echo.Dump(logM, false))
	_, err = logM.Insert()
	config.UpdateField(nil, `exported`, int(time.Now().Unix()), db.Cond{`id`: config.Id})
	return err
}

func (m *Mappings) get(result *exec.Recv, data echo.Store, mapping *Mapping, savedTables ...echo.Store) (interface{}, bool, error) {
	var newData interface{}
	if len(mapping.FromTable) > 0 && len(savedTables) > 0 {
		oldData, ok := savedTables[0].Get(mapping.FromTable).(echo.Store)
		if !ok {
			return newData, ok, nil
		}
		newData = oldData.Get(mapping.FromField)
	} else if mapping.FromParent > 0 {
		parent := result.Parents(mapping.FromParent)
		switch res := parent.Result().(type) {
		case []interface{}:
			i, err := strconv.Atoi(mapping.FromField)
			if err != nil {
				return newData, true, err
			}
			if i < len(res) {
				newData = res[i]
			}
		case map[string]interface{}:
			newData, _ = res[mapping.FromField]
		}
	} else {
		newData = data.Get(mapping.FromField)
	}
	var err error
	if len(mapping.FromPipe) > 0 {
		newData, err = emptyPipeItem.CallFilter(newData, mapping.FromPipe)
	}
	return newData, true, err
}
