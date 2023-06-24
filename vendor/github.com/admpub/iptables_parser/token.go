package iptables_parser

type Token int

const (
	// Special tokens
	ILLEGAL Token = iota
	EOF
	WS

	// Literals
	IDENT // main

	// Misc characters
	COLON     // :
	HASHTAG   // #
	QUOTATION // "
	BACKSLASH // \
	NOT       // !
	COMMA
	NEWLINE

	// Keywords
	SRC
	DEST
	COUNTER
	HEADER
	COMMENT
	COMMENTLINE
	APPEND
	FLAG
	DEFAULT
	LINE
)
