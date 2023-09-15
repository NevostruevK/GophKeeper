package service

import (
	"context"
	"testing"

	"github.com/NevostruevK/GophKeeper/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestService(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	service, _, server, err := StartTestService()
	require.NoError(t, err)
	defer func() {
		err = StopTestService(service, server)
		require.NoError(t, err)
	}()

	t.Run("Store Entry err", func(t *testing.T) {
		ds, err := service.StoreEntry(ctx, models.TEXT, "some title", models.NewText([]byte("some text")))
		require.Error(t, err)
		assert.Nil(t, ds)
	})
	t.Run("Load Specs err", func(t *testing.T) {
		user := models.NewUser("Load_Specs_err_login", "some_password")
		_, err := service.Register(ctx, user)
		require.NoError(t, err)
		in, err := service.LoadSpecs(ctx, models.TEXT)
		require.NoError(t, err)
		assert.Equal(t, 0, len(in))
		in, err = service.LoadSpecs(ctx, models.NOTIMPLEMENT)
		require.NoError(t, err)
		assert.Equal(t, 0, len(in))
	})
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
	t.Run("Get Data err (not found)", func(t *testing.T) {
		user := models.NewUser("Get Data err (not found)", "Get Data err")
		_, err := service.Register(ctx, user)
		require.NoError(t, err)
		spec := models.NewSpec(models.CARD, "some title")
		_, err = service.GetData(ctx, *spec)
		require.Error(t, err)
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
