package mysql

import "regexp"

var (
	reCollate       = regexp.MustCompile(` COLLATE ([^ ]+)`)
	reCharacter     = regexp.MustCompile(` CHARACTER SET ([^ ]+)`)
	reInnoDBComment = regexp.MustCompile(`(?:(.+); )?InnoDB free: .*`)

	reGrantOn       = regexp.MustCompile(`GRANT (.*) ON (.*) TO `)
	reGrantBrackets = regexp.MustCompile(` *([^(,]*[^ ,(])( *\\([^)]+\\))?`)
	reGrantOption   = regexp.MustCompile(` WITH GRANT OPTION`)
	reGrantIdent    = regexp.MustCompile(` IDENTIFIED BY PASSWORD '([^']+)`)
)
