package mysql

import "regexp"

var (
	reCollate       = regexp.MustCompile(` COLLATE ([^ ]+)`)
	reCharacter     = regexp.MustCompile(` CHARACTER SET ([^ ]+)`)
	reInnoDBComment = regexp.MustCompile(`(?:(.+); )?InnoDB free: .*`)
)
