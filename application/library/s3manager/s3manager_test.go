package s3manager_test

import (
	"context"
	"os"
	"testing"

	"github.com/admpub/nging/v3/application/dbschema"
	"github.com/admpub/nging/v3/application/library/s3manager/s3client"
	"github.com/stretchr/testify/assert"
)

func TestStat(t *testing.T) {
	return
	cfg := &dbschema.NgingCloudStorage{
		Key:      os.Getenv(`S3_KEY`),
		Secret:   os.Getenv(`S3_SECRET`),
		Secure:   `Y`,
		Region:   os.Getenv(`S3_REGION`),
		Bucket:   os.Getenv(`S3_BUCKET`),
		Endpoint: os.Getenv(`S3_ENDPOINT`),
		Baseurl:  os.Getenv(`S3_BASEURL`),
	}
	mgr, err := s3client.New(cfg, 1024000)
	if err != nil {
		panic(err)
	}
	exists, err := mgr.Exists(context.Background(), ``)
	if err != nil {
		panic(err)
	}
	assert.Equal(t, false, exists)
}
