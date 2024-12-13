package paymentUsecase

import (
	"context"
	"errors"
	"log"

	"github.com/topten1222/hello_sekai/modules/item"
	itemPb "github.com/topten1222/hello_sekai/modules/item/itemPb"
	"github.com/topten1222/hello_sekai/modules/payment"
	"github.com/topten1222/hello_sekai/modules/payment/paymentRepository"
)

type (
	PaymentUsecaseService interface {
		FindItemsInIds(context.Context, string, []*payment.ItemServiceReqDatum) error
		GetOffset(context.Context) (int64, error)
		UpserOffset(context.Context, int64) error
	}

	paymentUsecase struct {
		paymentRepo paymentRepository.PaymentRepositoryService
	}
)

func NewPaymentUsecase(paymentRepo paymentRepository.PaymentRepositoryService) PaymentUsecaseService {
	return &paymentUsecase{paymentRepo: paymentRepo}
}

func (u *paymentUsecase) FindItemsInIds(pctx context.Context, grpcUrl string, req []*payment.ItemServiceReqDatum) error {
	setIds := make(map[string]bool)
	for _, v := range req {
		if !setIds[v.ItemId] {
			setIds[v.ItemId] = true
		}
	}

	itemData, err := u.paymentRepo.FindItemsInIds(pctx, grpcUrl, &itemPb.FindItemsInIdsReq{
		Ids: func() []string {
			itemIds := make([]string, 0)
			for k, _ := range setIds {
				itemIds = append(itemIds, k)
			}
			return itemIds
		}(),
	})
	if err != nil {
		log.Printf("Error: FindItemsInIds: %v", err)
		return errors.New("error: FindItemsInIds")
	}

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

	for i := range req {
		if _, ok := itemMaps[req[i].ItemId]; !ok {
			return errors.New("item not found")
		}
		req[i].Price = itemMaps[req[i].ItemId].Price
	}

	return nil
}

func (u *paymentUsecase) GetOffset(pctx context.Context) (int64, error) {
	return u.paymentRepo.GetOffset(pctx)
}

func (u *paymentUsecase) UpserOffset(pctx context.Context, offset int64) error {
	return u.paymentRepo.UpserOffset(pctx, offset)
}
