package storage

import (
	"context"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/cmgsj/blob/pkg/blob"
	blobv1 "github.com/cmgsj/blob/pkg/gen/proto/blob/v1"
)

var _ blob.Storage = (*InMemoryStorage)(nil)

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{}
}

type InMemoryStorage struct {
	blobs sync.Map
}

func (s *InMemoryStorage) ListBlobs(ctx context.Context, path string) ([]string, error) {
	var blobs []string

	prefix := cleanBlobPrefix(path)

	s.blobs.Range(func(key, value any) bool {
		name := key.(string)

		if strings.HasPrefix(name, prefix) {
			blobs = append(blobs, name)
		}

		return true
	})

	slices.Sort(blobs)

	return blobs, nil
}

func (s *InMemoryStorage) GetBlob(ctx context.Context, name string) (*blobv1.Blob, error) {
	v, ok := s.blobs.Load(name)
	if !ok {
		return nil, blob.ErrBlobNotFound
	}

	return v.(*blobv1.Blob), nil
}

func (s *InMemoryStorage) WriteBlob(ctx context.Context, name string, content []byte) error {
	b := &blobv1.Blob{
		Name:      name,
		Content:   content,
		UpdatedAt: time.Now().Unix(),
	}

	s.blobs.Store(name, b)

	return nil
}

func (s *InMemoryStorage) RemoveBlob(ctx context.Context, name string) error {
	_, ok := s.blobs.Load(name)
	if !ok {
		return blob.ErrBlobNotFound
	}

	s.blobs.Delete(name)

	return nil
}
