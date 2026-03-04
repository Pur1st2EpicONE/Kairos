package postgres

import (
	"Kairos/internal/models"
	"context"
	"database/sql"
)

func (c *CoreStorage) GetBookingForConfirm(tx *sql.Tx, ctx context.Context, userID int64, eventUUID string) (*models.Booking, error) {

	booking := new(models.Booking)

	row := tx.QueryRowContext(ctx, `

	SELECT b.id, b.user_id, b.event_id, b.status, b.created_at, b.expires_at
	FROM bookings b
	JOIN events e 
	ON b.event_id = e.id
	WHERE b.user_id = $1 
	AND e.uuid = $2 
	ORDER BY b.id DESC
	LIMIT 1
	FOR UPDATE OF b`,

		userID, eventUUID)
	err := row.Scan(
		&booking.ID,
		&booking.UserID,
		&booking.EventID,
		&booking.Status,
		&booking.CreatedAt,
		&booking.ExpiresAt,
	)
	if err != nil {
		return nil, err
	}

	return booking, nil

}
