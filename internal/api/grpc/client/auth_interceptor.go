package client

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// AuthInterceptor is a client interceptor for authentication
type AuthInterceptor struct {
	authMethods map[string]bool
	accessToken string
}

// NewAuthInterceptor returns a new auth interceptor
func NewAuthInterceptor(
	authMethods map[string]bool,
	refreshDuration time.Duration,
) *AuthInterceptor {
	interceptor := &AuthInterceptor{
		authMethods: authMethods,
	}
	return interceptor
}

func (interceptor *AuthInterceptor) SetToken(accessToken string) {
	interceptor.accessToken = accessToken
}

// Unary returns a client interceptor to authenticate unary RPC
func (interceptor *AuthInterceptor) Unary() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		log.Printf("client --> unary interceptor: %s", method)

		if interceptor.authMethods[method] {
			ctx, cancel := context.WithTimeout(ctx, LoginTimeOut)
			defer cancel()
			return invoker(interceptor.attachToken(ctx), method, req, reply, cc, opts...)
		}

		ctx, cancel := context.WithTimeout(ctx, KeeperClientTimeOut)
		defer cancel()
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func (interceptor *AuthInterceptor) attachToken(ctx context.Context) context.Context {
	return metadata.AppendToOutgoingContext(ctx, "authorization", interceptor.accessToken)
}
