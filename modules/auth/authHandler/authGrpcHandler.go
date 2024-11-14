package authHandler

import (
	"context"

	authPb "github.com/topten1222/hello_sekai/modules/auth/authPb"
	authusecase "github.com/topten1222/hello_sekai/modules/auth/authUsecase"
)

type (
	AuthGrpHandler struct {
		authPb.UnimplementedAuthGrpcServiceServer
		authUsecase authusecase.AuthusecaseService
	}
)

func NewAuthGrpcHandler(authUsecase authusecase.AuthusecaseService) *AuthGrpHandler {
	return &AuthGrpHandler{authUsecase: authUsecase}
}

func (g *AuthGrpHandler) CredentialSearch(ctx context.Context, req *authPb.AccessTokenSearchReq) (*authPb.AccessTokenSearchRes, error) {
	return nil, nil
}

func (g *AuthGrpHandler) RolesCount(ctx context.Context, req *authPb.RoleCountReq) (*authPb.RoleCountRes, error) {
	return nil, nil
}
