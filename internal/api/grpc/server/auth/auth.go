// Server gRPC for authentification users.
package auth

import (
	"context"
	"errors"

	"github.com/NevostruevK/GophKeeper/internal/models"
	"github.com/NevostruevK/GophKeeper/internal/storage"
	pb "github.com/NevostruevK/GophKeeper/proto"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// AuthServer gRPC authentification server.
type AuthServer struct {
	pb.UnimplementedAuthServiceServer
	storage    UserStore
	jwtManager *JWTManager
}

// NewAuthServer returns AuthServer.
func NewAuthServer(userStore UserStore, jwtManager *JWTManager) pb.AuthServiceServer {
	return &AuthServer{
		storage:    userStore,
		jwtManager: jwtManager,
	}
}

// Login login user and returns the access token.
func (s *AuthServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	user := models.NewUser(req.Login, req.Password)

	id, err := s.storage.GetUser(ctx, *user)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) || errors.Is(err, storage.ErrWrongPassword) {
			return nil, status.Errorf(codes.NotFound, "incorrect login/password")
		}
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return s.genTokenAndSend(id)
}

// Register register user and returns the access token.
func (s *AuthServer) Register(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	user := models.NewUser(req.Login, req.Password)
	id, err := s.storage.AddUser(ctx, *user)
	if err != nil {
		if errors.Is(err, storage.ErrDuplicateLogin) {
			return nil, status.Errorf(codes.AlreadyExists, "user with this login already exists")
		}
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return s.genTokenAndSend(id)
}

func (s *AuthServer) genTokenAndSend(id uuid.UUID) (*pb.LoginResponse, error) {
	token, err := s.jwtManager.Generate(id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot generate access token")
	}
	return &pb.LoginResponse{AccessToken: token}, nil
}
