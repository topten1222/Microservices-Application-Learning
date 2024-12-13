package paymentRepository

import (
	"context"
	"errors"
	"log"
	"time"

	itemPb "github.com/topten1222/hello_sekai/modules/item/itemPb"
	"github.com/topten1222/hello_sekai/modules/models"
	"github.com/topten1222/hello_sekai/pkg/grpccon"
	"github.com/topten1222/hello_sekai/pkg/jwtauth"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	PaymentRepositoryService interface {
		FindItemsInIds(context.Context, string, *itemPb.FindItemsInIdsReq) (*itemPb.FindItemsInIdsRes, error)
		GetOffset(context.Context) (int64, error)
		UpserOffset(context.Context, int64) error
	}

	paymentRepository struct {
		db *mongo.Client
	}
)

func NewPaymentRepository(db *mongo.Client) PaymentRepositoryService {
	return &paymentRepository{db: db}
}

func (r *paymentRepository) paymentDbConn(pctx context.Context) *mongo.Database {
	return r.db.Database("payment_db")
}

func (r *paymentRepository) FindItemsInIds(pctx context.Context, grpcUrl string, req *itemPb.FindItemsInIdsReq) (*itemPb.FindItemsInIdsRes, error) {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	jwtauth.SetApiKeyInContext(&ctx)
	conn, err := grpccon.NewGrpcClient(grpcUrl)
	if err != nil {
		return nil, errors.New("erro: grpc connection faild")
	}

	result, err := conn.Item().FindItemsInIds(ctx, req)
	if err != nil {
		return nil, errors.New("error: item not found")
	}
	if result == nil {
		return nil, errors.New("error: item not found")
	}
	if len(result.Items) == 0 {
		return nil, errors.New("error: item not found")
	}
	return result, nil
}

func (r *paymentRepository) GetOffset(pctx context.Context) (int64, error) {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.paymentDbConn(ctx)
	col := db.Collection("payment_queue")

	result := new(models.KafkaOffset)
	if err := col.FindOne(ctx, bson.M{}).Decode(result); err != nil {
		return -1, errors.New("error: get offset faild")
	}

	return result.Offset, nil
}

func (r *paymentRepository) UpserOffset(pctx context.Context, offset int64) error {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.paymentDbConn(ctx)
	col := db.Collection("payment_queue")

	result, err := col.UpdateOne(ctx, bson.M{}, bson.M{"$set": bson.M{"offset": offset}}, options.Update().SetUpsert(true))

	if err != nil {
		return errors.New("error: cannot update upsertoffset")
	}
	log.Printf("upsert offset: %v", result)

	return nil

}
