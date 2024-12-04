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
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	ItemRepositoryService interface {
		InsertOneItem(context.Context, *item.Item) (primitive.ObjectID, error)
		IsUniqueItem(context.Context, string) bool
		FindOneItem(context.Context, string) (*item.Item, error)
		FindManyItems(context.Context, primitive.D, []*options.FindOptions) ([]*item.ItemShowCase, error)
		CountItems(context.Context, primitive.D) (int64, error)
		UpdateOneItem(context.Context, string, primitive.M) error
		EnableOrDisableItem(context.Context, string, bool) error
	}

	itemRepository struct {
		db *mongo.Client
	}
)

func NewItemRepository(db *mongo.Client) ItemRepositoryService {
	return &itemRepository{db: db}
}

func (r *itemRepository) itemDbConn(_ context.Context) *mongo.Database {
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

func (r *itemRepository) FindManyItems(pctx context.Context, filter primitive.D, opts []*options.FindOptions) ([]*item.ItemShowCase, error) {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()
	db := r.itemDbConn(ctx)
	col := db.Collection("items")

	cursors, err := col.Find(ctx, filter, opts...)
	if err != nil {
		log.Printf("Error: find many items %s", err.Error())

		return make([]*item.ItemShowCase, 0), errors.New("error: Find Many items faild")
	}
	results := make([]*item.ItemShowCase, 0)
	for cursors.Next(ctx) {
		result := new(item.Item)
		if err := cursors.Decode(result); err != nil {
			log.Printf("Error: Find many item faild %s", err.Error())
			return make([]*item.ItemShowCase, 0), errors.New("error: find many item faild")
		}
		results = append(results, &item.ItemShowCase{
			ItemId:   result.Id.Hex(),
			Price:    result.Price,
			Title:    result.Title,
			Damage:   result.Damage,
			ImageUrl: result.ImageUrl,
		})
	}

	return results, nil
}

func (r *itemRepository) CountItems(pctx context.Context, filter primitive.D) (int64, error) {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.itemDbConn(ctx)
	col := db.Collection("items")

	count, err := col.CountDocuments(ctx, filter)
	if err != nil {
		log.Printf("Error: count items %s", err.Error())
		return -1, err
	}

	return count, nil
}

func (r *itemRepository) UpdateOneItem(pctx context.Context, itemId string, req primitive.M) error {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.itemDbConn(ctx)
	col := db.Collection("items")
	result, err := col.UpdateOne(ctx, bson.M{"_id": utils.ConvertToObjectId(itemId)}, bson.M{"$set": req})
	if err != nil {
		log.Printf("Error: Update one item %s", err.Error())
		return errors.New("error: update one item")
	}
	log.Printf("success update %v", result)
	return nil
}

func (r *itemRepository) EnableOrDisableItem(pctx context.Context, itemId string, status bool) error {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.itemDbConn(ctx)
	col := db.Collection("items")

	result, err := col.UpdateOne(ctx, bson.M{"_id": utils.ConvertToObjectId(itemId)}, bson.M{"$set": bson.M{"usage_status": status}})
	if err != nil {
		log.Printf("Error: Enable or Disable item %s", err.Error())
		return errors.New("error: Enable or Disable item faild")
	}
	log.Printf("success update %v", result)
	return nil
}
