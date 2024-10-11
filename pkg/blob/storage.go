package blob

import (
	"context"
	"sort"
	"strings"
	"sync"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	blobv1 "github.com/cmgsj/blob/pkg/gen/proto/blob/v1"
)

var (
	ErrNotFound = status.Error(codes.NotFound, "blob not found")
)

type Storage interface {
	ListBlobs(ctx context.Context, path string) ([]string, error)
	GetBlob(ctx context.Context, name string) (*blobv1.Blob, error)
	WriteBlob(ctx context.Context, name string, content []byte) error
	RemoveBlob(ctx context.Context, name string) error
}

func NewInMemoryStorage() Storage {
	return &storage{
		m: make(map[string]*blobv1.Blob),
	}
}

type storage struct {
	mu sync.RWMutex
	m  map[string]*blobv1.Blob
}

func (s *storage) ListBlobs(ctx context.Context, path string) ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	path = strings.Trim(path, "/")
	var blobs []string
	for name := range s.m {
		if strings.HasPrefix(name, path) {
			blobs = append(blobs, name)
		}
	}
	sort.Strings(blobs)
	return blobs, nil
}

func (s *storage) GetBlob(ctx context.Context, name string) (*blobv1.Blob, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	blob, ok := s.m[name]
	if !ok {
		return nil, ErrNotFound
	}
	return blob, nil
}

func (s *storage) WriteBlob(ctx context.Context, name string, content []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	blob, ok := s.m[name]
	if !ok {
		blob = &blobv1.Blob{
			Name: name,
		}
	}
	blob.Content = content
	blob.UpdatedAt = time.Now().Unix()
	s.m[name] = blob
	return nil
}

func (s *storage) RemoveBlob(ctx context.Context, name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.m[name]
	if !ok {
		return ErrNotFound
	}
	delete(s.m, name)
	return nil
}
