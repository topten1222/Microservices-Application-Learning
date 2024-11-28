package middlewareHandler

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/topten1222/hello_sekai/config"
	"github.com/topten1222/hello_sekai/modules/middleware/middlewareUsecase"
	"github.com/topten1222/hello_sekai/pkg/response"
)

type (
	MiddlewareHandlerService interface {
		JwtAuthorization(echo.HandlerFunc) echo.HandlerFunc
		RbacAuthorization(echo.HandlerFunc, []int) echo.HandlerFunc
	}
	middlewareHandler struct {
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

func (h *middlewareHandler) JwtAuthorization(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		accessToken := strings.TrimPrefix(c.Request().Header.Get("Authorization"), "Bearer ")
		newCtx, err := h.middlewareUsecase.JwtAuthorization(c, h.cfg, accessToken)
		if err != nil {
			return response.ErrResponse(c, http.StatusBadRequest, err.Error())
		}
		return next(newCtx)
	}
}

func (h *middlewareHandler) RbacAuthorization(next echo.HandlerFunc, expected []int) echo.HandlerFunc {
	return func(c echo.Context) error {
		newCtx, err := h.middlewareUsecase.RbacAuthorization(c, h.cfg, expected)
		if err != nil {
			return response.ErrResponse(c, http.StatusBadRequest, err.Error())
		}
		return next(newCtx)
	}
}
