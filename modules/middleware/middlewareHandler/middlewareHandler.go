package middlewareHandler

import (
	"github.com/topten1222/hello_sekai/config"
	"github.com/topten1222/hello_sekai/modules/middleware/middlewareUsecase"
)

type (
	MiddlewareHandlerService interface{}
	middlewareHandler        struct {
		cfg               *config.Config
		middlewareUsecase middlewareUsecase.MiddlewareUsecaseService
	}
)

func NewMiddlewareHandler(cfg *config.Config, middlewareUsecase middlewareUsecase.MiddlewareUsecaseService) MiddlewareHandlerService {
	return &middlewareHandler{
		cfg:               cfg,
		middlewareUsecase: middlewareUsecase,
	}
}
