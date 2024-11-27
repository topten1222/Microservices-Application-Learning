package authUsecase

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/topten1222/hello_sekai/config"
	"github.com/topten1222/hello_sekai/modules/auth"
	authrepository "github.com/topten1222/hello_sekai/modules/auth/authRepository"
	"github.com/topten1222/hello_sekai/modules/player"
	playerPb "github.com/topten1222/hello_sekai/modules/player/playerPb"
	"github.com/topten1222/hello_sekai/pkg/jwtauth"
	"github.com/topten1222/hello_sekai/utils"
)

type (
	AuthusecaseService interface {
		Login(context.Context, *config.Config, *auth.PlayerLoginReq) (*auth.ProfileIntercepter, error)
		RefreshToken(context.Context, *config.Config, *auth.RefreshTokenReq) (*auth.ProfileIntercepter, error)
	}

	authusecase struct {
		authRepo authrepository.AuthrepositoryService
	}
)

func NewAuthUsecase(authRepo authrepository.AuthrepositoryService) AuthusecaseService {
	return &authusecase{authRepo: authRepo}
}

func (u *authusecase) Login(pctx context.Context, cfg *config.Config, req *auth.PlayerLoginReq) (*auth.ProfileIntercepter, error) {
	profile, err := u.authRepo.CredentialSearch(pctx, cfg.Grpc.PlayerUrl, &playerPb.CredentialSearchReq{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return nil, err
	}
	accessToken := jwtauth.NewAccessToken(cfg.Jwt.AccessSecretKey, cfg.Jwt.AccessDuaration, &jwtauth.Claims{
		Id:       "player:" + profile.Id,
		RoleCode: int(profile.RoleCode),
	}).SignToken()
	refreshToken := jwtauth.NewRefreshToken(cfg.Jwt.RefreshSecretKey, cfg.Jwt.RefreshDuaration, &jwtauth.Claims{
		Id:       "player:" + profile.Id,
		RoleCode: int(profile.RoleCode),
	}).SignToken()
	credentailId, err := u.authRepo.InsertOnePlayerCredentail(pctx, &auth.Credential{
		PlayerId:     "player:" + profile.Id,
		RoleCode:     int(profile.RoleCode),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		CreatedAt:    utils.LocalTime(),
		UpdatedAt:    utils.LocalTime(),
	})
	if err != nil {
		return nil, err
	}
	credentail, err := u.authRepo.FindOnePlayerCredentail(pctx, credentailId.Hex())
	if err != nil {
		return nil, err
	}
	loc, _ := time.LoadLocation("Asia/Bangkok")
	return &auth.ProfileIntercepter{
		PlayerProfile: &player.PlayerProfile{
			Id:        profile.Id,
			Email:     profile.Email,
			Username:  profile.Username,
			CreatedAt: utils.ConvertStringTimeToTime(profile.CreatedAt).In(loc),
			UpdatedAt: utils.ConvertStringTimeToTime(profile.UpdatedAt).In(loc),
		},
		Credential: &auth.CredentialRes{
			Id:           credentailId.Hex(),
			PlayerId:     credentail.PlayerId,
			RoleCode:     credentail.RoleCode,
			AccessToken:  credentail.AccessToken,
			RefreshToken: credentail.RefreshToken,
			CreatedAt:    credentail.CreatedAt.In(loc),
			UpdatedAt:    credentail.UpdatedAt.In(loc),
		},
	}, nil
}

func (u *authusecase) RefreshToken(pctx context.Context, cfg *config.Config, req *auth.RefreshTokenReq) (*auth.ProfileIntercepter, error) {
	claims, err := jwtauth.ParseToken(cfg.Jwt.RefreshSecretKey, req.RefreshToken)
	if err != nil {
		log.Printf("Error: Refresh Token: %s", err.Error())
		return nil, err
	}
	profile, err := u.authRepo.FindOnePlayerProfileToRefresh(pctx, cfg.Grpc.PlayerUrl, &playerPb.FindOnePlayerProfileToRefreshReq{
		PlayerId: strings.TrimPrefix(claims.Id, "player:"),
	})
	if err != nil {
		return nil, err
	}
	accessToken := jwtauth.NewAccessToken(cfg.Jwt.AccessSecretKey, cfg.Jwt.AccessDuaration, &jwtauth.Claims{
		Id:       profile.Id,
		RoleCode: int(profile.RoleCode),
	}).SignToken()
	refreshToken := jwtauth.ReloadToken(cfg.Jwt.RefreshSecretKey, claims.ExpiresAt.Unix(), &jwtauth.Claims{
		Id:       profile.Id,
		RoleCode: int(profile.RoleCode),
	})
	fmt.Println(accessToken)
	if err := u.authRepo.UpdateOnePlayerCredentail(pctx, req.CredentialId, &auth.UpdateRefreshTokenReq{
		PlayerId:     profile.Id,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		UpdatedAt:    utils.LocalTime(),
	}); err != nil {
		return nil, err
	}
	credentail, err := u.authRepo.FindOnePlayerCredentail(pctx, req.CredentialId)
	if err != nil {
		return nil, err
	}
	loc, _ := time.LoadLocation("Asia/Bangkok")

	return &auth.ProfileIntercepter{
		PlayerProfile: &player.PlayerProfile{
			Id:        "player: " + profile.Id,
			Email:     profile.Email,
			Username:  profile.Username,
			CreatedAt: utils.ConvertStringTimeToTime(profile.CreatedAt),
			UpdatedAt: utils.ConvertStringTimeToTime(profile.UpdatedAt),
		},
		Credential: &auth.CredentialRes{
			Id:           credentail.Id.Hex(),
			PlayerId:     credentail.PlayerId,
			RoleCode:     credentail.RoleCode,
			AccessToken:  credentail.AccessToken,
			RefreshToken: credentail.RefreshToken,
			CreatedAt:    credentail.CreatedAt.In(loc),
			UpdatedAt:    credentail.UpdatedAt.In(loc),
		},
	}, nil
}
