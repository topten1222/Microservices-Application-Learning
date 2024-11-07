package authUsecase

import authrepository "github.com/topten1222/hello_sekai/modules/auth/authRepository"

type (
	AuthusecaseService interface{}

	authusecase struct {
		authRepo authrepository.AuthrepositoryService
	}
)

func NewAuthUsecase(authRepo authrepository.AuthrepositoryService) AuthusecaseService {
	return &authusecase{authRepo: authRepo}
}
