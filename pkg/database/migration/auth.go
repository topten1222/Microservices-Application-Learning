package migration

import (
	"context"
	"log"

	"github.com/topten1222/hello_sekai/config"
	"github.com/topten1222/hello_sekai/modules/auth"
	"github.com/topten1222/hello_sekai/pkg/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func authDbConn(pctx context.Context, cfg *config.Config) *mongo.Database {
	return database.DbConnect(pctx, cfg).Database("auth_db")
}

func AuthMigrate(pctx context.Context, cfg *config.Config) {
	db := authDbConn(pctx, cfg)
	defer db.Client().Disconnect(pctx)

	col := db.Collection("auth")

	indexs, _ := col.Indexes().CreateMany(pctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "_id", Value: 1}}},
		{Keys: bson.D{{Key: "player_id", Value: 1}}},
		{Keys: bson.D{{Key: "refresh_token", Value: 1}}},
	})
	for _, index := range indexs {
		log.Printf("index:: %v", index)
	}

	col = db.Collection("roles")
	indexs, _ = col.Indexes().CreateMany(pctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "_id", Value: 1}}},
		{Keys: bson.D{{Key: "code", Value: 1}}},
	})

	for _, index := range indexs {
		log.Printf("index:: %v", index)
	}

	documents := func() []interface{} {
		roles := []*auth.Role{
			{
				Id:    primitive.NewObjectID(),
				Title: "Player",
				Code:  0,
			},
			{
				Id:    primitive.NewObjectID(),
				Title: "Admin",
				Code:  1,
			},
		}
		docs := make([]interface{}, 0)
		for _, r := range roles {
			docs = append(docs, r)
		}
		return docs
	}()
	result, err := col.InsertMany(pctx, documents)
	if err != nil {
		panic(err)
	}
	log.Println("Migrate auth complated: ", result)
}
