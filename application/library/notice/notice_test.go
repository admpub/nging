package notice

import (
	"testing"

	"github.com/webx-top/com"
)

func TestOpenMessage(t *testing.T) {
	OpenMessage(`testUser`, `testType`)
	user := DefaultUserNotices.User[`testUser`]
	if len(user.Notice.Types) != 1 {
		t.Errorf(`Size of types != %v`, 1)
	}
	if user.Notice.Types[`testType`] != true {
		t.Error(`Type of testType != true`)
	}
	com.Dump(DefaultUserNotices)
}

func TestSend(t *testing.T) {
	go func() {
		t.Log(string(RecvJSON(`testUser`)))
	}()
	Send(`testUser`, NewMessageWithValue(`testType`, `testTitle`, `testContent`))
	com.Dump(DefaultUserNotices)
}
