package utils

import (
	"fmt"
	"testing"

	"github.com/webx-top/com"
)

func TestNetInterface(t *testing.T) {
	ipv4, ipv6, err := GetNetInterface(``)
	if err != nil {
		panic(err)
	}
	fmt.Println(`ipv4:`, com.Dump(ipv4, false))
	fmt.Println(`ipv6:`, com.Dump(ipv6, false))
}
