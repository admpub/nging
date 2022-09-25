package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFloorNumber(t *testing.T) {
	n := FloorNumber(1, 20, 10)
	assert.Equal(t, 11, n)
	n = FloorNumber(2, 20, 10)
	assert.Equal(t, 31, n)
	n = FloorNumber(2, 20, 0)
	assert.Equal(t, 21, n)
}
