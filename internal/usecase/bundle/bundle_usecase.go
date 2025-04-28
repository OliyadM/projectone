package bundle

import (
	"context"
	"errors"
	"fmt"

	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/bundle"
	"go.mongodb.org/mongo-driver/mongo"
)

type bundleUsecase struct {
	bundleRepo bundle.Repository
}

func NewBundleUsecase(bundleRepo bundle.Repository) bundle.Usecase {
	return &bundleUsecase{
		bundleRepo: bundleRepo,
	}
}

func (u *bundleUsecase) CreateBundle(ctx context.Context, supplierID string, b *bundle.Bundle) error {
	if b.SupplierID != supplierID {
		return errors.New("unauthorized: supplier ID mismatch")
	}
	return u.bundleRepo.CreateBundle(ctx, b)
}

func (u *bundleUsecase) ListBundles(ctx context.Context, supplierID string) ([]*bundle.Bundle, error) {
	bundles, err := u.bundleRepo.ListBundles(ctx, supplierID)
	if err != nil {
		return nil, err
	}
	return bundles, nil
}

func (u *bundleUsecase) DeleteBundle(ctx context.Context, supplierID string, bundleID string) error {
	// Fetch the bundle to verify ownership
	bundle, err := u.bundleRepo.GetBundleByID(ctx, bundleID)
	if err != nil {
		return err
	}
	if bundle == nil {
		return errors.New("bundle not found")
	}

	// Verify the bundle belongs to the supplier
	if bundle.SupplierID != supplierID {
		return errors.New("unauthorized: you can only delete your own bundles")
	}

	// Delete (deactivate) the bundle
	err = u.bundleRepo.DeleteBundle(ctx, bundleID)
	if err == mongo.ErrNoDocuments {
		return errors.New("bundle not found or already purchased")
	}
	if err != nil {
		return err
	}

	return nil
}

func (u *bundleUsecase) GetBundleByID(ctx context.Context, supplierID string, id string) (*bundle.Bundle, error) { // Added
	bundle, err := u.bundleRepo.GetBundleByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if bundle == nil {
		return nil, errors.New("bundle not found")
	}
	if bundle.SupplierID != supplierID {
		return nil, errors.New("unauthorized: you can only view your own bundles")
	}
	return bundle, nil
}

func (u *bundleUsecase) UpdateBundle(ctx context.Context, supplierID string, id string, updatedData map[string]interface{}) error { // Added
	// Fetch the bundle to verify ownership and status
	bundle, err := u.GetBundleByID(ctx, supplierID, id)
	if err != nil {
		return err
	}

	// Check if the bundle is editable (must be "available")
	if bundle.Status != "available" {
		return errors.New("cannot update bundle: bundle must be in 'available' status")
	}

	// Update the bundle in the repository
	return u.bundleRepo.UpdateBundle(ctx, id, updatedData)
}

func (uc *bundleUsecase) ListAvailableBundles(ctx context.Context) ([]*bundle.Bundle, error) {
	return uc.bundleRepo.ListAvailableBundles(ctx)
}
func (u *bundleUsecase) DecreaseRemainingItemCount(ctx context.Context, bundleID string) error {
	b, err := u.bundleRepo.GetBundleByID(ctx, bundleID)
	if err != nil {
		return err
	}
	if b.RemainingItemCount <= 0 {
		return errors.New("bundle is fully unpacked")
	}

	newCount := b.RemainingItemCount - 1
	return u.bundleRepo.UpdateBundle(ctx, bundleID, map[string]interface{}{
		"remaining_item_count": newCount,
	})
}
func (u *bundleUsecase) GetBundlePublicByID(ctx context.Context, bundleID string) (*bundle.Bundle, error) {
	return u.bundleRepo.GetBundleByID(ctx, bundleID)
}

func (u *bundleUsecase) GetBundleByTitle(ctx context.Context, title string) (*bundle.Bundle, error) {
	if title == "" {
		return nil, errors.New("title cannot be empty")
	}

	bundle, err := u.bundleRepo.GetBundleByTitle(ctx, title)
	if err != nil {
		return nil, fmt.Errorf("failed to get bundle by title: %w", err)
	}

	return bundle, nil
}
