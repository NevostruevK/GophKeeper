// Keeper gRPC server.
package keeper

import (
	pb "github.com/NevostruevK/GophKeeper/proto"
)

// KeeperServer gRPC data server.
type KeeperServer struct {
	pb.UnimplementedKeeperServer
	storage DataStore
}

// NewKeeperServer returns KeeperServer.
func NewKeeperServer(dataStore DataStore) pb.KeeperServer {
	return &KeeperServer{storage: dataStore}
}
