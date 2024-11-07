package playerHandler

import "github.com/topten1222/hello_sekai/modules/player/playerUsecase"

type (
	playerGrpcHandlerService struct {
		playerUsecase playerUsecase.PlayerUsecaseService
	}
)

func NewPlayerGrpcHandler(playerUsecase playerUsecase.PlayerUsecaseService) *playerGrpcHandlerService {
	return &playerGrpcHandlerService{playerUsecase: playerUsecase}
}
