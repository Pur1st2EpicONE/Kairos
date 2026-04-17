package postgres

import (
	"Kairos/internal/models"
	"context"
	"fmt"

	"github.com/wb-go/wbf/retry"
)

// GetUserByLogin retrieves a user by their unique login name.
// Uses the configured retry strategy. Returns the user or an error
// (e.g., sql.ErrNoRows if login does not exist).
func (s *AuthStorage) GetUserByLogin(ctx context.Context, login string) (models.User, error) {

	row, err := s.db.QueryRowWithRetry(ctx, retry.Strategy(s.config.QueryRetryStrategy), `

    SELECT id, username, password
    FROM users
    WHERE username = $1`,

		login)
	if err != nil {
		return models.User{}, fmt.Errorf("failed to query row %w", err)
	}

	var user models.User
	err = row.Scan(&user.ID, &user.Login, &user.Password)
	if err != nil {
		return models.User{}, err
	}

	return user, nil

}
