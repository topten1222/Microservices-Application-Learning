package inventoryHandler

import (
	"github.com/topten1222/hello_sekai/config"
	"github.com/topten1222/hello_sekai/modules/inventory/inventoryUsecase"
)

type (
	InventoryHttpHandlerService interface{}

	inventoryHttpHandler struct {
		cfg                 *config.Config
		inventoryUsecaseSer inventoryUsecase.InventoryUsecaseService
	}
)

func NewInventoryHttpHanlder(cfg *config.Config, inventoryUsecaseSer inventoryUsecase.InventoryUsecaseService) InventoryHttpHandlerService {
	return inventoryHttpHandler{inventoryUsecaseSer: inventoryUsecaseSer}
}
