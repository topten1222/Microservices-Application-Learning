package authRepository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/topten1222/hello_sekai/modules/auth"
	playerPb "github.com/topten1222/hello_sekai/modules/player/playerPb"
	"github.com/topten1222/hello_sekai/pkg/grpccon"
	"github.com/topten1222/hello_sekai/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	AuthrepositoryService interface {
		InsertOnePlayerCredentail(context.Context, *auth.Credential) (primitive.ObjectID, error)
		CredentialSearch(context.Context, string, *playerPb.CredentialSearchReq) (*playerPb.PlayerProfile, error)
		FindOnePlayerCredentail(context.Context, string) (*auth.Credential, error)
		FindOnePlayerProfileToRefresh(context.Context, string, *playerPb.FindOnePlayerProfileToRefreshReq) (*playerPb.PlayerProfile, error)
		UpdateOnePlayerCredentail(context.Context, string, *auth.UpdateRefreshTokenReq) error
		DeleteOnePlayerCredentail(context.Context, string) (int64, error)
		FindOneAccessToken(context.Context, string) (*auth.Credential, error)
		RolesCount(context.Context) (int64, error)
	}

	authrepository struct {
		db *mongo.Client
	}
)

func NewAuthRepository(db *mongo.Client) AuthrepositoryService {
	return &authrepository{db: db}
}

func (r *authrepository) authDbConn(pctx context.Context) *mongo.Database {
	return r.db.Database("auth_db")
}

func (r *authrepository) CredentialSearch(pctx context.Context, grpcUrl string, req *playerPb.CredentialSearchReq) (*playerPb.PlayerProfile, error) {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()
	conn, err := grpccon.NewGrpcClient(grpcUrl)
	if err != nil {
		log.Printf("Error: grpc connection faild: %s", err.Error())
		return nil, errors.New("error: grpc connection faild")

	}
	fmt.Println(req)
	result, err := conn.Player().CredentialSearch(ctx, req)
	if err != nil {
		log.Printf("Error: CredentialSearch failed: %s", err.Error())
		return nil, errors.New("Error invalid email password")
	}

	return result, nil
}

func (r *authrepository) FindOnePlayerProfileToRefresh(pctx context.Context, grpcUrl string, req *playerPb.FindOnePlayerProfileToRefreshReq) (*playerPb.PlayerProfile, error) {
	ctx, cancel := context.WithTimeout(pctx, 30*time.Second)
	defer cancel()

	conn, err := grpccon.NewGrpcClient(grpcUrl)
	if err != nil {
		log.Printf("Error: grpc connection failed: %s", err.Error())
		return nil, errors.New("error: grpc connection failed")

	}
	result, err := conn.Player().FindOnePlayerProfileToRefresh(ctx, req)
	if err != nil {
		log.Printf("Error: FindOnePlayerProfileToRefresh failed: %s", err.Error())
		return nil, errors.New("error: find one player profile to refresh failed")
	}
	return result, nil
}

func (r authrepository) InsertOnePlayerCredentail(ctx context.Context, req *auth.Credential) (primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	db := r.authDbConn(ctx)
	col := db.Collection("auth")

	result, err := col.InsertOne(ctx, req)
	if err != nil {
		log.Printf("Error: InsertOnePlayerCredentail failed: %s", err.Error())
		return primitive.NilObjectID, err
	}
	return result.InsertedID.(primitive.ObjectID), nil
}

func (r *authrepository) FindOnePlayerCredentail(ctx context.Context, credentialId string) (*auth.Credential, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	db := r.authDbConn(ctx)
	col := db.Collection("auth")

	result := new(auth.Credential)

	if err := col.FindOne(ctx, bson.M{"_id": utils.ConvertToObjectId(credentialId)}).Decode(result); err != nil {
		log.Printf("Error: FindOnePlayerCredential failed: %s", err.Error())
		return nil, errors.New("error: find one player credential failed")
	}

	return result, nil
}

func (r *authrepository) UpdateOnePlayerCredentail(pctx context.Context, credentialId string, req *auth.UpdateRefreshTokenReq) error {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()
	db := r.authDbConn(ctx)
	col := db.Collection("auth")

	_, err := col.UpdateOne(
		ctx,
		bson.M{"_id": utils.ConvertToObjectId(credentialId)},
		bson.M{
			"$set": bson.M{
				"player_id":     req.PlayerId,
				"access_token":  req.AccessToken,
				"refresh_token": req.RefreshToken,
				"updated_at":    req.UpdatedAt,
			},
		},
	)
	if err != nil {
		log.Printf("Error UpdateOnePlayerCredentail failed: %s", err.Error())
		return errors.New("error: update one player credential failed")
	}
	return nil
}

func (r *authrepository) DeleteOnePlayerCredentail(pctx context.Context, credentialId string) (int64, error) {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.authDbConn(ctx)
	col := db.Collection("auth")
	result, err := col.DeleteOne(ctx, bson.M{
		"_id": utils.ConvertToObjectId(credentialId),
	})
	if err != nil {
		log.Printf("Error: DeleteOnePlayerCredentail failed: %s", err.Error())
		return -1, errors.New("error: delete one player credential")
	}
	log.Printf("Result : %s", result)
	return result.DeletedCount, nil
}

func (r *authrepository) FindOneAccessToken(pctx context.Context, accessToken string) (*auth.Credential, error) {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.authDbConn(ctx)
	col := db.Collection("auth")

	credentail := new(auth.Credential)
	if err := col.FindOne(ctx, bson.M{
		"access_token": accessToken,
	}).Decode(credentail); err != nil {
		return nil, err
	}

	return credentail, nil
}

func (r *authrepository) RolesCount(pctx context.Context) (int64, error) {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.authDbConn(ctx)
	col := db.Collection("roles")

	count, err := col.CountDocuments(ctx, bson.M{})
	if err != nil {
		log.Printf("Error: RolesCount Faild %s", err)
		return -1, err
	}
	return count, nil
}
