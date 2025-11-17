package blob

import (
	"errors"

	"connectrpc.com/connect"

	"github.com/cmgsj/blob/pkg/blob/storage"
)

func storageError(err error) error {
	if errors.Is(err, storage.ErrBlobNotFound) {
		return connect.NewError(connect.CodeNotFound, err)
	}

	return err
}
