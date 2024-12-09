package inventoryHandler

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/topten1222/hello_sekai/config"
	"github.com/topten1222/hello_sekai/modules/inventory"
	"github.com/topten1222/hello_sekai/modules/inventory/inventoryUsecase"
	"github.com/topten1222/hello_sekai/pkg/request"
	"github.com/topten1222/hello_sekai/pkg/response"
)

type (
	InventoryHttpHandlerService interface {
		FindPlayerItems(echo.Context) error
	}

	inventoryHttpHandler struct {
		cfg                 *config.Config
		inventoryUsecaseSer inventoryUsecase.InventoryUsecaseService
	}
)

func NewInventoryHttpHanlder(cfg *config.Config, inventoryUsecaseSer inventoryUsecase.InventoryUsecaseService) InventoryHttpHandlerService {
	return &inventoryHttpHandler{cfg: cfg, inventoryUsecaseSer: inventoryUsecaseSer}
}

func (h *inventoryHttpHandler) FindPlayerItems(c echo.Context) error {
	ctx := context.Background()

	wrapper := request.ContextWrapper(c)

	req := new(inventory.InventorySearchReq)
	playerId := c.Param("player_id")
	if err := wrapper.Bind(req); err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}
	res, err := h.inventoryUsecaseSer.FindPlayerItems(ctx, h.cfg, playerId, req)
	if err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	return response.SuccessResponse(c, http.StatusOK, res)
}
