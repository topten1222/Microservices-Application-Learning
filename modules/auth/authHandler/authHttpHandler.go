package authHandler

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/topten1222/hello_sekai/config"
	"github.com/topten1222/hello_sekai/modules/auth"
	authusecase "github.com/topten1222/hello_sekai/modules/auth/authUsecase"
	"github.com/topten1222/hello_sekai/pkg/request"
	"github.com/topten1222/hello_sekai/pkg/response"
)

type (
	AuthHttpHandlerService interface {
		Login(echo.Context) error
		RefreshToken(echo.Context) error
	}

	authHttpHandler struct {
		cfg         *config.Config
		authusecase authusecase.AuthusecaseService
	}
)

func NewAuthHttpHandler(cfg *config.Config, authusecase authusecase.AuthusecaseService) AuthHttpHandlerService {
	return &authHttpHandler{cfg: cfg, authusecase: authusecase}
}

func (h *authHttpHandler) Login(c echo.Context) error {
	ctx := context.Background()
	wrapper := request.ContextWrapper(c)
	req := new(auth.PlayerLoginReq)
	if err := wrapper.Bind(req); err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}
	res, err := h.authusecase.Login(ctx, h.cfg, req)
	if err != nil {
		return response.ErrResponse(c, http.StatusInternalServerError, err.Error())

	}

	return response.SuccessResponse(c, http.StatusOK, res)

}

func (h *authHttpHandler) RefreshToken(c echo.Context) error {
	ctx := context.Background()
	wrapper := request.ContextWrapper(c)
	req := new(auth.RefreshTokenReq)
	if err := wrapper.Bind(req); err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}
	res, err := h.authusecase.RefreshToken(ctx, h.cfg, req)
	if err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}
	return response.SuccessResponse(c, http.StatusOK, res)
}
