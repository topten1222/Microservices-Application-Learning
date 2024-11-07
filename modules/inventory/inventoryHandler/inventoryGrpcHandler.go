package inventoryHandler

import "github.com/topten1222/hello_sekai/modules/inventory/inventoryUsecase"

type (
	inventoryGrpcHanlderService struct {
		inventoryUsecase inventoryUsecase.InventoryUsecaseService
	}
)

func NewInventoryGrpcHandler(inventoryUsecase inventoryUsecase.InventoryUsecaseService) *inventoryGrpcHanlderService {
	return &inventoryGrpcHanlderService{inventoryUsecase: inventoryUsecase}
}
