package request

import (
	"testing"

	"github.com/admpub/nging/v5/application/library/codec"
	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	pass, _ := codec.DefaultSM2EncryptHex(`12345678`)
	rawPwd, _ := codec.DefaultSM2DecryptHex(pass)
	assert.Equal(t, rawPwd, `12345678`)
	data := &Login{
		User: `test`,
		Pass: pass,
		Code: `12345678`,
	}
	err := echoContext.Validate(data)
	assert.NoError(t, err)
}
