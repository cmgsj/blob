package interceptors

import "google.golang.org/grpc"

type FullInterceptor struct {
	clientInterceptor
	serverInterceptor
}

type ClientInterceptor struct {
	clientInterceptor
}

type clientInterceptor struct {
	unary  grpc.UnaryClientInterceptor
	stream grpc.StreamClientInterceptor
}

func (i clientInterceptor) UnaryClientInterceptor() grpc.UnaryClientInterceptor   { return i.unary }
func (i clientInterceptor) StreamClientInterceptor() grpc.StreamClientInterceptor { return i.stream }

type ServerInterceptor struct {
	serverInterceptor
}

type serverInterceptor struct {
	unary  grpc.UnaryServerInterceptor
	stream grpc.StreamServerInterceptor
}

func (i serverInterceptor) UnaryServerInterceptor() grpc.UnaryServerInterceptor   { return i.unary }
func (i serverInterceptor) StreamServerInterceptor() grpc.StreamServerInterceptor { return i.stream }
