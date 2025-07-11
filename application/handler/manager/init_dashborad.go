package manager

import (
	"html/template"

	"github.com/webx-top/com"
	"github.com/webx-top/echo"

	"github.com/coscms/webcore/library/backend"
	"github.com/coscms/webcore/library/dashboard"
	"github.com/coscms/webcore/library/httpserver"
	"github.com/coscms/webcore/library/notice"
	"github.com/coscms/webcore/library/role"
	"github.com/coscms/webcore/model"
)

func init() {
	httpserver.Backend.Dashboard.Cards.Add(0, (&dashboard.Card{
		IconName:  `fa-user`,
		IconColor: `success`,
		Short:     `USERS`,
		Name:      echo.T(`用户数量`),
		Summary:   ``,
	}).SetContentGenerator(func(ctx echo.Context) interface{} {
		//用户统计
		userMdl := model.NewUser(ctx)
		userCount, _ := userMdl.Count(nil)
		user := backend.User(ctx)
		if user == nil || !role.IsFounder(user) {
			return userCount
		}
		onlineUsers := notice.OnlineUserCount()
		if !notice.IsOnline(user.Username) {
			onlineUsers++
		}
		return template.HTML(com.String(userCount) + ` <a class="label label-success" href="` + backend.URLFor(`/manager/user`) + `?online=Y">` + ctx.T(`在线`) + `:` + com.String(onlineUsers) + `</a>`)
	}))
	httpserver.Backend.Dashboard.TopButtons.Add(0,
		&dashboard.Button{
			Tmpl: `manager/topbutton/language`,
		},
		&dashboard.Button{
			Tmpl: `manager/topbutton/upgrade`,
		},
		&dashboard.Button{
			Tmpl: `manager/topbutton/official-site`,
		},
		&dashboard.Button{
			Tmpl: `manager/topbutton/logs`,
		},
	)
	//httpserver.Backend.Dashboard.GlobalFooters.Add(-1,)
}
