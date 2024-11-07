package server

import (
	"github.com/topten1222/hello_sekai/modules/payment/paymentHandler"
	"github.com/topten1222/hello_sekai/modules/payment/paymentRepository"
	"github.com/topten1222/hello_sekai/modules/payment/paymentUsecase"
)

func (s *server) paymentService() {
	repo := paymentRepository.NewPaymentRepository(s.db)
	usecase := paymentUsecase.NewPaymentUsecase(repo)
	paymentHandler.NewPaymentHttpHandler(s.cfg, usecase)
	paymentHandler.NewPaymentGrpcHandler(usecase)
	paymentHandler.NewPaymentQueueHandler(s.cfg, usecase)

	s.app.Group("/payment")
}
