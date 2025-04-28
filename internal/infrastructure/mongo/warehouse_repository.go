package mongo

import (
	"context"

	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/warehouse"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoRepository struct {
	collection *mongo.Collection
}

func NewMongoWarehouseRepository(db *mongo.Database) warehouse.Repository {
	return &mongoRepository{
		collection: db.Collection("warehouses"),
	}
}

func (r *mongoRepository) AddItem(ctx context.Context, item *warehouse.WarehouseItem) error {
	_, err := r.collection.InsertOne(ctx, item)
	return err
}

func (r *mongoRepository) GetItemsByReseller(ctx context.Context, resellerID string) ([]*warehouse.WarehouseItem, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"reseller_id": resellerID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var items []*warehouse.WarehouseItem
	if err = cursor.All(ctx, &items); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *mongoRepository) GetItemsByBundle(ctx context.Context, bundleID string) ([]*warehouse.WarehouseItem, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"bundle_id": bundleID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var items []*warehouse.WarehouseItem
	if err = cursor.All(ctx, &items); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *mongoRepository) MarkItemAsListed(ctx context.Context, itemID string) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": itemID}, bson.M{"$set": bson.M{"status": "listed"}})
	return err
}

func (r *mongoRepository) MarkItemAsSkipped(ctx context.Context, itemID string) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": itemID}, bson.M{"$set": bson.M{"status": "skipped"}})
	return err
}

func (r *mongoRepository) DeleteItem(ctx context.Context, itemID string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": itemID})
	return err
}
func (r *mongoRepository) HasResellerReceivedBundle(ctx context.Context, resellerID string, bundleID string) (bool, error) {
	filter := bson.M{
		"reseller_id": resellerID,
		"bundle_id":   bundleID,
		"status":      bson.M{"$in": []string{"arrived", "listed"}},
	}
	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
func (r *mongoRepository) CountByStatus(ctx context.Context, status string) (int, error) {
	filter := bson.M{"status": status}
	count, err := r.collection.CountDocuments(ctx, filter)
	return int(count), err
}
