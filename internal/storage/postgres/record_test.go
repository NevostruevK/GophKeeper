package postgres_test

import (
	"context"
	"testing"

	"github.com/NevostruevK/GophKeeper/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStorage_AddRecord(t *testing.T) {
	data := models.Data([]byte("test data"))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	t.Run("Add record ok", func(t *testing.T) {
		id, err := addUser(ctx)
		require.NoError(t, err)
		r := models.NewRecord(models.TEXT, "some_title", data)
		ds, err := testStorage.AddRecord(ctx, id, r)
		require.NoError(t, err)
		require.NotNil(t, ds)
		assert.NotEqual(t, ds.ID, uuid.Nil)
		assert.NotEqual(t, id, ds.ID)
	})

	t.Run("Add record for unknown user error", func(t *testing.T) {
		r := models.NewRecord(models.TEXT, "some_title", data)
		_, err := testStorage.AddRecord(ctx, uuid.New(), r)
		assert.Error(t, err)
	})
}

func TestStorage_GetRecords(t *testing.T) {
	titles := []string{"title1", "title2", "title3"}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	t.Run("Get titles ok", func(t *testing.T) {
		id, err := addUser(ctx)
		require.NoError(t, err)

		require.NoError(t, addRecordsTitles(ctx, id, titles, models.TEXT))

		specs, err := testStorage.GetSpecs(ctx, id)
		require.NoError(t, err)

		assert.ElementsMatch(t, titles, getSpecsTitles(specs))
	})
	t.Run("Get title for unknown user nil", func(t *testing.T) {
		specs, err := testStorage.GetSpecs(ctx, uuid.New())
		require.NoError(t, err)
		assert.Empty(t, specs)
	})
}

func TestStorage_GetMetasForType(t *testing.T) {
	titlesText := []string{"title_text_1", "title_text_2", "title_text_3"}
	titlesFile := []string{"title_file_1", "title_file_2", "title_file_3"}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	t.Run("Get titles ok", func(t *testing.T) {
		id, err := addUser(ctx)
		require.NoError(t, err)

		require.NoError(t, addRecordsTitles(ctx, id, titlesText, models.TEXT))
		require.NoError(t, addRecordsTitles(ctx, id, titlesFile, models.FILE))

		specs, err := testStorage.GetSpecsOfType(ctx, id, models.TEXT)
		require.NoError(t, err)

		assert.ElementsMatch(t, titlesText, getSpecsTitles(specs))

		specs, err = testStorage.GetSpecsOfType(ctx, id, models.FILE)
		require.NoError(t, err)

		assert.ElementsMatch(t, titlesFile, getSpecsTitles(specs))
	})
	t.Run("Get title for unknown user nil", func(t *testing.T) {
		specs, err := testStorage.GetSpecsOfType(ctx, uuid.New(), models.FILE)
		require.NoError(t, err)
		assert.Empty(t, specs)
	})
}

func addRecordsTitles(ctx context.Context, userID uuid.UUID, titles []string, mType models.MType) error {
	data := models.Data([]byte("test data"))
	for _, title := range titles {
		r := models.NewRecord(mType, title, data)
		_, err := testStorage.AddRecord(ctx, userID, r)
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
