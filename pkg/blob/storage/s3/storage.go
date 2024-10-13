package s3

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/url"
	"slices"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"

	"github.com/cmgsj/blob/pkg/blob/storage"
	"github.com/cmgsj/blob/pkg/blob/storage/util"
	blobv1 "github.com/cmgsj/blob/pkg/gen/proto/blob/v1"
)

var _ storage.Storage = (*Storage)(nil)

type Storage struct {
	s3Client     *s3.Client
	bucket       string
	objectPrefix string
}

type StorageOptions struct {
	URI       string
	AccessKey string
	SecretKey string
}

func NewStorage(ctx context.Context, opts StorageOptions) (*Storage, error) {
	u, err := url.Parse(opts.URI)
	if err != nil {
		return nil, err
	}

	if u.Host == "" {
		return nil, fmt.Errorf("invalid google cloud storage uri %q: host is required", opts.URI)
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
			return nil, fmt.Errorf("invalid google cloud storage uri %q: bucket is required", opts.URI)
		}

		bucket = path[2]

		if len(path) > 3 {
			objectPrefix = strings.Join(path[3:], "/")
		}

		endpoint = fmt.Sprintf("%s://%s/%s/%s/", u.Scheme, u.Host, path[0], path[1])

	default:
		return nil, fmt.Errorf("invalid google cloud storage uri %q: unknown scheme", opts.URI)
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

	_, err = s3Client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		return nil, err
	}

	return &Storage{
		s3Client:     s3Client,
		bucket:       bucket,
		objectPrefix: util.BlobObjectPrefix(objectPrefix),
	}, nil
}

func (s *Storage) ListBlobs(ctx context.Context, path string) ([]string, error) {
	path = util.BlobPath(s.objectPrefix, path)

	objects, err := s.s3Client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String(s.bucket),
		Prefix: aws.String(path),
	})
	if err != nil {
		return nil, err
	}

	var blobNames []string

	for _, object := range objects.Contents {
		blobNames = append(blobNames, util.BlobPath(strings.TrimPrefix(*object.Key, s.objectPrefix)))
	}

	slices.Sort(blobNames)

	return blobNames, nil
}

func (s *Storage) GetBlob(ctx context.Context, name string) (*blobv1.Blob, error) {
	path := util.BlobPath(s.objectPrefix, name)

	object, err := s.s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		var e *types.NoSuchKey

		if errors.As(err, &e) {
			return nil, storage.ErrBlobNotFound
		}

		return nil, err
	}
	defer object.Body.Close()

	content, err := io.ReadAll(object.Body)
	if err != nil {
		return nil, err
	}

	return &blobv1.Blob{
		Name:      name,
		Content:   content,
		UpdatedAt: object.LastModified.Unix(),
	}, nil
}

func (s *Storage) WriteBlob(ctx context.Context, name string, content []byte) error {
	path := util.BlobPath(s.objectPrefix, name)

	_, err := s.s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(path),
		Body:   bytes.NewReader(content),
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) RemoveBlob(ctx context.Context, name string) error {
	path := util.BlobPath(s.objectPrefix, name)

	_, err := s.s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		var e *types.NoSuchKey

		if errors.As(err, &e) {
			return storage.ErrBlobNotFound
		}

		return err
	}

	return nil
}
