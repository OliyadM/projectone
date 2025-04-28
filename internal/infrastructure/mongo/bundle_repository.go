package mongo

import (
	"context"
	"errors" // Added
	"fmt"

	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/bundle"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type BundleRepository struct {
	collection *mongo.Collection
}

func NewBundleRepository(db *mongo.Database) *BundleRepository {
	return &BundleRepository{
		collection: db.Collection("bundles"),
	}
}

func (r *BundleRepository) CreateBundle(ctx context.Context, b *bundle.Bundle) error {
	_, err := r.collection.InsertOne(ctx, b)
	return err
}

func (r *BundleRepository) GetBundleByID(ctx context.Context, id string) (*bundle.Bundle, error) {
	var bundle bundle.Bundle
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&bundle)
	if err == mongo.ErrNoDocuments {
		return nil, errors.New("bundle not found") // Updated: Return a specific error message
	}
	if err != nil {
		return nil, err
	}
	return &bundle, nil
}

func (r *BundleRepository) ListBundles(ctx context.Context, supplierID string) ([]*bundle.Bundle, error) {
	var bundles []*bundle.Bundle
	cursor, err := r.collection.Find(ctx, bson.M{"supplierid": supplierID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var b bundle.Bundle
		if err := cursor.Decode(&b); err != nil {
			return nil, err
		}
		bundles = append(bundles, &b)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return bundles, nil
}

func (r *BundleRepository) ListAvailableBundles(ctx context.Context) ([]*bundle.Bundle, error) {
	var bundles []*bundle.Bundle
	cursor, err := r.collection.Find(ctx, bson.M{"status": "available"})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var b bundle.Bundle
		if err := cursor.Decode(&b); err != nil {
			return nil, err
		}
		bundles = append(bundles, &b)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return bundles, nil
}

func (r *BundleRepository) ListPurchasedByReseller(ctx context.Context, resellerID string) ([]*bundle.Bundle, error) {
	fmt.Printf("🔍 Fetching purchased bundles for reseller: %s\n", resellerID)
	
	var bundles []*bundle.Bundle
	cursor, err := r.collection.Find(ctx, bson.M{"resellerid": resellerID, "status": "purchased"})
	if err != nil {
		fmt.Printf("❌ Error querying bundles: %v\n", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var b bundle.Bundle
		if err := cursor.Decode(&b); err != nil {
			fmt.Printf("❌ Error decoding bundle: %v\n", err)
			return nil, err
		}
		bundles = append(bundles, &b)
	}

	if err := cursor.Err(); err != nil {
		fmt.Printf("❌ Cursor error: %v\n", err)
		return nil, err
	}

	fmt.Printf("✅ Found %d purchased bundles for reseller %s\n", len(bundles), resellerID)
	return bundles, nil
}

func (r *BundleRepository) UpdateBundleStatus(ctx context.Context, id string, status string) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"status": status}})
	return err
}

func (r *BundleRepository) MarkAsPurchased(ctx context.Context, bundleID string, resellerID string) error {
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": bundleID},
		bson.M{"$set": bson.M{"status": "purchased", "resellerid": resellerID}},
	)
	return err
}

func (r *BundleRepository) DeleteBundle(ctx context.Context, bundleID string) error {
	// Update the bundle's status to "deactivated"
	result, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": bundleID, "status": "available"},
		bson.M{"$set": bson.M{"status": "deactivated"}},
	)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

func (r *BundleRepository) UpdateBundle(ctx context.Context, id string, updatedData map[string]interface{}) error { // Added
	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": updatedData})
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("bundle not found")
	}
	return nil
}
func (r *BundleRepository) DecreaseBundleQuantity(ctx context.Context, bundleID string) error {
	update := bson.M{
		"$inc": bson.M{"quantity": -1},
	}
	filter := bson.M{
		"_id":      bundleID,
		"quantity": bson.M{"$gt": 0}, // Only decrease if quantity > 0
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("bundle not found or quantity already zero")
	}
	return nil
}
func (r *BundleRepository) CountBundles(ctx context.Context) (int, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{})
	return int(count), err
}

func (r *BundleRepository) GetBundleByTitle(ctx context.Context, title string) (*bundle.Bundle, error) {
	var bundle bundle.Bundle
	err := r.collection.FindOne(ctx, bson.M{"title": title}).Decode(&bundle)
	if err == mongo.ErrNoDocuments {
		return nil, errors.New("bundle not found")
	}
	if err != nil {
		return nil, err
	}
	return &bundle, nil
}
