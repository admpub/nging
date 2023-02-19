package sqlite

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/admpub/mysql-schema-sync/internal"
)

var _ internal.Comparer = NewCompare()

func NewCompare() *Compare {
	return &Compare{}
}

type Compare struct {
}

func (c *Compare) AlterData(sc *internal.SchemaSync, table string) (*internal.TableAlterData, error) {
	alter := internal.NewAlterData(table)

	sschema, err := sc.SourceDb.GetTableSchema(table)
	if err != nil {
		return nil, err
	}
	dschema, err := sc.DestDb.GetTableSchema(table)
	if err != nil {
		return nil, err
	}

	alter.SchemaDiff = internal.NewSchemaDiff(table, ParseSchema(sschema), ParseSchema(dschema))

	if sschema == dschema {
		return alter, nil
	}
	if len(sschema) == 0 {
		alter.Type = internal.AlterTypeDrop
		alter.SQL = fmt.Sprintf("drop table `%s`;", table)
		return alter, nil
	}
	if len(dschema) == 0 {
		alter.Type = internal.AlterTypeCreate
		if sc.Config.SQLPreprocessor() != nil {
			sschema = sc.Config.SQLPreprocessor()(sschema)
		}
		alter.SQL = sschema
		return alter, nil
	}

	diff := c.getSchemaDiff(sc, alter)
	if len(diff) > 0 {
		alter.Type = internal.AlterTypeAlter
		alter.SQL = diff
	}

	return alter, nil
}

func (c *Compare) getSchemaDiff(sc *internal.SchemaSync, alter *internal.TableAlterData) string {
	source := alter.SchemaDiff.Source
	dest := alter.SchemaDiff.Dest
	table := alter.Table
	var alterIndexLines []string
	var hasFieldChanged bool
	var hasIndexChanged bool
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
		if !has || !isSameSchemaItem(dt, destDt) {
			hasFieldChanged = true
			log.Println("trace check column.alter ", fmt.Sprintf("%s.%s", table, name), "changeField=", name)
			break
		}
		log.Println("trace check column.alter ", fmt.Sprintf("%s.%s", table, name), "not change")
	}

	//多余的字段暂不删除

	//比对索引
	for indexName, idx := range source.IndexAll {
		if sc.Config.IsIgnoreIndex(table, indexName) {
			log.Printf("ignore index %s.%s", table, indexName)
			continue
		}
		dIdx, has := dest.IndexAll[indexName]
		log.Println("trace indexName---->[", fmt.Sprintf("%s.%s", table, indexName), "] dest_has:", has, "\ndest_idx:", dIdx, "\nsource_idx:", idx)
		if has && idx.SQL != dIdx.SQL { // 字段索引不同则先删除，后面会自动创建
			alterIndexLines = append(alterIndexLines, "DROP INDEX `"+indexName+"`")
			log.Println("trace check index.alter ", fmt.Sprintf("%s.%s", table, indexName), "changeIndex=", indexName)
			hasIndexChanged = true
			continue
		}
		log.Println("trace check index.alter ", fmt.Sprintf("%s.%s", table, indexName), "not change")
	}

	//drop index
	if sc.Config.Drop {
		for indexName := range dest.IndexAll {
			if sc.Config.IsIgnoreIndex(table, indexName) {
				log.Printf("ignore index %s.%s", table, indexName)
				continue
			}
			_, has := source.IndexAll[indexName]
			if !has {
				alterIndexLines = append(alterIndexLines, "DROP INDEX `"+indexName+"`")
				log.Println("trace check index.drop ", fmt.Sprintf("%s.%s", table, indexName), "dropIndex=", indexName)
				continue
			}
			log.Println("trace check index.drop ", fmt.Sprintf("%s.%s", table, indexName), " not change")
		}
	}

	//比对外键（外键信息在 CREATE 信息中，因为采用了创建新表删除旧表的方式来更改表结构，所以不用再针对外键进行额外操作）
	for foreignName, idx := range source.ForeignAll {
		if sc.Config.IsIgnoreForeignKey(table, foreignName) {
			log.Printf("ignore foreignName %s.%s", table, foreignName)
			continue
		}
		dIdx, has := dest.ForeignAll[foreignName]
		log.Println("trace foreignName---->[", fmt.Sprintf("%s.%s", table, foreignName), "] dest_has:", has, "\ndest_idx:", dIdx, "\nsource_idx:", idx)
		if !has || idx.SQL != dIdx.SQL {
			hasFieldChanged = true
			break
		}
		log.Println("trace check foreignKey.alter ", fmt.Sprintf("%s.%s", table, foreignName), "not change")
	}

	if !hasFieldChanged && !hasIndexChanged && len(alterIndexLines) == 0 {
		return ``
	}
	tableName := alter.Table
	alters := alterIndexLines
	if hasFieldChanged {
		ddlFieldsDef := strings.TrimSuffix(source.SchemaRaw, `;`)
		// 此分支下的操作包含了对 hasIndexChanged 的处理
		var sameFields []string
		for name := range source.Fields {
			if _, has := dest.Fields[name]; has {
				sameFields = append(sameFields, name)
			}
		}
		if len(sameFields) == 0 {
			if sc.Config.Drop {
				alters = append(alters,
					"DROP TABLE `"+tableName+"`",
					ddlFieldsDef,
				)
			}
		} else {
			tempTable := "_" + tableName + "_old_" + time.Now().Local().Format("20060102_150405")
			newTableFields := "`" + strings.Join(sameFields, "`,`") + "`"
			firstRename := true
			if firstRename {
				alters = append(alters,
					"ALTER TABLE `"+tableName+"` RENAME TO `"+tempTable+"`",
					ddlFieldsDef,
					"INSERT INTO `"+tableName+"` ("+newTableFields+") SELECT "+newTableFields+" FROM `"+tempTable+"`",
					"DROP TABLE `"+tempTable+"`",
				)
			} else {
				alters = append(alters,
					"CREATE TABLE `"+tempTable+"` AS SELECT * FROM `"+tableName+"`",
					"DROP TABLE `"+tableName+"`",
					ddlFieldsDef,
					"INSERT INTO `"+tableName+"` ("+newTableFields+") SELECT "+newTableFields+" FROM `"+tempTable+"`",
				)
			}
		}
	} else if hasIndexChanged {
		alters = append(alters, ParseIndexDDL(source.SchemaRaw)...)
	}
	if len(alters) == 0 {
		return ``
	}
	queryies := []string{
		"SAVEPOINT alter_column_" + tableName,
		"PRAGMA foreign_keys = 0",
		"PRAGMA triggers = NO",
	}
	queryies = append(queryies, alters...)
	/// @todo add views
	queryies = append(queryies, "PRAGMA triggers = YES")
	queryies = append(queryies, "PRAGMA foreign_keys = 1")
	queryies = append(queryies, "RELEASE alter_column_"+tableName)
	return strings.Join(queryies, ";\n") + ";"
}
