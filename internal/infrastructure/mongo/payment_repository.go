package mongo

import (
	"context"
	"time"

	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/payment"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoPaymentRepository struct {
	collection *mongo.Collection
}

func NewMongoPaymentRepository(db *mongo.Database) payment.Repository {
	return &mongoPaymentRepository{
		collection: db.Collection("payments"),
	}
}

func (repo *mongoPaymentRepository) RecordPayment(ctx context.Context, p *payment.Payment) error {
	if p.ID == "" {
		p.ID = primitive.NewObjectID().Hex()
	}
	if p.CreatedAt == "" {
		p.CreatedAt = time.Now().Format(time.RFC3339)
	}
	_, err := repo.collection.InsertOne(ctx, p)
	return err
}

func (repo *mongoPaymentRepository) GetPaymentsByUser(ctx context.Context, userID string) ([]*payment.Payment, error) {
	filter := bson.M{"fromuserid": userID}
	cursor, err := repo.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var payments []*payment.Payment
	if err = cursor.All(ctx, &payments); err != nil {
		return nil, err
	}
	return payments, nil
}

func (repo *mongoPaymentRepository) GetPaymentsByType(ctx context.Context, userID string, pType payment.PaymentType) ([]*payment.Payment, error) {
	filter := bson.M{
		"fromuserid": userID,
		"type":       pType,
	}
	cursor, err := repo.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var payments []*payment.Payment
	if err = cursor.All(ctx, &payments); err != nil {
		return nil, err
	}
	return payments, nil
}
func (repo *mongoPaymentRepository) GetAllPlatformFees(ctx context.Context) (float64, float64, error) {
	pipeline := mongo.Pipeline{
		bson.D{
			{Key: "$match", Value: bson.D{
				{Key: "status", Value: bson.D{
					{Key: "$in", Value: bson.A{"paid", "Paid"}},
				}},
			}},
		},
		bson.D{
			{Key: "$group", Value: bson.D{
				{Key: "_id", Value: nil},
				{Key: "totalSales", Value: bson.D{{Key: "$sum", Value: "$amount"}}},
				{Key: "platformFees", Value: bson.D{{Key: "$sum", Value: "$platformfee"}}},
			}},
		},
	}

	cursor, err := repo.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, 0, err
	}
	defer cursor.Close(ctx)

	var result []struct {
		TotalSales   float64 `bson:"totalSales"`
		PlatformFees float64 `bson:"platformFees"`
	}

	if err := cursor.All(ctx, &result); err != nil || len(result) == 0 {
		return 0, 0, err
	}

	return result[0].TotalSales, result[0].PlatformFees, nil
}
