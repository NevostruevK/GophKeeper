package keeper

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/NevostruevK/GophKeeper/internal/api/grpc/server/auth"
	"github.com/NevostruevK/GophKeeper/internal/models"
	pb "github.com/NevostruevK/GophKeeper/proto"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var ErrExtractUserID = errors.New(`failed to extract user ID `)

func (s *KeeperServer) GetID(ctx context.Context) (uuid.UUID, error) {
	val := ctx.Value(auth.KeyUserID)
	id, ok := val.(uuid.UUID)
	if !ok {
		log.Println(ErrExtractUserID, val)
		return uuid.Nil, ErrExtractUserID
	}
	return id, nil
}

func (s *KeeperServer) GetSpecs(ctx context.Context, req *pb.GetSpecsRequest) (*pb.Specs, error) {
	id, err := s.GetID(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	specs, err := s.storage.GetSpecs(ctx, id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &pb.Specs{Specs: models.SpecsToProto(specs)}, nil
}

func (s *KeeperServer) GetSpecsOfType(ctx context.Context, req *pb.GetSpecsOfTypeRequest) (*pb.Specs, error) {
	id, err := s.GetID(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	specs, err := s.storage.GetSpecsOfType(ctx, id, models.ProtoToDataType(req.Type))
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &pb.Specs{Specs: models.SpecsToProto(specs)}, nil
}

func (s *KeeperServer) GetData(ctx context.Context, req *pb.DataSpec) (*pb.Data, error) {
	op := "KeeperServer.GetData: "
	id, err := uuid.Parse(req.Id)
	if err != nil {
		msg := fmt.Sprintf("%s failed with error %v for id: %s", op, err, req.Id)
		log.Println(msg)
		return nil, status.Errorf(codes.Internal, msg)
	}
	data, err := s.storage.GetData(ctx, &models.DataSpec{ID: id, DataSize: int(req.DataSize)})
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &pb.Data{Data: data}, nil
}

func (s *KeeperServer) AddRecord(ctx context.Context, req *pb.Record) (*pb.DataSpec, error) {
	id, err := s.GetID(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	r := models.NewRecord(
		models.ProtoToDataType(req.Type),
		req.Title,
		req.Data,
	)
	ds, err := s.storage.AddRecord(ctx, id, r)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return ds.ToProto(), nil
}
