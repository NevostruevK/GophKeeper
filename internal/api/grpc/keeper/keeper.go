package keeper

import (
	pb "github.com/NevostruevK/GophKeeper/proto"
)

type KeeperServer struct {
	pb.UnimplementedKeeperServer
	storage DataStore
}

func NewKeeperServer(dataStore DataStore) pb.KeeperServer {
	return &KeeperServer{storage: dataStore}
}
