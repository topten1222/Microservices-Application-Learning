package server

import (
	"github.com/topten1222/hello_sekai/modules/inventory/inventoryHandler"
	"github.com/topten1222/hello_sekai/modules/inventory/inventoryRepository"
	"github.com/topten1222/hello_sekai/modules/inventory/inventoryUsecase"
)

func (s *server) inventoryService() {
	repo := inventoryRepository.NewInventoryRepository(s.db)
	usecase := inventoryUsecase.NewInventoryUsecasee(repo)
	inventoryHandler.NewInventoryHttpHanlder(s.cfg, usecase)
	inventoryHandler.NewInventoryGrpcHandler(usecase)
	inventoryHandler.NewInventoryQueueHandler(s.cfg, usecase)

	s.app.Group("/inventory")
}
