package itemHandler

import (
	"context"

	itemPb "github.com/topten1222/hello_sekai/modules/item/itemPb"
	"github.com/topten1222/hello_sekai/modules/item/itemUsecase"
)

type (
	itemGrpcHandler struct {
		itemUsecase itemUsecase.ItemUsecaseService
		itemPb.ItemGrpcServiceServer
	}
)

func NewItemGrpcHandler(itemUsecase itemUsecase.ItemUsecaseService) *itemGrpcHandler {
	return &itemGrpcHandler{itemUsecase: itemUsecase}
}

func (g *itemGrpcHandler) FindItemsInIds(ctx context.Context, req *itemPb.FindItemsInIdsReq) (*itemPb.FindItemsInIdsRes, error) {
	return nil, nil
}
