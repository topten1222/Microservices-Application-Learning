package authHandler

import (
	authusecase "github.com/topten1222/hello_sekai/modules/auth/authUsecase"
)

type (
	AuthGrpHandler struct {
		authUsecase authusecase.AuthusecaseService
	}
)

func NewAuthGrpcHandler(authUsecase authusecase.AuthusecaseService) *AuthGrpHandler {
	return &AuthGrpHandler{authUsecase: authUsecase}
}
