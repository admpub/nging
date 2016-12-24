package result

import "encoding/gob"

func init() {
	gob.Register([]Resulter{})
}

type Resulter interface {
	GetSQL() string
	GetBeginTime() string
	GetElapsedTime() string
	GetAffected() int64
	GetError() string
}
