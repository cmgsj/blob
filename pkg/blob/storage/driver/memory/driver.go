package memory

import (
	"context"
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/cmgsj/blob/pkg/blob/storage/driver"
	blobv1 "github.com/cmgsj/blob/pkg/gen/proto/blob/v1"
)

var _ driver.Driver = (*Driver)(nil)

type Driver struct {
	m sync.Map
}

func NewDriver() *Driver {
	return &Driver{}
}

func (d *Driver) Bucket() string {
	return ""
}

func (d *Driver) ObjectPrefix() string {
	return ""
}

func (d *Driver) BucketExists(ctx context.Context, bucket string) (bool, error) {
	return true, nil
}

func (s *Driver) ListObjects(ctx context.Context, path string) ([]string, error) {
	var objectNames []string

	s.m.Range(func(key, value any) bool {
		name := key.(string)

		if strings.HasPrefix(name, path) {
			objectNames = append(objectNames, name)
		}

		return true
	})

	return objectNames, nil
}

func (s *Driver) GetObject(ctx context.Context, name string) ([]byte, int64, error) {
	v, ok := s.m.Load(name)
	if !ok {
		return nil, 0, errObjectNotFound
	}

	b := v.(*blobv1.Blob)

	return b.Content, b.UpdatedAt, nil
}

func (s *Driver) WriteObject(ctx context.Context, name string, content []byte) error {
	s.m.Store(name, &blobv1.Blob{
		Name:      name,
		Content:   content,
		UpdatedAt: time.Now().Unix(),
	})

	return nil
}

func (s *Driver) RemoveObject(ctx context.Context, name string) error {
	_, ok := s.m.Load(name)
	if !ok {
		return errObjectNotFound
	}

	s.m.Delete(name)

	return nil
}

var errObjectNotFound = errors.New("object not found")

func (s *Driver) IsObjectNotFound(err error) bool {
	return errors.Is(err, errObjectNotFound)
}
