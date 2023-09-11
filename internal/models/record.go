package models

import (
	"fmt"

	pb "github.com/NevostruevK/GophKeeper/proto"
	"github.com/google/uuid"
)

const (
	// TitleSize maximum size for record's field Title.
	TitleSize = 256
)

// MType types for record's data.
type MType string

const (
	// PAIR type for Pair.
	PAIR MType = "PAIR"
	// TEXT type for Text.
	TEXT MType = "TEXT"
	// FILE type for File.
	FILE MType = "FILE"
	// CARD type for Card.
	CARD MType = "CARD"
	// NOTIMPLEMENT type for not implemented data.
	NOTIMPLEMENT MType = "NOT IMPLEMENTED"
)

const (
	sPAIR = "PAIR"
	sTEXT = "TEXT"
	sFILE = "FILE"
	sCARD = "CARD"
)

// StringToMType converts string to MType.
func StringToMType(typ string) MType {
	switch typ {
	case sCARD:
		return CARD
	case sFILE:
		return FILE
	case sTEXT:
		return TEXT
	case sPAIR:
		return PAIR
	}
	return NOTIMPLEMENT
}

// DataTypeToProto converts MType to pb.DataType for gRPC transaction.
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
}

// ProtoToDataType converts pb.DataType to MType.
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
}

// Spec extention data specification.
type Spec struct {
	Type     MType  // data type
	Title    string // title for data
	DataSpec        // data size and data ID
}

// NewSpec returns Spec.
func NewSpec(typ MType, title string) *Spec {
	return &Spec{
		Type:  typ,
		Title: title,
	}
}

// ToProto converts Spec to pb.Spec for gRPC transaction.
func (s *Spec) ToProto() *pb.Spec {
	return &pb.Spec{
		Id:       s.ID.String(),
		Type:     DataTypeToProto(s.Type),
		Tytle:    s.Title,
		DataSize: int64(s.DataSize),
	}
}

// SpecsToProto converts slice of Specs to slice of pb.Spec for gRPC transaction.
func SpecsToProto(specs []Spec) []*pb.Spec {
	specsPB := make([]*pb.Spec, len(specs))
	for i, spec := range specs {
		specsPB[i] = spec.ToProto()
	}
	return specsPB
}

// ProtoToSpecs converts slice of pb.Spec to slice of Spec.
func ProtoToSpecs(pbSpecs []*pb.Spec) ([]Spec, error) {
	specs := make([]Spec, len(pbSpecs))
	for i, s := range pbSpecs {
		id, err := uuid.Parse(s.Id)
		if err != nil {
			return nil, err
		}
		specs[i] = Spec{
			DataSpec: DataSpec{id, int(s.DataSize)},
			Type:     ProtoToDataType(s.Type),
			Title:    s.Tytle,
		}
	}
	return specs, nil
}

// DataSpec data specification.
type DataSpec struct {
	ID       uuid.UUID // data ID
	DataSize int       // data size
}

// ToProto converts DataSpec to pb.DataSpec for gRPC transaction.
func (ds DataSpec) ToProto() *pb.DataSpec {
	return &pb.DataSpec{Id: ds.ID.String(), DataSize: uint64(ds.DataSize)}
}

// ProtoToDataSpec converts pb.DataSpec to DataSpec.
func ProtoToDataSpec(dsProto *pb.DataSpec) (*DataSpec, error) {
	id, err := uuid.Parse(dsProto.Id)
	if err != nil {
		return nil, fmt.Errorf("ProtoToDataSpec failed with error %v for id %s", err, dsProto.Id)
	}
	return &DataSpec{ID: id, DataSize: int(dsProto.DataSize)}, nil
}

// Data type for data.
type Data []byte

// Record record for storing.
type Record struct {
	Type  MType  // record type
	Title string // record title
	Data         // record data
}

// NewRecord returns Record.
func NewRecord(typ MType, title string, data []byte) *Record {
	return &Record{typ, title, data}
}

// ToSpec converts Record to Spec.
func (r *Record) ToSpec(ds DataSpec) *Spec {
	return &Spec{
		DataSpec: ds,
		Type:     r.Type,
		Title:    r.Title,
	}
}

// ToProto converts Record to pb.Record for gRPC transaction.
func (r *Record) ToProto() *pb.Record {
	return &pb.Record{
		Type:  DataTypeToProto(r.Type),
		Title: r.Title,
		Data:  r.Data,
	}
}
