package ip2region

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/webx-top/com"
)

func TestIPWAN(t *testing.T) {
	wan, err := GetWANIP(true)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(`IPv4:`, wan.IPv4)
	fmt.Println(`IPv6:`, wan.IPv6)
	r := regexp.MustCompile(IPv6Rule)
	matches := r.FindAllStringSubmatch(` 	2001:3CA1:010F:001A:121B:0000:0000:0010`, 1)
	com.Dump(matches)
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
