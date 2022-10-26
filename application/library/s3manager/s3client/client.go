package s3client

import (
	"crypto/tls"
	"net/http"

	"github.com/admpub/nging/v5/application/dbschema"
	"github.com/admpub/nging/v5/application/library/s3manager"
	minio "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func Connect(m *dbschema.NgingCloudStorage) (client *minio.Client, err error) {
	isSecure := m.Secure == `Y`
	options := &minio.Options{
		Creds:  credentials.NewStaticV4(m.Key, m.Secret, ""),
		Secure: isSecure,
		Region: m.Region,
	}
	if isSecure {
		options.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}
	client, err = minio.New(m.Endpoint, options)
	return
}

func New(m *dbschema.NgingCloudStorage, editableMaxSize int) *s3manager.S3Manager {
	return s3manager.New(Connect, m, editableMaxSize)
}
