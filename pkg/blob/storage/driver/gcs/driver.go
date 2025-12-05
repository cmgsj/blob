package gcs

import (
	"context"
	"errors"
	"io"

	"cloud.google.com/go/storage"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"

	"github.com/cmgsj/blob/pkg/blob/storage/driver"
)

const DriverType = "gcs"

var _ driver.Driver = (*Driver)(nil)

type Driver struct {
	gcsClient    *storage.Client
	bucket       string
	objectPrefix string
}

type DriverOptions struct {
	URI string
}

func NewDriver(ctx context.Context, o DriverOptions) (*Driver, error) {
	uri, err := driver.ParseURI(DriverType, o.URI)
	if err != nil {
		return nil, err
	}

	credentials, err := google.FindDefaultCredentials(ctx)
	if err != nil {
		return nil, err
	}

	gcsClient, err := storage.NewClient(
		ctx,
		option.WithCredentials(credentials),
		option.WithEndpoint(uri.Host),
	)
	if err != nil {
		return nil, err
	}

	return &Driver{
		gcsClient:    gcsClient,
		bucket:       uri.Bucket,
		objectPrefix: uri.ObjectPrefix,
	}, nil
}

func (d *Driver) Bucket() string {
	return d.bucket
}

func (d *Driver) ObjectPrefix() string {
	return d.objectPrefix
}

func (d *Driver) BucketExists(ctx context.Context, bucket string) (bool, error) {
	_, err := d.gcsClient.Bucket(bucket).Attrs(ctx)
	if err != nil {
		if errors.Is(err, storage.ErrBucketNotExist) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (d *Driver) ListObjects(ctx context.Context, path string) ([]string, error) {
	it := d.gcsClient.Bucket(d.bucket).Objects(ctx, &storage.Query{
		Prefix: path,
	})

	var objectNames []string

	for {
		attrs, err := it.Next()
		if err != nil {
			if errors.Is(err, iterator.Done) {
				break
			}

			return nil, err
		}

		objectNames = append(objectNames, attrs.Name)
	}

	return objectNames, nil
}

func (d *Driver) GetObject(ctx context.Context, name string) ([]byte, error) {
	reader, err := d.gcsClient.Bucket(d.bucket).Object(name).NewReader(ctx)
	if err != nil {
		return nil, err
	}

	defer func() { _ = reader.Close() }()

	content, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	return content, nil
}

func (d *Driver) PutObject(ctx context.Context, name string, content []byte) error {
	writer := d.gcsClient.Bucket(d.bucket).Object(name).NewWriter(ctx)

	_, err := writer.Write(content)
	if err != nil {
		return err
	}

	return writer.Close()
}

func (d *Driver) DeleteObject(ctx context.Context, name string) error {
	return d.gcsClient.Bucket(d.bucket).Object(name).Delete(ctx)
}

func (d *Driver) IsObjectNotFound(err error) bool {
	return errors.Is(err, storage.ErrObjectNotExist)
}
