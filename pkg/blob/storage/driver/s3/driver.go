package s3

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/url"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"

	"github.com/cmgsj/blob/pkg/blob/storage/driver"
)

var _ driver.Driver = (*Driver)(nil)

type Driver struct {
	s3Client     *s3.Client
	bucket       string
	objectPrefix string
}

type DriverOptions struct {
	URI       string
	AccessKey string
	SecretKey string
}

func NewDriver(ctx context.Context, opts DriverOptions) (*Driver, error) {
	u, err := url.Parse(opts.URI)
	if err != nil {
		return nil, err
	}

	if u.Host == "" {
		return nil, fmt.Errorf("invalid s3 uri %q: host is required", opts.URI)
	}

	var bucket string
	var objectPrefix string
	var endpoint string

	switch u.Scheme {
	case "s3":
		bucket = u.Host
		objectPrefix = u.Path

	case "http", "https":
		path := strings.Split(strings.Trim(u.Path, "/"), "/")

		if len(path) < 3 {
			return nil, fmt.Errorf("invalid s3 uri %q: bucket is required", opts.URI)
		}

		bucket = path[2]

		if len(path) > 3 {
			objectPrefix = strings.Join(path[3:], "/")
		}

		endpoint = fmt.Sprintf("%s://%s/%s/%s/", u.Scheme, u.Host, path[0], path[1])

	default:
		return nil, fmt.Errorf("invalid s3 uri %q: unknown scheme", opts.URI)
	}

	s3Client := s3.New(s3.Options{
		BaseEndpoint: aws.String(endpoint),
		Credentials: aws.CredentialsProviderFunc(func(ctx context.Context) (aws.Credentials, error) {
			return aws.Credentials{
				AccessKeyID:     opts.AccessKey,
				SecretAccessKey: opts.SecretKey,
			}, nil
		}),
	})

	return &Driver{
		s3Client:     s3Client,
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
	_, err := d.s3Client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		var e *types.NoSuchBucket
		if errors.As(err, &e) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (d *Driver) ListObjects(ctx context.Context, path string) ([]string, error) {
	objects, err := d.s3Client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String(d.bucket),
		Prefix: aws.String(path),
	})
	if err != nil {
		return nil, err
	}

	var objectNames []string

	for _, object := range objects.Contents {
		objectNames = append(objectNames, *object.Key)
	}

	return objectNames, nil
}

func (d *Driver) GetObject(ctx context.Context, name string) ([]byte, int64, error) {
	object, err := d.s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(d.bucket),
		Key:    aws.String(name),
	})
	if err != nil {
		return nil, 0, err
	}
	defer object.Body.Close()

	content, err := io.ReadAll(object.Body)
	if err != nil {
		return nil, 0, err
	}

	return content, object.LastModified.Unix(), nil
}

func (d *Driver) WriteObject(ctx context.Context, name string, content []byte) error {
	_, err := d.s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(d.bucket),
		Key:    aws.String(name),
		Body:   bytes.NewReader(content),
	})
	return err
}

func (d *Driver) RemoveObject(ctx context.Context, name string) error {
	_, err := d.s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(d.bucket),
		Key:    aws.String(name),
	})
	return err
}

func (d *Driver) IsObjectNotFound(err error) bool {
	var e *types.NoSuchKey
	return errors.As(err, &e)
}
