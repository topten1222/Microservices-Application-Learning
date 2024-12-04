package itemUsecase

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/topten1222/hello_sekai/modules/item"
	"github.com/topten1222/hello_sekai/modules/item/itemRepository"
	"github.com/topten1222/hello_sekai/modules/models"
	"github.com/topten1222/hello_sekai/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	ItemUsecaseService interface {
		CreateItem(context.Context, *item.CreateItemReq) (any, error)
		FindOneItem(context.Context, string) (*item.ItemShowCase, error)
		FindManyItems(context.Context, string, *item.ItemSearchReq) (*models.PaginateRes, error)
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
	return u.FindOneItem(pctx, itemId.Hex())
}

func (u *itemUsecase) FindOneItem(pctx context.Context, itemId string) (*item.ItemShowCase, error) {
	result, err := u.itemRepo.FindOneItem(pctx, itemId)
	if err != nil {
		return nil, err
	}
	return &item.ItemShowCase{
		ItemId:   fmt.Sprintf("item:%s", result.Id.Hex()),
		Title:    result.Title,
		Price:    result.Price,
		Damage:   result.Damage,
		ImageUrl: result.ImageUrl,
	}, nil
}

func (u *itemUsecase) FindManyItems(pctx context.Context, basePaginationUrl string, req *item.ItemSearchReq) (*models.PaginateRes, error) {
	findItemFilter := bson.D{}
	findItemOpts := make([]*options.FindOptions, 0)
	counItemsFilter := bson.D{}

	if req.Start != "" {
		req.Start = strings.TrimPrefix(req.Start, "item:")
		findItemFilter = append(findItemFilter, bson.E{
			"_id", bson.D{{"$gt", utils.ConvertToObjectId(req.Start)}}},
		)
	}
	if req.Title != "" {
		findItemFilter = append(findItemFilter, bson.E{
			"title", primitive.Regex{Pattern: req.Title, Options: "i"},
		})
		counItemsFilter = append(counItemsFilter, bson.E{
			"title", primitive.Regex{Pattern: req.Title, Options: "i"},
		})
	}
	findItemFilter = append(findItemFilter, bson.E{"usage_status", true})
	counItemsFilter = append(counItemsFilter, bson.E{"usage_status", true})

	findItemOpts = append(findItemOpts, options.Find().SetSort(bson.D{{"_d", 1}}))
	findItemOpts = append(findItemOpts, options.Find().SetLimit(int64(req.Limit)))

	results, err := u.itemRepo.FindManyItems(pctx, findItemFilter, findItemOpts)
	if err != nil {
		return nil, err
	}
	total, err := u.itemRepo.CountItems(pctx, counItemsFilter)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return &models.PaginateRes{
			Data:  make([]*item.ItemShowCase, 0),
			Tolal: total,
			Limit: req.Limit,
			First: models.FirstPaginate{
				Herf: fmt.Sprintf("%s?limit=%d&title=%s", basePaginationUrl, req.Limit, req.Title),
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
			Herf: fmt.Sprintf("%s?limit=%d&title=%s", basePaginationUrl, req.Limit, req.Title),
		},
		Next: models.NextPaginate{
			Start: results[len(results)-1].ItemId,
			Herf:  fmt.Sprintf("%s?limit=%d&title=%s&start=%s", basePaginationUrl, req.Limit, req.Title, results[len(results)-1].ItemId),
		},
	}, nil
}
