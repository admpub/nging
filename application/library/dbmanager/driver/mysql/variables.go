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
package mysql

import (
	"regexp"
	"strings"
)

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
	reFieldTypeBlob      = regexp.MustCompile("blob|bytea|raw|file")
	reFieldLengthInvalid = regexp.MustCompile("[^-0-9,+()[\\]]")
	reFieldLengthNumber  = regexp.MustCompile("^[0-9].*")
	reFieldEnumValue     = regexp.MustCompile(`'((?:[^']|'')*)'`)
	reFieldTextValue     = regexp.MustCompile(`text|lob`)

	pgsqlFieldDefaultValue = regexp.MustCompile("^[a-z]+\\(('[^']*')+\\)$")

	//以下数据来自客户端
	reGrantColumn           = regexp.MustCompile(`^([^() ]+)\s*(\([^)]*\))?$`)
	reGrantOptionValue      = regexp.MustCompile(`(GRANT OPTION)\([^)]*\)`)
	reNotWord               = regexp.MustCompile(`[^a-zA-Z0-9_]+`)
	reOnlyWord              = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	reOnlyNumber            = regexp.MustCompile(`^[0-9]+(\.[0-9]+)?$`)
	reOnlyFloatOrEmpty      = regexp.MustCompile(`^[0-9]*\.[0-9]*$`)
	reChineseAndPunctuation = regexp.MustCompile(`[\x80-\xFF]`)
	reSQLCondOrder          = regexp.MustCompile("^((COUNT\\(DISTINCT |[A-Z0-9_]+\\()(`(?:[^`]|``)+`|\"(?:[^\"]|\"\")+\")\\)|COUNT\\(\\*\\))$")
	reSQLFunction           = regexp.MustCompile("^(COUNT\\((\\*|(DISTINCT )?`(?:[^`]|``)+`)\\)|(AVG|GROUP_CONCAT|MAX|MIN|SUM)\\(`(?:[^`]|``)+`\\))$")

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

	operators     = []string{"=", "<", ">", "<=", ">=", "!=", "LIKE", "LIKE %%", "REGEXP", "IN", "IS NULL", "NOT LIKE", "NOT REGEXP", "NOT IN", "IS NOT NULL", "SQL"} ///< @var array operators used in select
	functions     = []string{"char_length", "date", "from_unixtime", "lower", "round", "sec_to_time", "time_to_sec", "upper"}                                         ///< @var array functions used in select
	grouping      = []string{"avg", "count", "count distinct", "group_concat", "max", "min", "sum"}                                                                   ///< @var array grouping functions used in select
	editFunctions = []map[string]string{                                                                                                                              ///< @var array of array("$type|$type2" => "$function/$function2") functions used in editing, [0] - edit and insert, [1] - edit only
		map[string]string{
			"char":      "md5/sha1/password/encrypt/uuid", //! JavaScript for disabling maxlength
			"binary":    "md5/sha1",
			"date|time": "now",
		}, map[string]string{
			"(^|[^o])int|float|double|decimal": "+/-", // not point
			"date":      "+ interval/- interval",
			"time":      "addtime/subtime",
			"char|text": "concat",
		},
	}

	pattern      = "`(?:[^`]|``)+`"
	reQuotedCol  = regexp.MustCompile(pattern)
	reForeignKey = regexp.MustCompile(`CONSTRAINT (` + pattern + `) FOREIGN KEY ?\(((?:` + pattern + `,? ?)+)\) REFERENCES (` + pattern + `)(?:\.(` + pattern + `))? \(((?:` + pattern + `,? ?)+)\)(?: ON DELETE (` + OnActions + `))?(?: ON UPDATE (` + OnActions + `))?`)

	trans = map[string]string{
		":": ":1",
		"]": ":2",
		"[": ":3",
	}

	reFunctionAddOrSubOr = regexp.MustCompile(`^([+-]|\|\|)$`)
	reFunctionInterval   = regexp.MustCompile(`^[+-] interval$`)
	reSQLValue           = regexp.MustCompile(`^(\d+|'[0-9.: -]') [A-Z_]+$`)
	reFieldName          = regexp.MustCompile(`^([\w(]+)(` + strings.Replace(regexp.QuoteMeta(quoteCol(`_`)), "_", ".*", -1) + `)([ \w)]+)$`)
	reNotSpaceOrDashOrAt = regexp.MustCompile(`[^ -@]`)
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
