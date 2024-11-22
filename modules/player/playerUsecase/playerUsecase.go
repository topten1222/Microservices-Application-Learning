package playerUsecase

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/topten1222/hello_sekai/modules/player"
	playerPb "github.com/topten1222/hello_sekai/modules/player/playerPb"
	"github.com/topten1222/hello_sekai/modules/player/playerRepository"
	"github.com/topten1222/hello_sekai/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type (
	PlayerUsecaseService interface {
		CreatePlayer(context.Context, *player.CreatePlayerReq) (*player.PlayerProfile, error)
		FindOnePlayerProfile(context.Context, string) (*player.PlayerProfile, error)
		AddPlayerMoney(context.Context, *player.CreatePlayerTransectionReq) (*player.PlayerSavingAccount, error)
		GetPlayerSavingAccount(context.Context, string) (*player.PlayerSavingAccount, error)
		FindOnePlayerCredentail(context.Context, string, string) (*playerPb.PlayerProfile, error)
	}

	playerUsecase struct {
		playerRepository playerRepository.PlayerRepositoryService
	}
)

func NewPlayerUsecase(playerRepository playerRepository.PlayerRepositoryService) PlayerUsecaseService {
	return &playerUsecase{playerRepository: playerRepository}
}

func (u *playerUsecase) CreatePlayer(pctx context.Context, req *player.CreatePlayerReq) (*player.PlayerProfile, error) {
	if !u.playerRepository.IsUniquePlayer(pctx, req.Email, req.Username) {
		return nil, errors.New("email or username is already taken")
	}
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("error: failed to hash password")
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
		return nil, errors.New("error: failed to create player")
	}

	return u.FindOnePlayerProfile(pctx, playerId.Hex())
}

func (u *playerUsecase) FindOnePlayerProfile(pctx context.Context, playerId string) (*player.PlayerProfile, error) {
	result, err := u.playerRepository.FindOnePlayerProfile(pctx, playerId)
	if err != nil {
		return nil, err
	}
	loc, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		return nil, errors.New("error: faild load location")
	}
	return &player.PlayerProfile{
		Id:        result.Id.Hex(),
		Email:     result.Email,
		Username:  result.Username,
		CreatedAt: result.CreatedAt.In(loc),
		UpdatedAt: result.UpdatedAt.In(loc),
	}, nil
}

func (u *playerUsecase) AddPlayerMoney(pctx context.Context, playerTransaction *player.CreatePlayerTransectionReq) (*player.PlayerSavingAccount, error) {
	if err := u.playerRepository.InsertOnePlayerTransaction(pctx, &player.PlayerTransaction{
		PlayerId:  playerTransaction.PlayerId,
		Amount:    playerTransaction.Amount,
		CreatedAt: utils.LocalTime(),
	}); err != nil {
		return nil, err
	}
	return u.playerRepository.GetPlayerSavingAccount(pctx, playerTransaction.PlayerId)

}

func (u *playerUsecase) GetPlayerSavingAccount(pctx context.Context, playerId string) (*player.PlayerSavingAccount, error) {
	return u.playerRepository.GetPlayerSavingAccount(pctx, playerId)
}

func (u *playerUsecase) FindOnePlayerCredentail(pctx context.Context, email, password string) (*playerPb.PlayerProfile, error) {
	result, err := u.playerRepository.FindOnePlayerCredential(pctx, email)
	if err != nil {
		return nil, err
	}
	fmt.Println(result.Password)
	fmt.Println(password)

	if err := bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(password)); err != nil {
		log.Printf("Error: invalida password %s", err.Error())
		return nil, errors.New("error: invalid password")
	}
	loc, _ := time.LoadLocation("Asia/Bangkok")
	return &playerPb.PlayerProfile{
		Id:        result.Id.Hex(),
		Email:     result.Email,
		Username:  result.Username,
		CreatedAt: result.CreatedAt.In(loc).String(),
		UpdatedAt: result.UpdatedAt.In(loc).String(),
	}, nil
}
