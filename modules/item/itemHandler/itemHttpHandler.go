package itemHandler

import (
	"context"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/topten1222/hello_sekai/config"
	"github.com/topten1222/hello_sekai/modules/item"
	"github.com/topten1222/hello_sekai/modules/item/itemUsecase"
	"github.com/topten1222/hello_sekai/pkg/request"
	"github.com/topten1222/hello_sekai/pkg/response"
)

type (
	ItemHttpHandlerService interface {
		CreateItem(echo.Context) error
		FindOneItem(echo.Context) error
	}

	itemHttpHandler struct {
		cfg         *config.Config
		itemUsecase itemUsecase.ItemUsecaseService
	}
)

func NewItemHttpHandler(cfg *config.Config, itemUsecase itemUsecase.ItemUsecaseService) ItemHttpHandlerService {
	return &itemHttpHandler{cfg: cfg, itemUsecase: itemUsecase}
}

func (h *itemHttpHandler) CreateItem(c echo.Context) error {
	ctx := context.Background()
	wrapper := request.ContextWrapper(c)

	req := new(item.CreateItemReq)
	if err := wrapper.Bind(req); err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	res, err := h.itemUsecase.CreateItem(ctx, req)
	if err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())

	}
	return response.SuccessResponse(c, http.StatusCreated, res)

}

func (h *itemHttpHandler) FindOneItem(c echo.Context) error {
	ctx := context.Background()

	itemId := strings.TrimPrefix(c.Param("item_id"), "item:")

	res, err := h.itemUsecase.FindOneItem(ctx, itemId)
	if err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}
	return response.SuccessResponse(c, http.StatusOK, res)

}
