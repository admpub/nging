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

	"github.com/nging-plugins/dbmanager/pkg/library/dbmanager/driver/mysql/utils"
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
	reFulltextKey   = regexp.MustCompile("FULLTEXT KEY `([^`]+)`[ ]*\\([^)]+\\) /\\*[^ ]* WITH ([^*]+) \\*/")

	reView                = regexp.MustCompile("^.+?\\s+AS\\s+")
	reField               = regexp.MustCompile("^([^( ]+)(?:\\((.+)\\))?( unsigned)?( zerofill)?$")
	reFieldOnUpdate       = regexp.MustCompile("^(?i)(?:DEFAULT_GENERATED )?on update (.+)")
	reFieldDefault        = regexp.MustCompile("char|set")
	reFieldPrivilegeDelim = regexp.MustCompile(", *")

	reFriendlyName = regexp.MustCompile("(?i)[^a-z0-9_]")
	//((?<!o)int(?!er)|numeric|real|float|double|decimal|money)
	reFieldTypeNumber    = regexp.MustCompile("(^|[^o])int(?:er)?|numeric|real|float|double|decimal|money")
	reCSVText            = regexp.MustCompile("[\"\n,;\t]")
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
	reSQLFunction           = regexp.MustCompile("^(COUNT\\((\\*|(DISTINCT )?`(?:[^`]|``)+`)\\)|(AVG|GROUP_CONCAT|MAX|MIN|SUM|CHAR_LENGTH|DATE|FROM_UNIXTIME|LOWER|ROUND|SEC_TO_TIME|TIME_TO_SEC|UPPER)\\(`(?:[^`]|``)+`\\))$")

	//Charsets MySQL 支持的字符集
	Charsets = utils.Charsets
	// UnsignedTags 无符号标签
	UnsignedTags = []string{"unsigned", "zerofill", "unsigned zerofill"}
	// EnumLength 枚举选项
	EnumLength           = "'(?:''|[^'\\\\]|\\\\.)*'"
	reEnumLength         = regexp.MustCompile(EnumLength)
	reContainsEnumLength = regexp.MustCompile("^\\s*\\(?\\s*" + EnumLength + "(?:\\s*,\\s*" + EnumLength + ")*\\s*\\)?\\s*$")
	// OnActions used in foreignKeys()
	OnActions = "RESTRICT|NO ACTION|CASCADE|SET NULL|SET DEFAULT"
	// PartitionTypes 分区类型
	PartitionTypes = []string{`HASH`, `LINEAR HASH`, `KEY`, `LINEAR KEY`, `RANGE`, `LIST`}

	typeGroups = []*FieldTypeGroup{
		&FieldTypeGroup{
			Types: []*FieldType{
				{"tinyint", 3},
				{"smallint", 5},
				{"mediumint", 8},
				{"int", 10},
				{"bigint", 20},
				{"decimal", 66},
				{"float", 12},
				{"double", 21},
			},
			Label: "数字",
			Type:  "Number",
		},
		&FieldTypeGroup{
			Types: []*FieldType{
				{"date", 10},
				{"datetime", 19},
				{"timestamp", 19},
				{"time", 10},
				{"year", 4},
			},
			Label: "日期和时间",
			Type:  "Datetime",
		},
		&FieldTypeGroup{
			Types: []*FieldType{
				{"char", 255},
				{"varchar", 65535},
				{"tinytext", 255},
				{"text", 65535},
				{"mediumtext", 16777215},
				{"longtext", 4294967295},
			},
			Label: "字符串",
			Type:  "String",
		},
		&FieldTypeGroup{
			Types: []*FieldType{
				{"enum", 65535},
				{"set", 64},
			},
			Label: "列表",
			Type:  "List",
		},
		&FieldTypeGroup{
			Types: []*FieldType{
				{"bit", 20},
				{"binary", 255},
				{"varbinary", 65535},
				{"tinyblob", 255},
				{"blob", 65535},
				{"mediumblob", 16777215},
				{"longblob", 4294967295},
			},
			Label: "二进制",
			Type:  "Binary",
		},
		&FieldTypeGroup{
			Types: []*FieldType{
				{"geometry", 0},
				{"point", 0},
				{"linestring", 0},
				{"polygon", 0},
				{"multipoint", 0},
				{"multilinestring", 0},
				{"multipolygon", 0},
				{"geometrycollection", 0},
			},
			Label: "几何图形",
			Type:  "Geometry",
		},
	}
	types = map[string]uint64{} ///< @var array ($type,$maximum_unsigned_length, ...)

	operators     = []string{"=", "<", ">", "<=", ">=", "!=", "LIKE", "LIKE %%", "REGEXP", "IN", "IS NULL", "NOT LIKE", "NOT REGEXP", "NOT IN", "IS NOT NULL", "SQL"} ///< @var array operators used in select
	functions     = []string{"char_length", "date", "from_unixtime", "lower", "round", "sec_to_time", "time_to_sec", "upper"}                                         ///< @var array functions used in select
	grouping      = []string{"avg", "count", "count distinct", "group_concat", "max", "min", "sum"}                                                                   ///< @var array grouping functions used in select
	editFunctions = []map[string]string{
		///< @var array of array("$type|$type2" => "$function/$function2") functions used in editing, [0] - edit and insert, [1] - edit only
		map[string]string{
			"char":      "md5/sha1/password/encrypt/uuid", //! JavaScript for disabling maxlength
			"binary":    "md5/sha1",
			"date|time": "now",
		}, map[string]string{
			"(^|[^o])int|float|double|decimal": "+/-", // not point
			"date":                             "+ interval/- interval",
			"time":                             "addtime/subtime",
			"char|text":                        "concat",
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
