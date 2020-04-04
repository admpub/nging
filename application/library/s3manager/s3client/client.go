package s3client

import (
	"net/http"
	"crypto/tls"

	minio "github.com/minio/minio-go"
	"github.com/admpub/nging/application/library/s3manager"
	"github.com/admpub/nging/application/dbschema"
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
	if m.Secure != `Y` {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client.SetCustomTransport(tr)
	}
	return client, nil
}

func New(m *dbschema.NgingCloudStorage, editableMaxSize int64) (*s3manager.S3Manager,error) {
	client, err := Connect(m)
	if err != nil {
		return nil, err
	}
	mgr := s3manager.New(client, m.Bucket, editableMaxSize)
	return mgr, err
}