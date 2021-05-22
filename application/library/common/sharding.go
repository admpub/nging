package common

import (
	"math"
	"time"

	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/param"
)

// IDSharding 按照ID进行分片
func IDSharding(id uint64, shardingNum float64) uint64 {
	return uint64(math.Ceil(float64(id) / shardingNum))
}

// MD5Sharding 按照MD5进行分片
func MD5Sharding(str interface{}, length ...int) string {
	v := com.Md5(param.AsString(str))
	if len(length) == 0 || length[0] < 1 {
		return v[0:1]
	}
	if length[0] >= 32 {
		return v
	}
	return v[0:length[0]]
}

// MonthSharding 按照日期月进行分片
func MonthSharding(ctx echo.Context) string {
	return GetNowTime(ctx).Format("2006_01")
}

// YearSharding 按照日期年进行分片
func YearSharding(ctx echo.Context) string {
	return GetNowTime(ctx).Format("2006")
}

// GetNowTime 获取当前时间(同一个context中只获取一次)
func GetNowTime(ctx echo.Context) time.Time {
	t, y := ctx.Internal().Get(`time.now`).(time.Time)
	if !y || t.IsZero() {
		t = time.Now().Local()
		ctx.Internal().Set(`time.now`, t)
	}
	return t
}
