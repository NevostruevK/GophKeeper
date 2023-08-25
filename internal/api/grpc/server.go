package grpc

import (
	"context"
	"log"
	"net"

	"github.com/NevostruevK/GophKeeper/internal/api/grpc/auth"
	pb "github.com/NevostruevK/GophKeeper/proto"
	"google.golang.org/grpc"
)

type Server struct {
	*grpc.Server
}

type ServerOptions []grpc.ServerOption

func NewServerOptions(jwtManager *auth.JWTManager, useTLS bool) ServerOptions {
	options := ServerOptions{}
	return options
}

func NewServer(authServer pb.AuthServiceServer, options ServerOptions) *Server {

	grpcServer := grpc.NewServer(options...)

	pb.RegisterAuthServiceServer(grpcServer, authServer)
	//	reflection.Register(grpcServer)

	return &Server{grpcServer}
}

func (s *Server) Start(address string) error {
	listen, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}
	return s.Serve(listen)
}

func (s *Server) Shutdown(ctx context.Context) {
	s.GracefulStop()
}
