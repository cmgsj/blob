package storage

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/url"
	"slices"
	"strings"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"

	"github.com/cmgsj/blob/pkg/blob"
	blobv1 "github.com/cmgsj/blob/pkg/gen/proto/blob/v1"
)

const GoogleCloudStorageFolderName = "blobs"

var _ blob.Storage = (*GoogleCloudStorage)(nil)

func NewGoogleCloudStorage(ctx context.Context, uri string, opts ...option.ClientOption) (*GoogleCloudStorage, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}

	if u.Scheme != "gs" || u.Host == "" {
		return nil, fmt.Errorf("invalid uri %q", uri)
	}

	storageClient, err := storage.NewClient(ctx, opts...)
	if err != nil {
		return nil, err
	}

	bucket := storageClient.Bucket(u.Host)

	_, err = bucket.Attrs(ctx)
	if err != nil {
		return nil, err
	}

	folder := GoogleCloudStorageFolderName

	if u.Path != "" {
		folder = joinBlobPrefix(u.Path, folder)
	}

	return &GoogleCloudStorage{
		bucket: bucket,
		folder: folder,
	}, nil
}

type GoogleCloudStorage struct {
	bucket *storage.BucketHandle
	folder string
}

func (s *GoogleCloudStorage) ListBlobs(ctx context.Context, path string) ([]string, error) {
	var blobs []string

	it := s.bucket.Objects(ctx, &storage.Query{
		Prefix: joinBlobPrefix(s.folder, path),
	})

	for {
		attrs, err := it.Next()
		if err != nil {
			if err == iterator.Done {
				break
			}
			return nil, googleCloudStorageError(err)
		}

		if strings.HasPrefix(attrs.Name, path) {
			blobs = append(blobs, attrs.Name)
		}
	}

	slices.Sort(blobs)

	return blobs, nil
}

func (s *GoogleCloudStorage) GetBlob(ctx context.Context, name string) (*blobv1.Blob, error) {
	reader, err := s.bucket.Object(name).NewReader(ctx)
	if err != nil {
		return nil, googleCloudStorageError(err)
	}

	content, err := io.ReadAll(reader)
	if err != nil {
		return nil, googleCloudStorageError(err)
	}

	err = reader.Close()
	if err != nil {
		return nil, googleCloudStorageError(err)
	}

	attrs, err := s.bucket.Object(name).Attrs(ctx)
	if err != nil {
		return nil, googleCloudStorageError(err)
	}

	b := &blobv1.Blob{
		Name:      name,
		Content:   content,
		UpdatedAt: attrs.Updated.Unix(),
	}

	return b, nil
}

func (s *GoogleCloudStorage) WriteBlob(ctx context.Context, name string, content []byte) error {
	writer := s.bucket.Object(name).NewWriter(ctx)

	_, err := writer.Write(content)
	if err != nil {
		return googleCloudStorageError(err)
	}

	err = writer.Close()
	if err != nil {
		return googleCloudStorageError(err)
	}

	return nil
}

func (s *GoogleCloudStorage) RemoveBlob(ctx context.Context, name string) error {
	err := s.bucket.Object(name).Delete(ctx)
	if err != nil {
		return googleCloudStorageError(err)
	}

	return nil
}

func googleCloudStorageError(err error) error {
	if errors.Is(err, storage.ErrObjectNotExist) {
		return blob.ErrBlobNotFound
	}

	return err
}
