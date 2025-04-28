package authinfra

import (
	"golang.org/x/crypto/bcrypt"
)

type passwordService struct{}

func NewPasswordService() *passwordService {
	return &passwordService{}
}

func (p *passwordService) HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashed), err
}

func (p *passwordService) CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
