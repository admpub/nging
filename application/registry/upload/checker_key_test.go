package upload

import (
	"net/url"
	"strings"
	"testing"

	"github.com/webx-top/echo/testing/test"
)

func TestUploadURL(t *testing.T) {
	urls := BackendUploadURL(`/manager/upload/movie`, `refid`, `123`)
	values, err := url.ParseQuery(strings.SplitN(urls, `?`, 2)[1])
	if err != nil {
		t.Fatal(err)
	}
	//com.Dump(values)
	test.True(t, strings.HasPrefix(urls, `/manager/upload//manager/upload/movie?refid=123&time=`))
	token := values.Get(`token`)
	values.Del(`token`)
	test.Eq(t, token, Token(values))
}
