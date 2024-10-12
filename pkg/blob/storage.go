package blob

import (
	"context"

	blobv1 "github.com/cmgsj/blob/pkg/gen/proto/blob/v1"
)

type Storage interface {
	ListBlobs(ctx context.Context, path string) ([]string, error)
	GetBlob(ctx context.Context, name string) (*blobv1.Blob, error)
	WriteBlob(ctx context.Context, name string, content []byte) error
	RemoveBlob(ctx context.Context, name string) error
}
