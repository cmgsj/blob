package handler

import (
	"context"

	"connectrpc.com/connect"

	"github.com/cmgsj/blob/pkg/blob/storage"
	apiv1 "github.com/cmgsj/blob/pkg/proto/blob/api/v1"
	"github.com/cmgsj/blob/pkg/proto/blob/api/v1/apiv1connect"
)

var _ apiv1connect.BlobServiceHandler = (*Handler)(nil)

type Handler struct {
	storage *storage.Storage
}

func NewHandler(storage *storage.Storage) *Handler {
	return &Handler{
		storage: storage,
	}
}

func (h *Handler) ListBlobs(ctx context.Context, request *connect.Request[apiv1.ListBlobsRequest]) (*connect.Response[apiv1.ListBlobsResponse], error) {
	names, err := h.storage.ListBlobs(ctx, request.Msg.GetPath())
	if err != nil {
		return nil, storageError(err)
	}

	response := &apiv1.ListBlobsResponse{}

	response.SetNames(names)

	return connect.NewResponse(response), nil
}

func (h *Handler) GetBlob(ctx context.Context, request *connect.Request[apiv1.GetBlobRequest]) (*connect.Response[apiv1.GetBlobResponse], error) {
	blob, err := h.storage.GetBlob(ctx, request.Msg.GetName())
	if err != nil {
		return nil, storageError(err)
	}

	response := &apiv1.GetBlobResponse{}

	response.SetBlob(blob)

	return connect.NewResponse(response), nil
}

func (h *Handler) SetBlob(ctx context.Context, request *connect.Request[apiv1.SetBlobRequest]) (*connect.Response[apiv1.SetBlobResponse], error) {
	err := h.storage.PutBlob(ctx, request.Msg.GetName(), request.Msg.GetContent())
	if err != nil {
		return nil, storageError(err)
	}

	response := &apiv1.SetBlobResponse{}

	return connect.NewResponse(response), nil
}

func (h *Handler) DeleteBlob(ctx context.Context, request *connect.Request[apiv1.DeleteBlobRequest]) (*connect.Response[apiv1.DeleteBlobResponse], error) {
	err := h.storage.DeleteBlob(ctx, request.Msg.GetName())
	if err != nil {
		return nil, storageError(err)
	}

	response := &apiv1.DeleteBlobResponse{}

	return connect.NewResponse(response), nil
}
