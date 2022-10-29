package request

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	data := &Register{
		InvitationCode:       `test1234567890abc`,
		Username:             `test`,
		Email:                `123@webx.top`,
		Password:             `12345678`,
		ConfirmationPassword: `12345678`,
	}
	result := testValidator.Validate(data)
	assert.NoError(t, result.Error())
}
