package migration

import (
	"context"

	"github.com/labstack/gommon/log"
	"github.com/topten1222/hello_sekai/config"
	"github.com/topten1222/hello_sekai/pkg/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func paymentDbConn(pctx context.Context, cfg *config.Config) *mongo.Database {
	return database.DbConnect(pctx, cfg).Database("payment_db")
}

func PaymentMigrate(pctx context.Context, cfg *config.Config) {
	db := paymentDbConn(pctx, cfg)
	defer db.Client().Disconnect(pctx)
	col := db.Collection("payment_queue")

	results, err := col.InsertOne(pctx, bson.M{"offset": -1})
	if err != nil {
		log.Printf("Error inser %v", err)
	}
	log.Printf("Migrate payment complete: ", results)
}
