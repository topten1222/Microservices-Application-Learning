package server

import (
	"log"

	"github.com/topten1222/hello_sekai/modules/auth/authHandler"
	authPb "github.com/topten1222/hello_sekai/modules/auth/authPb"
	"github.com/topten1222/hello_sekai/modules/auth/authRepository"
	"github.com/topten1222/hello_sekai/modules/auth/authUsecase"
	"github.com/topten1222/hello_sekai/pkg/grpccon"
)

func (s *server) authService() {
	authRepo := authRepository.NewAuthRepository(s.db)
	authUsecase := authUsecase.NewAuthUsecase(authRepo)
	httpHandler := authHandler.NewAuthHttpHandler(s.cfg, authUsecase)
	authHandler.NewAuthGrpcHandler(authUsecase)

	//grpc

	go func() {
		grpcServer, list := grpccon.NewGrpcServer(&s.cfg.Jwt, s.cfg.Grpc.AuthUrl)
		authPb.RegisterAuthGrpcServiceServer(grpcServer, authHandler.NewAuthGrpcHandler(authUsecase))
		log.Printf("Auth Grpc server listening on %s", s.cfg.Grpc.AuthUrl)
		grpcServer.Serve(list)
	}()

	auth := s.app.Group("/auth_v1")

	auth.GET("/", s.healthCheckService)
	auth.POST("/auth/login", httpHandler.Login)
	auth.POST("/auth/refresh-token", httpHandler.RefreshToken)
	auth.POST("/auth/logout", httpHandler.Logout)
	// Health Check

}
