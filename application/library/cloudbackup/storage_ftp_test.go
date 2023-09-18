package cloudbackup

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/webx-top/echo/defaults"
)

func TestStorageFTP(t *testing.T) {
	c := NewStorageFTP(`127.0.0.1:25`, `test`, `test123456`).(*StorageFTP)
	err := c.Connect()
	assert.NoError(t, err)
	defer c.Close()
	err = c.MkdirAll(`/mkdirall/1/2`)
	assert.NoError(t, err)
	ctx := defaults.NewMockContext()
	fp, err := os.Open(`./storage_ftp_test.go`)
	assert.NoError(t, err)
	defer fp.Close()
	err = c.Put(ctx, fp, `/mkdirall/1/2/3/storage_ftp_test.go`, 0)
	assert.NoError(t, err)
	err = c.Put(ctx, fp, `/mkdirall2/1/2/storage_ftp_test.go`, 0)
	assert.NoError(t, err)
}
