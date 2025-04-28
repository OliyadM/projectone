package mongo

import (
	"context"
	"errors"

	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/cartitem"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type CartItemRepository struct {
	collection *mongo.Collection
}

func NewCartItemRepository(db *mongo.Database) *CartItemRepository {
	return &CartItemRepository{
		collection: db.Collection("cartitems"),
	}
}

// CreateCartItem adds an item to the user's cart.
// It checks for duplicates before inserting.
func (r *CartItemRepository) CreateCartItem(ctx context.Context, item *cartitem.CartItem) error {
	// Check if the item already exists in the user's cart.
	var existing cartitem.CartItem
	err := r.collection.FindOne(ctx, bson.M{
		"userid":    item.UserID,
		"listingid": item.ListingID,
	}).Decode(&existing)
	if err == nil {
		return errors.New("item already in cart")
	}
	if err != mongo.ErrNoDocuments {
		return err
	}

	// Insert the new cart item.
	_, err = r.collection.InsertOne(ctx, item)
	return err
}

// GetCartItems retrieves all cart items for the given user.
func (r *CartItemRepository) GetCartItems(ctx context.Context, userID string) ([]*cartitem.CartItem, error) {
	var items []*cartitem.CartItem
	cursor, err := r.collection.Find(ctx, bson.M{"userid": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var item cartitem.CartItem
		if err := cursor.Decode(&item); err != nil {
			return nil, err
		}
		items = append(items, &item)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

// DeleteCartItem removes a specific item from the user's cart.
func (r *CartItemRepository) DeleteCartItem(ctx context.Context, userID string, listingID string) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{
		"userid":    userID,
		"listingid": listingID,
	})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

// ClearCart removes all items from the user's cart.
func (r *CartItemRepository) ClearCart(ctx context.Context, userID string) error {
	_, err := r.collection.DeleteMany(ctx, bson.M{"userid": userID})
	return err
}
