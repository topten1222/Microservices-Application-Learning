package migration

import (
	"context"
	"log"

	"github.com/topten1222/hello_sekai/config"
	"github.com/topten1222/hello_sekai/pkg/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func inventoryDbConn(pctx context.Context, cfg *config.Config) *mongo.Database {
	return database.DbConnect(pctx, cfg).Database("inventory_db")
}

func InventoryMigrate(pctx context.Context, cfg *config.Config) {
	db := inventoryDbConn(pctx, cfg)
	defer db.Client().Disconnect(pctx)
	col := db.Collection("players_inventory")
	indexs, _ := col.Indexes().CreateMany(pctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "player_id", Value: 1}, {Key: "item_id", Value: 1}}},
	})

	for _, index := range indexs {
		log.Printf("Create Index : %v", index)
	}

	col = db.Collection("players_inventory_queue")

	results, err := col.InsertOne(pctx, bson.M{"offset": -1})
	if err != nil {
		log.Printf("Error insert: %v", err)
	}
	log.Printf("Migrate auth completed: %v", results)
}
