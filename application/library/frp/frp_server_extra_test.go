package frp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServerExtra(t *testing.T) {
	cfg := NewServerConfigExtra()
	cfg.Extra = []byte(`{"a":1,"b":2}`)
	str := cfg.String()
	err := cfg.Parse(str)
	if err != nil {
		panic(err)
	}
	assert.Equal(t, 1, cfg.unmarshaledExtra.Int(`a`))
	assert.Equal(t, 2, cfg.unmarshaledExtra.Int(`b`))
}
