package blob

import (
	"context"

	"connectrpc.com/connect"

	"github.com/cmgsj/blob/pkg/blob/storage"
	apiv1 "github.com/cmgsj/blob/pkg/proto/blob/api/v1"
	"github.com/cmgsj/blob/pkg/proto/blob/api/v1/apiv1connect"
)

var _ apiv1connect.BlobServiceHandler = (*Service)(nil)

type Service struct {
	storage *storage.Storage
}

func NewService(storage *storage.Storage) *Service {
	return &Service{
		storage: storage,
	}
}

func (s *Service) ListBlobs(ctx context.Context, request *connect.Request[apiv1.ListBlobsRequest]) (*connect.Response[apiv1.ListBlobsResponse], error) {
	names, err := s.storage.ListBlobs(ctx, request.Msg.GetPath())
	if err != nil {
		return nil, storageError(err)
	}

	response := &apiv1.ListBlobsResponse{}

	response.SetNames(names)

	return connect.NewResponse(response), nil
}

func (s *Service) GetBlob(ctx context.Context, request *connect.Request[apiv1.GetBlobRequest]) (*connect.Response[apiv1.GetBlobResponse], error) {
	blob, err := s.storage.GetBlob(ctx, request.Msg.GetName())
	if err != nil {
		return nil, storageError(err)
	}

	response := &apiv1.GetBlobResponse{}

	response.SetBlob(blob)

	return connect.NewResponse(response), nil
}

func (s *Service) SetBlob(ctx context.Context, request *connect.Request[apiv1.SetBlobRequest]) (*connect.Response[apiv1.SetBlobResponse], error) {
	err := s.storage.PutBlob(ctx, request.Msg.GetName(), request.Msg.GetContent())
	if err != nil {
		return nil, storageError(err)
	}

	response := &apiv1.SetBlobResponse{}

	return connect.NewResponse(response), nil
}

func (s *Service) DeleteBlob(ctx context.Context, request *connect.Request[apiv1.DeleteBlobRequest]) (*connect.Response[apiv1.DeleteBlobResponse], error) {
	err := s.storage.DeleteBlob(ctx, request.Msg.GetName())
	if err != nil {
		return nil, storageError(err)
	}

	response := &apiv1.DeleteBlobResponse{}

	return connect.NewResponse(response), nil
}
