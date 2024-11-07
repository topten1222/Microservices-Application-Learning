package paymentUsecase

import "github.com/topten1222/hello_sekai/modules/payment/paymentRepository"

type (
	PaymentUsecaseService interface{}

	paymentUsecase struct {
		paymentRepo paymentRepository.PaymentRepositoryService
	}
)

func NewPaymentUsecase(paymentRepo paymentRepository.PaymentRepositoryService) PaymentUsecaseService {
	return &paymentUsecase{paymentRepo: paymentRepo}
}
