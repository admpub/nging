package pipe

import (
	"github.com/admpub/nging/application/registry/upload/driver"
	uploadClient "github.com/webx-top/client/upload"
)

type PipeFunc func(storer driver.Storer, results uploadClient.Results, recv interface{}) error

var pipes = map[string]PipeFunc{}

func Register(pipeName string, pipeFunc PipeFunc) {
	pipes[pipeName] = pipeFunc
}

func Get(pipeName string) PipeFunc {
	fn, ok := pipes[pipeName]
	if ok {
		return fn
	}
	return nil
}
