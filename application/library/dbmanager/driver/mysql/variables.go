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
	reFieldOnUpdate       = regexp.MustCompile("^on update (.+)")
	reFieldDefault        = regexp.MustCompile("char|set")
	reFieldPrivilegeDelim = regexp.MustCompile(", *")

	reFieldTypeNumber    = regexp.MustCompile("(^|[^o])int|float|double|decimal")
	reFieldTypeText      = regexp.MustCompile("char|text|enum|set")
	reFieldTypeBit       = regexp.MustCompile("^([0-9]+|b'[0-1]+')\\$")
	reFieldLengthInvalid = regexp.MustCompile("[^-0-9,+()[\\]]")
	reFieldLengthNumber  = regexp.MustCompile("^[0-9].*")

	pgsqlFieldDefaultValue = regexp.MustCompile("^[a-z]+\\(('[^']*')+\\)\\$")

	//以下数据来自客户端
	reGrantColumn      = regexp.MustCompile(`^([^() ]+)\s*(\([^)]*\))?$`)
	reGrantOptionValue = regexp.MustCompile(`(GRANT OPTION)\([^)]*\)`)
	reNotWord          = regexp.MustCompile(`[^a-zA-Z0-9_]+`)

	UnsignedTags = []string{"unsigned", "zerofill", "unsigned zerofill"}
	EnumLength   = "'(?:''|[^'\\\\]|\\\\.)*'"
)
