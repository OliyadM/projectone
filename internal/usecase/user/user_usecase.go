package userusecase

import (
	"context"

	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/user"
)

type userUsecase struct {
	repo user.Repository
}

func NewUserUsecase(repo user.Repository) user.Usecase {
	return &userUsecase{repo: repo}
}

func (uc *userUsecase) GetByID(ctx context.Context, id string) (*user.User, error) {
	return uc.repo.GetByID(ctx, id)
}

func (uc *userUsecase) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	return uc.repo.GetUserByEmail(ctx, email)
}

func (uc *userUsecase) ListByRole(ctx context.Context, role user.Role) ([]*user.User, error) {
	return uc.repo.ListUsersByRole(ctx, role)
}

func (uc *userUsecase) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	return uc.repo.UpdateUser(ctx, id, updates)
}

func (uc *userUsecase) Delete(ctx context.Context, id string) error {
	return uc.repo.DeleteUser(ctx, id)
}
func (u *userUsecase) GetBlacklistedUsers(ctx context.Context) ([]*user.User, error) {
	return u.repo.GetBlacklistedUsers(ctx)
}
