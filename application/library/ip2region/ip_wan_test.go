package ip2region

import (
	"fmt"
	"testing"

	"github.com/admpub/ip2region/v2/binding/golang/ip2region"
	"github.com/stretchr/testify/assert"
)

func TestIPWAN(t *testing.T) {
	wan, err := GetWANIP(0, 4)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(`IPv4:`, wan.IP)
	wan, err = GetWANIP(0, 6) //在不支持IPv6的环境下会抛panic
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(`IPv6:`, wan.IP)
	matches := ipv6Regexp.FindAllStringSubmatch(` 	2001:3CA1:010F:001A:121B:0000:0000:0010`, 1)
	expected := [][]string{
		[]string{
			"2001:3CA1:010F:001A:121B:0000:0000:0010",
			"2001:3CA1:010F:001A:121B:0000:0000:0010",
			"", "",
		},
	}
	assert.Equal(t, expected, matches)
	//panic(``)
}

func TestStringify(t *testing.T) {
	info := ip2region.IpInfo{
		Country:  `中国`,
		Region:   `东北`,
		Province: `山东省`,
		City:     `济南市`,
		ISP:      `联通`,
	}
	result := Stringify(info, false)
	assert.Equal(t, `中国东北山东省济南市联通`, result)

	result2 := Stringify(info)
	assert.Equal(t, `{"国家":"中国","地区":"东北","省份":"山东省","城市":"济南市","线路":"联通"}`, result2)
}
