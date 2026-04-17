package impl

import (
	"Kairos/internal/errs"
	"Kairos/internal/models"
	"context"

	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

const bcryptMaxLen = 72 // Maximum password length accepted by bcrypt

// CreateUser registers a new user.
// It validates password length, hashes the password with bcrypt,
// and stores the user in the database. Returns the generated user ID
// or an error (e.g., ErrPasswordTooLong, ErrUserAlreadyExists).
func (a *AuthService) CreateUser(ctx context.Context, user models.User) (int64, error) {

	if len(user.Password) > bcryptMaxLen {
		return 0, errs.ErrPasswordTooLong
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		a.logger.LogError("service — failed to hash user password", err, "layer", "service.impl")
		return 0, err
	}
	user.Password = string(passwordHash)

	userID, err := a.storage.CreateUser(ctx, user)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return 0, errs.ErrUserAlreadyExists
		}
		a.logger.LogError("service — failed to create new user", err, "layer", "service.impl")
		return 0, err
	}

	return userID, nil

}
