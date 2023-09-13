package auth_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/NevostruevK/GophKeeper/internal/api/grpc/server"
	"github.com/NevostruevK/GophKeeper/internal/api/grpc/server/auth"
	"github.com/NevostruevK/GophKeeper/internal/models"
	pb "github.com/NevostruevK/GophKeeper/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	"github.com/NevostruevK/GophKeeper/internal/storage/memory"
)

func Test_Register(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	address := "127.0.0.1:8082"

	conn, client, err := startClient(address)
	require.NoError(t, err)
	defer conn.Close()

	server := startServer(address)
	defer server.Shutdown(ctx)
	t.Run("register ok", func(t *testing.T) {
		user := models.NewUser(newLogPass())

		_, err := client.Register(
			ctx,
			&pb.LoginRequest{Login: user.Login, Password: user.Password},
		)
		require.NoError(t, err)
	})
	t.Run("register the same login err", func(t *testing.T) {
		user := models.NewUser(newLogPass())
		_, err := client.Register(
			ctx,
			&pb.LoginRequest{Login: user.Login, Password: user.Password},
		)
		require.NoError(t, err)
		_, err = client.Register(
			ctx,
			&pb.LoginRequest{Login: user.Login, Password: "any password"},
		)
		require.Error(t, err)
		e, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, codes.AlreadyExists, e.Code())
	})
}

func Test_Login(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	address := "127.0.0.1:8083"

	conn, client, err := startClient(address)
	require.NoError(t, err)
	defer conn.Close()

	server := startServer(address)
	defer server.Shutdown(ctx)
	t.Run("login ok", func(t *testing.T) {
		user := models.NewUser(newLogPass())

		_, err := client.Register(
			ctx,
			&pb.LoginRequest{Login: user.Login, Password: user.Password},
		)
		require.NoError(t, err)

		_, err = client.Login(
			ctx,
			&pb.LoginRequest{Login: user.Login, Password: user.Password},
		)
		require.NoError(t, err)
	})
	t.Run("wrong login err", func(t *testing.T) {
		_, err = client.Login(
			ctx,
			&pb.LoginRequest{Login: "wrong login", Password: "any password"},
		)
		require.Error(t, err)
		e, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, codes.NotFound, e.Code())
	})
}

func startServer(address string) *server.Server {
	jwtManager := auth.NewJWTManager("test_secret_key", time.Hour)
	st := memory.NewUserStore()
	authServer := auth.NewAuthServer(st, jwtManager)
	server := server.NewServer(authServer, nil, nil)

	go server.Start(address)
	return server
}

func startClient(address string) (*grpc.ClientConn, pb.AuthServiceClient, error) {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, err
	}
	return conn, pb.NewAuthServiceClient(conn), nil
}

func genLoginPassword() func() (string, string) {
	var num int
	return func() (string, string) {
		num++
		login := fmt.Sprintf("test_login_%d", num)
		password := fmt.Sprintf("test_password_%d", num)
		return login, password
	}
}

var newLogPass = genLoginPassword()
