package auth_test

import (
	"context"
	"fmt"
	"log"
	"net"
	"testing"
	"time"

	"github.com/NevostruevK/GophKeeper/internal/api/grpc/server"
	"github.com/NevostruevK/GophKeeper/internal/api/grpc/server/auth"
	"github.com/NevostruevK/GophKeeper/internal/models"
	"github.com/NevostruevK/GophKeeper/internal/storage/memory"
	pb "github.com/NevostruevK/GophKeeper/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Test_Register(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	conn, client, err := startClient()
	require.NoError(t, err)
	defer conn.Close()

	server := startServer()
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

	server := startServer()
	defer server.Shutdown(ctx)

	conn, client, err := startClient()
	require.NoError(t, err)
	defer conn.Close()

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

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func bufDialer(ctx context.Context, address string) (net.Conn, error) {
	return lis.Dial()
}

func startServer() *server.Server {
	lis = bufconn.Listen(bufSize)
	jwtManager := auth.NewJWTManager("test_secret_key", time.Hour)
	st := memory.NewUserStore()
	authServer := auth.NewAuthServer(st, jwtManager)
	server := server.NewServer(authServer, nil, nil)
	go func() {
		if err := server.Serve(lis); err != nil {
			log.Fatal(err)
		}
	}()
	return server
}

func startClient() (*grpc.ClientConn, pb.AuthServiceClient, error) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
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
