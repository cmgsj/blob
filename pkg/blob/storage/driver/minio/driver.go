package minio

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/url"
	"strings"

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
	URI string
}

func NewDriver(ctx context.Context, opts DriverOptions) (*Driver, error) {
	u, err := url.Parse(opts.URI)
	if err != nil {
		return nil, err
	}

	if u.Host == "" {
		return nil, fmt.Errorf("invalid minio uri %q: host is required", opts.URI)
	}

	var (
		bucket       string
		objectPrefix string
		endpoint     string
	)

	switch u.Scheme {
	case "minio":
		bucket = u.Host
		objectPrefix = u.Path

	case "http", "https":
		path := strings.Split(strings.Trim(u.Path, "/"), "/")

		if len(path) < 3 {
			return nil, fmt.Errorf("invalid minio uri %q: bucket is required", opts.URI)
		}

		bucket = path[2]

		if len(path) > 3 {
			objectPrefix = strings.Join(path[3:], "/")
		}

		endpoint = fmt.Sprintf("%s://%s/%s/%s/", u.Scheme, u.Host, path[0], path[1])

	default:
		return nil, fmt.Errorf("invalid minio uri %q: unknown scheme", opts.URI)
	}

	creds := credentials.NewEnvMinio()

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds: creds,
	})
	if err != nil {
		return nil, err
	}

	return &Driver{
		minioClient:  minioClient,
		bucket:       bucket,
		objectPrefix: objectPrefix,
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

	objectNames := make([]string, 0, len(objects))

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

	defer func() { _ = object.Close() }()

	content, err := io.ReadAll(object)
	if err != nil {
		return nil, err
	}

	return content, nil
}

func (d *Driver) PutObject(ctx context.Context, name string, content []byte) error {
	_, err := d.minioClient.PutObject(ctx, d.bucket, name, bytes.NewReader(content), int64(len(content)), minio.PutObjectOptions{})

	return err
}

func (d *Driver) DeleteObject(ctx context.Context, name string) error {
	return d.minioClient.RemoveObject(ctx, d.bucket, name, minio.RemoveObjectOptions{})
}

func (d *Driver) IsObjectNotFound(err error) bool {
	return minio.ToErrorResponse(err).Code == "NoSuchKey"
}
