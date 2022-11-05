package request

import (
	"strings"

	"github.com/admpub/nging/v5/application/library/config"
	"github.com/admpub/nging/v5/application/library/config/subconfig/sdb"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
)

var _ echo.AfterValidate = &Setup{}

type Setup struct {
	Type       string `validate:"required"`
	User       string
	Password   string
	Host       string
	Database   string
	Charset    string `validate:"omitempty,alphanum"`
	AdminUser  string `validate:"username"`
	AdminPass  string `validate:"required,min=8,max=64"`
	AdminEmail string `validate:"required,email"`
}

func (s *Setup) AfterValidate(c echo.Context) error {
	if _, ok := config.DBInstallers[s.Type]; !ok {
		return c.NewError(code.InvalidParameter, `不支持的数据库类型: %v`, s.Type).SetZone(`type`)
	}
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
				return c.NewError(code.InvalidParameter, `字符集参数无效`).SetZone(`charset`)
			}
		}
	}
	return nil
}
