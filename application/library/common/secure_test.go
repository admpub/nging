package common

import (
	"testing"
	"github.com/webx-top/echo/testing/test"
)

func TestSecure(t *testing.T) {
	s:=`<p>test<a href="http://www.admpub.com">link</a>test</p>`
	test.Eq(t,`<p>test<a href="http://www.admpub.com" rel="nofollow">link</a>test</p>`,RemoveXSS(s))
	s=`<p>test<a href="http://www.admpub.com"><img src="http://www.admpub.com/test" />link</a>test</p>`
	test.Eq(t,`<p>test<img src="http://www.admpub.com/test"/>linktest</p>`,RemoveXSS(s, true))
}