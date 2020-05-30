/*
   Nging is a toolbox for webmasters
   Copyright (C) 2018-present  Wenhui Shen <swh@admpub.com>

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU Affero General Public License as published
   by the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Affero General Public License for more details.

   You should have received a copy of the GNU Affero General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package manager

import (
	"github.com/webx-top/echo"

	"github.com/admpub/nging/application/handler"
	_ "github.com/admpub/nging/application/handler/manager/file"
	"github.com/admpub/nging/application/registry/upload/helper"
)

func init() {
	handler.RegisterToGroup(`/manager`, func(g echo.RouteRegister) {
		handler.Echo().Route(`GET,HEAD`, helper.UploadURLPath+`:type/*`, File) //显示上传文件夹下的静态文件
		g.Route(`GET,POST`, `/user`, User)
		g.Route(`GET,POST`, `/role`, Role)
		g.Route(`GET,POST`, `/user_add`, UserAdd)
		g.Route(`GET,POST`, `/user_edit`, UserEdit)
		g.Route(`GET,POST`, `/user_delete`, UserDelete)
		g.Route(`GET,POST`, `/role_add`, RoleAdd)
		g.Route(`GET,POST`, `/role_edit`, RoleEdit)
		g.Route(`GET,POST`, `/role_delete`, RoleDelete)
		g.Route(`GET,POST`, `/invitation`, Invitation)
		g.Route(`GET,POST`, `/invitation_add`, InvitationAdd)
		g.Route(`GET,POST`, `/invitation_edit`, InvitationEdit)
		g.Route(`GET,POST`, `/invitation_delete`, InvitationDelete)
		g.Route(`GET,POST`, `/verification`, Verification)
		g.Route(`GET,POST`, `/verification_delete`, VerificationDelete)
		g.Route(`GET`, `/clear_cache`, ClearCache)
		g.Route(`GET,POST`, `/settings`, Settings)
		g.Route(`POST`, `/upload/:type`, Upload) //文件上传
		g.Route(`GET,POST`, `/crop`, Crop)       //裁剪图片
		g.Route(`GET,POST`, `/uploaded_file`, UploadedFile)
		g.Route(`GET,POST`, `/alert_recipient`, AlertRecipient)
		g.Route(`GET,POST`, `/alert_recipient_add`, AlertRecipientAdd)
		g.Route(`GET,POST`, `/alert_recipient_edit`, AlertRecipientEdit)
		g.Route(`GET,POST`, `/alert_recipient_test`, AlertRecipientTest)
		g.Route(`GET,POST`, `/alert_recipient_edit`, AlertRecipientEdit)
		g.Route(`GET,POST`, `/alert_recipient_delete`, AlertRecipientDelete)
		g.Route(`GET,POST`, `/kv`, KvIndex)
		g.Route(`GET,POST`, `/kv_add`, KvAdd)
		g.Route(`GET,POST`, `/kv_edit`, KvEdit)
		g.Route(`GET,POST`, `/kv_delete`, KvDelete)
	})

}
