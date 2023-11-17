package request

import (
	"github.com/docker/docker/api/types"
)

type Login struct {
	types.AuthConfig
}
