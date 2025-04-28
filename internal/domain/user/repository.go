package user

import "context"

type Repository interface {
	CreateUser(ctx context.Context, u *User) error
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetByID(ctx context.Context, id string) (*User, error)
	ListUsersByRole(ctx context.Context, role Role) ([]*User, error)
	UpdateUser(ctx context.Context, id string, updates map[string]interface{}) error
	DeleteUser(ctx context.Context, id string) error
	FindUserByUsername(ctx context.Context, username string) (*User, error)
	UpdateTrustData(ctx context.Context, user *User) error
	GetBlacklistedUsers(ctx context.Context) ([]*User, error)
	CountActiveUsers(ctx context.Context) (int, error)
}
