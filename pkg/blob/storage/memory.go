package storage

import (
	"context"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/cmgsj/blob/pkg/blob"
	blobv1 "github.com/cmgsj/blob/pkg/gen/proto/blob/v1"
)

func InMemory() blob.Storage {
	return &memory{
		blobs: make(map[string]*blobv1.Blob),
	}
}

type memory struct {
	mu    sync.RWMutex
	blobs map[string]*blobv1.Blob
}

func (s *memory) ListBlobs(ctx context.Context, path string) ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	path = strings.Trim(path, "/")

	var blobs []string

	for name := range s.blobs {
		if strings.HasPrefix(name, path) {
			blobs = append(blobs, name)
		}
	}

	sort.Strings(blobs)

	return blobs, nil
}

func (s *memory) GetBlob(ctx context.Context, name string) (*blobv1.Blob, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	b, ok := s.blobs[name]
	if !ok {
		return nil, blob.ErrBlobNotFound
	}

	return b, nil
}

func (s *memory) WriteBlob(ctx context.Context, name string, content []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	b, ok := s.blobs[name]
	if !ok {
		b = &blobv1.Blob{
			Name: name,
		}
	}

	b.Content = content
	b.UpdatedAt = time.Now().Unix()

	s.blobs[name] = b

	return nil
}

func (s *memory) RemoveBlob(ctx context.Context, name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.blobs[name]
	if !ok {
		return blob.ErrBlobNotFound
	}

	delete(s.blobs, name)

	return nil
}
