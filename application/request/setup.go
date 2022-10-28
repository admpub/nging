package request

import (
	"strings"

	"github.com/admpub/nging/v5/application/library/config/subconfig/sdb"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
)

var _ echo.AfterValidate = &Setup{}

type Setup struct {
	Type       string `validate:"oneof='mysql' 'sqlite'"`
	User       string `validate:""`
	Password   string `validate:""`
	Host       string `validate:"omitempty,hostname"`
	Database   string `validate:"omitempty,alphanum_"`
	Charset    string `validate:""`
	AdminUser  string `validate:"username"`
	AdminPass  string `validate:"required,min=8,max=64"`
	AdminEmail string `validate:"required,email"`
}

func (s *Setup) AfterValidate(c echo.Context) error {
	if s.Type == `sqlite` {
		s.User = ``
		s.Password = ``
		s.Host = ``
		if !strings.HasSuffix(s.Database, `.db`) {
			s.Database += `.db`
		}
	} else {
		if len(s.Charset) == 0 {
			s.Charset = sdb.MySQLDefaultCharset
		} else {
			if !com.InSlice(s.Charset, sdb.MySQLSupportCharsetList) {
				return c.NewError(code.InvalidParameter, c.T(`字符集参数无效`)).SetZone(`charset`)
			}
		}
	}
	return nil
}
