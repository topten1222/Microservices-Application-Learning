package authRepository

import (
	"context"
	"errors"
	"log"
	"time"

	playerPb "github.com/topten1222/hello_sekai/modules/player/playerPb"
	"github.com/topten1222/hello_sekai/pkg/grpccon"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	AuthrepositoryService interface{}

	authrepository struct {
		db *mongo.Client
	}
)

func NewAuthRepository(db *mongo.Client) AuthrepositoryService {
	return &authrepository{db: db}
}

func (r *authrepository) authDbConn(pctx context.Context) *mongo.Database {
	return r.db.Database("auth_db")
}

func (r *authrepository) CredentialSearch(pctx context.Context, grpcUrl string, req *playerPb.CredentialSearchReq) (*playerPb.PlayerProfile, error) {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()
	conn, err := grpccon.NewGrpcClient(grpcUrl)
	if err != nil {
		log.Printf("Error: grpc connection faild: %s", err.Error())
		return nil, errors.New("error: grpc connection faild")

	}
	result, err := conn.Player().CredentialSearch(ctx, req)
	if err != nil {
		log.Printf("Error: CredentialSearch failed: %s", err.Error())
		return nil, err
	}

	return result, nil
}
