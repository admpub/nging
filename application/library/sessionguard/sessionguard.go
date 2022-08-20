package sessionguard

import (
	"encoding/json"

	"github.com/admpub/log"
	"github.com/admpub/nging/v4/application/dbschema"
	"github.com/admpub/nging/v4/application/library/ip2region"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
)

func GetLastLoginInfo(ctx echo.Context, ownerType string, ownerId uint64, sessionId string) (*dbschema.NgingLoginLog, error) {
	m := dbschema.NewNgingLoginLog(ctx)
	err := m.Get(func(r db.Result) db.Result {
		return r.OrderBy(`-created`)
	}, db.And(
		db.Cond{`owner_id`: ownerId},
		db.Cond{`owner_type`: ownerType},
		db.Cond{`session_id`: sessionId},
		db.Cond{`success`: `Y`},
	))
	return m, err
}

// Validate 验证 session 环境是否安全，避免 cookie 和 session id 被窃取
// 在非匿名模式下 UserAgent 和 IP 归属地与登录时的一致
func Validate(ctx echo.Context, lastIP string, ownerType string, ownerId uint64) bool {
	k := `backend.Anonymous`
	if ownerType != `user` {
		k = `frontend.Anonymous`
	}
	if echo.Bool(k) {
		return true
	}
	info, err := GetLastLoginInfo(ctx, ownerType, ownerId, ctx.Session().MustID())
	if err != nil {
		log.Errorf(`failed to GetLastLoginInfo: %v`, err)
		return false
	}
	if info.UserAgent != ctx.Request().UserAgent() {
		return false
	}
	currentIP := ctx.RealIP()
	if lastIP == currentIP {
		return true
	}
	if len(info.IpLocation) == 0 {
		return false
	}
	ipLoc := map[string]string{}
	err = json.Unmarshal([]byte(info.IpLocation), &ipLoc)
	if err != nil {
		log.Errorf(`failed to unmarshal IpLocation: %v`, err)
		return false
	}
	ipInfo, err := ip2region.IPInfo(currentIP)
	if err != nil {
		log.Errorf(`failed to get IPInfo: %v`, err)
		return false
	}
	if ipInfo.Country == `0` {
		ipInfo.Country = ``
	}
	if ipInfo.Region == `0` {
		ipInfo.Region = ``
	}
	if ipInfo.Province == `0` {
		ipInfo.Province = ``
	}
	if ipInfo.City == `0` {
		ipInfo.City = ``
	}
	return ipInfo.Country == ipLoc[`国家`] &&
		ipInfo.Region == ipLoc[`地区`] &&
		ipInfo.Province == ipLoc[`省份`] &&
		ipInfo.City == ipLoc[`城市`]
}
