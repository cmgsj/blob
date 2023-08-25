package interceptors

import (
	"google.golang.org/grpc"
)

type FullInterceptor interface {
	ServerInterceptor
	ClientInterceptor
}

type ServerInterceptor interface {
	UnaryServerInterceptor() grpc.UnaryServerInterceptor
	StreamServerInterceptor() grpc.StreamServerInterceptor
}

type ClientInterceptor interface {
	UnaryClientInterceptor() grpc.UnaryClientInterceptor
	StreamClientInterceptor() grpc.StreamClientInterceptor
}

type interceptor struct {
	us grpc.UnaryServerInterceptor
	ss grpc.StreamServerInterceptor
	uc grpc.UnaryClientInterceptor
	sc grpc.StreamClientInterceptor
}

func (i interceptor) UnaryServerInterceptor() grpc.UnaryServerInterceptor   { return i.us }
func (i interceptor) StreamServerInterceptor() grpc.StreamServerInterceptor { return i.ss }
func (i interceptor) UnaryClientInterceptor() grpc.UnaryClientInterceptor   { return i.uc }
func (i interceptor) StreamClientInterceptor() grpc.StreamClientInterceptor { return i.sc }
