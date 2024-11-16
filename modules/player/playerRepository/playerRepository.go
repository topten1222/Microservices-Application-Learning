package playerRepository

import (
	"context"
	"log"
	"time"

	"github.com/topten1222/hello_sekai/modules/player"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	PlayerRepositoryService interface {
		IsUniquePlayer(context.Context, string, string) bool
		InsertOnePlayer(context.Context, *player.Player) (primitive.ObjectID, error)
	}
	playerRepository struct {
		db *mongo.Client
	}
)

func NewPlayerRepository(db *mongo.Client) PlayerRepositoryService {
	return &playerRepository{db: db}
}

func (r *playerRepository) playerDbConn(pctx context.Context) *mongo.Database {
	return r.db.Database("player_db")
}

func (r *playerRepository) IsUniquePlayer(pctx context.Context, email, username string) bool {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.playerDbConn(ctx)
	col := db.Collection("players")

	player := new(player.Player)
	if err := col.FindOne(ctx, bson.M{
		"$or": []bson.M{
			{"username": username},
			{"email": email},
		},
	}).Decode(player); err != nil {
		log.Printf("Error: IsUniquePLayer: %s", err.Error())
		return true
	}
	return false
}

func (r *playerRepository) InsertOnePlayer(pctx context.Context, player *player.Player) (primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()
	db := r.playerDbConn(ctx)
	col := db.Collection("players")

	playerId, err := col.InsertOne(ctx, player)
	if err != nil {
		log.Printf("Error: InsertOnePlayer: %s", err.Error())
		return primitive.NilObjectID, err
	}
	return playerId.InsertedID.(primitive.ObjectID), nil
}
