package user

import "context"

type Usecase interface {
	GetByID(ctx context.Context, id string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	ListByRole(ctx context.Context, role Role) ([]*User, error)
	Update(ctx context.Context, id string, updates map[string]interface{}) error
	Delete(ctx context.Context, id string) error
	GetBlacklistedUsers(ctx context.Context) ([]*User, error)
}
