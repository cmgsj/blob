package googlecloud

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/url"
	"slices"
	"strings"

	gcs "cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"

	"github.com/cmgsj/blob/pkg/blob/storage"
	"github.com/cmgsj/blob/pkg/blob/storage/util"
	blobv1 "github.com/cmgsj/blob/pkg/gen/proto/blob/v1"
)

const BlobStorageFolder = "blobs"

var _ storage.Storage = (*Storage)(nil)

type Storage struct {
	bucket *gcs.BucketHandle
	folder string
}

func NewStorage(ctx context.Context, uri string, opts ...option.ClientOption) (*Storage, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}

	if u.Host == "" {
		return nil, fmt.Errorf("invalid google cloud storage uri %q: host is required", uri)
	}

	var bucketName string
	var bucketPrefix string

	switch u.Scheme {
	case "gs":
		bucketName = u.Host
		bucketPrefix = u.Path

	case "http", "https":
		path := strings.Split(strings.Trim(u.Path, "/"), "/")

		if len(path) < 3 {
			return nil, fmt.Errorf("invalid google cloud storage uri %q: bucket is required", uri)
		}

		bucketName = path[2]

		if len(path) > 3 {
			bucketPrefix = strings.Join(path[3:], "/")
		}

		endpoint := fmt.Sprintf("%s://%s/%s/%s/", u.Scheme, u.Host, path[0], path[1])

		opts = append(opts, option.WithEndpoint(endpoint))

	default:
		return nil, fmt.Errorf("invalid google cloud storage uri %q: unknown scheme", uri)
	}

	gcsClient, err := gcs.NewClient(ctx, opts...)
	if err != nil {
		return nil, err
	}

	bucket := gcsClient.Bucket(bucketName)

	_, err = bucket.Attrs(ctx)
	if err != nil {
		return nil, err
	}

	folder := BlobStorageFolder

	if bucketPrefix != "" {
		folder = util.BlobPath(bucketPrefix, folder)
	}

	return &Storage{
		bucket: bucket,
		folder: folder,
	}, nil
}

func (s *Storage) ListBlobs(ctx context.Context, path string) ([]string, error) {
	path = util.BlobPath(s.folder, path)

	it := s.bucket.Objects(ctx, &gcs.Query{
		Prefix: path,
	})

	var blobNames []string

	for {
		attrs, err := it.Next()
		if err != nil {
			if err == iterator.Done {
				break
			}
			return nil, err
		}

		blobNames = append(blobNames, util.BlobPath(strings.TrimPrefix(attrs.Name, s.folder)))
	}

	slices.Sort(blobNames)

	return blobNames, nil
}

func (s *Storage) GetBlob(ctx context.Context, name string) (*blobv1.Blob, error) {
	path := util.BlobPath(s.folder, name)

	reader, err := s.bucket.Object(path).NewReader(ctx)
	if err != nil {
		if errors.Is(err, gcs.ErrObjectNotExist) {
			return nil, storage.ErrBlobNotFound
		}

		return nil, err
	}

	content, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	err = reader.Close()
	if err != nil {
		return nil, err
	}

	attrs, err := s.bucket.Object(path).Attrs(ctx)
	if err != nil {
		if errors.Is(err, gcs.ErrObjectNotExist) {
			return nil, storage.ErrBlobNotFound
		}

		return nil, err
	}

	return &blobv1.Blob{
		Name:      name,
		Content:   content,
		UpdatedAt: attrs.Updated.Unix(),
	}, nil
}

func (s *Storage) WriteBlob(ctx context.Context, name string, content []byte) error {
	path := util.BlobPath(s.folder, name)

	writer := s.bucket.Object(path).NewWriter(ctx)

	_, err := writer.Write(content)
	if err != nil {
		return err
	}

	err = writer.Close()
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) RemoveBlob(ctx context.Context, name string) error {
	path := util.BlobPath(s.folder, name)

	err := s.bucket.Object(path).Delete(ctx)
	if err != nil {
		if errors.Is(err, gcs.ErrObjectNotExist) {
			return storage.ErrBlobNotFound
		}

		return nil
	}

	return nil
}

func IsStorage(uri string) bool {
	u, err := url.Parse(uri)
	if err != nil {
		return false
	}

	switch u.Scheme {
	case "gs", "http", "https":
		return true
	}

	return false
}
