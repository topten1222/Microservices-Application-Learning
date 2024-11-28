package middlewareRepository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	authPb "github.com/topten1222/hello_sekai/modules/auth/authPb"
	"github.com/topten1222/hello_sekai/pkg/grpccon"
	"github.com/topten1222/hello_sekai/pkg/jwtauth"
)

type (
	MiddlewareRepositoryHandlerService interface {
		AccessTokenSearch(context.Context, string, string) error
		RolesCount(context.Context, string) (int64, error)
	}

	middlewaRerepository struct{}
)

func NewMiddlewareRepository() MiddlewareRepositoryHandlerService {
	return &middlewaRerepository{}
}

func (r *middlewaRerepository) AccessTokenSearch(pctx context.Context, grpcUrl, accessToken string) error {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()
	conn, err := grpccon.NewGrpcClient(grpcUrl)
	if err != nil {
		log.Printf("Error: grpce connection faild %s", err)

		return err
	}
	result, err := conn.Auth().CredentialSearch(ctx, &authPb.AccessTokenSearchReq{
		AccessToken: accessToken,
	})
	if err != nil {
		log.Printf("Error: access token connection faild %s", err)
		return err
	}
	if result == nil {
		return errors.New("error: access token is invalid")
	}
	if !result.InValid {
		return errors.New("error: access token is invalid")
	}
	return nil
}

func (r *middlewaRerepository) RolesCount(pctx context.Context, grpcUrl string) (int64, error) {
	ctx, cancel := context.WithTimeout(pctx, 30*time.Second)
	defer cancel()

	conn, err := grpccon.NewGrpcClient(grpcUrl)
	if err != nil {
		log.Printf("Error: gRPC connection failed: %s", err.Error())
		return -1, errors.New("error: gRPC connection failed")
	}

	jwtauth.SetApiKeyInContext(&ctx)
	result, err := conn.Auth().RolesCount(ctx, &authPb.RoleCountReq{})
	fmt.Println("Result ROLE COUNT::: ", result)
	if err != nil {
		log.Printf("Error: CredentialSearch failed: %s", err.Error())
		return -1, errors.New("error: email or password is incorrect")
	}

	if result == nil {
		log.Printf("Error: roles count failed")
		return -1, errors.New("error: roles count failed")
	}

	return result.Count, nil
}
