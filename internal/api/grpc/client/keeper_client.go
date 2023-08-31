package client

import (
	"context"
	"time"

	"github.com/NevostruevK/GophKeeper/internal/models"
	pb "github.com/NevostruevK/GophKeeper/proto"
	"google.golang.org/grpc"
)

const KeeperClientTimeOut = time.Second

// KeeperClient клиент для работы с данными.
type KeeperClient struct {
	conn    *grpc.ClientConn
	service pb.KeeperClient
}

// NewAuthClient конструктор для создания клиента данных.
func NewKeeperClient(conn *grpc.ClientConn) *KeeperClient {
	service := pb.NewKeeperClient(conn)
	return &KeeperClient{conn, service}
}

func (c *KeeperClient) GetSpecs(ctx context.Context) ([]models.Spec, error) {
	specs, err := c.service.GetSpecs(ctx, &pb.GetSpecsRequest{})
	if err != nil {
		return nil, err
	}
	return models.ProtoToSpecs(specs.Specs)
}

func (c *KeeperClient) GetSpecsOfType(ctx context.Context, typ models.MType) ([]models.Spec, error) {
	specs, err := c.service.GetSpecsOfType(ctx, &pb.GetSpecsOfTypeRequest{Type: models.DataTypeToProto(typ)})
	if err != nil {
		return nil, err
	}
	return models.ProtoToSpecs(specs.Specs)
}

func (c *KeeperClient) GetData(ctx context.Context, ds models.DataSpec) (models.Data, error) {
	data, err := c.service.GetData(ctx, &pb.DataSpec{
		Id:       ds.ID.String(),
		DataSize: uint64(ds.DataSize),
	})
	if err != nil {
		return nil, err
	}
	return data.Data, nil
}

func (c *KeeperClient) AddRecord(ctx context.Context, r *models.Record) (*models.DataSpec, error) {
	ds, err := c.service.AddRecord(ctx, &pb.Record{
		Title: r.Title,
		Type:  models.DataTypeToProto(r.Type),
		Data:  r.Data,
	})
	if err != nil {
		return nil, err
	}
	return models.ProtoToDataSpec(ds)
}

func (client *KeeperClient) Close() error {
	return client.conn.Close()
}
