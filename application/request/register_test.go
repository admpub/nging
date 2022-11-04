package request

import (
	"testing"

	"github.com/admpub/nging/v5/application/library/codec"
	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	pass, _ := codec.DefaultSM2EncryptHex(`12345678`)
	data := &Register{
		InvitationCode:       `test1234567890abc`,
		Username:             `test`,
		Email:                `123@webx.top`,
		Password:             pass,
		ConfirmationPassword: `12345678`,
	}
	err := echoContext.Validate(data)
	assert.NoError(t, err)
}
