package driver

import "context"

type Driver interface {
	Bucket() string
	ObjectPrefix() string
	BucketExists(ctx context.Context, bucket string) (bool, error)
	ListObjects(ctx context.Context, path string) ([]string, error)
	GetObject(ctx context.Context, name string) ([]byte, int64, error)
	WriteObject(ctx context.Context, name string, content []byte) error
	RemoveObject(ctx context.Context, name string) error
	IsObjectNotFound(err error) bool
}
