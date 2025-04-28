package auth

import (
	"context"

	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/user"
	"github.com/golang-jwt/jwt/v5"
)

type PasswordService interface {
	HashPassword(password string) (string, error)
	CheckPasswordHash(password, hash string) bool
}

type JWTService interface {
	GenerateToken(userID, username, role string) (string, error)
	ParseToken(token string) (*jwt.Token, jwt.MapClaims, error)
}
type AuthUsecase interface {
	Login(ctx context.Context, creds LoginCredentials) (*LoginResult, error)
	Register(ctx context.Context, user user.User) (*LoginResult, error)
}
