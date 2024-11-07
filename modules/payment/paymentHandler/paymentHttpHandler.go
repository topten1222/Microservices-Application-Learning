package paymentHandler

import (
	"github.com/topten1222/hello_sekai/config"
	"github.com/topten1222/hello_sekai/modules/payment/paymentUsecase"
)

type (
	PaymentHttpHandlerService interface{}

	paymentHttpHandler struct {
		cfg            *config.Config
		paymentUsecase paymentUsecase.PaymentUsecaseService
	}
)

func NewPaymentHttpHandler(cfg *config.Config, paymentUsecase paymentUsecase.PaymentUsecaseService) PaymentHttpHandlerService {
	return &paymentHttpHandler{cfg: cfg, paymentUsecase: paymentUsecase}
}
