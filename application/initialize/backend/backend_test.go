package backend

import (
	"testing"

	"github.com/webx-top/echo/testing/test"
)

func TestMakeSubdomains(t *testing.T) {
	r := MakeSubdomains(`www.coscms.com,www.coscms.com:8181`, []string{})
	test.Eq(t, `www.coscms.com,www.coscms.com:8181,www.coscms.com:0`, r)
	r = MakeSubdomains(`www.coscms.com,www.coscms.com:8181`, DefaultLocalHostNames)
	test.Eq(t, `www.coscms.com,www.coscms.com:8181,www.coscms.com:0,127.0.0.1:0,localhost:0`, r)
}
