package cmder

import (
	"errors"
	"fmt"

	"github.com/admpub/nging/v4/application/library/config/cmder"
)

var (
	ErrUnsupported             = errors.New("unsupported")
	ErrNoAvailibaleConfigFound = errors.New(`No available configurations found`)
)

func GetServer() (*FRPServer, error) {
	cm := cmder.Get(`frpserver`)
	if cm == nil {
		return nil, fmt.Errorf("frpserver: %w", ErrUnsupported)
	}
	return cm.(*FRPServer), nil
}

func GetClient() (*FRPClient, error) {
	cm := cmder.Get(`frpclient`)
	if cm == nil {
		return nil, fmt.Errorf("frpclient: %w", ErrUnsupported)
	}
	return cm.(*FRPClient), nil
}
