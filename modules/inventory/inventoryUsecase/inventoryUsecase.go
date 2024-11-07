package inventoryUsecase

import (
	"github.com/topten1222/hello_sekai/modules/inventory/inventoryRepository"
)

type (
	InventoryUsecaseService interface{}

	inventoryUsecase struct {
		inventoryRepo inventoryRepository.InventoryRepositoryService
	}
)

func NewInventoryUsecasee(inventoryRepo inventoryRepository.InventoryRepositoryService) InventoryUsecaseService {
	return inventoryUsecase{inventoryRepo: inventoryRepo}
}
