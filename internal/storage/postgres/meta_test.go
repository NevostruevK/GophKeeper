package postgres_test

import (
	"context"
	"testing"

	"github.com/NevostruevK/GophKeeper/internal/models"
	storage "github.com/NevostruevK/GophKeeper/internal/storage/postgres"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStorage_AddMeta(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	st, err := newStorage(ctx)
	require.NoError(t, err)

	ids := idsDB{make([]uuid.UUID, 0, 4)}
	defer func() {
		require.NoError(t, deleteFromDB(ctx, st, ids.ids))
	}()

	t.Run("Add meta ok", func(t *testing.T) {
		user, err := addUser(ctx, st, &ids)
		require.NoError(t, err)
		meta := models.NewMeta(user.ID, models.TEXT, "some_title")
		id, err := st.AddMeta(ctx, meta)
		require.NoError(t, err)
		assert.Equal(t, meta.ID, id)
	})

	t.Run("Add meta for unknown user error", func(t *testing.T) {
		meta := models.NewMeta(uuid.New(), models.TEXT, "some_title")
		id, err := st.AddMeta(ctx, meta)
		assert.Error(t, err)
		assert.Equal(t, uuid.Nil, id)
	})
}

func TestStorage_GetMetas(t *testing.T) {
	titles := []string{"title1", "title2", "title3"}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	st, err := newStorage(ctx)
	require.NoError(t, err)

	ids := idsDB{make([]uuid.UUID, 0, 4)}
	defer func() {
		require.NoError(t, deleteFromDB(ctx, st, ids.ids))
	}()
	t.Run("Get titles ok", func(t *testing.T) {
		user, err := addUser(ctx, st, &ids)
		require.NoError(t, err)

		require.NoError(t, addMetaTitles(ctx, st, user.ID, titles, models.TEXT))

		metas, err := st.GetMetas(ctx, user.ID)
		require.NoError(t, err)

		assert.ElementsMatch(t, titles, getMetaTitles(metas))
	})
	t.Run("Get title for unknown user nil", func(t *testing.T) {
		metas, err := st.GetMetas(ctx, uuid.New())
		require.NoError(t, err)
		assert.Empty(t, metas)
	})
}

func TestStorage_GetMetasForType(t *testing.T) {
	titlesText := []string{"title_text_1", "title_text_2", "title_text_3"}
	titlesFile := []string{"title_file_1", "title_file_2", "title_file_3"}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	st, err := newStorage(ctx)
	require.NoError(t, err)

	ids := idsDB{make([]uuid.UUID, 0, 4)}
	defer func() {
		require.NoError(t, deleteFromDB(ctx, st, ids.ids))
	}()
	t.Run("Get titles ok", func(t *testing.T) {
		user, err := addUser(ctx, st, &ids)
		require.NoError(t, err)

		require.NoError(t, addMetaTitles(ctx, st, user.ID, titlesText, models.TEXT))
		require.NoError(t, addMetaTitles(ctx, st, user.ID, titlesFile, models.FILE))

		metas, err := st.GetMetasForType(ctx, user.ID, models.TEXT)
		require.NoError(t, err)

		assert.ElementsMatch(t, titlesText, getMetaTitles(metas))

		metas, err = st.GetMetasForType(ctx, user.ID, models.FILE)
		require.NoError(t, err)

		assert.ElementsMatch(t, titlesFile, getMetaTitles(metas))
	})
	t.Run("Get title for unknown user nil", func(t *testing.T) {
		metas, err := st.GetMetasForType(ctx, uuid.New(), models.FILE)
		require.NoError(t, err)
		assert.Empty(t, metas)
	})
}

func addMetaTitles(ctx context.Context, st *storage.Storage, userID uuid.UUID, titles []string, mType models.MType) error {
	for _, title := range titles {
		meta := models.NewMeta(userID, mType, title)
		_, err := st.AddMeta(ctx, meta)
		if err != nil {
			return err
		}
	}
	return nil
}

func getMetaTitles(metas []models.Meta) []string {
	titles := make([]string, len(metas))
	for i, meta := range metas {
		titles[i] = meta.Title
	}
	return titles
}
