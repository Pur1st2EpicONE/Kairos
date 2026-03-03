package impl

import (
	"Kairos/internal/models"
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func (a *AuthService) CreateUser(ctx context.Context, user models.User) (int64, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, fmt.Errorf("failed to hash user password: %w", err)
	}
	user.Password = string(passwordHash)
	return a.storage.CreateUser(ctx, user)
}
