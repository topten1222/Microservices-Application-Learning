package server

import (
	"github.com/topten1222/hello_sekai/modules/player/playerHandler"
	"github.com/topten1222/hello_sekai/modules/player/playerRepository"
	"github.com/topten1222/hello_sekai/modules/player/playerUsecase"
)

func (s *server) playerService() {
	repo := playerRepository.NewPlayerRepository(s.db)
	usecase := playerUsecase.NewPlayerUsecase(repo)
	playerHandler.NewPlayerHttpHandler(s.cfg, usecase)
	playerHandler.NewPlayerGrpcHandler(usecase)
	playerHandler.NewPlayerQueueHandler(s.cfg, usecase)

	s.app.Group("/player")
}
