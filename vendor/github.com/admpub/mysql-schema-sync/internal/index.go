package internal

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"
)

// DbIndex db index
type DbIndex struct {
	IndexType      indexType
	Name           string
	SQL            string
	RelationTables []string //相关联的表
}

type indexType string

const (
	indexTypePrimary    indexType = "PRIMARY"
	indexTypeIndex                = "INDEX"
	indexTypeForeignKey           = "FOREIGN KEY"
)

func (idx *DbIndex) alterAddSQL(drop bool) string {
	alterSQL := []string{}
	if drop {
		dropSQL := idx.alterDropSQL()
		if dropSQL != "" {
			alterSQL = append(alterSQL, dropSQL)
		}
	}

	switch idx.IndexType {
	case indexTypePrimary:
		alterSQL = append(alterSQL, "ADD "+idx.SQL)
	case indexTypeIndex, indexTypeForeignKey:
		alterSQL = append(alterSQL, fmt.Sprintf("ADD %s", idx.SQL))
	default:
		log.Fatalln("unknow indexType", idx.IndexType)
	}
	return strings.Join(alterSQL, ",\n")
}

func (idx *DbIndex) String() string {
	bs, _ := json.MarshalIndent(idx, "  ", " ")
	return string(bs)
}

func (idx *DbIndex) alterDropSQL() string {
	switch idx.IndexType {
	case indexTypePrimary:
		return "DROP PRIMARY KEY"
	case indexTypeIndex:
		return fmt.Sprintf("DROP INDEX `%s`", idx.Name)
	case indexTypeForeignKey:
		return fmt.Sprintf("DROP FOREIGN KEY `%s`", idx.Name)
	default:
		log.Fatalln("unknow indexType", idx.IndexType)
	}
	return ""
}

func (idx *DbIndex) addRelationTable(table string) {
	table = strings.TrimSpace(table)
	if table != "" {
		idx.RelationTables = append(idx.RelationTables, table)
	}
}

//匹配索引字段
var indexReg = regexp.MustCompile(`^([A-Z]+\s)?KEY\s`)

//匹配外键
var foreignKeyReg = regexp.MustCompile("^CONSTRAINT `(.+)` FOREIGN KEY.+ REFERENCES `(.+)` ")

func parseDbIndexLine(line string) *DbIndex {
	line = strings.TrimSpace(line)
	idx := &DbIndex{
		SQL:            line,
		RelationTables: []string{},
	}
	if strings.HasPrefix(line, "PRIMARY") {
		idx.IndexType = indexTypePrimary
		idx.Name = "PRIMARY KEY"
		return idx
	}

	//  UNIQUE KEY `idx_a` (`a`) USING HASH COMMENT '注释',
	//  FULLTEXT KEY `c` (`c`)
	//  PRIMARY KEY (`d`)
	//  KEY `idx_e` (`e`),
	if indexReg.MatchString(line) {
		arr := strings.Split(line, "`")
		idx.IndexType = indexTypeIndex
		idx.Name = arr[1]
		return idx
	}

	//CONSTRAINT `busi_table_ibfk_1` FOREIGN KEY (`repo_id`) REFERENCES `repo_table` (`repo_id`)
	foreignMatches := foreignKeyReg.FindStringSubmatch(line)
	if len(foreignMatches) > 0 {
		idx.IndexType = indexTypeForeignKey
		idx.Name = foreignMatches[1]
		idx.addRelationTable(foreignMatches[2])
		return idx
	}

	log.Fatalln("db_index parse failed,unsupport,line:", line)
	return nil
}
