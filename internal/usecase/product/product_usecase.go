package productusecase

import (
	"context"
	"errors"

	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/bundle"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/product"
)

type productUsecase struct {
	repo       product.Repository
	bundleRepo bundle.Repository
}

func NewProductUsecase(repo product.Repository, bundleRepo bundle.Repository) product.Usecase {
	return &productUsecase{
		repo:       repo,
		bundleRepo: bundleRepo,
	}
}
func (uc *productUsecase) AddProduct(ctx context.Context, p *product.Product) error {
	if p.BundleID != "" {
		// Fetch bundle
		b, err := uc.bundleRepo.GetBundleByID(ctx, p.BundleID)
		if err != nil {
			return err
		}

		// Check available quantity
		if b.Quantity <= 0 {
			return errors.New("bundle is out of stock")
		}

		// Add product first
		if err := uc.repo.AddProduct(ctx, p); err != nil {
			return err
		}

		// Decrement bundle quantity
		updates := map[string]interface{}{
			"quantity": b.Quantity - 1,
		}
		return uc.bundleRepo.UpdateBundle(ctx, b.ID, updates)
	}

	// If not tied to a bundle, allow adding normally
	return uc.repo.AddProduct(ctx, p)
}

func (uc *productUsecase) GetProductByID(ctx context.Context, id string) (*product.Product, error) {
	return uc.repo.GetProductByID(ctx, id)
}

func (uc *productUsecase) GetProductByTitle(ctx context.Context, title string) (*product.Product, error) {
	if title == "" {
		return nil, errors.New("title cannot be empty")
	}
	return uc.repo.GetProductByTitle(ctx, title)
}

func (uc *productUsecase) ListProductsByReseller(ctx context.Context, resellerID string, page, limit int) ([]*product.Product, error) {
	return uc.repo.ListProductsByReseller(ctx, resellerID, page, limit)
}

func (uc *productUsecase) ListAvailableProducts(ctx context.Context, page, limit int) ([]*product.Product, error) {
	return uc.repo.ListAvailableProducts(ctx, page, limit)
}

func (uc *productUsecase) DeleteProduct(ctx context.Context, id string) error {
	return uc.repo.DeleteProduct(ctx, id)
}

func (uc *productUsecase) UpdateProduct(ctx context.Context, id string, updates map[string]interface{}) error {
	return uc.repo.UpdateProduct(ctx, id, updates)
}
