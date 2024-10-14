package minio

import (
	"bytes"
	"context"
	"io"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/cmgsj/blob/pkg/blob/storage/driver"
)

var _ driver.Driver = (*Driver)(nil)

type Driver struct {
	minioClient  *minio.Client
	bucket       string
	objectPrefix string
}

type DriverOptions struct {
	Address      string
	AccessKey    string
	SecretKey    string
	Bucket       string
	ObjectPrefix string
	Secure       bool
}

func NewDriver(ctx context.Context, opts DriverOptions) (*Driver, error) {
	minioClient, err := minio.New(opts.Address, &minio.Options{
		Creds:  credentials.NewStatic(opts.AccessKey, opts.SecretKey, "", credentials.SignatureDefault),
		Secure: opts.Secure,
	})
	if err != nil {
		return nil, err
	}

	return &Driver{
		minioClient:  minioClient,
		bucket:       opts.Bucket,
		objectPrefix: opts.ObjectPrefix,
	}, nil
}

func (d *Driver) Bucket() string {
	return d.bucket
}

func (d *Driver) ObjectPrefix() string {
	return d.objectPrefix
}

func (d *Driver) BucketExists(ctx context.Context, bucket string) (bool, error) {
	return d.minioClient.BucketExists(ctx, bucket)
}

func (d *Driver) ListObjects(ctx context.Context, path string) ([]string, error) {
	objects := d.minioClient.ListObjects(ctx, d.bucket, minio.ListObjectsOptions{
		Prefix:    path,
		Recursive: true,
	})

	var objectNames []string

	for object := range objects {
		objectNames = append(objectNames, object.Key)
	}

	return objectNames, nil
}

func (d *Driver) GetObject(ctx context.Context, name string) ([]byte, error) {
	object, err := d.minioClient.GetObject(ctx, d.bucket, name, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	defer object.Close()

	content, err := io.ReadAll(object)
	if err != nil {
		return nil, err
	}

	return content, nil
}

func (d *Driver) WriteObject(ctx context.Context, name string, content []byte) error {
	_, err := d.minioClient.PutObject(ctx, d.bucket, name, bytes.NewReader(content), int64(len(content)), minio.PutObjectOptions{})
	return err
}

func (d *Driver) RemoveObject(ctx context.Context, name string) error {
	return d.minioClient.RemoveObject(ctx, d.bucket, name, minio.RemoveObjectOptions{})
}

func (d *Driver) IsObjectNotFound(err error) bool {
	return minio.ToErrorResponse(err).Code == "NoSuchKey"
}
