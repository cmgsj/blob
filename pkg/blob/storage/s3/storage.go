package s3

import (
	"context"

	"github.com/cmgsj/blob/pkg/blob/storage"
	"github.com/cmgsj/blob/pkg/blob/storage/minio"
)

var _ storage.Storage = (*Storage)(nil)

type Storage = minio.Storage

type StorageOptions = minio.StorageOptions

func NewStorage(ctx context.Context, opts StorageOptions) (*Storage, error) {
	return minio.NewStorage(ctx, opts)
}
