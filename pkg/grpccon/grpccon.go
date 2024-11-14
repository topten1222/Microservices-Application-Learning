package grpccon

import (
	"log"
	"net"

	"github.com/topten1222/hello_sekai/config"
	authPb "github.com/topten1222/hello_sekai/modules/auth/authPb"
	inventoryPb "github.com/topten1222/hello_sekai/modules/inventory/inventoryPb"
	itemPb "github.com/topten1222/hello_sekai/modules/item/itemPb"
	playerPb "github.com/topten1222/hello_sekai/modules/player/playerPb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type (
	GrpcClientFactoryHandler interface {
		Auth() authPb.AuthGrpcServiceClient
		Player() playerPb.PlayerGrpcServiceClient
		Item() itemPb.ItemGrpcServiceClient
		Inventory() inventoryPb.InventoryGrpcServiceClient
	}

	grpcClientFactory struct {
		client *grpc.ClientConn
	}

	grpcAuth struct{}
)

func (g *grpcClientFactory) Auth() authPb.AuthGrpcServiceClient {
	return authPb.NewAuthGrpcServiceClient(g.client)
}

func (g *grpcClientFactory) Player() playerPb.PlayerGrpcServiceClient {
	return playerPb.NewPlayerGrpcServiceClient(g.client)
}

func (g *grpcClientFactory) Item() itemPb.ItemGrpcServiceClient {
	return itemPb.NewItemGrpcServiceClient(g.client)
}

func (g *grpcClientFactory) Inventory() inventoryPb.InventoryGrpcServiceClient {
	return inventoryPb.NewInventoryGrpcServiceClient(g.client)
}

func NewGrpcClient(host string) (GrpcClientFactoryHandler, error) {
	opts := make([]grpc.DialOption, 0)
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	clientConn, err := grpc.Dial(host, opts...)
	if err != nil {
		return nil, err
	}
	return &grpcClientFactory{client: clientConn}, nil
}

func NewGrpcServer(cfg *config.Jwt, host string) (*grpc.Server, net.Listener) {
	opts := make([]grpc.ServerOption, 0)
	grpcServer := grpc.NewServer(opts...)
	lit, err := net.Listen("tcp", host)
	if err != nil {
		log.Fatal("Error:: ", err)
	}
	return grpcServer, lit
}
