package paymentHandler

import "github.com/topten1222/hello_sekai/modules/payment/paymentUsecase"

type (
	paymentGrpcHandler struct {
		paymentUsecase paymentUsecase.PaymentUsecaseService
	}
)

func NewPaymentGrpcHandler(paymentUsecase paymentUsecase.PaymentUsecaseService) *paymentGrpcHandler {
	return &paymentGrpcHandler{paymentUsecase: paymentUsecase}
}
