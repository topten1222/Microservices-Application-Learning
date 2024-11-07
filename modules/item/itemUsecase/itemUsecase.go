package itemUsecase

import "github.com/topten1222/hello_sekai/modules/item/itemRepository"

type (
	ItemUsecaseService interface{}

	itemUsecase struct {
		itemRepo itemRepository.ItemRepositoryService
	}
)

func NewItemUsecase(itemRepo itemRepository.ItemRepositoryService) ItemUsecaseService {
	return &itemUsecase{itemRepo: itemRepo}
}
