package bundle

import "context"

type Usecase interface {
	CreateBundle(ctx context.Context, supplierID string, bundle *Bundle) error
	ListBundles(ctx context.Context, supplierID string) ([]*Bundle, error)
	DeleteBundle(ctx context.Context, supplierID string, bundleID string) error
	GetBundleByID(ctx context.Context, supplierID string, id string) (*Bundle, error)                         // Added
	UpdateBundle(ctx context.Context, supplierID string, id string, updatedData map[string]interface{}) error // Added
	ListAvailableBundles(ctx context.Context) ([]*Bundle, error)
	DecreaseRemainingItemCount(ctx context.Context, bundleID string) error
	GetBundlePublicByID(ctx context.Context, bundleID string) (*Bundle, error)
	GetBundleByTitle(ctx context.Context, title string) (*Bundle, error)
}
