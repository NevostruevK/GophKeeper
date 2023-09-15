package keeper_test

import (
	"context"
	"fmt"
	"log"
	"net"
	"testing"
	"time"

	"github.com/NevostruevK/GophKeeper/internal/api/grpc/client"
	"github.com/NevostruevK/GophKeeper/internal/api/grpc/server"
	"github.com/NevostruevK/GophKeeper/internal/api/grpc/server/auth"
	"github.com/NevostruevK/GophKeeper/internal/api/grpc/server/keeper"
	"github.com/NevostruevK/GophKeeper/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/NevostruevK/GophKeeper/internal/storage/memory"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

func TestGRPC_GetSpecs(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	server := startServer()
	defer server.Shutdown(ctx)

	client, err := newClient(bufDialer)

	require.NoError(t, err)
	defer client.Close()

	outRecords := []models.Record{
		*models.NewRecord(models.TEXT, "title text", []byte("text data")),
		*models.NewRecord(models.PAIR, "title pair", []byte("pair data")),
		*models.NewRecord(models.FILE, "title file", []byte("file data")),
		*models.NewRecord(models.CARD, "title card", []byte("card data")),
	}
	outSpecs := make([]models.Spec, len(outRecords))
	outTypes := make([]models.MType, len(outRecords))

	_, err = client.Auth.Register(ctx, "add_record_login", "some_password")
	require.NoError(t, err)
	for i, record := range outRecords {
		ds, err := client.Keeper.AddRecord(ctx, &record)
		require.NoError(t, err)
		outSpecs[i] = *record.ToSpec(*ds)
		outTypes[i] = record.Type
	}

	t.Run("get specs ok", func(t *testing.T) {
		inSpecs, err := client.Keeper.GetSpecs(ctx)
		require.NoError(t, err)
		assert.ElementsMatch(t, outSpecs, inSpecs)
	})
	t.Run("get specs of types ok", func(t *testing.T) {
		for i, spec := range outSpecs {
			inSpecs, err := client.Keeper.GetSpecsOfType(ctx, outTypes[i])
			require.NoError(t, err)
			require.Equal(t, 1, len(inSpecs))
			assert.Equal(t, spec, inSpecs[0])
		}
	})
	t.Run("get data ok", func(t *testing.T) {
		_, err := client.Auth.Register(ctx, "get data ok login", "some_password")
		require.NoError(t, err)
		ds, err := client.Keeper.AddRecord(ctx, &outRecords[0])
		require.NoError(t, err)
		data, err := client.Keeper.GetData(ctx, *ds)
		require.NoError(t, err)
		assert.Equal(t, outRecords[0].Data, data)
	})
	t.Run("get data err", func(t *testing.T) {
		ds := models.DataSpec{ID: uuid.New(), DataSize: 0}
		_, err := client.Keeper.GetData(ctx, ds)
		require.Error(t, err)
	})
}

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func bufDialer(ctx context.Context, address string) (net.Conn, error) {
	return lis.Dial()
}

func startServer() *server.Server {
	lis = bufconn.Listen(bufSize)
	dataStore := memory.NewDataStore()
	userStore := memory.NewUserStore()
	keeperServer := keeper.NewKeeperServer(dataStore)
	jwtManager := auth.NewJWTManager("test_secret_key", time.Hour)
	options, _ := server.NewServerOptions(jwtManager, false)
	authServer := auth.NewAuthServer(userStore, jwtManager)
	s := server.NewServer(authServer, keeperServer, options)
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatal(err)
		}
	}()
	return s
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
