package postgres

import (
	"context"
	"database/sql"
)

func (c *CoreStorage) CancelBooking(tx *sql.Tx, ctx context.Context, bookingID int64) (int64, error) {

	var eventID int64
	err := tx.QueryRowContext(ctx, `

	UPDATE bookings
	SET status = 'expired'
	WHERE id = $1 AND status = 'pending'
	RETURNING event_id`,

		bookingID).Scan(&eventID)
	if err != nil {
		return 0, err
	}

	return eventID, nil

}
