package inventoryHandler

import (
	"github.com/topten1222/hello_sekai/config"
	"github.com/topten1222/hello_sekai/modules/inventory/inventoryUsecase"
)

type (
	InventoryQueueHandlerService interface{}

	inventoryQueueHandler struct {
		cfg              *config.Config
		inventoryUsecase inventoryUsecase.InventoryUsecaseService
	}
)

func NewInventoryQueueHandler(cfg *config.Config, inventoryUsecase inventoryUsecase.InventoryUsecaseService) InventoryQueueHandlerService {
	return &inventoryQueueHandler{cfg: cfg, inventoryUsecase: inventoryUsecase}
}
