package blob

import (
	"context"

	"github.com/cmgsj/blob/pkg/blob/storage"
	apiv1 "github.com/cmgsj/blob/pkg/proto/blob/api/v1"
)

var _ apiv1.BlobServiceServer = (*Server)(nil)

type Server struct {
	storage *storage.Storage
}

func NewServer(storage *storage.Storage) *Server {
	return &Server{
		storage: storage,
	}
}

func (s *Server) ListBlobs(ctx context.Context, request *apiv1.ListBlobsRequest) (*apiv1.ListBlobsResponse, error) {
	names, err := s.storage.ListBlobs(ctx, request.GetPath())
	if err != nil {
		return nil, storageError(err)
	}

	response := &apiv1.ListBlobsResponse{}

	response.SetNames(names)

	return response, nil
}

func (s *Server) GetBlob(ctx context.Context, request *apiv1.GetBlobRequest) (*apiv1.GetBlobResponse, error) {
	blob, err := s.storage.GetBlob(ctx, request.GetName())
	if err != nil {
		return nil, storageError(err)
	}

	response := &apiv1.GetBlobResponse{}

	response.SetBlob(blob)

	return response, nil
}

func (s *Server) SetBlob(ctx context.Context, request *apiv1.SetBlobRequest) (*apiv1.SetBlobResponse, error) {
	err := s.storage.PutBlob(ctx, request.GetName(), request.GetContent())
	if err != nil {
		return nil, storageError(err)
	}

	response := &apiv1.SetBlobResponse{}

	return response, nil
}

func (s *Server) DeleteBlob(ctx context.Context, request *apiv1.DeleteBlobRequest) (*apiv1.DeleteBlobResponse, error) {
	err := s.storage.DeleteBlob(ctx, request.GetName())
	if err != nil {
		return nil, storageError(err)
	}

	response := &apiv1.DeleteBlobResponse{}

	return response, nil
}
