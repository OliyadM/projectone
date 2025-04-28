package mongo

import (
	"context"
	"fmt"
	"time"

	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/product"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoProductRepository struct {
	collection *mongo.Collection
}

func NewMongoProductRepository(db *mongo.Database) product.Repository {
	return &mongoProductRepository{
		collection: db.Collection("products"),
	}
}

func (r *mongoProductRepository) AddProduct(ctx context.Context, p *product.Product) error {
	p.CreatedAt = time.Now().Format(time.RFC3339)
	_, err := r.collection.InsertOne(ctx, p)
	return err
}

func (r *mongoProductRepository) GetProductByID(ctx context.Context, id string) (*product.Product, error) {
	var p product.Product
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&p)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *mongoProductRepository) ListProductsByReseller(ctx context.Context, resellerID string, page, limit int) ([]*product.Product, error) {
	var products []*product.Product
	skip := (page - 1) * limit
	opts := options.Find().SetSkip(int64(skip)).SetLimit(int64(limit))

	resellerObjectID, err := primitive.ObjectIDFromHex(resellerID)
	if err != nil {
		return nil, fmt.Errorf("invalid reseller ID: %w", err)
	}

	cursor, err := r.collection.Find(ctx, bson.M{"reseller_id": resellerObjectID}, opts)
	if err != nil {
		return nil, fmt.Errorf("database query failed: %w", err)
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var p product.Product
		if err := cursor.Decode(&p); err != nil {
			return nil, fmt.Errorf("failed to decode product: %w", err)
		}
		products = append(products, &p)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return products, nil
}

func (r *mongoProductRepository) GetProductByTitle(ctx context.Context, title string) (*product.Product, error) {
	var p product.Product
	err := r.collection.FindOne(ctx, bson.M{"title": title}).Decode(&p)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("product not found: %w", err)
		}
		return nil, fmt.Errorf("failed to find product: %w", err)
	}
	return &p, nil
}

func (r *mongoProductRepository) ListAvailableProducts(ctx context.Context, page, limit int) ([]*product.Product, error) {
	var products []*product.Product
	skip := (page - 1) * limit
	opts := options.Find().SetSkip(int64(skip)).SetLimit(int64(limit))

	cursor, err := r.collection.Find(ctx, bson.M{"status": "available"}, opts)
	if err != nil {
		return nil, fmt.Errorf("database query failed: %w", err)
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var p product.Product
		if err := cursor.Decode(&p); err != nil {
			return nil, fmt.Errorf("failed to decode product: %w", err)
		}
		products = append(products, &p)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return products, nil
}

func (r *mongoProductRepository) DeleteProduct(ctx context.Context, id string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *mongoProductRepository) UpdateProduct(ctx context.Context, id string, updates map[string]interface{}) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": updates})
	return err
}
func (r *mongoProductRepository) GetProductsByBundleID(ctx context.Context, bundleID string) ([]*product.Product, error) {
	var products []*product.Product

	cursor, err := r.collection.Find(ctx, bson.M{"bundle_id": bundleID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var p product.Product
		if err := cursor.Decode(&p); err != nil {
			return nil, err
		}
		products = append(products, &p)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

// GetSoldProductsByReseller returns all products that are marked as sold for a specific reseller
func (r *mongoProductRepository) GetSoldProductsByReseller(ctx context.Context, resellerID string) ([]*product.Product, error) {
	collection := r.collection
	
	// Convert string ID to ObjectID
	objID, err := primitive.ObjectIDFromHex(resellerID)
	if err != nil {
		return nil, fmt.Errorf("invalid reseller ID: %w", err)
	}

	// Find all products that are sold and belong to this reseller
	filter := bson.M{
		"reseller_id": objID,
		"status": "sold",
	}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find sold products: %w", err)
	}
	defer cursor.Close(ctx)

	var products []*product.Product
	if err := cursor.All(ctx, &products); err != nil {
		return nil, fmt.Errorf("failed to decode products: %w", err)
	}

	return products, nil
}
