package upgrade

import (
	"encoding/json"

	"github.com/admpub/nging/v5/application/dbschema"
	"github.com/admpub/nging/v5/application/library/common"
	"github.com/admpub/nging/v5/application/library/role"
	"github.com/admpub/nging/v5/application/model"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/defaults"
)

func init() {
	type Role struct {
		Id           uint   `db:"id,omitempty,pk" bson:"id,omitempty" comment:"ID" json:"id" xml:"id"`
		PermCmd      string `db:"perm_cmd" bson:"perm_cmd" comment:"指令集权限(多个用“,”隔开)" json:"perm_cmd" xml:"perm_cmd"`
		PermAction   string `db:"perm_action" bson:"perm_action" comment:"操作权限(多个用“,”隔开)" json:"perm_action" xml:"perm_action"`
		PermBehavior string `db:"perm_behavior" bson:"perm_behavior" comment:"行为权限(多个用“,”隔开)" json:"perm_behavior" xml:"perm_behavior"`
	}
	echo.OnCallback(`nging.upgrade.db.before`, func(data echo.Event) error {
		installedSchemaVer := data.Context.Float64(`installedSchemaVer`)
		if installedSchemaVer >= 5 {
			return nil
		}
		ctx := defaults.NewMockContext()
		m := dbschema.NewNgingUserRole(ctx)
		recv := []*Role{}
		_, err := m.NewParam().SetRecv(&recv).List()
		if err != nil {
			return err
		}
		b, err := json.Marshal(recv)
		if err != nil {
			return err
		}
		err = common.WriteCache(`upgrade`, `db.lt5.nging_user_role.json`, b)
		return err
	})
	echo.OnCallback(`nging.upgrade.db.after`, func(data echo.Event) error {
		installedSchemaVer := data.Context.Float64(`installedSchemaVer`)
		if installedSchemaVer >= 5 {
			return nil
		}
		b, err := common.ReadCache(`upgrade`, `db.lt5.nging_user_role.json`)
		if err != nil {
			return err
		}
		recv := []*Role{}
		err = json.Unmarshal(b, &recv)
		if err != nil {
			return err
		}
		ctx := defaults.NewMockContext()
		rpM := model.NewUserRolePermission(ctx)
		for _, row := range recv {
			rpM.RoleId = row.Id
			if len(row.PermAction) > 0 {
				rpM.Type = role.RolePermissionTypePage
				rpM.Permission = row.PermAction
				_, err = rpM.Add()
				if err != nil {
					break
				}
			}
			if len(row.PermCmd) > 0 {
				rpM.Type = role.RolePermissionTypeCommand
				rpM.Permission = row.PermCmd
				_, err = rpM.Add()
				if err != nil {
					break
				}
			}
			if len(row.PermBehavior) > 0 {
				rpM.Type = role.RolePermissionTypeBehavior
				rpM.Permission = row.PermBehavior
				_, err = rpM.Add()
				if err != nil {
					break
				}
			}
		}
		if err == nil {
			common.RemoveCache(`upgrade`, `db.lt5.nging_user_role.json`)
		}
		return err
	})
}
