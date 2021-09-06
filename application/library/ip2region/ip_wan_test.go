package ip2region

import (
	"fmt"
	"testing"

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
