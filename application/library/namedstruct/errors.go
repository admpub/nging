package namedstruct

import "errors"

var (
	ErrNotExist     = errors.New(`the struct not exist`)
	ErrNameConflict = errors.New(`name conflict! failed to register, this name is already registered`)
)
