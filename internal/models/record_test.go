package models_test

import (
	"testing"

	"github.com/NevostruevK/GophKeeper/internal/models"
	pb "github.com/NevostruevK/GophKeeper/proto"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	sPAIR = "PAIR"
	sTEXT = "TEXT"
	sFILE = "FILE"
	sCARD = "CARD"
)

func TestModels_StringToMType(t *testing.T) {
	tests := []struct {
		name string
		arg  string
		want models.MType
	}{
		{
			name: sPAIR,
			arg:  sPAIR,
			want: models.PAIR,
		},
		{
			name: sTEXT,
			arg:  sTEXT,
			want: models.TEXT,
		},
		{
			name: sFILE,
			arg:  sFILE,
			want: models.FILE,
		},
		{
			name: sCARD,
			arg:  sCARD,
			want: models.CARD,
		},
		{
			name: "any string",
			arg:  "any string",
			want: models.NOTIMPLEMENT,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := models.StringToMType(tt.arg); got != tt.want {
				t.Errorf("StringToType(%s) got %v, want %v", tt.arg, got, tt.want)
			}
		})
	}
}

func TestModels_DataTypeToProto(t *testing.T) {
	tests := []struct {
		name string
		arg  models.MType
		want pb.DataType
	}{
		{
			name: sPAIR,
			arg:  models.PAIR,
			want: pb.DataType_PAIR,
		},
		{
			name: sTEXT,
			arg:  models.TEXT,
			want: pb.DataType_TEXT,
		},
		{
			name: sFILE,
			arg:  models.FILE,
			want: pb.DataType_FILE,
		},
		{
			name: sCARD,
			arg:  models.CARD,
			want: pb.DataType_CARD,
		},
		{
			name: "any string",
			arg:  "any string",
			want: pb.DataType_TEXT,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := models.DataTypeToProto(tt.arg); got != tt.want {
				t.Errorf("DataTypeToProto(%s) got %v, want %v", tt.arg, got, tt.want)
			}
		})
	}
}

func TestModels_ProtoToDataType(t *testing.T) {
	tests := []struct {
		name string
		arg  pb.DataType
		want models.MType
	}{
		{
			name: sPAIR,
			arg:  pb.DataType_PAIR,
			want: models.PAIR,
		},
		{
			name: sTEXT,
			arg:  pb.DataType_TEXT,
			want: models.TEXT,
		},
		{
			name: sFILE,
			arg:  pb.DataType_FILE,
			want: models.FILE,
		},
		{
			name: sCARD,
			arg:  pb.DataType_CARD,
			want: models.CARD,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := models.ProtoToDataType(tt.arg); got != tt.want {
				t.Errorf("ProtoToDataType(%s) got %v, want %v", tt.arg, got, tt.want)
			}
		})
	}
}

func TestModelsSpec_ToProto(t *testing.T) {
	spec := models.NewSpec(models.TEXT, "some title")
	spec.ID = uuid.New()
	spec.DataSize = len(spec.Title)
	actual := spec.ToProto()
	expected := &pb.Spec{
		Id:       spec.ID.String(),
		Type:     pb.DataType_TEXT,
		Tytle:    spec.Title,
		DataSize: int64(spec.DataSize),
	}
	assert.Equal(t, expected, actual)
}

func TestModels_SpecsToProto(t *testing.T) {
	outSpecs := []models.Spec{
		*models.NewSpec(models.TEXT, "some text title"),
		*models.NewSpec(models.PAIR, "some pair title"),
		*models.NewSpec(models.FILE, "some file title"),
		*models.NewSpec(models.CARD, "some card title"),
	}
	for i, spec := range outSpecs {
		outSpecs[i].ID = uuid.New()
		outSpecs[i].DataSize = len(spec.Title)
	}
	t.Run("Proto to Specs ok", func(t *testing.T) {
		pbSpecs := models.SpecsToProto(outSpecs)
		require.NotNil(t, pbSpecs)
		assert.Equal(t, len(outSpecs), len(pbSpecs))
		inSpecs, err := models.ProtoToSpecs(pbSpecs)
		require.NoError(t, err)
		assert.ElementsMatch(t, outSpecs, inSpecs)
	})
	t.Run("Proto to Specs err", func(t *testing.T) {
		pbSpecs := models.SpecsToProto(outSpecs)
		require.NotNil(t, pbSpecs)
		assert.Equal(t, len(outSpecs), len(pbSpecs))
		pbSpecs[0].Id = "Bad UUID"
		_, err := models.ProtoToSpecs(pbSpecs)
		require.Error(t, err)
	})
}

func TestModelsDataSpec_ToProto(t *testing.T) {
	t.Run("DataSpec to proto ok", func(t *testing.T) {
		outDs := models.DataSpec{ID: uuid.New(), DataSize: 123}
		pbDataSpec := outDs.ToProto()
		require.NotNil(t, pbDataSpec)
		inDs, err := models.ProtoToDataSpec(pbDataSpec)
		require.NoError(t, err)
		assert.Equal(t, outDs, *inDs)
	})
	t.Run("DataSpec to proto err", func(t *testing.T) {
		outDs := models.DataSpec{ID: uuid.New(), DataSize: 123}
		pbDataSpec := outDs.ToProto()
		require.NotNil(t, pbDataSpec)
		pbDataSpec.Id = "Bad UUID"
		_, err := models.ProtoToDataSpec(pbDataSpec)
		require.Error(t, err)
	})
}

func TestModelsRecord_ToSpec(t *testing.T) {
	ds := models.DataSpec{ID: uuid.New(), DataSize: 123}
	tests := []struct {
		name string
		obj  *models.Record
		want models.Spec
	}{
		{
			name: sPAIR,
			obj:  models.NewRecord(models.PAIR, "title pair", []byte("data pair")),
			want: models.Spec{models.PAIR, "title pair", ds},
		},
		{
			name: sTEXT,
			obj:  models.NewRecord(models.TEXT, "title text", []byte("data text")),
			want: models.Spec{models.TEXT, "title text", ds},
		},
		{
			name: sFILE,
			obj:  models.NewRecord(models.FILE, "title file", []byte("data file")),
			want: models.Spec{models.FILE, "title file", ds},
		},
		{
			name: sCARD,
			obj:  models.NewRecord(models.CARD, "title card", []byte("data card")),
			want: models.Spec{models.CARD, "title card", ds},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.obj.ToSpec(ds); *got != tt.want {
				t.Errorf("record ToSpec got %v, want %v", got, tt.want)
			}
		})
	}
}
