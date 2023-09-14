// Server gRPC.
package server

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"

	"log"
	"net"
	"time"

	"github.com/NevostruevK/GophKeeper/internal/api/grpc/server/auth"
	pb "github.com/NevostruevK/GophKeeper/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

// Server gRPC server.
type Server struct {
	*grpc.Server
}

// ServerOptions настройки gRPC сервера.
type ServerOptions []grpc.ServerOption

// NewServerOptions настройка менеджера аутентификации и включение TLS.
func NewServerOptions(jwtManager *auth.JWTManager, enableTLS bool) (ServerOptions, error) {
	options := ServerOptions{}
	interceptor := auth.NewAuthInterceptor(jwtManager)
	options = append(options, grpc.ChainUnaryInterceptor(logInterceptor, interceptor.Unary()))

	if enableTLS {
		tlsCredentials, err := loadTLSCredentials()
		if err != nil {
			return options, fmt.Errorf("cannot load TLS credentials: %w", err)
		}

		options = append(options, grpc.Creds(tlsCredentials))
	}
	return options, nil
}

// NewServer возвращае gRPC Server.
func NewServer(authServer pb.AuthServiceServer, keeperServer pb.KeeperServer, options ServerOptions) *Server {

	grpcServer := grpc.NewServer(options...)

	pb.RegisterAuthServiceServer(grpcServer, authServer)
	pb.RegisterKeeperServer(grpcServer, keeperServer)
	reflection.Register(grpcServer)

	return &Server{grpcServer}
}

// Start start gRPC server.
func (s *Server) Start(address string) error {
	listen, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("cannot start server: %w", err)
	}
	return s.Serve(listen)
}

// Shutdown stop gRPC server.
func (s *Server) Shutdown(ctx context.Context) {
	s.GracefulStop()
}

func logInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	start := time.Now()
	res, err := handler(ctx, req)

	if err != nil {
		log.Printf("[ERROR] server gRPC: %35s,  %v", info.FullMethod, err)
	} else {
		log.Printf("[INFO] server gRPC: %35s,  %v", info.FullMethod, time.Since(start))
	}
	return res, err
}

func loadTLSCredentials() (credentials.TransportCredentials, error) {
	const (
		serverCertFile   = "../../cert/server-cert.pem"
		serverKeyFile    = "../../cert/server-key.pem"
		clientCACertFile = "../../cert/ca-cert.pem"
	)
	pemClientCA, err := os.ReadFile(clientCACertFile)
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemClientCA) {
		return nil, fmt.Errorf("failed to add client CA's certificate")
	}

	serverCert, err := tls.LoadX509KeyPair(serverCertFile, serverKeyFile)
	if err != nil {
		return nil, err
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certPool,
	}

	return credentials.NewTLS(config), nil
}
