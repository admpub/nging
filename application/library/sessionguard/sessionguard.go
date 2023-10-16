package sessionguard

import (
	"encoding/json"

	"github.com/admpub/log"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"

	"github.com/admpub/nging/v5/application/dbschema"

	ip2regionparser "github.com/admpub/ip2region/v2/binding/golang/ip2region"
	"github.com/admpub/nging/v5/application/library/common"
	"github.com/admpub/nging/v5/application/library/ip2region"
)

func GetLastLoginInfo(ctx echo.Context, ownerType string, ownerId uint64, sessionId string) (*dbschema.NgingLoginLog, error) {
	m := dbschema.NewNgingLoginLog(ctx)
	err := m.Get(func(r db.Result) db.Result {
		return r.OrderBy(`-created`)
	}, db.And(
		db.Cond{`owner_type`: ownerType},
		db.Cond{`owner_id`: ownerId},
		db.Cond{`session_id`: sessionId},
		db.Cond{`success`: `Y`},
	))
	return m, err
}

// Validate 验证 session 环境是否安全，避免 cookie 和 session id 被窃取
// 在非匿名模式下 UserAgent 和 IP 归属地与登录时的一致
func Validate(ctx echo.Context, lastIP string, ownerType string, ownerId uint64) bool {
	if common.IsAnonymousMode(ownerType) {
		return true
	}
	env := GetEnvFromSession(ctx, ownerType)
	if env != nil {
		return validateEnv(ctx, lastIP, env.UserAgent, func() *ip2regionparser.IpInfo {
			return &env.Location
		})
	}
	info, err := GetLastLoginInfo(ctx, ownerType, ownerId, ctx.Session().MustID())
	if err != nil {
		log.Errorf(`failed to GetLastLoginInfo: %v`, err)
		return false
	}
	return validateEnv(ctx, lastIP, info.UserAgent, func() *ip2regionparser.IpInfo {
		if len(info.IpLocation) == 0 {
			return nil
		}
		ipLoc := map[string]string{}
		err = json.Unmarshal([]byte(info.IpLocation), &ipLoc)
		if err != nil {
			log.Errorf(`failed to unmarshal IpLocation: %v`, err)
			return nil
		}
		oldLocation := &ip2regionparser.IpInfo{
			Country:  ipLoc[`国家`],
			Region:   ipLoc[`地区`],
			Province: ipLoc[`省份`],
			City:     ipLoc[`城市`],
		}
		return oldLocation
	})
}

func validateEnv(ctx echo.Context, lastIP string, oldUserAgent string, oldLocationGetter func() *ip2regionparser.IpInfo) bool {
	if oldUserAgent != ctx.Request().UserAgent() {
		return false
	}
	currentIP := ctx.RealIP()
	if lastIP == currentIP {
		return true
	}
	ipInfo, err := ip2region.IPInfo(currentIP)
	if err != nil {
		log.Errorf(`failed to get IPInfo: %v`, err)
		return false
	}
	ip2region.ClearZero(&ipInfo)
	oldLocation := oldLocationGetter()
	if oldLocation == nil {
		return false
	}
	return ipInfo.Country == oldLocation.Country &&
		ipInfo.Region == oldLocation.Region &&
		ipInfo.Province == oldLocation.Province &&
		ipInfo.City == oldLocation.City
}
