package ip2region

import "testing"

func TestIPWAN(t *testing.T) {
	wan, err := GetWANIP()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(wan.IPv4)
	t.Log(wan.IPv6)
}
