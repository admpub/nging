package setup

import (
	"github.com/webx-top/echo"

	"github.com/admpub/nging/v4/application/library/role"
	"github.com/nging-plugins/servermanager/application/model"
)

func init() {
	role.PermCommandList = getCommandList
}

func getCommandList(ctx echo.Context) ([]interface{}, error) {
	cmdM := model.NewCommand(ctx)
	_, err := cmdM.ListByOffset(nil, nil, 0, -1, `disabled`, `N`)
	if err != nil {
		return nil, err
	}
	rows := cmdM.Objects()
	list := make([]interface{}, len(rows))
	for index, row := range rows {
		list[index] = row
	}
	return list, err
}
