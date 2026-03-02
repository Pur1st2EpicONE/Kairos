package impl

import (
	"Kairos/internal/models"
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func (a *AuthService) CreateUser(ctx context.Context, user models.User) (int, error) {
	passwordHash, err := hashPassword(user.Password)
	if err != nil {
		return 0, fmt.Errorf("failed to hash user password: %w", err)
	}
	user.Password = string(passwordHash)
	return a.storage.CreateUser(ctx, user)
}

func hashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}
