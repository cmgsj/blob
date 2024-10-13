package memory

import (
	"context"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/cmgsj/blob/pkg/blob/storage"
	"github.com/cmgsj/blob/pkg/blob/storage/util"
	blobv1 "github.com/cmgsj/blob/pkg/gen/proto/blob/v1"
)

var _ storage.Storage = (*Storage)(nil)

type Storage struct {
	m sync.Map
}

func NewStorage() *Storage {
	return &Storage{}
}

func (s *Storage) ListBlobs(ctx context.Context, path string) ([]string, error) {
	prefix := util.BlobNamePrefix(path)

	var blobNames []string

	s.m.Range(func(key, value any) bool {
		name := key.(string)

		if strings.HasPrefix(name, prefix) {
			blobNames = append(blobNames, name)
		}

		return true
	})

	slices.Sort(blobNames)

	return blobNames, nil
}

func (s *Storage) GetBlob(ctx context.Context, name string) (*blobv1.Blob, error) {
	v, ok := s.m.Load(name)
	if !ok {
		return nil, storage.ErrBlobNotFound
	}

	return v.(*blobv1.Blob), nil
}

func (s *Storage) WriteBlob(ctx context.Context, name string, content []byte) error {
	s.m.Store(name, &blobv1.Blob{
		Name:      name,
		Content:   content,
		UpdatedAt: time.Now().Unix(),
	})

	return nil
}

func (s *Storage) RemoveBlob(ctx context.Context, name string) error {
	_, ok := s.m.Load(name)
	if !ok {
		return storage.ErrBlobNotFound
	}

	s.m.Delete(name)

	return nil
}

func IsStorage(uri string) bool {
	switch uri {
	case ":memory:":
		return true
	}

	return false
}
