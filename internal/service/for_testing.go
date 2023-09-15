package service

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/NevostruevK/GophKeeper/internal/api/grpc/client"
	"github.com/NevostruevK/GophKeeper/internal/api/grpc/server"
	"github.com/NevostruevK/GophKeeper/internal/api/grpc/server/auth"
	"github.com/NevostruevK/GophKeeper/internal/api/grpc/server/keeper"
	"github.com/NevostruevK/GophKeeper/internal/storage/memory"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

func StartTestService() (*Service, *client.Client, *server.Server, error) {
	const bufSize = 1024 * 1024

	var lis *bufconn.Listener

	bufDialer := func(ctx context.Context, address string) (net.Conn, error) {
		return lis.Dial()
	}

	lis = bufconn.Listen(bufSize)
	dataStorage := memory.NewDataStore()
	userSoorage := memory.NewUserStore()
	keeperServer := keeper.NewKeeperServer(dataStorage)
	jwtManager := auth.NewJWTManager("test_secret_key", time.Hour)
	options, err := server.NewServerOptions(jwtManager, false)
	if err != nil {
		return nil, nil, nil, err
	}
	authServer := auth.NewAuthServer(userSoorage, jwtManager)
	server := server.NewServer(authServer, keeperServer, options)
	go func() {
		if err := server.Serve(lis); err != nil {
			log.Fatal(err)
		}
	}()
	client, err := newClient(bufDialer)
	if err != nil {
		return nil, nil, nil, err
	}
	service := NewService(client)
	return service, client, server, nil
}

func StopTestService(service *Service, server *server.Server) error {
	server.Shutdown(context.TODO())
	return service.client.Close()
}

func newClient(dial func(ctx context.Context, address string) (net.Conn, error)) (*client.Client, error) {
	transportOption := grpc.WithTransportCredentials(insecure.NewCredentials())

	ctx := context.Background()
	conn, err := grpc.DialContext(
		ctx,
		"bufnet",
		grpc.WithContextDialer(dial),
		transportOption)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, fmt.Errorf("cannot dial server: %v", err)
	}

	interceptor := client.NewAuthInterceptor(client.AuthMethods(), time.Hour)
	authClient := client.NewAuthClient(conn, interceptor)

	conn, err = grpc.DialContext(
		ctx,
		"bufnet",
		grpc.WithContextDialer(dial),
		transportOption,
		grpc.WithUnaryInterceptor(interceptor.Unary()))
	if err != nil {
		return nil, err
	}

	KeeperClient := client.NewKeeperClient(conn)

	return &client.Client{Auth: authClient, Keeper: KeeperClient}, nil
}
