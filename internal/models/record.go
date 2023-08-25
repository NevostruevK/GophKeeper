package models

import (
	pb "github.com/NevostruevK/GophKeeper/proto"
	"github.com/google/uuid"
)

const (
	TitleSize = 256
)

type MType string

const (
	PAIR MType = "PAIR"
	TEXT MType = "TEXT"
	FILE MType = "FILE"
	CARD MType = "CARD"
)

func DataTypeToProto(typ MType) pb.DataType {
	switch typ {
	case PAIR:
		return pb.DataType_PAIR
	case TEXT:
		return pb.DataType_TEXT
	case FILE:
		return pb.DataType_FILE
	default:
		return pb.DataType_CARD
	}
	//TODO ошибка не совпал тип
}

func ProtoToDataType(typ pb.DataType) MType {
	switch typ {
	case pb.DataType_PAIR:
		return PAIR
	case pb.DataType_TEXT:
		return TEXT
	case pb.DataType_FILE:
		return FILE
	default:
		return CARD
	}
	//TODO ошибка не совпал тип
}

type Spec struct {
	ID             uuid.UUID
	Type           MType
	Title          string
	DataSize       int
	HasDescription bool
}

func NewSpec(typ MType, title string) *Spec {
	return &Spec{
		ID:    uuid.New(),
		Type:  typ,
		Title: title,
	}
}

func (s *Spec) ToProto() *pb.Spec {
	return &pb.Spec{
		Id:             s.ID.String(),
		Type:           DataTypeToProto(s.Type),
		Tytle:          s.Title,
		DataSize:       int64(s.DataSize),
		HasDescription: s.HasDescription,
	}
}

func SpecsToProto(specs []Spec) []*pb.Spec {
	specsPB := make([]*pb.Spec, len(specs))
	for i, spec := range specs {
		specsPB[i] = spec.ToProto()
	}
	return specsPB
}

type DataSpec struct {
	ID       uuid.UUID
	DataSize int
}

type Data []byte

type Record struct {
	Type  MType
	Title string
	Data
	Description []byte
}

func NewRecord(typ MType, title string, data, description []byte) *Record {
	return &Record{typ, title, data, description}
}

func (r *Record) ToSpec(id uuid.UUID) *Spec {
	return &Spec{
		ID:             id,
		Type:           r.Type,
		Title:          r.Title,
		DataSize:       len(r.Data),
		HasDescription: r.Description != nil,
	}
}
