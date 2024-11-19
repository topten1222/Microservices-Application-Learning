package playerRepository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/topten1222/hello_sekai/modules/player"
	"github.com/topten1222/hello_sekai/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	PlayerRepositoryService interface {
		IsUniquePlayer(context.Context, string, string) bool
		InsertOnePlayer(context.Context, *player.Player) (primitive.ObjectID, error)
		FindOnePlayerProfile(context.Context, string) (*player.PlayerProfileBson, error)
		InsertOnePlayerTransaction(context.Context, *player.PlayerTransaction) error
		GetPlayerSavingAccount(context.Context, string) (*player.PlayerSavingAccount, error)
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

func (r *playerRepository) FindOnePlayerProfile(pctx context.Context, playerId string) (*player.PlayerProfileBson, error) {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.playerDbConn(ctx)
	col := db.Collection("players")

	result := new(player.PlayerProfileBson)
	err := col.FindOne(
		ctx,
		bson.M{"_id": utils.ConvertToObjectId(playerId)},
		options.FindOne().SetProjection(bson.M{
			"_id":        1,
			"email":      1,
			"username":   1,
			"created_at": 1,
			"updated_at": 1,
		})).Decode(result)
	if err != nil {
		log.Printf("Error: FindOnePlayerProfile: %s", err.Error())
		return nil, errors.New("Error Player Not Found")
	}
	return result, nil
}

func (r *playerRepository) InsertOnePlayerTransaction(pctx context.Context, playerTransaction *player.PlayerTransaction) error {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()
	db := r.playerDbConn(ctx)
	col := db.Collection("player_transections")
	result, err := col.InsertOne(ctx, playerTransaction)
	if err != nil {
		log.Printf("Erro: InsertOnePlayerTransaction: %s", err.Error())
		return errors.New("error: insert one player transaction failed")
	}
	log.Printf("Result: InsertOnePlayerTransaction: %v", result.InsertedID.(primitive.ObjectID))
	return nil
}

func (r *playerRepository) GetPlayerSavingAccount(pctx context.Context, playerId string) (*player.PlayerSavingAccount, error) {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()
	db := r.playerDbConn(ctx)
	col := db.Collection("player_transections")
	fmt.Println(playerId)
	filter := bson.A{
		bson.D{{Key: "$match", Value: bson.D{{Key: "player_id", Value: playerId}}}},
		bson.D{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$player_id"},
			{Key: "balance", Value: bson.D{{Key: "$sum", Value: "$amount"}}},
		}}},
		bson.D{
			{Key: "$project", Value: bson.D{
				{Key: "player_id", Value: "$_id"},
				{Key: "_id", Value: 0},
				{Key: "balance", Value: 1},
			},
			}},
	}
	cursor, err := col.Aggregate(ctx, filter)
	if err != nil {
		log.Printf("Error: GetPlayerSavingAccount: %s", err.Error())

		return nil, errors.New("error: faild to get player saving account")
	}

	result := new(player.PlayerSavingAccount)
	for cursor.Next(ctx) {
		if err := cursor.Decode(result); err != nil {
			log.Printf("Error: GetPlayerSavingAccount: %s", err.Error())
			return nil, errors.New("error: faild to get player saving account")

		}

	}
	fmt.Println("ererer", result)
	return result, nil
}
