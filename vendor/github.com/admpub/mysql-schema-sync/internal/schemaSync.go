package internal

import (
	"fmt"
	"log"
	"strings"
)

// SchemaSync 配置文件
type SchemaSync struct {
	Config   *Config
	SourceDb DBOperator
	DestDb   DBOperator
	Comparer Comparer
}

// NewSchemaSync 对一个配置进行同步
// dbOperators: 0: source; 1: destination
func NewSchemaSync(config *Config, dbOperators ...DBOperator) *SchemaSync {
	s := new(SchemaSync)
	s.Config = config
	s.Comparer = NewMyCompare()
	switch len(dbOperators) {
	case 2:
		s.SourceDb = dbOperators[0]
		s.DestDb = dbOperators[1]
	case 1:
		s.SourceDb = dbOperators[0]
	}
	if s.SourceDb == nil {
		s.SourceDb = NewMyDb(config.SourceDSN, "source")
	}
	if s.DestDb == nil {
		s.DestDb = NewMyDb(config.DestDSN, "dest")
	}
	return s
}

func (sc *SchemaSync) SetComparer(comparer Comparer) *SchemaSync {
	sc.Comparer = comparer
	return sc
}

// GetNewTableNames 获取所有新增加的表名
func (sc *SchemaSync) GetNewTableNames() []string {
	sourceTables := sc.SourceDb.GetTableNames()
	destTables := sc.DestDb.GetTableNames()

	var newTables []string

	for _, name := range sourceTables {
		if !inStringSlice(name, destTables) {
			newTables = append(newTables, name)
		}
	}
	return newTables
}

func (sc *SchemaSync) getAlterDataByTable(table string) *TableAlterData {
	return sc.Comparer.AlterData(sc, table)
}

// SyncSQL4Dest sync schema change
func (sc *SchemaSync) SyncSQL4Dest(sqlStr string, sqls []string) error {
	log.Println("Exec_SQL_START:")
	log.Println(">>>>>>")
	log.Println(sqlStr)
	log.Println("<<<<<<")
	sqlStr = strings.TrimSpace(sqlStr)
	if sqlStr == "" {
		log.Println("sql_is_empty,skip")
		return nil
	}
	t := NewMyTimer()
	ret, err := sc.DestDb.Exec(sqlStr)

	//how to enable allowMultiQueries?
	if err != nil && len(sqls) > 1 {
		log.Println("exec_mut_query failed,err=", err, ",now exec sqls foreach")
		tx, errTx := sc.DestDb.Begin()
		if errTx == nil {
			for _, sql := range sqls {
				ret, err = tx.Exec(sql)
				log.Println("query_one:[", sql, "]", err)
				if err != nil {
					break
				}
			}
			if err == nil {
				err = tx.Commit()
			} else {
				tx.Rollback()
			}
		} else {
			err = errTx
		}
	}
	t.Stop()
	if err != nil {
		log.Println("EXEC_SQL_FAIELD", err)
		return err
	}
	log.Println("EXEC_SQL_SUCCESS,used:", t.UsedSecond())
	var affected int64
	if ret != nil {
		affected, err = ret.RowsAffected()
	}
	log.Println("EXEC_SQL_RET:", affected, err)
	return err
}

func (sc *SchemaSync) Close() error {
	if sc.DestDb != nil {
		if err := sc.DestDb.Close(); err != nil {
			return err
		}
	}
	if sc.SourceDb != nil {
		if err := sc.SourceDb.Close(); err != nil {
			return err
		}
	}
	return nil
}

// CheckSchemaDiff 执行最终的diff
func CheckSchemaDiff(cfg *Config, dbOperators ...DBOperator) *Statics {
	statics := newStatics(cfg)
	sc := NewSchemaSync(cfg, dbOperators...)
	if cfg.comparer != nil {
		sc.SetComparer(cfg.comparer)
	}

	defer (func() {
		statics.timer.Stop()
		statics.sendMailNotice()
		sc.Close()
	})()

	newTables := sc.SourceDb.GetTableNames()
	log.Println("source db table total:", len(newTables))

	changedTables := make(map[string][]*TableAlterData)

	for index, table := range newTables {
		log.Printf("Index : %d Table : %s\n", index, table)
		if !cfg.CheckMatchTables(table) || cfg.IsSkipTables(table) {
			log.Println("Table:", table, "skip")
			continue
		}

		sd := sc.getAlterDataByTable(table)

		if sd.Type != AlterTypeNo {
			fmt.Println(sd)
			fmt.Println("")
			relationTables := sd.SchemaDiff.RelationTables()
			//fmt.Println("relationTables:",table,relationTables)

			//将所有有外键关联的单独放
			groupKey := "multi"
			if len(relationTables) == 0 {
				groupKey = "single_" + table
			}
			if _, has := changedTables[groupKey]; !has {
				changedTables[groupKey] = make([]*TableAlterData, 0)
			}
			changedTables[groupKey] = append(changedTables[groupKey], sd)
		} else {
			log.Println("table:", table, "not change,", sd)
		}
	}

	log.Println("trace changedTables:", changedTables)

	countSuccess := 0
	countFailed := 0
	canRunTypePref := "single"
	//先执行单个表的
run_sync:
	for typeName, sds := range changedTables {
		if !strings.HasPrefix(typeName, canRunTypePref) {
			continue
		}
		log.Println("runSyncType:", typeName)
		var sqls []string
		var sts []*tableStatics
		for _, sd := range sds {
			sql := strings.TrimRight(sd.SQL, ";")
			sqls = append(sqls, sql)

			st := statics.newTableStatics(sd.Table, sd)
			sts = append(sts, st)
		}

		sql := strings.Join(sqls, ";\n")
		if !strings.HasSuffix(sql, `;`) {
			sql += `;`
		}
		var ret error

		if sc.Config.Sync {

			ret = sc.SyncSQL4Dest(sql, sqls)
			if ret == nil {
				countSuccess++
			} else {
				countFailed++
			}
		}
		for _, st := range sts {
			st.alterRet = ret
			st.schemaAfter = sc.DestDb.GetTableSchema(st.table)
			st.timer.Stop()
		}

	} //end for

	//最后再执行多个表的alter
	if canRunTypePref == "single" {
		canRunTypePref = "multi"
		goto run_sync
	}

	if sc.Config.Sync {
		log.Println("execute_all_sql_done,success_total:", countSuccess, "failed_total:", countFailed)
	}

	return statics
}
