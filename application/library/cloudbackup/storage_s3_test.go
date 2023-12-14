package cloudbackup

import (
	"os"
	"testing"

	"github.com/admpub/godotenv"
	"github.com/admpub/nging/v5/application/dbschema"
	"github.com/stretchr/testify/assert"
	"github.com/webx-top/echo/defaults"
)

func TestStorageS3(t *testing.T) {
	storageS3Test = true
	godotenv.Load()
	c := NewStorageS3(dbschema.NgingCloudStorage{
		Type:     `oss`,
		Key:      os.Getenv(`STORAGE_S3_KEY`),
		Secret:   os.Getenv(`STORAGE_S3_SECRET`),
		Bucket:   `nging`,
		Endpoint: `oss-cn-hongkong.aliyuncs.com`,
		Region:   `cn-hongkong`,
		Secure:   `Y`,
	}).(*StorageS3)
	err := c.Connect()
	assert.NoError(t, err)
	defer c.Close()
	ctx := defaults.NewMockContext()
	fp, err := os.Open(`./storage_ftp_test.go`)
	assert.NoError(t, err)
	defer fp.Close()
	stat, err := fp.Stat()
	assert.NoError(t, err)
	err = c.Put(ctx, fp, `/cloudbackuptest/1/2/3/storage_ftp_test.go`, stat.Size())
	assert.NoError(t, err)
	err = c.Put(ctx, fp, `/cloudbackuptest/1/2/storage_ftp_test.go`, stat.Size())
	assert.NoError(t, err)

	err = c.Restore(ctx, `/cloudbackuptest`, `./testdata/s3-restored`, nil)
	assert.NoError(t, err)
}
