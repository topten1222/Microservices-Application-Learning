package itemRepository

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type (
	ItemRepositoryService interface{}

	itemRepositroy struct {
		db *mongo.Client
	}
)

func NewItemRepository(db *mongo.Client) ItemRepositoryService {
	return &itemRepositroy{db: db}
}

func (r *itemRepositroy) itemDbConn(pctx context.Context) *mongo.Database {
	return r.db.Database("item_db")
}
