package keeper

import (
	"context"

	"github.com/NevostruevK/GophKeeper/internal/models"
	pb "github.com/NevostruevK/GophKeeper/proto"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *KeeperServer) GetID() uuid.UUID {
	return uuid.New()
}

func (s *KeeperServer) GetSpecs(ctx context.Context, req *pb.GetSpecsRequest) (*pb.Specs, error) {
	id := s.GetID()
	specs, err := s.storage.GetSpecs(ctx, id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &pb.Specs{Specs: models.SpecsToProto(specs)}, nil
}

func (s *KeeperServer) GetSpecsOfType(ctx context.Context, req *pb.GetSpecsOfTypeRequest) (*pb.Specs, error) {
	return nil, nil
}

func (s *KeeperServer) GetData(ctx context.Context, req *pb.DataSpec) (*pb.Data, error) {
	return nil, nil
}

func (s *KeeperServer) GetDescription(ctx context.Context, req *pb.RecordID) (*pb.Data, error) {
	return nil, nil
}

func (s *KeeperServer) AddRecord(ctx context.Context, req *pb.Record) (*pb.RecordID, error) {
	id := s.GetID()

	r := models.NewRecord(
		models.ProtoToDataType(req.Type),
		req.Title,
		req.Data,
		req.Description,
	)
	id, err := s.storage.AddRecord(ctx, id, r)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &pb.RecordID{Id: id.String()}, nil
}

/*
rpc GetSpecs(GetSpecsRequest) returns (Specs);
rpc GetSpecsOfType(GetSpecsOfTypeRequest) returns (Specs);
rpc GetData(DataSpec) returns (Data);
rpc GetDescription(RecordID) returns (Data);
rpc AddRecord(AddRecordRequest) returns (RecordID);
*/
