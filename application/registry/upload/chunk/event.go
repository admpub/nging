package chunk

import (
	"github.com/admpub/events"
	"github.com/admpub/log"
	"github.com/admpub/nging/v5/application/dbschema"
	"github.com/webx-top/echo"
)

func onUserLogoutOrLogin(ev events.Event) error {
	user := ev.Context.Get(`user`).(*dbschema.NgingUser)
	if user == nil {
		return nil
	}
	err := CleanFileByOwner(`user`, uint64(user.Id))
	if err != nil {
		log.Errorf(`failed to CleanFileByOwner(user, %d): %v`, user.Id, err.Error())
	}
	return nil
}

func init() {
	echo.OnCallback(`nging.user.logout.success`, onUserLogoutOrLogin)
	echo.OnCallback(`nging.user.login.success`, onUserLogoutOrLogin)
}
