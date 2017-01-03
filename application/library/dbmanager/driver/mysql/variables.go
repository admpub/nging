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

import "regexp"

var (
	//数据来自数据库
	reCollate       = regexp.MustCompile(` COLLATE ([^ ]+)`)
	reCharacter     = regexp.MustCompile(` CHARACTER SET ([^ ]+)`)
	reInnoDBComment = regexp.MustCompile(`(?:(.+); )?InnoDB free: .*`)
	reGrantOn       = regexp.MustCompile(`GRANT (.*) ON (.*) TO `)
	reGrantBrackets = regexp.MustCompile(` *([^(,]*[^ ,(])( *\([^)]+\))?`)
	reGrantOption   = regexp.MustCompile(` WITH GRANT OPTION`)
	reGrantIdent    = regexp.MustCompile(` IDENTIFIED BY PASSWORD '([^']+)`)

	reView                = regexp.MustCompile("^.+?\\s+AS\\s+")
	reField               = regexp.MustCompile("^([^( ]+)(?:\\((.+)\\))?( unsigned)?( zerofill)?$")
	reFieldOnUpdate       = regexp.MustCompile("^(?i)on update (.+)")
	reFieldDefault        = regexp.MustCompile("char|set")
	reFieldPrivilegeDelim = regexp.MustCompile(", *")

	reFieldTypeNumber    = regexp.MustCompile("(^|[^o])int|float|double|decimal")
	reFieldTypeText      = regexp.MustCompile("char|text|enum|set")
	reFieldTypeBit       = regexp.MustCompile("^([0-9]+|b'[0-1]+')$")
	reFieldLengthInvalid = regexp.MustCompile("[^-0-9,+()[\\]]")
	reFieldLengthNumber  = regexp.MustCompile("^[0-9].*")

	pgsqlFieldDefaultValue = regexp.MustCompile("^[a-z]+\\(('[^']*')+\\)$")

	//以下数据来自客户端
	reGrantColumn      = regexp.MustCompile(`^([^() ]+)\s*(\([^)]*\))?$`)
	reGrantOptionValue = regexp.MustCompile(`(GRANT OPTION)\([^)]*\)`)
	reNotWord          = regexp.MustCompile(`[^a-zA-Z0-9_]+`)
	reOnlyWord         = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)

	UnsignedTags   = []string{"unsigned", "zerofill", "unsigned zerofill"}
	EnumLength     = "'(?:''|[^'\\\\]|\\\\.)*'"
	OnActions      = "RESTRICT|NO ACTION|CASCADE|SET NULL|SET DEFAULT" ///< @var string used in foreignKeys()
	PartitionTypes = []string{`HASH`, `LINEAR HASH`, `KEY`, `LINEAR KEY`, `RANGE`, `LIST`}

	typeGroups = []*FieldTypeGroup{
		&FieldTypeGroup{
			Types: []*FieldType{
				&FieldType{"tinyint", 3},
				&FieldType{"smallint", 5},
				&FieldType{"mediumint", 8},
				&FieldType{"int", 10},
				&FieldType{"bigint", 20},
				&FieldType{"decimal", 66},
				&FieldType{"float", 12},
				&FieldType{"double", 21},
			},
			Label: "数字",
			Type:  "Number",
		},
		&FieldTypeGroup{
			Types: []*FieldType{
				&FieldType{"date", 10},
				&FieldType{"datetime", 19},
				&FieldType{"timestamp", 19},
				&FieldType{"time", 10},
				&FieldType{"year", 4},
			},
			Label: "日期和时间",
			Type:  "Datetime",
		},
		&FieldTypeGroup{
			Types: []*FieldType{
				&FieldType{"char", 255},
				&FieldType{"varchar", 65535},
				&FieldType{"tinytext", 255},
				&FieldType{"text", 65535},
				&FieldType{"mediumtext", 16777215},
				&FieldType{"longtext", 4294967295},
			},
			Label: "字符串",
			Type:  "String",
		},
		&FieldTypeGroup{
			Types: []*FieldType{
				&FieldType{"enum", 65535},
				&FieldType{"set", 64},
			},
			Label: "列表",
			Type:  "List",
		},
		&FieldTypeGroup{
			Types: []*FieldType{
				&FieldType{"bit", 20},
				&FieldType{"binary", 255},
				&FieldType{"varbinary", 65535},
				&FieldType{"tinyblob", 255},
				&FieldType{"blob", 65535},
				&FieldType{"mediumblob", 16777215},
				&FieldType{"longblob", 4294967295},
			},
			Label: "二进制",
			Type:  "Binary",
		},
		&FieldTypeGroup{
			Types: []*FieldType{
				&FieldType{"geometry", 0},
				&FieldType{"point", 0},
				&FieldType{"linestring", 0},
				&FieldType{"polygon", 0},
				&FieldType{"multipoint", 0},
				&FieldType{"multilinestring", 0},
				&FieldType{"multipolygon", 0},
				&FieldType{"geometrycollection", 0},
			},
			Label: "几何图形",
			Type:  "Geometry",
		},
	}
	types = map[string]uint64{} ///< @var array ($type,$maximum_unsigned_length, ...)

	pattern      = "`(?:[^`]|``)+`"
	reQuotedCol  = regexp.MustCompile(pattern)
	reForeignKey = regexp.MustCompile(`CONSTRAINT (` + pattern + `) FOREIGN KEY ?\(((?:` + pattern + `,? ?)+)\) REFERENCES (` + pattern + `)(?:\.(` + pattern + `))? \(((?:` + pattern + `,? ?)+)\)(?: ON DELETE (` + OnActions + `))?(?: ON UPDATE (` + OnActions + `))?`)
)

func init() {
	for _, group := range typeGroups {
		for _, typ := range group.Types {
			types[typ.Name] = typ.MaxUnsignedLength
		}
	}
}

type FieldType struct {
	Name              string
	MaxUnsignedLength uint64
}
type FieldTypeGroup struct {
	Types []*FieldType
	Label string
	Type  string
}

func (f *FieldTypeGroup) IsString(typeName string) bool {
	return reFieldTypeText.MatchString(typeName)
}

func (f *FieldTypeGroup) IsNumeric(typeName string) bool {
	return reFieldTypeNumber.MatchString(typeName)
}
