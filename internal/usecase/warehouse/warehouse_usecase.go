package warehouse_usecase

import (
	"context"

	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/bundle"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/warehouse"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/models"
)

type WarehouseUseCase interface {
	GetWarehouseItems(ctx context.Context, resellerID string) ([]*warehouse.WarehouseItem, error)
}

type warehouseUseCaseImpl struct {
	warehouseRepo warehouse.Repository
	bundleRepo    bundle.Repository
}

func NewWarehouseUseCase(warehouseRepo warehouse.Repository, bundleRepo bundle.Repository) warehouse.WarehouseUseCase {
	return &warehouseUseCaseImpl{
		warehouseRepo: warehouseRepo,
		bundleRepo:    bundleRepo,
	}
}
func (uc *warehouseUseCaseImpl) GetWarehouseItems(ctx context.Context, resellerID string) ([]*models.WarehouseItemResponse, error) {
	// Get the warehouse items by reseller
	items, err := uc.warehouseRepo.GetItemsByReseller(ctx, resellerID)
	if err != nil {
		return nil, err
	}

	// Initialize responses as empty slice
	responses := make([]*models.WarehouseItemResponse, 0)
	for _, item := range items {
		bundle, err := uc.bundleRepo.GetBundleByID(ctx, item.BundleID)
		if err != nil || bundle == nil {
			continue // skip if bundle not found
		}

		responses = append(responses, &models.WarehouseItemResponse{
			ID:             item.ID,
			ResellerID:     item.ResellerID,
			BundleID:       item.BundleID,
			Status:         item.Status,
			CreatedAt:      item.CreatedAt,
			Title:          bundle.Title,
			SampleImage:    bundle.SampleImage,
			DeclaredRating: float64(bundle.DeclaredRating),
			RemainingItems: bundle.RemainingItemCount,
		})
	}

	return responses, nil
}
