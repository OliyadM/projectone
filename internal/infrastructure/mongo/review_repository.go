package mongo

import (
	"context"
	
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/review"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ReviewRepository struct {
	collection *mongo.Collection
}

func NewReviewRepository(db *mongo.Database) *ReviewRepository {
	return &ReviewRepository{
		collection: db.Collection("reviews"),
	}
}

func (r *ReviewRepository) CreateReview(ctx context.Context, rev *review.Review) error {
	_, err := r.collection.InsertOne(ctx, rev)
	return err
}

func (r *ReviewRepository) GetReviewByUserAndProduct(ctx context.Context, userID, productID string) (*review.Review, error) {
	var rev review.Review
	err := r.collection.FindOne(ctx, bson.M{"user_id": userID, "product_id": productID}).Decode(&rev)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &rev, nil
}

func (r *ReviewRepository) GetReviewsByReseller(ctx context.Context, resellerID string) ([]*review.Review, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"reseller_id": resellerID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var reviews []*review.Review
	if err = cursor.All(ctx, &reviews); err != nil {
		return nil, err
	}

	return reviews, nil
}
