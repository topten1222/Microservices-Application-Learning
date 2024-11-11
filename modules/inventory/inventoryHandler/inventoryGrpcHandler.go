package inventoryHandler

import (
	"context"

	inventoryPb "github.com/topten1222/hello_sekai/modules/inventory/inventoryPb"
	"github.com/topten1222/hello_sekai/modules/inventory/inventoryUsecase"
)

type (
	inventoryGrpcHanlderService struct {
		inventoryUsecase inventoryUsecase.InventoryUsecaseService
		inventoryPb.UnimplementedInventoryGrpcServiceServer
	}
)

func NewInventoryGrpcHandler(inventoryUsecase inventoryUsecase.InventoryUsecaseService) *inventoryGrpcHanlderService {
	return &inventoryGrpcHanlderService{inventoryUsecase: inventoryUsecase}
}

func (g *inventoryGrpcHanlderService) IsAvailableToSell(ctx context.Context, req *inventoryPb.IsAvailableToSellToReq) (*inventoryPb.IsAvailableToSellToRes, error) {
	return nil, nil
}
