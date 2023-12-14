package cloudbackup

import (
	"os"
	"testing"

	"github.com/admpub/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/webx-top/echo/defaults"
)

func TestStorageWebDAV(t *testing.T) {
	godotenv.Load()
	c := NewStorageWebDAV(`http://192.168.50.132:65000`, os.Getenv(`STORAGE_WEBDAV_USERNAME`), os.Getenv(`STORAGE_WEBDAV_PASSWORD`)).(*StorageWebDAV)
	err := c.Connect()
	assert.NoError(t, err)
	defer c.Close()
	ctx := defaults.NewMockContext()

	fp, err := os.Open(`./storage_ftp_test.go`)
	assert.NoError(t, err)
	stat, err := fp.Stat()
	assert.NoError(t, err)
	err = c.Put(ctx, fp, `/sdb9/cloudbackuptest/1/2/3/storage_ftp_test.go`, stat.Size())
	assert.NoError(t, err)
	fp.Close()

	fp, err = os.Open(`./storage_ftp_test.go`)
	assert.NoError(t, err)
	stat, err = fp.Stat()
	assert.NoError(t, err)
	err = c.Put(ctx, fp, `/sdb9/cloudbackuptest/1/2/storage_ftp_test.go`, stat.Size())
	assert.NoError(t, err)
	fp.Close()

	err = c.Restore(ctx, `/sdb9/cloudbackuptest`, `./testdata/webdav-restored`, nil)
	assert.NoError(t, err)
}
