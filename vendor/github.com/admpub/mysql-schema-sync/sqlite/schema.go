package sqlite

import (
	"regexp"
	"strings"

	"github.com/admpub/mysql-schema-sync/internal"
)

var (
	indexReg      = regexp.MustCompile("^CREATE (?:UNIQUE )?INDEX [`\"]([^`\"]+)[`\"]")
	index2Reg     = regexp.MustCompile("^CONSTRAINT [`\"]([^`\"]+)[`\"] UNIQUE")
	foreignKeyReg = regexp.MustCompile("^CONSTRAINT [`\"]([^`\"]+)[`\"] FOREIGN KEY.+ REFERENCES [`\"]([^`\"]+)[`\"] ")
)

func parseDbIndexLine(line string) *internal.DbIndex {
	line = strings.TrimSpace(line)
	idx := &internal.DbIndex{
		SQL:            line,
		RelationTables: []string{},
	}
	//CREATE UNIQUE INDEX `UNQ_forever_process3_name` ON `forever_process3`(`name`,`name2`)
	//CREATE INDEX `forever_process3_en` ON `forever_process3`(`en`)
	indexMatches := indexReg.FindStringSubmatch(line)
	if len(indexMatches) > 0 {
		idx.IndexType = internal.IndexTypeIndex
		idx.Name = indexMatches[1]
		return idx
	}
	indexMatches = index2Reg.FindStringSubmatch(line)
	if len(indexMatches) > 0 {
		idx.IndexType = internal.IndexTypeIndex
		idx.Name = indexMatches[1]
		return idx
	}

	//PRIMARY KEY ("id"),
	//CONSTRAINT "wid" FOREIGN KEY ("wid") REFERENCES "forever_process" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION,
	//CONSTRAINT "idx" UNIQUE ("idx") ON CONFLICT FAIL
	foreignMatches := foreignKeyReg.FindStringSubmatch(line)
	if len(foreignMatches) > 0 {
		idx.IndexType = internal.IndexTypeForeignKey
		idx.Name = foreignMatches[1]
		idx.AddRelationTable(foreignMatches[2])
		return idx
	}
	return nil
}

// ParseSchema parse table's schema
func ParseSchema(schema string) *internal.MySchema {
	schema = strings.TrimSpace(schema)
	lines := strings.Split(schema, "\n")
	mys := internal.NewSchema(schema)
	var hasPrimaryKey bool
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		line = strings.TrimRight(line, ",")
		line = strings.TrimRight(line, " ")
		var quote string
		switch line[0] {
		case '`':
			quote = "`"
		case '"':
			quote = `"`
		}
		//println(`[`, line, `]`)
		var idx *internal.DbIndex
		var isPK bool
		if len(quote) > 0 {
			index := strings.Index(line[1:], quote)
			name := line[1 : index+1]
			mys.Fields[name] = line
			if !hasPrimaryKey && strings.HasSuffix(line, ` PRIMARY KEY`) {
				isPK = true
			}
		} else if !hasPrimaryKey && strings.HasPrefix(line, `PRIMARY KEY `) {
			isPK = true
		}
		if isPK {
			idx = &internal.DbIndex{
				SQL:            line,
				RelationTables: []string{},
				IndexType:      internal.IndexTypePrimary,
				Name:           "PRIMARY KEY",
			}
			hasPrimaryKey = true
		} else {
			idx = parseDbIndexLine(line)
		}
		if idx == nil {
			continue
		}
		switch idx.IndexType {
		case internal.IndexTypeForeignKey:
			mys.ForeignAll[idx.Name] = idx
		default:
			mys.IndexAll[idx.Name] = idx
		}
	}
	return mys

}

var backets = regexp.MustCompile(`\([^()]+(,[^()]+)+\)`)

func FormatSchema(schema string) string {
	tempCleaned := strings.TrimRight(schema, "\r\n")
	if !strings.Contains(tempCleaned, "\n") {
		p := strings.Index(schema, `(`)
		var replaced bool
		if p > 0 {
			schema2 := backets.ReplaceAllStringFunc(schema[p+1:], func(s string) string {
				replaced = true
				return strings.ReplaceAll(s, `,`, `#comma#`)
			})
			schema = schema[:p] + "(\n" + schema2
		}
		p = strings.LastIndex(schema, `)`)
		if p > 0 {
			schema = schema[:p] + "\n" + schema[p:]
		}
		lines := strings.Split(schema, ",")
		schema = strings.Join(lines, ",\n")
		if replaced {
			schema = strings.ReplaceAll(schema, `#comma#`, `,`)
		}
	}
	return schema
}
