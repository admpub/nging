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
	"github.com/admpub/nging/application/handler"
	"github.com/admpub/nging/application/registry/upload/helper"
	"github.com/webx-top/echo"
)

func init() {
	handler.RegisterToGroup(`/manager`, func(g echo.RouteRegister) {
		handler.Echo().Route(`GET,HEAD`, helper.UploadURLPath+`:type/`+`*`, File)
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
	})
}
