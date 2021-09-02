package s3manager

import (
	"io"

	minio "github.com/minio/minio-go/v7"
)

// S3 is a client to interact with S3 storage.
type S3 interface {
	GetObject(bucketName, objectName string, opts minio.GetObjectOptions) (*minio.Object, error)
	ListBuckets() ([]minio.BucketInfo, error)
	ListObjectsV2(bucketName, objectPrefix string, recursive bool, doneCh <-chan struct{}) <-chan minio.ObjectInfo
	MakeBucket(bucketName, location string) error
	PutObject(bucketName, objectName string, reader io.Reader, objectSize int64, opts minio.PutObjectOptions) (int64, error)
	RemoveBucket(bucketName string) error
	RemoveObject(bucketName, objectName string) error
}
