package middlewareRepository

import (
	"context"
	"errors"
	"log"
	"time"

	authPb "github.com/topten1222/hello_sekai/modules/auth/authPb"
	"github.com/topten1222/hello_sekai/pkg/grpccon"
)

type (
	MiddlewareRepositoryHandlerService interface {
		AccessTokenSearch(pctx context.Context, grpcUrl, accessToken string) error
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
