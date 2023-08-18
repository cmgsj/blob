package interceptors

import (
	"context"
	"log/slog"

	"google.golang.org/grpc"
)

func ServerUnaryLogger(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	slog.Debug("inbound server unary grpc", "method", info.FullMethod)
	resp, err := handler(ctx, req)
	if err != nil {
		slog.Error("outbound server unary grpc", "method", info.FullMethod, "error", err)
	} else {
		slog.Debug("outbound server unary grpc", "method", info.FullMethod)
	}
	return resp, err
}

func ServerStreamLogger(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	slog.Debug("inbound server stream grpc", "method", info.FullMethod, "client", info.IsClientStream, "server", info.IsServerStream)
	err := handler(srv, ss)
	if err != nil {
		slog.Error("outbound server stream grpc", "method", info.FullMethod, "client", info.IsClientStream, "server", info.IsServerStream, "error", err)
	} else {
		slog.Debug("outbound server stream grpc", "method", info.FullMethod, "client", info.IsClientStream, "server", info.IsServerStream)
	}
	return err
}

func ClientUnaryLogger(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	slog.Debug("outbound client unary grpc", "method", method)
	err := invoker(ctx, method, req, reply, cc, opts...)
	if err != nil {
		slog.Debug("inbound client unary grpc", "method", method, "error", err)
	} else {
		slog.Debug("inbound client unary grpc", "method", method)
	}
	return err
}

func ClientStreamLogger(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	slog.Debug("outbound client stream grpc", "method", method, "client", desc.ClientStreams, "server", desc.ServerStreams)
	stream, err := streamer(ctx, desc, cc, method, opts...)
	if err != nil {
		slog.Error("inbound client stream grpc", "method", method, "client", desc.ClientStreams, "server", desc.ServerStreams, "error", err)
	} else {
		slog.Debug("inbound client stream grpc", "method", method, "client", desc.ClientStreams, "server", desc.ServerStreams)
	}
	return stream, err
}
