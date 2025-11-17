package memory

import (
	"context"
	"errors"
	"strings"
	"sync"

	"github.com/cmgsj/blob/pkg/blob/storage/driver"
)

var errObjectNotFound = errors.New("object not found")

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

func (d *Driver) ListObjects(ctx context.Context, path string) ([]string, error) {
	var objectNames []string

	d.m.Range(func(key, value any) bool {
		name := key.(string)

		if strings.HasPrefix(name, path) {
			objectNames = append(objectNames, name)
		}

		return true
	})

	return objectNames, nil
}

func (d *Driver) GetObject(ctx context.Context, name string) ([]byte, error) {
	v, ok := d.m.Load(name)
	if !ok {
		return nil, errObjectNotFound
	}

	return v.([]byte), nil
}

func (d *Driver) PutObject(ctx context.Context, name string, content []byte) error {
	d.m.Store(name, content)

	return nil
}

func (d *Driver) DeleteObject(ctx context.Context, name string) error {
	_, ok := d.m.Load(name)
	if !ok {
		return errObjectNotFound
	}

	d.m.Delete(name)

	return nil
}

func (d *Driver) IsObjectNotFound(err error) bool {
	return errors.Is(err, errObjectNotFound)
}
