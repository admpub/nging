package selfupdate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRestart(t *testing.T) {
	err := Restart(nil, `/Users/hank/go/src/github.com/admpub/nging/nging`)
	assert.NoError(t, err)
}
