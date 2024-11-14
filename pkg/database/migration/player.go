package migration

import (
	"context"
	"log"

	"github.com/topten1222/hello_sekai/config"
	"github.com/topten1222/hello_sekai/modules/player"
	"github.com/topten1222/hello_sekai/pkg/database"
	"github.com/topten1222/hello_sekai/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func playerDbConn(pctx context.Context, cfg *config.Config) *mongo.Database {
	return database.DbConnect(pctx, cfg).Database("player_db")
}

func PlayerMigrate(pctx context.Context, cfg *config.Config) {
	db := playerDbConn(pctx, cfg)
	defer db.Client().Disconnect(pctx)

	col := db.Collection("players")

	indexs, _ := col.Indexes().CreateMany(pctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "_id", Value: 1}}},
		{Keys: bson.D{{Key: "email", Value: 1}}},
	})

	for _, index := range indexs {
		log.Printf("Create Index : %v", index)
	}

	documents := func() []interface{} {
		players := []*player.Player{
			{
				Email:    "player0001@sekai.com",
				Password: "123456",
				Username: "player001",
				PlayerRoles: []player.PlayerRole{
					{
						RoleTitle: "player",
						RoleCode:  0,
					},
				},
				CreatedAt: utils.LocalTime(),
			},
			{
				Email:    "player0002@sekai.com",
				Password: "123456",
				Username: "player002",
				PlayerRoles: []player.PlayerRole{
					{
						RoleTitle: "player",
						RoleCode:  0,
					},
				},
				CreatedAt: utils.LocalTime(),
			},
			{
				Email:    "player0003@sekai.com",
				Password: "123456",
				Username: "player003",
				PlayerRoles: []player.PlayerRole{
					{
						RoleTitle: "player",
						RoleCode:  0,
					},
				},
				CreatedAt: utils.LocalTime(),
			},
			{
				Email:    "admin001@sekai.com",
				Password: "123456",
				Username: "admin001",
				PlayerRoles: []player.PlayerRole{
					{
						RoleTitle: "Player",
						RoleCode:  0,
					},
					{
						RoleTitle: "Admin",
						RoleCode:  1,
					},
				},
				CreatedAt: utils.LocalTime(),
			},
		}
		docs := make([]interface{}, 0)
		for _, r := range players {
			docs = append(docs, r)
		}
		return docs
	}()

	results, err := col.InsertMany(pctx, documents)
	if err != nil {
		log.Printf("Error insert:: %v", err)
	}
	log.Printf("Migrate player complete: %v", results)

	playerTransactions := make([]interface{}, 0)
	for _, p := range results.InsertedIDs {
		playerTransactions = append(playerTransactions, &player.PlayerTransaction{
			PlayerId:  "player:" + p.(primitive.ObjectID).Hex(),
			Amount:    1000,
			CreatedAt: utils.LocalTime(),
		})
	}

	col = db.Collection("player_transections")
	indexs, _ = col.Indexes().CreateMany(pctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "_id", Value: 1}}},
		{Keys: bson.D{{Key: "player_id", Value: 1}}},
	})
	for _, index := range indexs {
		log.Printf("Create Index : %v", index)
	}
	results, err = col.InsertMany(pctx, playerTransactions)
	if err != nil {
		log.Printf("Inser Player Err: %v", err)
	}
	log.Printf("Migrate player_transection %v", results)

	col = db.Collection("player_transaction_queue")
	result, err := col.InsertOne(pctx, bson.M{"offset": 1})
	if err != nil {
		log.Printf("Insert One Err: %v", err)
	}
	log.Printf("Migrate player transecation queue %v", result)
}
