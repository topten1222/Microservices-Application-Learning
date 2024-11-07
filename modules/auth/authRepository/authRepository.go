package authRepository

import (
	"context"

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
