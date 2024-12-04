package itemRepository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/topten1222/hello_sekai/modules/item"
	"github.com/topten1222/hello_sekai/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	ItemRepositoryService interface {
		InsertOneItem(context.Context, *item.Item) (primitive.ObjectID, error)
		IsUniqueItem(context.Context, string) bool
		FindOneItem(context.Context, string) (*item.Item, error)
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

func (r *itemRepository) FindOneItem(pctx context.Context, itemId string) (*item.Item, error) {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.itemDbConn(ctx)
	col := db.Collection("items")

	result := new(item.Item)
	fmt.Println(itemId)
	if err := col.FindOne(ctx, bson.M{"_id": utils.ConvertToObjectId(itemId)}).Decode(result); err != nil {
		log.Printf("Error: find one item %s", err.Error())
		return nil, errors.New("error: item not found")
	}
	return result, nil
}
