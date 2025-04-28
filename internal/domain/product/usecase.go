package product

import "context"

type Usecase interface {
	AddProduct(ctx context.Context, p *Product) error
	GetProductByID(ctx context.Context, id string) (*Product, error)
	GetProductByTitle(ctx context.Context, title string) (*Product, error)
	ListProductsByReseller(ctx context.Context, resellerID string, page, limit int) ([]*Product, error)
	ListAvailableProducts(ctx context.Context, page, limit int) ([]*Product, error)
	DeleteProduct(ctx context.Context, id string) error
	UpdateProduct(ctx context.Context, id string, updates map[string]interface{}) error
}
