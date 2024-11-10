package server

import (
	"github.com/topten1222/hello_sekai/modules/auth/authHandler"
	"github.com/topten1222/hello_sekai/modules/auth/authRepository"
	"github.com/topten1222/hello_sekai/modules/auth/authUsecase"
)

func (s *server) authService() {
	authRepo := authRepository.NewAuthRepository(s.db)
	authUsecase := authUsecase.NewAuthUsecase(authRepo)
	authHandler.NewAuthHttpHandler(s.cfg, authUsecase)
	authHandler.NewAuthGrpcHandler(authUsecase)

	auth := s.app.Group("/auth_v1")

	auth.GET("/", s.healthCheckService)
	// Health Check

}
