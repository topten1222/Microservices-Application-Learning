package middlewaRerepository

type (
	MiddlewareRepositoryHandlerService interface{}

	middlewaRerepository struct{}
)

func NewMiddlewareRepository() MiddlewareRepositoryHandlerService {
	return &middlewaRerepository{}
}
