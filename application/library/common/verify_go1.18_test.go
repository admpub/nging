package common

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIf(t *testing.T) {
	v := If(true, 1, 2)
	assert.Equal(t, `int`, fmt.Sprintf("%T", v))
	assert.Equal(t, 1, v)

	v2 := If(true, uint64(1), 2)
	assert.Equal(t, `uint64`, fmt.Sprintf("%T", v2))
	assert.Equal(t, uint64(1), v2)
}
