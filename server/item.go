package server

import (
	"log"

	"github.com/topten1222/hello_sekai/modules/item/itemHandler"
	itemPb "github.com/topten1222/hello_sekai/modules/item/itemPb"
	"github.com/topten1222/hello_sekai/modules/item/itemRepository"
	"github.com/topten1222/hello_sekai/modules/item/itemUsecase"
	"github.com/topten1222/hello_sekai/pkg/grpccon"
)

func (s *server) itemService() {
	repo := itemRepository.NewItemRepository(s.db)
	usecase := itemUsecase.NewItemUsecase(repo)
	itemHandler.NewItemHttpHandler(s.cfg, usecase)
	itemHandler.NewItemGrpcHandler(usecase)

	go func() {
		grpcServer, list := grpccon.NewGrpcServer(&s.cfg.Jwt, s.cfg.Grpc.ItemUrl)
		itemPb.RegisterItemGrpcServiceServer(grpcServer, itemHandler.NewItemGrpcHandler(usecase))
		log.Printf("Item grpc server listening on %s", s.cfg.Grpc.ItemUrl)
		grpcServer.Serve(list)
	}()

	item := s.app.Group("/item_v1")
	item.GET("/", s.healthCheckService)
}
