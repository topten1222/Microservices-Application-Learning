package middlewareUsecase

import middlewarerepository "github.com/topten1222/hello_sekai/modules/middleware/middlewareRepository"

type (
	MiddlewareUsecaseService interface{}
	middlewareUsecase        struct {
		middlewarerepository middlewarerepository.MiddlewareRepositoryHandlerService
	}
)

func NewMiddlewareUsecase(middlewareRepository middlewarerepository.MiddlewareRepositoryHandlerService) MiddlewareUsecaseService {
	return &middlewareUsecase{middlewarerepository: middlewareRepository}
}
