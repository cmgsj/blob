package interceptors

import (
	"github.com/cmgsj/blob/pkg/authn"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
)

func NewAuthenticator() ServerInterceptor {
	authenticator := authn.NewAuthenticator()
	return interceptor{
		us: auth.UnaryServerInterceptor(authenticator.Authenticate),
		ss: auth.StreamServerInterceptor(authenticator.Authenticate),
	}
}
