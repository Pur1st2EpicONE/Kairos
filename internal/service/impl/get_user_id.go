package impl

import (
	"Kairos/internal/errs"
	"Kairos/internal/models"
	"context"
	"database/sql"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func (a *AuthService) GetUserId(ctx context.Context, user models.User) (int, error) {

	if user.Login == "" || user.Password == "" {
		fmt.Println("empty login or pass")
		return 0, errs.ErrInvalidCredentials
	}

	user2, err := a.storage.GetUserByLogin(ctx, user.Login)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			fmt.Println(err)
			return 0, errs.ErrInvalidCredentials
		}
		return 0, fmt.Errorf("failed to get user: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user2.Password), []byte(user.Password)); err != nil {
		fmt.Println("CompareHashAndPassword", err)
		return 0, errs.ErrInvalidCredentials
	}

	return user.ID, nil
}
