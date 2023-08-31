package models

import (
	"fmt"

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
	ID       uuid.UUID
	Type     MType
	Title    string
	DataSize int
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
		Id:       s.ID.String(),
		Type:     DataTypeToProto(s.Type),
		Tytle:    s.Title,
		DataSize: int64(s.DataSize),
	}
}

func SpecsToProto(specs []Spec) []*pb.Spec {
	specsPB := make([]*pb.Spec, len(specs))
	for i, spec := range specs {
		specsPB[i] = spec.ToProto()
	}
	return specsPB
}

func ProtoToSpecs(pbSpecs []*pb.Spec) ([]Spec, error) {
	specs := make([]Spec, len(pbSpecs))
	for i, s := range pbSpecs {
		id, err := uuid.Parse(s.Id)
		if err != nil {
			return nil, err
		}
		specs[i] = Spec{
			ID:       id,
			Type:     ProtoToDataType(s.Type),
			Title:    s.Tytle,
			DataSize: int(s.DataSize),
		}
	}
	return specs, nil
}

type DataSpec struct {
	ID       uuid.UUID
	DataSize int
}

func (ds DataSpec) ToProto() *pb.DataSpec {
	return &pb.DataSpec{Id: ds.ID.String(), DataSize: uint64(ds.DataSize)}
}

func ProtoToDataSpec(dsProto *pb.DataSpec) (*DataSpec, error) {
	id, err := uuid.Parse(dsProto.Id)
	if err != nil {
		return nil, fmt.Errorf("ProtoToDataSpec failed with error %v for id %s", err, dsProto.Id)
	}
	return &DataSpec{ID: id, DataSize: int(dsProto.DataSize)}, nil
}

type Data []byte

type Record struct {
	Type  MType
	Title string
	Data
}

func NewRecord(typ MType, title string, data []byte) *Record {
	return &Record{typ, title, data}
}

func (r *Record) ToSpec(ds *DataSpec) *Spec {
	return &Spec{
		ID:       ds.ID,
		Type:     r.Type,
		Title:    r.Title,
		DataSize: ds.DataSize,
	}
}

func (r *Record) ToProto() *pb.Record {
	return &pb.Record{
		Type:  DataTypeToProto(r.Type),
		Title: r.Title,
		Data:  r.Data,
	}
}
