package playerUsecase

import (
	"context"
	"errors"

	"github.com/topten1222/hello_sekai/modules/player"
	"github.com/topten1222/hello_sekai/modules/player/playerRepository"
	"github.com/topten1222/hello_sekai/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type (
	PlayerUsecaseService interface {
		CreatePlayer(context.Context, *player.CreatePlayerReq) (string, error)
	}

	playerUsecase struct {
		playerRepository playerRepository.PlayerRepositoryService
	}
)

func NewPlayerUsecase(playerRepository playerRepository.PlayerRepositoryService) PlayerUsecaseService {
	return &playerUsecase{playerRepository: playerRepository}
}

func (u *playerUsecase) CreatePlayer(pctx context.Context, req *player.CreatePlayerReq) (string, error) {
	if !u.playerRepository.IsUniquePlayer(pctx, req.Email, req.Username) {
		return "", errors.New("email or username is already taken")
	}
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.New("error: failed to hash password")
	}
	playerId, err := u.playerRepository.InsertOnePlayer(pctx, &player.Player{
		Id:       primitive.NewObjectID(),
		Email:    req.Email,
		Password: string(hashPassword),
		Username: req.Username,
		PlayerRoles: []player.PlayerRole{
			{
				RoleTitle: "Player",
				RoleCode:  0,
			},
		},
		CreatedAt: utils.LocalTime(),
		UpdatedAt: utils.LocalTime(),
	})
	if err != nil {
		return "", errors.New("error: failed to create player")
	}
	return playerId.Hex(), nil
}
