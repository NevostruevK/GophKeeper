package tui

import (
	"context"
	"testing"

	"github.com/NevostruevK/GophKeeper/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadForm_setChoice(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	service, server, client, err := startService()
	require.NoError(t, err)
	defer func() {
		err = stopService(server, client)
		require.NoError(t, err)
	}()
	ch := make(chan any)
	go func() {
		err := run(ctx, service, "", "", ch)
		require.NoError(t, err)
	}()
	defer app.Stop()
	<-ch
	loadFormTest := newLoadForm()
	t.Run("Unauthenticated err", func(t *testing.T) {
		loadFormTest.choice.SetCurrentOption(0)
		actual := messager.GetText(false)
		expected := errUnauthenticated
		assert.Equal(t, expected, actual)
	})
	t.Run("there are no entries ok", func(t *testing.T) {
		_, err := srv.Register(ctx, models.NewUser("there are no entries ok", "there are no entries ok"))
		require.NoError(t, err)
		loadFormTest.choice.SetCurrentOption(0)
		actual := messager.GetText(false)
		expected := mesThereAreNoEntries
		assert.Equal(t, expected, actual)
	})
	t.Run("found 1 specs ok", func(t *testing.T) {
		_, err := srv.Register(ctx, models.NewUser("found 1 specs ok", "found 1 specs ok"))
		require.NoError(t, err)
		_, err = srv.StoreEntry(ctx, models.TEXT, "some title", models.NewText([]byte("some text")))
		require.NoError(t, err)
		loadFormTest.choice.SetCurrentOption(0)
		actual := messager.GetText(false)
		expected := "found 1 specs\n"
		assert.Equal(t, expected, actual)
	})
}
