package sqlite

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/admpub/mysql-schema-sync/internal"
	"github.com/webx-top/com"
)

var _ internal.Comparer = NewCompare()

func NewCompare() *Compare {
	return &Compare{}
}

type Compare struct {
}

func (c *Compare) AlterData(sc *internal.SchemaSync, table string) *internal.TableAlterData {
	alter := internal.NewAlterData(table)

	sschema := sc.SourceDb.GetTableSchema(table)
	dschema := sc.DestDb.GetTableSchema(table)

	alter.SchemaDiff = internal.NewSchemaDiff(table, ParseSchema(sschema), ParseSchema(dschema))

	if sschema == dschema {
		return alter
	}
	if len(sschema) == 0 {
		alter.Type = internal.AlterTypeDrop
		alter.SQL = fmt.Sprintf("drop table `%s`;", table)
		return alter
	}
	if len(dschema) == 0 {
		alter.Type = internal.AlterTypeCreate
		if sc.Config.SQLPreprocessor() != nil {
			sschema = sc.Config.SQLPreprocessor()(sschema)
		}
		alter.SQL = sschema
		return alter
	}

	diff := c.getSchemaDiff(sc, alter)
	if len(diff) > 0 {
		alter.Type = internal.AlterTypeAlter
		alter.SQL = diff
	}

	return alter
}

func (c *Compare) getSchemaDiff(sc *internal.SchemaSync, alter *internal.TableAlterData) string {
	source := alter.SchemaDiff.Source
	dest := alter.SchemaDiff.Dest
	table := alter.Table
	var hasChanges bool
	//比对字段
	for name, dt := range source.Fields {
		if sc.Config.IsIgnoreField(table, name) {
			log.Printf("ignore column %s.%s", table, name)
			continue
		}
		if sc.Config.SQLPreprocessor() != nil {
			dt = sc.Config.SQLPreprocessor()(dt)
		}
		destDt, has := dest.Fields[name]
		if has {
			if !isSameSchemaItem(dt, destDt) {
				com.Dump(map[string]interface{}{`src.`: dt, `dest`: destDt})
				hasChanges = true
			}
		} else {
			hasChanges = true
		}
		if hasChanges {
			log.Println("trace check column.alter ", fmt.Sprintf("%s.%s", table, name), "changeField=", name)
			break
		}
		log.Println("trace check column.alter ", fmt.Sprintf("%s.%s", table, name), "not change")
	}

	//源库已经删除的字段
	if !hasChanges && sc.Config.Drop {
		for name := range dest.Fields {
			if sc.Config.IsIgnoreField(table, name) {
				log.Printf("ignore column %s.%s", table, name)
				continue
			}
			if _, has := source.Fields[name]; !has {
				log.Println("trace check column.drop ", fmt.Sprintf("%s.%s", table, name), "dropField=", name)
				hasChanges = true
				break
			}
			log.Println("trace check column.drop ", fmt.Sprintf("%s.%s", table, name), "not change")
		}
	}

	//多余的字段暂不删除

	if !hasChanges {
		//比对索引
		for indexName, idx := range source.IndexAll {
			if sc.Config.IsIgnoreIndex(table, indexName) {
				log.Printf("ignore index %s.%s", table, indexName)
				continue
			}
			dIdx, has := dest.IndexAll[indexName]
			log.Println("trace indexName---->[", fmt.Sprintf("%s.%s", table, indexName), "] dest_has:", has, "\ndest_idx:", dIdx, "\nsource_idx:", idx)
			if has {
				if idx.SQL != dIdx.SQL {
					com.Dump(map[string]interface{}{`src.`: idx.SQL, `dest`: dIdx.SQL})
					hasChanges = true
				}
			} else {
				hasChanges = true
			}
			if hasChanges {
				log.Println("trace check index.alter ", fmt.Sprintf("%s.%s", table, indexName), "changeIndex=", indexName)
				break
			}
			log.Println("trace check index.alter ", fmt.Sprintf("%s.%s", table, indexName), "not change")
		}
	}

	//drop index
	if !hasChanges && sc.Config.Drop {
		for indexName := range dest.IndexAll {
			if sc.Config.IsIgnoreIndex(table, indexName) {
				log.Printf("ignore index %s.%s", table, indexName)
				continue
			}
			if _, has := source.IndexAll[indexName]; !has {
				hasChanges = true
			}
			if hasChanges {
				log.Println("trace check index.drop ", fmt.Sprintf("%s.%s", table, indexName), "dropIndex=", indexName)
				break
			}
			log.Println("trace check index.drop ", fmt.Sprintf("%s.%s", table, indexName), " not change")
		}
	}

	if !hasChanges {
		//比对外键
		for foreignName, idx := range source.ForeignAll {
			if sc.Config.IsIgnoreForeignKey(table, foreignName) {
				log.Printf("ignore foreignName %s.%s", table, foreignName)
				continue
			}
			dIdx, has := dest.ForeignAll[foreignName]
			log.Println("trace foreignName---->[", fmt.Sprintf("%s.%s", table, foreignName), "] dest_has:", has, "\ndest_idx:", dIdx, "\nsource_idx:", idx)
			if has {
				if idx.SQL != dIdx.SQL {
					com.Dump(map[string]interface{}{`src.`: idx.SQL, `dest`: dIdx.SQL})
					hasChanges = true
				}
			} else {
				hasChanges = true
			}
			if hasChanges {
				log.Println("trace check foreignKey.alter ", fmt.Sprintf("%s.%s", table, foreignName), "changeForeign=", foreignName)
				break
			}
			log.Println("trace check foreignKey.alter ", fmt.Sprintf("%s.%s", table, foreignName), "not change")
		}
	}

	//drop 外键
	if !hasChanges && sc.Config.Drop {
		for foreignName, dIdx := range dest.ForeignAll {
			if sc.Config.IsIgnoreForeignKey(table, foreignName) {
				log.Printf("ignore foreignName %s.%s", table, foreignName)
				continue
			}
			if _, has := source.ForeignAll[foreignName]; !has {
				log.Println("trace foreignName --->[", fmt.Sprintf("%s.%s", table, foreignName), "]", "didx:", dIdx)
				hasChanges = true
			}
			if hasChanges {
				log.Println("trace check foreignKey.drop ", fmt.Sprintf("%s.%s", table, foreignName), "dropForeign=", foreignName)
			} else {
				log.Println("trace check foreignKey.drop ", fmt.Sprintf("%s.%s", table, foreignName), "not change")
			}
		}
	}

	if !hasChanges {
		return ``
	}
	var sameFields []string
	for name := range source.Fields {
		if _, has := dest.Fields[name]; has {
			sameFields = append(sameFields, name)
		}
	}
	if len(sameFields) == 0 {
		com.Dump(sameFields)
		panic(table)
	}

	tableName := alter.Table
	ddlFieldsDef := strings.TrimSuffix(source.SchemaRaw, `;`)
	tempTable := "_" + tableName + "_old_" + time.Now().Local().Format("20060102_150405")
	newTableFields := "`" + strings.Join(sameFields, "`,`") + "`"
	queryies := []string{
		"SAVEPOINT alter_column_" + tableName,
		"PRAGMA foreign_keys = 0",
		"PRAGMA triggers = NO",
	}
	firstRename := true
	var alters []string
	if firstRename {
		alters = []string{
			"ALTER TABLE `" + tableName + "` RENAME TO `" + tempTable + "`",
			ddlFieldsDef,
			"INSERT INTO `" + tableName + "` (" + newTableFields + ") SELECT " + newTableFields + " FROM `" + tempTable + "`",
			"DROP TABLE `" + tempTable + "`",
		}
	} else {
		alters = []string{
			"CREATE TABLE `" + tempTable + "` AS SELECT * FROM `" + tableName + "`",
			"DROP TABLE `" + tableName + "`",
			ddlFieldsDef,
			"INSERT INTO `" + tableName + "` (" + newTableFields + ") SELECT " + newTableFields + " FROM `" + tempTable + "`",
		}
	}
	queryies = append(queryies, alters...)
	/// @todo add views
	queryies = append(queryies, "PRAGMA triggers = YES")
	queryies = append(queryies, "PRAGMA foreign_keys = 1")
	queryies = append(queryies, "RELEASE alter_column_"+tableName)
	return strings.Join(queryies, ";\n") + ";"
}
