package postgres

import (
	"Kairos/internal/models"
	"context"
	"database/sql"
)

func (c *CoreStorage) CreateBooking(tx *sql.Tx, ctx context.Context, booking *models.Booking) (int64, error) {

	var bookingID int64
	err := tx.QueryRowContext(ctx, `

    INSERT INTO bookings (user_id, event_id, status, created_at, expires_at)
    VALUES ($1, $2, $3, $4, $5)
    RETURNING id`,

		booking.UserID, booking.EventID, booking.Status, booking.CreatedAt, booking.ExpiresAt).Scan(&bookingID)
	if err != nil {
		return 0, err
	}

	return bookingID, nil

}
