package packer

import "errors"

var ErrNotFound = errors.New(`no package manager found`)
var ErrUnsupported = errors.New(`not supported`)
