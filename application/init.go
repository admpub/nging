/*
   Nging is a toolbox for webmasters
   Copyright (C) 2018-present Wenhui Shen <swh@admpub.com>

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
package application

import (
	_ "github.com/admpub/nging/v4/application/handler/index"
	_ "github.com/admpub/nging/v4/application/handler/manager"
	_ "github.com/admpub/nging/v4/application/handler/setup"
	_ "github.com/admpub/nging/v4/application/handler/tool"
	_ "github.com/admpub/nging/v4/application/handler/user"
	_ "github.com/admpub/nging/v4/application/initialize/backend"
	_ "github.com/admpub/nging/v4/application/library/cron/cmder"
	_ "github.com/admpub/nging/v4/application/listener"
)
