package playerHandler

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/topten1222/hello_sekai/config"
	"github.com/topten1222/hello_sekai/modules/player"
	"github.com/topten1222/hello_sekai/modules/player/playerUsecase"
	"github.com/topten1222/hello_sekai/pkg/request"
	"github.com/topten1222/hello_sekai/pkg/response"
)

type (
	PlayerHttpHandlerService interface {
		CreatePlayer(echo.Context) error
		FindOnePlayerProfile(echo.Context) error
		AddPlayerMoney(echo.Context) error
	}

	playerHttpHandler struct {
		cfg           *config.Config
		playerUsecase playerUsecase.PlayerUsecaseService
	}
)

func NewPlayerHttpHandler(cfg *config.Config, playerUsecase playerUsecase.PlayerUsecaseService) PlayerHttpHandlerService {
	return &playerHttpHandler{cfg: cfg, playerUsecase: playerUsecase}
}

func (h *playerHttpHandler) CreatePlayer(c echo.Context) error {
	ctx := context.Background()
	fmt.Println("erer222")
	wrapper := request.ContextWrapper(c)
	req := new(player.CreatePlayerReq)
	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(req); err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}
	if err := wrapper.Bind(req); err != nil {
		fmt.Println("eer")
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}
	fmt.Println("9898898")

	if err := validator.New().Struct(req); err != nil {
		fmt.Println("ttt")

		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}
	res, err := h.playerUsecase.CreatePlayer(ctx, req)
	if err != nil {
		return response.ErrResponse(c, http.StatusInternalServerError, err.Error())
	}
	return response.SuccessResponse(c, http.StatusOK, res)
}

func (h *playerHttpHandler) FindOnePlayerProfile(c echo.Context) error {
	ctx := context.Background()
	playerId := strings.TrimPrefix(c.Param("player_id"), "player:")
	if playerId == "" {
		return response.ErrResponse(c, http.StatusBadRequest, "player_id is required")
	}
	res, err := h.playerUsecase.FindOnePlayerProfile(ctx, playerId)
	if err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}
	return response.SuccessResponse(c, http.StatusOK, res)
}

func (h *playerHttpHandler) AddPlayerMoney(c echo.Context) error {
	ctx := context.Background()
	wrapper := request.ContextWrapper(c)
	req := new(player.CreatePlayerTransectionReq)
	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(req); err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}
	if err := wrapper.Bind(req); err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}
	if _, err := h.playerUsecase.FindOnePlayerProfile(ctx, req.PlayerId); err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}
	if err := h.playerUsecase.AddPlayerMoney(ctx, req); err != nil {
		return response.ErrResponse(c, http.StatusInternalServerError, err.Error())
	}
	return response.SuccessResponse(c, http.StatusOK, map[string]any{"message": "success"})
}
