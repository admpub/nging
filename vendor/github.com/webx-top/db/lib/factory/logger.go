package factory

import "github.com/admpub/log"

var Log Logger = log.GetLogger(`db`)

type Logger interface {
	Error(a ...interface{})
	Warn(a ...interface{})
	Info(a ...interface{})
}
