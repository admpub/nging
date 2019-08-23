package upload

import (
	"strings"
	"testing"

	"github.com/webx-top/echo/testing/test"
)

func TestUploadURL(t *testing.T) {
	test.True(t, strings.HasPrefix(BackendUploadURL(`/manager/upload/movie`, `refid`, `123`), `/manager/upload//manager/upload/movie?refid=123&time=`))
}
