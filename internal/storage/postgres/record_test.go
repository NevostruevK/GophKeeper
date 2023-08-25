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

func TestStorage_AddRecord(t *testing.T) {
	data := models.Data([]byte("test data"))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	st, err := newStorage(ctx)
	require.NoError(t, err)

	ids := idsDB{make([]uuid.UUID, 0, 4)}
	defer func() {
		require.NoError(t, deleteFromDB(ctx, st, ids.ids))
	}()

	t.Run("Add record ok", func(t *testing.T) {
		user, err := addUser(ctx, st, &ids)
		require.NoError(t, err)
		r := models.NewRecord(models.TEXT, "some_title", data, nil)
		ds, err := st.AddRecord(ctx, user.ID, r)
		require.NoError(t, err)
		require.NotNil(t, ds)
		assert.NotEqual(t, ds.ID, uuid.Nil)
		assert.NotEqual(t, user.ID, ds.ID)
	})

	t.Run("Add record for unknown user error", func(t *testing.T) {
		r := models.NewRecord(models.TEXT, "some_title", data, nil)
		_, err := st.AddRecord(ctx, uuid.New(), r)
		assert.Error(t, err)
	})
}

func TestStorage_GetRecords(t *testing.T) {
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

		require.NoError(t, addRecordsTitles(ctx, st, user.ID, titles, models.TEXT))

		specs, err := st.GetSpecs(ctx, user.ID)
		require.NoError(t, err)

		assert.ElementsMatch(t, titles, getSpecsTitles(specs))
	})
	t.Run("Get title for unknown user nil", func(t *testing.T) {
		specs, err := st.GetSpecs(ctx, uuid.New())
		require.NoError(t, err)
		assert.Empty(t, specs)
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

		require.NoError(t, addRecordsTitles(ctx, st, user.ID, titlesText, models.TEXT))
		require.NoError(t, addRecordsTitles(ctx, st, user.ID, titlesFile, models.FILE))

		specs, err := st.GetSpecsOfType(ctx, user.ID, models.TEXT)
		require.NoError(t, err)

		assert.ElementsMatch(t, titlesText, getSpecsTitles(specs))

		specs, err = st.GetSpecsOfType(ctx, user.ID, models.FILE)
		require.NoError(t, err)

		assert.ElementsMatch(t, titlesFile, getSpecsTitles(specs))
	})
	t.Run("Get title for unknown user nil", func(t *testing.T) {
		specs, err := st.GetSpecsOfType(ctx, uuid.New(), models.FILE)
		require.NoError(t, err)
		assert.Empty(t, specs)
	})
}

func addRecordsTitles(ctx context.Context, st *storage.Storage, userID uuid.UUID, titles []string, mType models.MType) error {
	data := models.Data([]byte("test data"))
	for _, title := range titles {
		r := models.NewRecord(mType, title, data, nil)
		_, err := st.AddRecord(ctx, userID, r)
		if err != nil {
			return err
		}
	}
	return nil
}

func getSpecsTitles(specs []models.Spec) []string {
	titles := make([]string, len(specs))
	for i, r := range specs {
		titles[i] = r.Title
	}
	return titles
}
