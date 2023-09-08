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

	server := startServer(address)
	defer server.Shutdown(ctx)

	client, err := client.NewClient(address, false)
	require.NoError(t, err)
	defer client.Close()
	t.Run("get specs ok", func(t *testing.T) {
		_, err := client.Auth.Register(ctx, "some_login", "some_password")
		require.NoError(t, err)
		titles := []string{"title1", "title2", "title3"}
		texts := [][]byte{[]byte("text1"), []byte("text2"), []byte("text3")}
		records := make([]models.Record, len(titles))
		outSpecs := make([]models.Spec, len(titles))
		for i, title := range titles {
			records[i] = *models.NewRecord(models.TEXT, title, texts[i])
			ds, err := client.Keeper.AddRecord(ctx, &records[i])
			require.NoError(t, err)
			outSpecs[i] = *records[i].ToSpec(*ds)
		}
		inSpecs, err := client.Keeper.GetSpecs(ctx)
		require.NoError(t, err)
		assert.ElementsMatch(t, outSpecs, inSpecs)
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
