package utils

import (
	"os"

	"github.com/sirupsen/logrus"
)

var (
	timestampFormat = "2006/01/0 10:10:10"
	Logger          = &logrus.Logger{
		Out: os.Stdout,
		Formatter: &logrus.TextFormatter{
			TimestampFormat: timestampFormat,
		},
		Hooks: make(logrus.LevelHooks),
		Level: logrus.InfoLevel,
	}
)
