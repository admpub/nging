package table

import (
	"errors"
)

var (
	ErrExistsFile       = errors.New("exists file")
	ErrInvalidFieldName = errors.New("Invalid field name")
)
