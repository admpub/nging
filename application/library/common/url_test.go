package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/webx-top/com"
)

func TestSortedURLValues(t *testing.T) {
	r := NewSortedURLValues(`a=b&b=100&a=c`)
	assert.Equal(t, `a`, r[0].Key)
	assert.Equal(t, []string{`b`, `c`}, r[0].Values)
	assert.Equal(t, `b`, r[1].Key)
	assert.Equal(t, []string{`100`}, r[1].Values)
	r.Del(`a`)
	assert.Equal(t, 1, len(r))
	r.Del(`b`)
	assert.Equal(t, 0, len(r))
	r.ParseQuery(`aa=1&ab=2&ac=3&ad=4`)
	com.Dump(r)
	assert.Equal(t, 4, len(r))
	assert.Equal(t, `aa`, r[0].Key)
	assert.Equal(t, []string{`1`}, r[0].Values)
	assert.Equal(t, `ab`, r[1].Key)
	assert.Equal(t, []string{`2`}, r[1].Values)
	assert.Equal(t, `ac`, r[2].Key)
	assert.Equal(t, []string{`3`}, r[2].Values)
	assert.Equal(t, `ad`, r[3].Key)
	assert.Equal(t, []string{`4`}, r[3].Values)
}
