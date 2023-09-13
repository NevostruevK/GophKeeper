package keeper_test

import (
	"context"
	"testing"
	"time"

	"github.com/NevostruevK/GophKeeper/internal/api/grpc/client"
	"github.com/NevostruevK/GophKeeper/internal/api/grpc/server"
	"github.com/NevostruevK/GophKeeper/internal/api/grpc/server/auth"
	"github.com/NevostruevK/GophKeeper/internal/api/grpc/server/keeper"
	"github.com/NevostruevK/GophKeeper/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/NevostruevK/GophKeeper/internal/storage/memory"
)

func TestGRPC_GetSpecs(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	address := "127.0.0.1:8081"
	client, err := client.NewClient(address, false)
	require.NoError(t, err)
	defer client.Close()

	outRecords := []models.Record{
		*models.NewRecord(models.TEXT, "title text", []byte("text data")),
		*models.NewRecord(models.PAIR, "title pair", []byte("pair data")),
		*models.NewRecord(models.FILE, "title file", []byte("file data")),
		*models.NewRecord(models.CARD, "title card", []byte("card data")),
	}

	t.Run("get specs ok", func(t *testing.T) {
		server := startServer(address)
		defer server.Shutdown(ctx)
		_, err := client.Auth.Register(ctx, "get_specs_ok_login", "some_password")
		require.NoError(t, err)
		outSpecs := make([]models.Spec, len(outRecords))
		for i, record := range outRecords {
			ds, err := client.Keeper.AddRecord(ctx, &record)
			require.NoError(t, err)
			outSpecs[i] = *record.ToSpec(*ds)
		}
		inSpecs, err := client.Keeper.GetSpecs(ctx)
		require.NoError(t, err)
		assert.ElementsMatch(t, outSpecs, inSpecs)
	})
	t.Run("get specs of types ok", func(t *testing.T) {
		server := startServer(address)
		defer server.Shutdown(ctx)
		_, err := client.Auth.Register(ctx, "get_specs_types_ok_login", "some_password")
		require.NoError(t, err)
		outSpecs := make([]models.Spec, len(outRecords))
		outTypes := make([]models.MType, len(outRecords))
		for i, record := range outRecords {
			ds, err := client.Keeper.AddRecord(ctx, &record)
			require.NoError(t, err)
			outSpecs[i] = *record.ToSpec(*ds)
			outTypes[i] = record.Type
		}
		for i, spec := range outSpecs {
			inSpecs, err := client.Keeper.GetSpecsOfType(ctx, outTypes[i])
			require.NoError(t, err)
			require.Equal(t, 1, len(inSpecs))
			assert.Equal(t, spec, inSpecs[0])
		}
	})
	t.Run("get data ok", func(t *testing.T) {
		server := startServer(address)
		defer server.Shutdown(ctx)
		_, err := client.Auth.Register(ctx, "get_specs_types_ok_login", "some_password")
		require.NoError(t, err)
		ds, err := client.Keeper.AddRecord(ctx, &outRecords[0])
		require.NoError(t, err)
		data, err := client.Keeper.GetData(ctx, *ds)
		require.NoError(t, err)
		assert.Equal(t, outRecords[0].Data, data)
	})
}

func startServer(address string) *server.Server {
	dataStore := memory.NewDataStore()
	userStore := memory.NewUserStore()
	keeperServer := keeper.NewKeeperServer(dataStore)
	jwtManager := auth.NewJWTManager("test_secret_key", time.Hour)
	options, _ := server.NewServerOptions(jwtManager, false)
	authServer := auth.NewAuthServer(userStore, jwtManager)
	s := server.NewServer(authServer, keeperServer, options)
	go s.Start(address)
	return s
}
