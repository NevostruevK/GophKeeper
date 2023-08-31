package auth

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type ctxKey string

const KeyUserID ctxKey = "userID"

var authMethods = map[string]bool{
	"/gophKeeper.AuthService/Login":    true,
	"/gophKeeper.AuthService/Register": true,
}

// AuthInterceptor is a server interceptor for authentication and authorization
type AuthInterceptor struct {
	jwtManager *JWTManager
}

// NewAuthInterceptor returns a new auth interceptor
func NewAuthInterceptor(jwtManager *JWTManager) *AuthInterceptor {
	return &AuthInterceptor{jwtManager}
}

func (interceptor *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		if authMethods[info.FullMethod] {
			return handler(ctx, req)
		}
		id, err := interceptor.authorize(ctx)
		if err != nil {
			return nil, err
		}
		ctx = context.WithValue(ctx, KeyUserID, id)
		return handler(ctx, req)
	}
}

func (interceptor *AuthInterceptor) authorize(ctx context.Context) (id uuid.UUID, err error) {

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return id, status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	values := md["authorization"]
	if len(values) == 0 {
		return id, status.Errorf(codes.Unauthenticated, "authorization token is not provided")
	}

	accessToken := values[0]
	claims, err := interceptor.jwtManager.Verify(accessToken)
	if err != nil {
		return id, status.Errorf(codes.Unauthenticated, "access token is invalid: %v", err)
	}

	return claims, nil
}
