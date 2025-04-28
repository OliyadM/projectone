package authinfra

import (
	"fmt"
	"time"

	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/auth"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type jwtService struct {
	secretKey string
}

func NewJWTService(secretKey string) *jwtService {
	return &jwtService{secretKey: secretKey}
}

func (s *jwtService) GenerateToken(userID, username, role string) (string, error) {
	if _, err := primitive.ObjectIDFromHex(userID); err != nil {
		return "", fmt.Errorf("invalid user ID format: %w", err)
	}

	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"role":     role,
		"exp":      time.Now().Add(72 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secretKey))
}

func (s *jwtService) ParseToken(tokenStr string) (*jwt.Token, jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(s.secretKey), nil
	})
	if err != nil || !token.Valid {
		return nil, nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, nil, jwt.ErrInvalidKeyType
	}

	return token, claims, nil
}

// Ensure jwtService implements auth.JWTService
var _ auth.JWTService = (*jwtService)(nil)
