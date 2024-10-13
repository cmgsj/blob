package gcs

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/url"
	"strings"

	gcs "cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"

	"github.com/cmgsj/blob/pkg/blob/storage/driver"
)

var _ driver.Driver = (*Driver)(nil)

type Driver struct {
	gcsClient    *gcs.Client
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
		return nil, fmt.Errorf("invalid google cloud storage uri %q: host is required", opts.URI)
	}

	var bucket string
	var objectPrefix string
	var clientOpts []option.ClientOption

	switch u.Scheme {
	case "gs":
		bucket = u.Host
		objectPrefix = u.Path

	case "http", "https":
		path := strings.Split(strings.Trim(u.Path, "/"), "/")

		if len(path) < 3 {
			return nil, fmt.Errorf("invalid google cloud storage uri %q: bucket is required", opts.URI)
		}

		bucket = path[2]

		if len(path) > 3 {
			objectPrefix = strings.Join(path[3:], "/")
		}

		endpoint := fmt.Sprintf("%s://%s/%s/%s/", u.Scheme, u.Host, path[0], path[1])

		clientOpts = append(clientOpts, option.WithEndpoint(endpoint))

	default:
		return nil, fmt.Errorf("invalid google cloud storage uri %q: unknown scheme", opts.URI)
	}

	gcsClient, err := gcs.NewClient(ctx, clientOpts...)
	if err != nil {
		return nil, err
	}

	return &Driver{
		gcsClient:    gcsClient,
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
	_, err := d.gcsClient.Bucket(bucket).Attrs(ctx)
	if err != nil {
		if errors.Is(err, gcs.ErrBucketNotExist) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (s *Driver) ListObjects(ctx context.Context, path string) ([]string, error) {
	it := s.gcsClient.Bucket(s.bucket).Objects(ctx, &gcs.Query{
		Prefix: path,
	})

	var objectNames []string

	for {
		attrs, err := it.Next()
		if err != nil {
			if err == iterator.Done {
				break
			}
			return nil, err
		}

		objectNames = append(objectNames, attrs.Name)
	}

	return objectNames, nil
}

func (s *Driver) GetObject(ctx context.Context, name string) ([]byte, int64, error) {
	reader, err := s.gcsClient.Bucket(s.bucket).Object(name).NewReader(ctx)
	if err != nil {
		return nil, 0, err
	}

	content, err := io.ReadAll(reader)
	if err != nil {
		return nil, 0, err
	}

	err = reader.Close()
	if err != nil {
		return nil, 0, err
	}

	attrs, err := s.gcsClient.Bucket(s.bucket).Object(name).Attrs(ctx)
	if err != nil {
		return nil, 0, err
	}

	return content, attrs.Updated.Unix(), nil
}

func (s *Driver) WriteObject(ctx context.Context, name string, content []byte) error {
	writer := s.gcsClient.Bucket(s.bucket).Object(name).NewWriter(ctx)

	_, err := writer.Write(content)
	if err != nil {
		return err
	}

	return writer.Close()
}

func (s *Driver) RemoveObject(ctx context.Context, name string) error {
	return s.gcsClient.Bucket(s.bucket).Object(name).Delete(ctx)
}

func (d *Driver) IsObjectNotFound(err error) bool {
	return errors.Is(err, gcs.ErrObjectNotExist)
}
