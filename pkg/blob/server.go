package blob

import (
	"context"

	blobv1 "github.com/cmgsj/blob/pkg/gen/proto/blob/v1"
)

var ServiceName = blobv1.BlobService_ServiceDesc.ServiceName

type Server struct {
	blobv1.UnimplementedBlobServiceServer
	storage Storage
}

func NewServer() blobv1.BlobServiceServer {
	return &Server{
		storage: NewInMemoryStorage(),
	}
}

func (s *Server) ListBlobs(ctx context.Context, req *blobv1.ListBlobsRequest) (*blobv1.ListBlobsResponse, error) {
	blobs, err := s.storage.ListBlobs(ctx, req.GetPath())
	if err != nil {
		return nil, err
	}

	return &blobv1.ListBlobsResponse{BlobNames: blobs, Count: uint64(len(blobs))}, nil
}

func (s *Server) GetBlob(ctx context.Context, req *blobv1.GetBlobRequest) (*blobv1.GetBlobResponse, error) {
	blob, err := s.storage.GetBlob(ctx, req.GetBlobName())
	if err != nil {
		return nil, err
	}

	return &blobv1.GetBlobResponse{Blob: blob}, nil
}

func (s *Server) WriteBlob(ctx context.Context, req *blobv1.WriteBlobRequest) (*blobv1.WriteBlobResponse, error) {
	err := s.storage.WriteBlob(ctx, req.GetBlobName(), req.GetContent())

	return &blobv1.WriteBlobResponse{}, err
}

func (s *Server) RemoveBlob(ctx context.Context, req *blobv1.RemoveBlobRequest) (*blobv1.RemoveBlobResponse, error) {
	err := s.storage.RemoveBlob(ctx, req.GetBlobName())

	return &blobv1.RemoveBlobResponse{}, err
}
