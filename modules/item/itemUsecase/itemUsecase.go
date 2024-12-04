package itemUsecase

import (
	"context"
	"errors"

	"github.com/topten1222/hello_sekai/modules/item"
	"github.com/topten1222/hello_sekai/modules/item/itemRepository"
	"github.com/topten1222/hello_sekai/utils"
)

type (
	ItemUsecaseService interface {
		CreateItem(context.Context, *item.CreateItemReq) (any, error)
	}

	itemUsecase struct {
		itemRepo itemRepository.ItemRepositoryService
	}
)

func NewItemUsecase(itemRepo itemRepository.ItemRepositoryService) ItemUsecaseService {
	return &itemUsecase{itemRepo: itemRepo}
}

func (u *itemUsecase) CreateItem(pctx context.Context, req *item.CreateItemReq) (any, error) {
	if !u.itemRepo.IsUniqueItem(pctx, req.Title) {
		return nil, errors.New("item name is not unique")
	}

	itemId, err := u.itemRepo.InsertOneItem(pctx, &item.Item{
		Title:       req.Title,
		Price:       req.Price,
		Damage:      req.Damage,
		UsageStatus: true,
		ImageUrl:    req.ImageUrl,
		CreatedAt:   utils.LocalTime(),
		UpdatedAt:   utils.LocalTime(),
	})
	if err != nil {
		return nil, err
	}
	return itemId.Hex(), nil
}
