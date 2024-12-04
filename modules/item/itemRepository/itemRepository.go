package itemRepository

import (
	"context"
	"log"
	"time"

	"github.com/topten1222/hello_sekai/modules/item"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	ItemRepositoryService interface {
		InsertOneItem(context.Context, *item.Item) (primitive.ObjectID, error)
		IsUniqueItem(context.Context, string) bool
	}

	itemRepository struct {
		db *mongo.Client
	}
)

func NewItemRepository(db *mongo.Client) ItemRepositoryService {
	return &itemRepository{db: db}
}

func (r *itemRepository) itemDbConn(pctx context.Context) *mongo.Database {
	return r.db.Database("item_db")
}

func (r *itemRepository) IsUniqueItem(pctx context.Context, title string) bool {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.itemDbConn(ctx)
	col := db.Collection("items")

	result := new(item.Item)

	if err := col.FindOne(ctx, bson.M{"title": title}).Decode(result); err != nil {
		log.Printf("Error: is unique player %s", err.Error())
		return true
	}
	return false
}

func (r *itemRepository) InsertOneItem(pctx context.Context, req *item.Item) (primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.itemDbConn(ctx)
	col := db.Collection("items")

	itemId, err := col.InsertOne(ctx, req)
	if err != nil {
		log.Printf("Error: insert one item %s", err.Error())
		return primitive.NilObjectID, err
	}

	return itemId.InsertedID.(primitive.ObjectID), nil
}
