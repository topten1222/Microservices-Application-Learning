package authHandler

import (
	"github.com/topten1222/hello_sekai/config"
	authusecase "github.com/topten1222/hello_sekai/modules/auth/authUsecase"
)

type (
	AuthHttpHandlerService interface{}

	authHttpHandler struct {
		cfg         *config.Config
		authusecase authusecase.AuthusecaseService
	}
)

func NewAuthHttpHandler(cfg *config.Config, authusecase authusecase.AuthusecaseService) AuthHttpHandlerService {
	return &authHttpHandler{cfg: cfg, authusecase: authusecase}
}
