package postgres

import (
	"Kairos/internal/errs"
	"Kairos/internal/models"
	"context"

	"github.com/lib/pq"
)

func (s *AuthStorage) CreateUser(ctx context.Context, user models.User) (int64, error) {

	var insertedID int64
	err := s.db.QueryRowContext(ctx, `

    INSERT INTO users (username, password)
    VALUES ($1, $2)
    RETURNING id;`,

		user.Login, user.Password).Scan(&insertedID)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return 0, errs.ErrUserAlreadyExists
		}
		return 0, err
	}

	return insertedID, nil

}
