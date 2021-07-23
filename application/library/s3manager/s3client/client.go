package s3client

import (
	"crypto/tls"
	"net/http"

	"github.com/admpub/nging/v3/application/dbschema"
	"github.com/admpub/nging/v3/application/library/s3manager"
	minio "github.com/minio/minio-go"
)

func Connect(m *dbschema.NgingCloudStorage) (client *minio.Client, err error) {
	isSecure := m.Secure == `Y`
	if len(m.Region) == 0 {
		client, err = minio.New(m.Endpoint, m.Key, m.Secret, isSecure)
	} else {
		client, err = minio.NewWithRegion(m.Endpoint, m.Key, m.Secret, isSecure, m.Region)
	}
	if err != nil {
		return client, err
	}
	if isSecure {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client.SetCustomTransport(tr)
	}
	return client, nil
}

func New(m *dbschema.NgingCloudStorage, editableMaxSize int64) (*s3manager.S3Manager, error) {
	client, err := Connect(m)
	if err != nil {
		return nil, err
	}
	mgr := s3manager.New(client, m, editableMaxSize)
	return mgr, err
}
