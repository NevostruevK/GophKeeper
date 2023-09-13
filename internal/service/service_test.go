package service

import (
	"context"
	"testing"
	"time"

	"github.com/NevostruevK/GophKeeper/internal/api/grpc/client"
	"github.com/NevostruevK/GophKeeper/internal/api/grpc/server"
	"github.com/NevostruevK/GophKeeper/internal/api/grpc/server/auth"
	"github.com/NevostruevK/GophKeeper/internal/api/grpc/server/keeper"
	"github.com/NevostruevK/GophKeeper/internal/config"
	"github.com/NevostruevK/GophKeeper/internal/models"
	"github.com/NevostruevK/GophKeeper/internal/storage/memory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestService(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	service, server, err := startService()
	require.NoError(t, err)
	defer func() {
		err = stopService(service, server)
		require.NoError(t, err)
	}()
	t.Run("Login error", func(t *testing.T) {
		user := models.NewUser("non_existent_login", "some_password")
		_, err := service.Login(ctx, user)
		assert.Error(t, err)
	})
	t.Run("Login ok", func(t *testing.T) {
		user := models.NewUser("test_login", "test_password")
		_, err := service.Register(ctx, user)
		require.NoError(t, err)
		_, err = service.Login(ctx, user)
		assert.NoError(t, err)
	})
	t.Run("Store Entry ok", func(t *testing.T) {
		user := models.NewUser("Store_Entry_ok_login", "some_password")
		_, err := service.Register(ctx, user)
		require.NoError(t, err)
		ds, err := service.StoreEntry(ctx, models.TEXT, "some title", models.NewText([]byte("some text")))
		assert.NoError(t, err)
		assert.NotNil(t, ds)
	})
	t.Run("Get Data ok", func(t *testing.T) {
		text := models.NewText([]byte("text for Get Data ok test"))
		user := models.NewUser("Get_data_ok_login", "some_password")
		_, err := service.Register(ctx, user)
		require.NoError(t, err)
		_, err = service.StoreEntry(ctx, models.TEXT, "some title", text)
		require.NoError(t, err)
		specs, err := service.LoadSpecs(ctx, models.TEXT)
		require.NoError(t, err)
		require.Equal(t, 1, len(specs))
		entry, err := service.GetData(ctx, specs[0])
		require.NoError(t, err)
		assert.Equal(t, text.String(), entry.String())
	})
	t.Run("Load Specs ok", func(t *testing.T) {
		out := make([]models.Spec, 0, 3)
		text := models.NewText([]byte("text for Load Specs ok"))
		pair := models.NewPair("pair_Load_Specs_ok_login", "pair_Load_Specs_ok_password")
		file := models.NewFile("file_Load_Specs_ok", []byte("file_Load_Specs_ok"))
		user := models.NewUser("Load_Specs_ok_login", "some_password")
		_, err := service.Register(ctx, user)
		require.NoError(t, err)
		ds, err := service.StoreEntry(ctx, models.TEXT, "some title for text", text)
		require.NoError(t, err)
		out = append(out, models.Spec{Type: models.TEXT, Title: "some title for text", DataSpec: *ds})
		ds, err = service.StoreEntry(ctx, models.PAIR, "some title for pair", pair)
		require.NoError(t, err)
		out = append(out, models.Spec{Type: models.PAIR, Title: "some title for pair", DataSpec: *ds})
		ds, err = service.StoreEntry(ctx, models.FILE, "some title for file", file)
		require.NoError(t, err)
		out = append(out, models.Spec{Type: models.FILE, Title: "some title for file", DataSpec: *ds})

		in, err := service.LoadSpecs(ctx, models.NOTIMPLEMENT)
		require.NoError(t, err)
		assert.ElementsMatch(t, out, in)
	})
}

func startService() (*Service, *server.Server, error) {
	cfg := config.Config{
		Address:       "127.0.0.1:8080",
		TokenKey:      "secretKeyForUserIdentification",
		EnableTLS:     false,
		TokenDuration: time.Hour,
	}
	dataStorage := memory.NewDataStore()
	userSoorage := memory.NewUserStore()
	keeperServer := keeper.NewKeeperServer(dataStorage)
	jwtManager := auth.NewJWTManager(cfg.TokenKey, cfg.TokenDuration)
	options, err := server.NewServerOptions(jwtManager, cfg.EnableTLS)
	if err != nil {
		return nil, nil, err
	}
	authServer := auth.NewAuthServer(userSoorage, jwtManager)
	server := server.NewServer(authServer, keeperServer, options)
	go server.Start(cfg.Address)
	client, err := client.NewClient(cfg.Address, cfg.EnableTLS)
	if err != nil {
		return nil, nil, err
	}
	service := NewService(client)
	return service, server, nil
}

func stopService(service *Service, server *server.Server) error {
	server.Shutdown(context.TODO())
	return service.client.Close()
}
