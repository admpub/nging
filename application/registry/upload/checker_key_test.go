package upload

import (
	"fmt"
	"net/url"
	"strings"
	"testing"

	"github.com/webx-top/echo/testing/test"
)

func TestUploadURL(t *testing.T) {
	urls := BackendUploadURL(`movie`, `refid`, `123`)
	values, err := url.ParseQuery(strings.SplitN(urls, `?`, 2)[1])
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(urls)
	//com.Dump(values)
	test.True(t, strings.HasPrefix(urls, `/manager/upload/movie?refid=123&time=`))
	token := values.Get(`token`)
	values.Del(`token`)
	test.Eq(t, token, Token(values))
	test.Eq(t, token, Token(`refid`, values.Get(`refid`), `time`, values.Get(`time`)))
}

func TestUploadURL2(t *testing.T) {
	urls := BackendUploadURL(`movie`, `refid`, ``)
	values, err := url.ParseQuery(strings.SplitN(urls, `?`, 2)[1])
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(urls)
	//com.Dump(values)
	test.True(t, strings.HasPrefix(urls, `/manager/upload/movie?refid=0&time=`))
	token := values.Get(`token`)
	values.Del(`token`)
	test.Eq(t, token, Token(values))
	test.Eq(t, token, Token(`refid`, values.Get(`refid`), `time`, values.Get(`time`)))
}
