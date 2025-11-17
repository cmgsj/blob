package blob

import (
	"context"

	"github.com/cmgsj/blob/pkg/blob/storage"
	blobv1 "github.com/cmgsj/blob/pkg/proto/blob/v1"
)

var _ blobv1.BlobServiceServer = (*Server)(nil)

type Server struct {
	storage *storage.Storage
}

func NewServer(storage *storage.Storage) *Server {
	return &Server{
		storage: storage,
	}
}

func (s *Server) ListBlobs(ctx context.Context, request *blobv1.ListBlobsRequest) (*blobv1.ListBlobsResponse, error) {
	names, err := s.storage.ListBlobs(ctx, request.GetPath())
	if err != nil {
		return nil, storageError(err)
	}

	response := &blobv1.ListBlobsResponse{}

	response.SetNames(names)

	return response, nil
}

func (s *Server) GetBlob(ctx context.Context, request *blobv1.GetBlobRequest) (*blobv1.GetBlobResponse, error) {
	blob, err := s.storage.GetBlob(ctx, request.GetName())
	if err != nil {
		return nil, storageError(err)
	}

	response := &blobv1.GetBlobResponse{}

	response.SetBlob(blob)

	return response, nil
}

func (s *Server) SetBlob(ctx context.Context, request *blobv1.SetBlobRequest) (*blobv1.SetBlobResponse, error) {
	err := s.storage.PutBlob(ctx, request.GetName(), request.GetContent())
	if err != nil {
		return nil, storageError(err)
	}

	response := &blobv1.SetBlobResponse{}

	return response, nil
}

func (s *Server) DeleteBlob(ctx context.Context, request *blobv1.DeleteBlobRequest) (*blobv1.DeleteBlobResponse, error) {
	err := s.storage.DeleteBlob(ctx, request.GetName())
	if err != nil {
		return nil, storageError(err)
	}

	response := &blobv1.DeleteBlobResponse{}

	return response, nil
}
