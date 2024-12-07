package inventoryRepository

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/topten1222/hello_sekai/modules/inventory"
	itemPb "github.com/topten1222/hello_sekai/modules/item/itemPb"
	"github.com/topten1222/hello_sekai/pkg/grpccon"
	"github.com/topten1222/hello_sekai/pkg/jwtauth"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	InventoryRepositoryService interface {
		FindItemInIds(context.Context, string, itemPb.FindItemsInIdsReq) (*itemPb.FindItemsInIdsRes, error)
		FindPlayerItems(context.Context, primitive.D, []*options.FindOptions) ([]*inventory.Inventory, error)
		CountPlayerItems(context.Context, string) (int64, error)
	}

	inventoryRepository struct {
		db *mongo.Client
	}
)

func NewInventoryRepository(db *mongo.Client) InventoryRepositoryService {
	return &inventoryRepository{db: db}
}

func (r *inventoryRepository) inventoryDbConn(pctx context.Context) *mongo.Database {
	return r.db.Database("inventory_db")
}

func (r *inventoryRepository) FindItemInIds(pctx context.Context, grpcUrl string, req itemPb.FindItemsInIdsReq) (*itemPb.FindItemsInIdsRes, error) {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	jwtauth.SetApiKeyInContext(&ctx)
	conn, err := grpccon.NewGrpcClient(grpcUrl)
	if err != nil {
		log.Printf("Error: grpc connection faild: %s", err.Error())
		return nil, errors.New("erorr: grpc connection faild")
	}

	result, err := conn.Item().FindItemsInIds(ctx, &req)
	if err != nil {
		log.Printf("Error: FindItensInIds Faild %s", err.Error())
		return nil, errors.New("error: find items in ids faild")
	}
	if result == nil {
		log.Printf("Error: FindItemsInIds faild %s", err.Error())
		return nil, errors.New("error: items not found")
	}
	if len(result.Items) == 0 {
		log.Printf("Error: FindItemsInIds faild %s", err.Error())
		return nil, errors.New("error: items not found")
	}

	return result, nil
}

func (r *inventoryRepository) FindPlayerItems(pctx context.Context, filter primitive.D, opts []*options.FindOptions) ([]*inventory.Inventory, error) {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	col := r.inventoryDbConn(ctx).Collection("players_inventory")

	cursors, err := col.Find(ctx, filter, opts...)
	if err != nil {
		log.Printf("Error: FindPlayerItems faild %s", err.Error())
		return nil, errors.New("error: find player items faild")
	}
	results := make([]*inventory.Inventory, 0)
	for cursors.Next(ctx) {
		result := new(inventory.Inventory)
		if err := cursors.Decode(result); err != nil {
			log.Printf("Error: findPlayerItems faild %s", err.Error())
			return nil, errors.New("error: player item not found")
		}
		results = append(results, result)
	}
	return results, nil
}

func (r *inventoryRepository) CountPlayerItems(pctx context.Context, playerId string) (int64, error) {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.inventoryDbConn(ctx)
	col := db.Collection("players_inventory")

	count, err := col.CountDocuments(ctx, bson.M{"player_id": playerId})
	if err != nil {
		log.Printf("Error: CountItems faild %s", err.Error())
		return -1, errors.New("error: count items faild")
	}
	return count, nil
}
