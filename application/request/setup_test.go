package request

import (
	"testing"

	"github.com/admpub/copier"
	"github.com/stretchr/testify/assert"
	"github.com/webx-top/echo/defaults"
	"github.com/webx-top/validator"
)

var testValidator *validator.Validate

func init() {
	testValidator = validator.New(defaults.NewMockContext(), `zh`)
}

func TestSetup(t *testing.T) {
	data := &Setup{
		Type:       `mysql`,
		User:       ``,
		Password:   ``,
		Host:       ``,
		Database:   ``,
		Charset:    ``,
		AdminUser:  `admin`,
		AdminPass:  `admin123`,
		AdminEmail: `test@admpub.com`,
	}
	result := testValidator.Validate(data)
	assert.NoError(t, result.Error())

	dataCopy := &Setup{}
	err := copier.Copy(dataCopy, data)
	assert.NoError(t, err)
	assert.Equal(t, data, dataCopy)

	dataCopy.Type = `errType`
	result = testValidator.Validate(dataCopy)
	assert.Error(t, result.Error())
	assert.Equal(t, `Type`, result.Field())

	dataCopy.Type = data.Type
	dataCopy.AdminUser = "admin\nadmin"
	result = testValidator.Validate(dataCopy)
	assert.Error(t, result.Error())
	assert.Equal(t, `AdminUser`, result.Field())

	dataCopy.AdminUser = data.AdminUser
	dataCopy.Database = "'"
	result = testValidator.Validate(dataCopy)
	assert.Error(t, result.Error())
	assert.Equal(t, `Database`, result.Field())
	dataCopy.Database = "`"
	result = testValidator.Validate(dataCopy)
	assert.Error(t, result.Error())
	assert.Equal(t, `Database`, result.Field())
}
