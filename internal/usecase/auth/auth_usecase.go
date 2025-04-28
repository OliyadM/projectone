package auth

import (
	"context"
	"errors"
	"time"

	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/auth"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/user"

	// "github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type authUsecase struct {
	userRepo        user.Repository
	passwordService auth.PasswordService
	jwtService      auth.JWTService
}

func NewAuthUsecase(
	userRepo user.Repository,
	passwordService auth.PasswordService,
	jwtService auth.JWTService,
) auth.AuthUsecase {
	return &authUsecase{
		userRepo:        userRepo,
		passwordService: passwordService,
		jwtService:      jwtService,
	}
}

func (uc *authUsecase) Login(ctx context.Context, creds auth.LoginCredentials) (*auth.LoginResult, error) {
	u, err := uc.userRepo.FindUserByUsername(ctx, creds.Username)
	if err != nil || !uc.passwordService.CheckPasswordHash(creds.Password, u.Password) {
		return nil, errors.New("invalid username or password")
	}

	if creds.Role != "" && creds.Role != string(u.Role) {
		return nil, errors.New("access denied: user is not a " + creds.Role)
	}

	token, err := uc.jwtService.GenerateToken(u.ID, u.Username, string(u.Role))
	if err != nil {
		return nil, err
	}

	return &auth.LoginResult{
		Token:    token,
		ID:       u.ID,
		Username: u.Username,
		Role:     string(u.Role),
	}, nil
}
func (uc *authUsecase) Register(ctx context.Context, newUser user.User) (*auth.LoginResult, error) {
	// Check if user already exists
	existing, _ := uc.userRepo.FindUserByUsername(ctx, newUser.Username)
	if existing != nil {
		return nil, errors.New("user already exists")
	}

	// Hash password
	hashed, err := uc.passwordService.HashPassword(newUser.Password)
	if err != nil {
		return nil, err
	}

	// Generate ObjectID
	objectID := primitive.NewObjectID()
	newUser.ID = objectID.Hex()
	newUser.Password = hashed
	newUser.CreatedAt = time.Now()

	// Set trust score
	if newUser.Role == "supplier" || newUser.Role == "reseller" {
		newUser.TrustScore = 100
	}

	// Default role
	if newUser.Role == "" {
		newUser.Role = "consumer"
	}

	// Save user
	if err := uc.userRepo.CreateUser(ctx, &newUser); err != nil {
		return nil, err
	}

	// Generate token
	token, err := uc.jwtService.GenerateToken(newUser.ID, newUser.Username, string(newUser.Role))
	if err != nil {
		return nil, err
	}

	// Return structured login result
	return &auth.LoginResult{
		Token:    token,
		ID:       newUser.ID,
		Username: newUser.Username,
		Role:     string(newUser.Role),
	}, nil
}
