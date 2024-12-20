package playerHandler

import (
	"context"

	playerPb "github.com/topten1222/hello_sekai/modules/player/playerPb"
	"github.com/topten1222/hello_sekai/modules/player/playerUsecase"
)

type (
	playerGrpcHandler struct {
		playerUsecase playerUsecase.PlayerUsecaseService
		playerPb.UnimplementedPlayerGrpcServiceServer
	}
)

func NewPlayerGrpcHandler(playerUsecase playerUsecase.PlayerUsecaseService) *playerGrpcHandler {
	return &playerGrpcHandler{playerUsecase: playerUsecase}
}

func (g *playerGrpcHandler) CredentialSearch(ctx context.Context, req *playerPb.CredentialSearchReq) (*playerPb.PlayerProfile, error) {
	return g.playerUsecase.FindOnePlayerCredentail(ctx, req.Email, req.Password)
}

func (g *playerGrpcHandler) FindOnePlayerProfileToRefresh(ctx context.Context, req *playerPb.FindOnePlayerProfileToRefreshReq) (*playerPb.PlayerProfile, error) {
	return g.playerUsecase.FindOnePlayerToRefresh(ctx, req.PlayerId)
}

func (g *playerGrpcHandler) GetPlayerSavingAccount(ctx context.Context, req *playerPb.GetPlayerSavingAccountReq) (*playerPb.PlayerSavingAccountRes, error) {
	return nil, nil
}
