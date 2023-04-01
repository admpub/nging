package ip2region

import (
	"fmt"
	"strings"

	"github.com/admpub/ip2region/v2/binding/golang/ip2region"
	syncOnce "github.com/admpub/once"
	"github.com/webx-top/echo"
)

var (
	region   *ip2region.Ip2Region
	dictFile string
	once     syncOnce.Once
)

func init() {
	dictFile = echo.Wd() + echo.FilePathSeparator + `data` + echo.FilePathSeparator + `ip2region` + echo.FilePathSeparator + `ip2region.xdb`
}

func SetDictFile(f string) {
	dictFile = f
	once.Reset()
}

func SetInstance(newInstance *ip2region.Ip2Region) {
	if region == nil {
		region = newInstance
	} else {
		oldRegion := *region
		*region = *newInstance
		oldRegion.Close()
	}
}

func initialize() (err error) {
	if region != nil {
		region.Close()
	}
	region, err = ip2region.New(dictFile)
	return
}

func IsInitialized() bool {
	return region != nil
}

func IPInfo(ip string) (info ip2region.IpInfo, err error) {
	if len(ip) == 0 {
		return
	}
	once.Do(func() {
		err = initialize()
	})
	if err != nil {
		return
	}
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf(`%v`, e)
		}
	}()
	info, err = region.MemorySearch(ip)
	return
}

func ClearZero(info *ip2region.IpInfo) {
	if info.Country == `0` {
		info.Country = ``
	}
	if info.Region == `0` {
		info.Region = ``
	}
	if info.Province == `0` {
		info.Province = ``
	}
	if info.City == `0` {
		info.City = ``
	}
	if info.ISP == `0` {
		info.ISP = ``
	}
}

func IsZero(str string) bool {
	return len(str) == 0 || str == `0`
}

func Stringify(info ip2region.IpInfo, jsonify ...bool) string {
	var (
		formats []string
		args    []interface{}
	)
	if !IsZero(info.Country) {
		formats = append(formats, `"国家":%q`)
		args = append(args, info.Country)
	}
	if !IsZero(info.Region) {
		formats = append(formats, `"地区":%q`)
		args = append(args, info.Region)
	}
	if !IsZero(info.Province) {
		formats = append(formats, `"省份":%q`)
		args = append(args, info.Province)
	}
	if !IsZero(info.City) {
		formats = append(formats, `"城市":%q`)
		args = append(args, info.City)
	}
	if !IsZero(info.ISP) {
		formats = append(formats, `"线路":%q`)
		args = append(args, info.ISP)
	}
	if len(jsonify) == 0 || jsonify[0] {
		return fmt.Sprintf(`{`+strings.Join(formats, `,`)+`}`, args...)
	}
	return fmt.Sprintf(strings.Repeat(`%s`, len(args)), args...)
}
