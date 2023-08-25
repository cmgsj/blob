package authn

import (
	"context"

	"github.com/golang-jwt/jwt/v5"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/metadata"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	AuthHeader = "X-API-Key"
	ClaimsKey  = "claims"
)

var (
	ErrUnauthenticated = status.Error(codes.Unauthenticated, "unauthenticated")
)

type Authenticator interface {
	Authenticate(ctx context.Context) (context.Context, error)
}

func NewAuthenticator() Authenticator {
	return &authorizer{}
}

type Claims struct {
	jwt.RegisteredClaims
}

type authorizer struct {
	key []byte
}

func (a *authorizer) Authenticate(ctx context.Context) (context.Context, error) {
	md := metadata.ExtractIncoming(ctx)

	tokenString := md.Get(AuthHeader)

	// _, _ = jwt.WithAudience(""), jwt.WithIssuer("")
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, a.keyFunc)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, err
	}

	md.Set(ClaimsKey, claims.Subject)

	return md.ToIncoming(ctx), nil
}

func (a *authorizer) keyFunc(_ *jwt.Token) (interface{}, error) { return a.key, nil }
