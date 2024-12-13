package inventoryUsecase

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/topten1222/hello_sekai/config"
	"github.com/topten1222/hello_sekai/modules/inventory"
	"github.com/topten1222/hello_sekai/modules/inventory/inventoryRepository"
	"github.com/topten1222/hello_sekai/modules/item"
	itemPb "github.com/topten1222/hello_sekai/modules/item/itemPb"
	"github.com/topten1222/hello_sekai/modules/models"
	"github.com/topten1222/hello_sekai/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	InventoryUsecaseService interface {
		FindPlayerItems(context.Context, *config.Config, string, *inventory.InventorySearchReq) (*models.PaginateRes, error)
		GetOffset(context.Context) (int64, error)
		UpserOffset(context.Context, int64) error
	}

	inventoryUsecase struct {
		inventoryRepo inventoryRepository.InventoryRepositoryService
	}
)

func NewInventoryUsecasee(inventoryRepo inventoryRepository.InventoryRepositoryService) InventoryUsecaseService {
	return &inventoryUsecase{inventoryRepo: inventoryRepo}
}

func (u *inventoryUsecase) FindPlayerItems(pctx context.Context, cfg *config.Config, playerId string, req *inventory.InventorySearchReq) (*models.PaginateRes, error) {

	findItemFilter := bson.D{}
	findItemsOpts := make([]*options.FindOptions, 0)
	findItemsOpts = append(findItemsOpts, options.Find().SetSort(bson.D{{Key: "_id", Value: 1}}))
	findItemsOpts = append(findItemsOpts, options.Find().SetLimit(int64(req.Limit)))

	if req.Start != "" {
		findItemFilter = append(findItemFilter, bson.E{Key: "_id", Value: bson.D{{Key: "$gt", Value: utils.ConvertToObjectId(req.Start)}}})
	}
	findItemFilter = append(findItemFilter, bson.E{Key: "player_id", Value: playerId})

	inventoryData, err := u.inventoryRepo.FindPlayerItems(pctx, findItemFilter, findItemsOpts)
	if err != nil {
		return nil, err
	}
	fmt.Println("Befor ItemData: ")
	fmt.Println("GRPC ITEM")

	itemData, err := u.inventoryRepo.FindItemInIds(pctx, cfg.Grpc.ItemUrl, &itemPb.FindItemsInIdsReq{
		Ids: func() []string {
			itemIds := make([]string, 0)
			for _, v := range inventoryData {
				itemIds = append(itemIds, v.ItemId)
			}
			fmt.Println("ItemIds:: ", itemIds)
			return itemIds
		}(),
	})
	if err != nil {
		log.Printf("Error Find Item Ids: %v", err.Error())
		return nil, err
	}

	fmt.Println("ItemData: ", itemData)

	itemMaps := make(map[string]*item.ItemShowCase)
	for _, v := range itemData.Items {
		itemMaps[v.Id] = &item.ItemShowCase{
			ItemId:   v.Id,
			Title:    v.Title,
			Price:    v.Price,
			ImageUrl: v.ImageUrl,
			Damage:   int(v.Damage),
		}
	}
	results := make([]*inventory.ItemInventory, 0)
	for _, v := range inventoryData {
		itemShowCase := &item.ItemShowCase{}
		itemId := strings.TrimPrefix(v.ItemId, "item:")
		if itemMaps[itemId] != nil {
			itemShowCase = &item.ItemShowCase{
				ItemId:   v.ItemId,
				Title:    itemMaps[itemId].Title,
				Price:    itemMaps[itemId].Price,
				Damage:   itemMaps[itemId].Damage,
				ImageUrl: itemMaps[itemId].ImageUrl,
			}
		}
		results = append(results, &inventory.ItemInventory{
			InventoryId:  v.Id.Hex(),
			PlayerId:     v.PlayerId,
			ItemShowCase: itemShowCase,
		})
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
		Data:  results,
		Tolal: total,
		Limit: req.Limit,
		First: models.FirstPaginate{
			Herf: fmt.Sprintf("%s/%s?limit=%d", cfg.Paginate.InventoryNextPageBaseUrl, playerId, req.Limit),
		},
		Next: models.NextPaginate{
			Start: results[len(results)-1].InventoryId,

			Herf: fmt.Sprintf("%s/%s?limit=%d&start=%s", cfg.Paginate.InventoryNextPageBaseUrl, playerId, req.Limit, results[len(results)-1].InventoryId),
		},
	}, nil
}

func (u *inventoryUsecase) GetOffset(pctx context.Context) (int64, error) {
	return u.inventoryRepo.GetOffset(pctx)
}

func (u *inventoryUsecase) UpserOffset(pctx context.Context, offset int64) error {
	return u.inventoryRepo.UpserOffset(pctx, offset)
}
