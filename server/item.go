package server

import (
	"github.com/topten1222/hello_sekai/modules/item/itemHandler"
	"github.com/topten1222/hello_sekai/modules/item/itemRepository"
	"github.com/topten1222/hello_sekai/modules/item/itemUsecase"
)

func (s *server) itemService() {
	repo := itemRepository.NewItemRepository(s.db)
	usecase := itemUsecase.NewItemUsecase(repo)
	itemHandler.NewItemHttpHandler(s.cfg, usecase)
	itemHandler.NewItemGrpcHandler(usecase)

	item := s.app.Group("/item_v1")
	item.GET("/", s.healthCheckService)
}
