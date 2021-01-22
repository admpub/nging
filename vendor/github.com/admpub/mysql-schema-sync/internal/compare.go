package internal

import (
	"fmt"
	"log"
	"strings"
)

var _ Comparer = NewMyCompare()

func NewMyCompare() *MyCompare {
	return &MyCompare{}
}

type MyCompare struct {
}

func (my *MyCompare) AlterData(sc *SchemaSync, table string) *TableAlterData {
	alter := NewAlterData(table)
	sschema := sc.SourceDb.GetTableSchema(table)
	dschema := sc.DestDb.GetTableSchema(table)

	alter.SchemaDiff = newSchemaDiff(table, sschema, dschema)

	if sschema == dschema {
		return alter
	}
	if len(sschema) == 0 {
		alter.Type = AlterTypeDrop
		alter.SQL = fmt.Sprintf("drop table `%s`;", table)
		return alter
	}
	if len(dschema) == 0 {
		alter.Type = AlterTypeCreate

		if sc.Config.SQLPreprocessor() != nil {
			sschema = sc.Config.SQLPreprocessor()(sschema)
		}

		alter.SQL = sschema + ";"
		return alter
	}

	diff := my.getSchemaDiff(sc, alter)
	if len(diff) > 0 {
		alter.Type = AlterTypeAlter
		alter.SQL = fmt.Sprintf("ALTER TABLE `%s`\n%s;", table, diff)
	}

	return alter
}

func (my *MyCompare) getSchemaDiff(sc *SchemaSync, alter *TableAlterData) string {
	sourceMyS := alter.SchemaDiff.Source
	destMyS := alter.SchemaDiff.Dest
	table := alter.Table

	var alterLines []string
	//比对字段
	for name, dt := range sourceMyS.Fields {
		if sc.Config.IsIgnoreField(table, name) {
			log.Printf("ignore column %s.%s", table, name)
			continue
		}
		if sc.Config.SQLPreprocessor() != nil {
			dt = sc.Config.SQLPreprocessor()(dt)
		}
		var alterSQL string
		if destDt, has := destMyS.Fields[name]; has {
			if !isSameSchemaItem(dt, destDt) {
				alterSQL = fmt.Sprintf("CHANGE `%s` %s", name, dt)
			}
		} else {
			alterSQL = "ADD " + dt
		}
		if alterSQL != "" {
			log.Println("trace check column.alter ", fmt.Sprintf("%s.%s", table, name), "alterSQL=", alterSQL)
			alterLines = append(alterLines, alterSQL)
		} else {
			log.Println("trace check column.alter ", fmt.Sprintf("%s.%s", table, name), "not change")
		}
	}

	//源库已经删除的字段
	if sc.Config.Drop {
		for name := range destMyS.Fields {
			if sc.Config.IsIgnoreField(table, name) {
				log.Printf("ignore column %s.%s", table, name)
				continue
			}
			if _, has := sourceMyS.Fields[name]; !has {
				alterSQL := fmt.Sprintf("drop `%s`", name)
				alterLines = append(alterLines, alterSQL)
				log.Println("trace check column.drop ", fmt.Sprintf("%s.%s", table, name), "alterSQL=", alterSQL)
			} else {
				log.Println("trace check column.drop ", fmt.Sprintf("%s.%s", table, name), "not change")
			}
		}
	}

	//多余的字段暂不删除

	//比对索引
	for indexName, idx := range sourceMyS.IndexAll {
		if sc.Config.IsIgnoreIndex(table, indexName) {
			log.Printf("ignore index %s.%s", table, indexName)
			continue
		}
		dIdx, has := destMyS.IndexAll[indexName]
		log.Println("trace indexName---->[", fmt.Sprintf("%s.%s", table, indexName), "] dest_has:", has, "\ndest_idx:", dIdx, "\nsource_idx:", idx)
		alterSQL := ""
		if has {
			if idx.SQL != dIdx.SQL {
				alterSQL = idx.AlterAddSQL(true)
			}
		} else {
			alterSQL = idx.AlterAddSQL(false)
		}
		if alterSQL != "" {
			alterLines = append(alterLines, alterSQL)
			log.Println("trace check index.alter ", fmt.Sprintf("%s.%s", table, indexName), "alterSQL=", alterSQL)
		} else {
			log.Println("trace check index.alter ", fmt.Sprintf("%s.%s", table, indexName), "not change")
		}
	}

	//drop index
	if sc.Config.Drop {
		for indexName, dIdx := range destMyS.IndexAll {
			if sc.Config.IsIgnoreIndex(table, indexName) {
				log.Printf("ignore index %s.%s", table, indexName)
				continue
			}
			var dropSQL string
			if _, has := sourceMyS.IndexAll[indexName]; !has {
				dropSQL = dIdx.AlterDropSQL()
			}

			if dropSQL != "" {
				alterLines = append(alterLines, dropSQL)
				log.Println("trace check index.drop ", fmt.Sprintf("%s.%s", table, indexName), "alterSQL=", dropSQL)
			} else {
				log.Println("trace check index.drop ", fmt.Sprintf("%s.%s", table, indexName), " not change")
			}
		}
	}

	//比对外键
	for foreignName, idx := range sourceMyS.ForeignAll {
		if sc.Config.IsIgnoreForeignKey(table, foreignName) {
			log.Printf("ignore foreignName %s.%s", table, foreignName)
			continue
		}
		dIdx, has := destMyS.ForeignAll[foreignName]
		log.Println("trace foreignName---->[", fmt.Sprintf("%s.%s", table, foreignName), "] dest_has:", has, "\ndest_idx:", dIdx, "\nsource_idx:", idx)
		alterSQL := ""
		if has {
			if idx.SQL != dIdx.SQL {
				alterSQL = idx.AlterAddSQL(true)
			}
		} else {
			alterSQL = idx.AlterAddSQL(false)
		}
		if alterSQL != "" {
			alterLines = append(alterLines, alterSQL)
			log.Println("trace check foreignKey.alter ", fmt.Sprintf("%s.%s", table, foreignName), "alterSQL=", alterSQL)
		} else {
			log.Println("trace check foreignKey.alter ", fmt.Sprintf("%s.%s", table, foreignName), "not change")
		}
	}

	//drop 外键
	if sc.Config.Drop {
		for foreignName, dIdx := range destMyS.ForeignAll {
			if sc.Config.IsIgnoreForeignKey(table, foreignName) {
				log.Printf("ignore foreignName %s.%s", table, foreignName)
				continue
			}
			var dropSQL string
			if _, has := sourceMyS.ForeignAll[foreignName]; !has {
				log.Println("trace foreignName --->[", fmt.Sprintf("%s.%s", table, foreignName), "]", "didx:", dIdx)
				dropSQL = dIdx.AlterDropSQL()

			}
			if dropSQL != "" {
				alterLines = append(alterLines, dropSQL)
				log.Println("trace check foreignKey.drop ", fmt.Sprintf("%s.%s", table, foreignName), "alterSQL=", dropSQL)
			} else {
				log.Println("trace check foreignKey.drop ", fmt.Sprintf("%s.%s", table, foreignName), "not change")
			}
		}
	}

	return strings.Join(alterLines, ",\n")
}
