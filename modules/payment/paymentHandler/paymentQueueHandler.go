package paymentHandler

import (
	"github.com/topten1222/hello_sekai/config"
	"github.com/topten1222/hello_sekai/modules/payment/paymentUsecase"
)

type (
	PaymentQueueHandlerService interface{}

	paymentQueueHandler struct {
		cfg            *config.Config
		paymentUsecase paymentUsecase.PaymentUsecaseService
	}
)

func NewPaymentQueueHandler(cfg *config.Config, paymentUsecase paymentUsecase.PaymentUsecaseService) PaymentQueueHandlerService {
	return &paymentQueueHandler{cfg: cfg, paymentUsecase: paymentUsecase}
}
