package itemHandler

import (
	"github.com/topten1222/hello_sekai/config"
	"github.com/topten1222/hello_sekai/modules/item/itemUsecase"
)

type (
	ItemHandlerService interface{}

	itemHandler struct {
		cfg         *config.Config
		itemUsecase itemUsecase.ItemUsecaseService
	}
)

func NewItemHttpHandler(cfg *config.Config, itemUsecase itemUsecase.ItemUsecaseService) ItemHandlerService {
	return &itemHandler{cfg: cfg, itemUsecase: itemUsecase}
}
