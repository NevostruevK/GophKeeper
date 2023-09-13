package client

import (
	"context"
	"time"

	pb "github.com/NevostruevK/GophKeeper/proto"
	"google.golang.org/grpc"
)

const LoginTimeOut = time.Second

// AuthClient is a client to call authentication RPC
type AuthClient struct {
	conn        *grpc.ClientConn
	service     pb.AuthServiceClient
	interceptor *AuthInterceptor
}

// NewAuthClient returns a new auth client
func NewAuthClient(conn *grpc.ClientConn, interceptor *AuthInterceptor) *AuthClient {
	service := pb.NewAuthServiceClient(conn)
	return &AuthClient{conn, service, interceptor}
}

// Login login user and returns the access token
func (client *AuthClient) Login(ctx context.Context, login, password string) (string, error) {
	req := &pb.LoginRequest{
		Login:    login,
		Password: password,
	}

	res, err := client.service.Login(ctx, req)
	if err != nil {
		return "", err
	}
	accessToken := res.GetAccessToken()
	client.interceptor.SetToken(accessToken)
	return accessToken, nil
}

// Register register user and returns the access token
func (client *AuthClient) Register(ctx context.Context, login, password string) (string, error) {
	req := &pb.LoginRequest{
		Login:    login,
		Password: password,
	}

	res, err := client.service.Register(ctx, req)
	if err != nil {
		return "", err
	}

	accessToken := res.GetAccessToken()
	client.interceptor.SetToken(accessToken)
	return accessToken, nil
}

// Close free resources.
func (client *AuthClient) Close() error {
	return client.conn.Close()
}
