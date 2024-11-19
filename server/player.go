package server

import (
	"log"

	"github.com/topten1222/hello_sekai/modules/player/playerHandler"
	playerPb "github.com/topten1222/hello_sekai/modules/player/playerPb"
	"github.com/topten1222/hello_sekai/modules/player/playerRepository"
	"github.com/topten1222/hello_sekai/modules/player/playerUsecase"
	"github.com/topten1222/hello_sekai/pkg/grpccon"
)

func (s *server) playerService() {
	repo := playerRepository.NewPlayerRepository(s.db)
	usecase := playerUsecase.NewPlayerUsecase(repo)
	httpHandler := playerHandler.NewPlayerHttpHandler(s.cfg, usecase)
	playerHandler.NewPlayerGrpcHandler(usecase)
	playerHandler.NewPlayerQueueHandler(s.cfg, usecase)
	go func() {
		grpcServer, list := grpccon.NewGrpcServer(&s.cfg.Jwt, s.cfg.Grpc.PlayerUrl)
		playerPb.RegisterPlayerGrpcServiceServer(grpcServer, playerHandler.NewPlayerGrpcHandler(usecase))
		log.Printf("Player grpc server listening on %s", s.cfg.Grpc.PlayerUrl)
		grpcServer.Serve(list)
	}()

	player := s.app.Group("/player_v1")
	player.GET("/", s.healthCheckService)
	player.POST("/player/register", httpHandler.CreatePlayer)
	player.GET("/player/:player_id", httpHandler.FindOnePlayerProfile)
	player.POST("/player/add-money", httpHandler.AddPlayerMoney)
	player.GET("/player/account/:player_id", httpHandler.GetPlayerSavingAccount)

}
