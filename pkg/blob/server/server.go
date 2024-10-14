package server

import (
	"context"

	"github.com/cmgsj/blob/pkg/blob/storage"
	blobv1 "github.com/cmgsj/blob/pkg/gen/proto/blob/v1"
)

var _ blobv1.BlobServiceServer = (*Server)(nil)

type Server struct {
	blobv1.UnimplementedBlobServiceServer
	storage *storage.Storage
}

func NewServer(storage *storage.Storage) *Server {
	return &Server{
		storage: storage,
	}
}

func (s *Server) ListBlobs(ctx context.Context, req *blobv1.ListBlobsRequest) (*blobv1.ListBlobsResponse, error) {
	blobNames, err := s.storage.ListBlobs(ctx, req.GetPath())
	if err != nil {
		return nil, storageError(err)
	}

	return &blobv1.ListBlobsResponse{
		BlobNames: blobNames,
	}, nil
}

func (s *Server) GetBlob(ctx context.Context, req *blobv1.GetBlobRequest) (*blobv1.GetBlobResponse, error) {
	blob, err := s.storage.GetBlob(ctx, req.GetBlobName())
	if err != nil {
		return nil, storageError(err)
	}

	return &blobv1.GetBlobResponse{
		Blob: blob,
	}, nil
}

func (s *Server) PutBlob(ctx context.Context, req *blobv1.PutBlobRequest) (*blobv1.PutBlobResponse, error) {
	err := s.storage.PutBlob(ctx, req.GetBlobName(), req.GetContent())
	if err != nil {
		return nil, storageError(err)
	}

	return &blobv1.PutBlobResponse{}, nil
}

func (s *Server) DeleteBlob(ctx context.Context, req *blobv1.DeleteBlobRequest) (*blobv1.DeleteBlobResponse, error) {
	err := s.storage.DeleteBlob(ctx, req.GetBlobName())
	if err != nil {
		return nil, storageError(err)
	}

	return &blobv1.DeleteBlobResponse{}, nil
}
