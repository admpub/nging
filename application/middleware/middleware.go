package middleware

import (
	"github.com/admpub/nging/v4/application/dbschema"
	"github.com/admpub/nging/v4/application/handler"
	"github.com/admpub/nging/v4/application/model"
	"github.com/webx-top/echo"
)

var Middlewares []interface{}

func Use(m ...interface{}) {
	Middlewares = append(Middlewares, m...)
}

func init() {
	handler.GetRoleList = func(c echo.Context) []*dbschema.NgingUserRole {
		user := handler.User(c)
		if user == nil {
			return nil
		}
		roleM := model.NewUserRole(c)
		return roleM.ListByUser(user)
	}
}
