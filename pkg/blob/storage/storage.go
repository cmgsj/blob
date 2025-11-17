package storage

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/cmgsj/blob/pkg/blob/storage/driver"
	apiv1 "github.com/cmgsj/blob/pkg/proto/blob/api/v1"
)

const ObjectPrefix = "blobs"

type Storage struct {
	driver       driver.Driver
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

	return &Storage{
		driver:       driver,
		objectPrefix: joinBlobPath(driver.ObjectPrefix(), ObjectPrefix),
	}, nil
}

func (s *Storage) ListBlobs(ctx context.Context, path string) ([]string, error) {
	path = joinBlobPath(s.objectPrefix, path)

	objects, err := s.driver.ListObjects(ctx, path)
	if err != nil {
		return nil, err
	}

	for i := range objects {
		objects[i] = joinBlobPath(strings.TrimPrefix(objects[i], s.objectPrefix))
	}

	slices.Sort(objects)

	return objects, nil
}

func (s *Storage) GetBlob(ctx context.Context, name string) (*apiv1.Blob, error) {
	path := joinBlobPath(s.objectPrefix, name)

	content, err := s.driver.GetObject(ctx, path)
	if err != nil {
		if s.driver.IsObjectNotFound(err) {
			return nil, ErrBlobNotFound
		}

		return nil, err
	}

	blob := &apiv1.Blob{}

	blob.SetName(name)
	blob.SetContent(content)

	return blob, nil
}

func (s *Storage) PutBlob(ctx context.Context, name string, content []byte) error {
	path := joinBlobPath(s.objectPrefix, name)

	err := s.driver.PutObject(ctx, path, content)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) DeleteBlob(ctx context.Context, name string) error {
	path := joinBlobPath(s.objectPrefix, name)

	err := s.driver.DeleteObject(ctx, path)
	if err != nil {
		if s.driver.IsObjectNotFound(err) {
			return ErrBlobNotFound
		}

		return err
	}

	return nil
}

func joinBlobPath(base string, elems ...string) string {
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
