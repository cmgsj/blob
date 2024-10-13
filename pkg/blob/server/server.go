package server

import (
	"context"

	"github.com/cmgsj/blob/pkg/blob/storage"
	blobv1 "github.com/cmgsj/blob/pkg/gen/proto/blob/v1"
)

type Server struct {
	blobv1.UnimplementedBlobServiceServer
	storage storage.Storage
}

func NewServer(storage storage.Storage) blobv1.BlobServiceServer {
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
		Count:     uint64(len(blobNames)),
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

func (s *Server) WriteBlob(ctx context.Context, req *blobv1.WriteBlobRequest) (*blobv1.WriteBlobResponse, error) {
	err := s.storage.WriteBlob(ctx, req.GetBlobName(), req.GetContent())
	if err != nil {
		return nil, storageError(err)
	}

	return &blobv1.WriteBlobResponse{}, nil
}

func (s *Server) RemoveBlob(ctx context.Context, req *blobv1.RemoveBlobRequest) (*blobv1.RemoveBlobResponse, error) {
	err := s.storage.RemoveBlob(ctx, req.GetBlobName())
	if err != nil {
		return nil, storageError(err)
	}

	return &blobv1.RemoveBlobResponse{}, nil
}
