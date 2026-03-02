package postgres

import (
	"Kairos/internal/models"
	"context"
)

func (s *AuthStorage) GetUserByLogin(ctx context.Context, login string) (models.User, error) {

	var user models.User
	err := s.db.QueryRowContext(ctx, `

    SELECT id, username, password
    FROM users
    WHERE username = $1`,

		login).Scan(&user.ID, &user.Login, &user.Password)
	if err != nil {
		return models.User{}, err
	}

	return user, nil

}
