package server

import (
	"log"

	"github.com/topten1222/hello_sekai/modules/inventory/inventoryHandler"
	inventoryPb "github.com/topten1222/hello_sekai/modules/inventory/inventoryPb"
	"github.com/topten1222/hello_sekai/modules/inventory/inventoryRepository"
	"github.com/topten1222/hello_sekai/modules/inventory/inventoryUsecase"
	"github.com/topten1222/hello_sekai/pkg/grpccon"
)

func (s *server) inventoryService() {
	repo := inventoryRepository.NewInventoryRepository(s.db)
	usecase := inventoryUsecase.NewInventoryUsecasee(repo)
	httpHandler := inventoryHandler.NewInventoryHttpHanlder(s.cfg, usecase)
	inventoryHandler.NewInventoryGrpcHandler(usecase)
	inventoryHandler.NewInventoryQueueHandler(s.cfg, usecase)

	go func() {
		grpcServer, list := grpccon.NewGrpcServer(&s.cfg.Jwt, s.cfg.Grpc.InventoryUrl)
		inventoryPb.RegisterInventoryGrpcServiceServer(grpcServer, inventoryHandler.NewInventoryGrpcHandler(usecase))
		log.Printf("Inventory grpc server listening on %s", s.cfg.Grpc.InventoryUrl)
		grpcServer.Serve(list)
	}()

	inventory := s.app.Group("/inventory_v1")
	inventory.GET("/", s.healthCheckService)
	inventory.GET("/inventory:/:player_id", httpHandler.FindPlayerItems)
}
