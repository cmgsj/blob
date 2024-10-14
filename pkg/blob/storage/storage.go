package storage

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/cmgsj/blob/pkg/blob/storage/driver"
	blobv1 "github.com/cmgsj/blob/pkg/gen/proto/blob/v1"
)

const ObjectPrefix = "blobs"

type Storage struct {
	driver       driver.Driver
	bucket       string
	objectPrefix string
}

func NewStorage(ctx context.Context, driver driver.Driver) (*Storage, error) {
	exists, err := driver.BucketExists(ctx, driver.Bucket())
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, fmt.Errorf("bucket %q does not exist", driver.Bucket())
	}

	objectPrefix := driver.ObjectPrefix()

	if objectPrefix == "" {
		objectPrefix = ObjectPrefix
	} else {
		objectPrefix = blobPath(objectPrefix, ObjectPrefix)
	}

	return &Storage{
		driver:       driver,
		bucket:       driver.Bucket(),
		objectPrefix: objectPrefix,
	}, nil
}

func (s *Storage) ListBlobs(ctx context.Context, path string) ([]string, error) {
	path = blobPath(s.objectPrefix, path)

	objects, err := s.driver.ListObjects(ctx, path)
	if err != nil {
		return nil, err
	}

	for i := range objects {
		objects[i] = blobPath(strings.TrimPrefix(objects[i], s.objectPrefix))
	}

	slices.Sort(objects)

	return objects, nil
}

func (s *Storage) GetBlob(ctx context.Context, name string) (*blobv1.Blob, error) {
	path := blobPath(s.objectPrefix, name)

	content, err := s.driver.GetObject(ctx, path)
	if err != nil {
		if s.driver.IsObjectNotFound(err) {
			return nil, ErrBlobNotFound
		}

		return nil, err
	}

	return &blobv1.Blob{
		Name:    name,
		Content: content,
	}, nil
}

func (s *Storage) WriteBlob(ctx context.Context, name string, content []byte) error {
	path := blobPath(s.objectPrefix, name)

	err := s.driver.WriteObject(ctx, path, content)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) RemoveBlob(ctx context.Context, name string) error {
	path := blobPath(s.objectPrefix, name)

	err := s.driver.RemoveObject(ctx, path)
	if err != nil {
		if s.driver.IsObjectNotFound(err) {
			return ErrBlobNotFound
		}

		return err
	}

	return nil
}

func blobPath(base string, elems ...string) string {
	var path []string

	base = strings.Trim(base, "/")

	if base != "" {
		path = append(path, base)
	}

	for _, elem := range elems {
		elem = strings.Trim(elem, "/")

		if elem != "" {
			path = append(path, elem)
		}
	}

	return strings.Join(path, "/")
}
