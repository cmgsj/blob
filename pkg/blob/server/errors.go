package server

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/cmgsj/blob/pkg/blob/storage"
)

var (
	ErrBlobNotFound = status.Error(codes.NotFound, "blob not found")
)

func storageError(err error) error {
	if errors.Is(err, storage.ErrBlobNotFound) {
		return ErrBlobNotFound
	}

	return err
}
