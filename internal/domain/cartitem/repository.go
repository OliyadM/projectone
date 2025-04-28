package cartitem

import "context"

type Repository interface {
	CreateCartItem(ctx context.Context, item *CartItem) error

	GetCartItems(ctx context.Context, userID string) ([]*CartItem, error)

	DeleteCartItem(ctx context.Context, userID string, listingID string) error

	ClearCart(ctx context.Context, userID string) error
}
