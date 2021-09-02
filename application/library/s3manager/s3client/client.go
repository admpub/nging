package s3client

import (
	"crypto/tls"
	"net/http"

	"github.com/admpub/nging/v3/application/dbschema"
	"github.com/admpub/nging/v3/application/library/s3manager"
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

func New(m *dbschema.NgingCloudStorage, editableMaxSize int64) (*s3manager.S3Manager, error) {
	client, err := Connect(m)
	if err != nil {
		return nil, err
	}
	mgr := s3manager.New(client, m, editableMaxSize)
	return mgr, err
}
