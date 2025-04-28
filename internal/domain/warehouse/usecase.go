package warehouse

import (
	"context"

	"github.com/Zeamanuel-Admasu/afro-vintage-backend/models"
)

type WarehouseUseCase interface {
	GetWarehouseItems(ctx context.Context, resellerID string) ([]*models.WarehouseItemResponse, error)
}
