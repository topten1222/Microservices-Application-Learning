package inventoryUsecase

import (
	"context"
	"fmt"

	"github.com/topten1222/hello_sekai/modules/inventory"
	"github.com/topten1222/hello_sekai/modules/inventory/inventoryRepository"
	"github.com/topten1222/hello_sekai/modules/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	InventoryUsecaseService interface {
		FindPlayerItems(context.Context, string, string, *inventory.InventorySearchReq) (*models.PaginateRes, error)
	}

	inventoryUsecase struct {
		inventoryRepo inventoryRepository.InventoryRepositoryService
	}
)

func NewInventoryUsecasee(inventoryRepo inventoryRepository.InventoryRepositoryService) InventoryUsecaseService {
	return &inventoryUsecase{inventoryRepo: inventoryRepo}
}

func (u *inventoryUsecase) FindPlayerItems(pctx context.Context, basePaginatateUrl, playerId string, req *inventory.InventorySearchReq) (*models.PaginateRes, error) {
	findItemsOpts := make([]*options.FindOptions, 0)
	findItemsOpts = append(findItemsOpts, options.Find().SetSort(bson.D{{"_id", 1}}))
	findItemsOpts = append(findItemsOpts, options.Find().SetLimit(int64(req.Limit)))

	results, err := u.inventoryRepo.FindPlayerItems(pctx, playerId, findItemsOpts)
	if err != nil {
		return nil, err
	}

	total, err := u.inventoryRepo.CountPlayerItems(pctx, playerId)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return &models.PaginateRes{
			Data:  make([]*inventory.ItemInventory, 0),
			Tolal: total,
			Limit: req.Limit,
			First: models.FirstPaginate{
				Herf: "",
			},
			Next: models.NextPaginate{
				Start: "",
				Herf:  "",
			},
		}, nil
	}

	return &models.PaginateRes{
		Data:  make([]*inventory.ItemInventory, 0),
		Tolal: total,
		Limit: req.Limit,
		First: models.FirstPaginate{
			Herf: fmt.Sprintf("%s?limit=%d", basePaginatateUrl, req.Limit),
		},
		Next: models.NextPaginate{
			Start: results[len(results)-1].ItemId,

			Herf: fmt.Sprintf("%s?limit=%d&start=%s", basePaginatateUrl, req.Limit, results[len(results)-1].ItemId),
		},
	}, nil
}
