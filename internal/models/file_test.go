package models_test

import (
	"fmt"
	"testing"

	"github.com/NevostruevK/GophKeeper/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestModelsFile(t *testing.T) {
	t.Run("test String", func(t *testing.T) {
		file := models.NewFile("file.txt", []byte("file data"))
		s := file.String()
		assert.Equal(t, "file.txt : 9", s)
	})
	t.Run("test Show", func(t *testing.T) {
		file := models.NewFile("file.txt", []byte("file data"))
		s := file.Show()
		assert.Equal(t, "file.txt : 9", s)
	})
}

func TestModelsFile_IsReadyForStorage(t *testing.T) {
	const (
		fileNameIsEmpty = "file name is empty"
		fileIsNotExist  = "file is not exist"
	)
	type result struct {
		bool
		string
	}
	tests := []struct {
		name string
		obj  models.File
		want result
	}{
		{
			name: "test ok",
			obj:  *models.NewFile("./test_data/file.tst", nil),
			want: result{true, ""},
		},
		{
			name: "test err (file name is empty)",
			obj:  *models.NewFile("", []byte("file data")),
			want: result{false, fileNameIsEmpty},
		},
		{
			name: "test err (file is not exist)",
			obj:  *models.NewFile("not-exist", []byte("file data")),
			want: result{false, fmt.Sprintf("file not-exist %s", fileIsNotExist)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if ok, s := tt.obj.IsReadyForStorage(); ok != tt.want.bool || s != tt.want.string {
				t.Errorf("IsReadyForStorage(%v) got (%v , %v), want %v", tt.obj, ok, s, tt.want)
			}
		})
	}
}
